package jnumber

import (
	"math/big"
	"strings"
)

// ParseBigInt returns the integer represented by the given japanese numerals.
func ParseBigInt(s string) (*big.Int, error) {
	if s == "" {
		return nil, ErrEmpty
	}
	initBigIntsOnce.Do(initBigInts)
	isNegative := strings.HasPrefix(s, "-")
	if isNegative {
		s = s[1:]
	}
	parser := bigIntParser{
		sum: big.NewInt(0),
	}
	err := parser.parse(s)
	if err != nil {
		return nil, err
	}
	if isNegative {
		return parser.sum.Neg(parser.sum), nil
	}
	return parser.sum, nil
}

// parser contains the state for the parsing process. The segment could be a slice
// for easier code but this way we can avoid an allocation.
type bigIntParser struct {
	// current sum
	sum *big.Int
	// holds digit until we know what to do with them
	segment [16]big.Int
	// position in the segment array for the next digit
	segmentIndex int
}

// append adds a digit to the current segment.
func (p *bigIntParser) append(digit *big.Int) {
	p.segment[p.segmentIndex].Set(digit)
	p.segmentIndex++
}

// clearSegment clears the current segment to start parsing a new segment.
func (p *bigIntParser) clearSegment() {
	p.segmentIndex = 0
}

// push integrates a new digit to the current segment.
func (p *bigIntParser) push(digit *big.Int) error {
	if p.segmentIndex == 0 {
		p.append(digit)
		return nil
	}
	lastIndex := p.segmentIndex - 1
	lastDigit := &p.segment[lastIndex]
	if lastDigit.Cmp(digit) < 0 {
		if digit.Cmp(&b十) < 0 {
			return ErrInvalidSequence
		} else if (digit.Cmp(&b百) == 0 || digit.Cmp(&b千) == 0) && lastDigit.Cmp(&b十) >= 0 {
			return ErrInvalidSequence
		}
		// Example: 二十 -> 2 * 10 -> 20
		lastDigit.Mul(lastDigit, digit)
		return nil
	} else if lastDigit.Cmp(digit) > 0 && p.segmentIndex < len(p.segment) && (digit.Cmp(&b十) >= 0 || lastDigit.Cmp(&b十) >= 0) {
		// Example: 十一 -> 20 + 1 -> 21
		// We don't know if the next digit needs be multiplied with this digit or
		// added to the last. Store for handling at the end of the segment.
		p.append(digit)
		return nil
	}
	return ErrInvalidSequence
}

// endSegmentWith is a combination of push() and endSegment() for digits >= 万 (10000).
func (p *bigIntParser) endSegmentWith(digit *big.Int) error {
	if p.segmentIndex == 0 {
		return ErrInvalidSequence
	}
	var multiplierSum big.Int
	var lastDigit *big.Int
	for i := 0; i < p.segmentIndex; i++ {
		segmentDigit := &p.segment[i]
		if lastDigit != nil && segmentDigit.Cmp(lastDigit) >= 0 {
			return ErrInvalidSequence
		}
		multiplierSum.Add(&multiplierSum, segmentDigit)
		lastDigit = segmentDigit
	}
	if multiplierSum.Cmp(&b万) >= 0 || multiplierSum.Cmp(digit) >= 0 {
		return ErrInvalidSequence
	}
	multiplierSum.Mul(&multiplierSum, digit)
	p.sum.Add(p.sum, &multiplierSum)
	p.clearSegment()
	return nil
}

func (p *bigIntParser) endSegment() error {
	if p.segmentIndex == 0 {
		return nil
	}
	var segmentSum big.Int
	var lastDigit *big.Int
	for i := 0; i < p.segmentIndex; i++ {
		segmentDigit := &p.segment[i]
		if lastDigit != nil && segmentDigit.Cmp(lastDigit) >= 0 {
			return ErrInvalidSequence
		}
		segmentSum.Add(&segmentSum, segmentDigit)
		lastDigit = segmentDigit
	}
	if segmentSum.Cmp(&b万) >= 0 {
		return ErrInvalidSequence
	}
	p.sum.Add(p.sum, &segmentSum)
	return nil
}

