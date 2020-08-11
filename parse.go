package jnumber

import (
	"math"
	"strings"
	"unicode/utf8"
)

// ParseInt returns the integer represented by the given japanese numerals.
func ParseInt(s string) (int64, error) {
	isNegative := strings.HasPrefix(s, "-")
	if isNegative {
		s = s[1:]
	}
	sum, err := ParseUint(s)
	if err != nil {
		return 0, err
	}
	if isNegative {
		if sum > -math.MinInt64 {
			return 0, ErrOverflow
		}
		return -int64(sum), nil
	}
	if sum > math.MaxInt64 {
		return 0, ErrOverflow
	}
	return int64(sum), nil
}

// ParseUint returns the unsigned integer represented by the given japanese numerals.
func ParseUint(s string) (uint64, error) {
	if s == "" {
		return 0, ErrEmpty
	}
	parser := parser{}
	err := parser.parse(s)
	if err != nil {
		return 0, err
	}
	return parser.sum, nil
}

// parser contains the state for the parsing process. The segment could be a slice
// for easier code but this way we can avoid an allocation.
type parser struct {
	// current sum
	sum uint64
	// holds digit until we know what to do with them
	segment [16]uint64
	// position in the segment array for the next digit
	segmentIndex int
}

// append adds a digit to the current segment.
func (p *parser) append(digit uint64) {
	p.segment[p.segmentIndex] = digit
	p.segmentIndex++
}

// clearSegment clears the current segment to start parsing a new segment.
func (p *parser) clearSegment() {
	p.segmentIndex = 0
}

// push integrates a new digit to the current segment.
func (p *parser) push(digit uint64) error {
	if p.segmentIndex == 0 {
		p.append(digit)
		return nil
	}
	lastIndex := p.segmentIndex - 1
	lastDigit := p.segment[lastIndex]
	if lastDigit < digit {
		if digit < i十 {
			return ErrInvalidSequence
		} else if (digit == i百 || digit == i千) && lastDigit >= i十 {
			return ErrInvalidSequence
		}
		// Example: 二十 -> 2 * 10 -> 20
		p.segment[lastIndex] = lastDigit * digit
		return nil
	} else if lastDigit > digit && p.segmentIndex < len(p.segment) && (digit >= i十 || lastDigit >= i十) {
		// Example: 十一 -> 20 + 1 -> 21
		// We don't know if the next digit needs be multiplied with this digit or
		// added to the last. Store for handling at the end of the segment.
		p.append(digit)
		return nil
	}
	return ErrInvalidSequence
}

// endSegmentWith is a combination of push() and endSegment() for digits >= 万 (10000).
func (p *parser) endSegmentWith(digit uint64) error {
	if p.segmentIndex == 0 {
		return ErrInvalidSequence
	}
	var multiplierSum uint64
	var lastDigit uint64
	for i := 0; i < p.segmentIndex; i++ {
		segmentDigit := p.segment[i]
		if lastDigit > 0 && segmentDigit >= lastDigit {
			return ErrInvalidSequence
		}
		multiplierSum += segmentDigit
		lastDigit = segmentDigit
	}
	if multiplierSum >= i万 || multiplierSum >= digit {
		return ErrInvalidSequence
	}
	if digit == i京 && multiplierSum > maxParseUintMultiplier {
		return ErrOverflow
	}
	oldSum := p.sum
	p.sum += multiplierSum * digit
	p.clearSegment()
	if p.sum < oldSum {
		return ErrOverflow
	}
	return nil
}

func (p *parser) endSegment() error {
	if p.segmentIndex == 0 {
		return nil
	}
	var segmentSum uint64
	var lastDigit uint64
	for i := 0; i < p.segmentIndex; i++ {
		segmentDigit := p.segment[i]
		if lastDigit > 0 && segmentDigit >= lastDigit {
			return ErrInvalidSequence
		}
		segmentSum += segmentDigit
		lastDigit = segmentDigit
	}
	if segmentSum >= i万 {
		return ErrInvalidSequence
	}
	oldSum := p.sum
	p.sum += segmentSum
	if p.sum < oldSum {
		return ErrOverflow
	}
	return nil
}

func (p *parser) parse(s string) error {
	n := len(s)
	i := 0
loop:
	for ; i < n-2; i += 3 {
		r := decodeUtf8Kanji(i, s)
		value, ok := ValueOf(r)
		if value > 0 && ok {
			var err error
			if value < i万 {
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
			// 10^20 - 10^68 overflows uint64
			// only the first kanji for multi kanji numbers
			case '垓', '秭', '穣', '溝',
				'澗', '正', '載', '極',
				'恒', '阿', '那', '不',
				'無':
				return ErrOverflow
			default:
				return checkUnexpectedRune(s[i:])
			}
		}
	}
	if i < n {
		return checkUnexpectedRune(s[i:])
	}
	return p.endSegment()
}

func checkUnexpectedRune(s string) error {
	r, _ := utf8.DecodeRuneInString(s)
	switch r {
	case utf8.RuneError:
		return ErrEncoding
	default:
		return &UnexpectedRuneError{r, 0}
	}
}

// All kanji we want consist of 3 bytes in utf-8 encoding. This may seem unsafe,
// but if we encounter an unexpected or invalid rune, ValueOf will catch those
// values and we can retrieve the real rune with utf8.DecodeRuneInString and
// return a proper error. This optimization alone reduces the total parse time by
// about 30% to 60% compared to a for-range-loop over the string.
// byte 1: 1110 xxxx
// byte 2: 10xx xxxx
// byte 3: 10xx xxxx
func decodeUtf8Kanji(i int, s string) rune {
	byte1 := rune(s[i])
	byte2 := rune(s[i+1])
	byte3 := rune(s[i+2])
	validation := (byte1&0b_1111_0000 | (byte2&0b_1100_0000)>>4 | (byte3&0b_1100_0000)>>6) ^ 0b_1110_1010
	return byte3&0b_0011_1111 | byte2&0b_0011_1111<<6 | byte1&0b_0000_1111<<12 | validation<<24
}