func (p *bigIntParser) parse(s string) error {
	n := len(s)
	i := 0
	var expectedRunes stack
loop:
	for ; i < n-2; i += 3 {
		r := decodeUtf8Kanji(i, s)
		if skip, err := expectedRunes.pop(r); err != nil {
			return err
		} else if skip {
			continue
		}
		value := bigIntValueOf(r)
		if value != nil && value.Sign() > 0 {
			var err error
			if value.Cmp(&b万) < 0 {
				err = p.push(value)
			} else {
				err = p.endSegmentWith(value)
			}
			if err != nil {
				return err
			}
		} else {
			switch r {
			case '零', '〇':
				if i == 0 {
					i += 3
					break loop
				}
				return ErrInvalidSequence
			default:
				return checkUnexpectedRune(s[i:])
			}
		}
		switch r {
		case '恒': // 恒河沙 10^52
			expectedRunes.push('沙', '河')
		case '阿': // 阿僧祇 10^56
			expectedRunes.push('祇', '僧')
		case '那': // 那由他 10^60
			expectedRunes.push('他', '由')
		case '不': // 不可思議 10^64
			expectedRunes.push('議', '思', '可')
		case '無': // 無量大数
			expectedRunes.push('数', '大', '量')
		}
	}
	if i < n {
		return checkUnexpectedRune(s[i:])
	}
	if !expectedRunes.empty() {
		return ErrEOF
	}
	return p.endSegment()
}

// stack stores the expected runes for multi kanji numerals.
type stack struct {
	runes []rune
}

// push adds the given runes to the stack.
func (s *stack) push(runes ...rune) {
	s.runes = append(s.runes, runes...)
}

// pop removes a rune from the stack and returns an error if the stack is not empty and the given rune does not match
// the top of the stack.
func (s *stack) pop(r rune) (skip bool, err error) {
	if len(s.runes) == 0 {
		return
	}
	last := len(s.runes) - 1
	if expected := s.runes[last]; expected != r {
		return false, &UnexpectedRuneError{r, expected}
	}
	s.runes = s.runes[:last]
	return true, nil
}

func (s *stack) empty() bool {
	return len(s.runes) == 0
}

// bigIntValueOf returns the value of the given japanese numeral. Expects only the first
// rune of multi kanji numerals. The result must treated as read-only.
func bigIntValueOf(r rune) *big.Int {
	switch r {
	case '零':
		return &b零
	case '〇':
		return &b零
	case '一':
		return &b一
	case '二':
		return &b二
	case '三':
		return &b三
	case '四':
		return &b四
	case '五':
		return &b五
	case '六':
		return &b六
	case '七':
		return &b七
	case '八':
		return &b八
	case '九':
		return &b九
	case '十':
		return &b十
	case '百':
		return &b百
	case '千':
		return &b千
	case '万':
		return &b万
	case '億':
		return &b億
	case '兆':
		return &b兆
	case '京':
		return &b京
	case '垓':
		return &b垓
	case '秭':
		return &b秭
	case '穣':
		return &b穣
	case '溝':
		return &b溝
	case '澗':
		return &b澗
	case '正':
		return &b正
	case '載':
		return &b載
	case '極':
		return &b極
	// first rune of multi kanji numerals
	case '恒':
		return &b恒河沙
	case '阿':
		return &b阿僧祇
	case '那':
		return &b那由他
	case '不':
		return &b不可思議
	case '無':
		return &b無量大数
	// formal numbers / daiji / 大字
	case '壱':
		return &b一
	case '弐':
		return &b二
	case '参':
		return &b三
	case '伍':
		return &b五
	case '拾':
		return &b十
	case '萬':
		return &b万
	default:
		return nil
	}
}
