package jnumber

import (
	"math"
	"math/bits"
	"strings"
	"unicode/utf8"
)

// ParseInt returns the integer represented by the given japanese numerals.
func ParseInt(s string) (int64, error) {
	abs := strings.TrimPrefix(s, negativePrefix)
	isNegative := s != abs
	sum, err := ParseUint(abs)
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
	n := len(s)
	if n == 0 {
		return 0, ErrEmpty
	}
	sum := uint64(0)
	segment := uint64(0)
	lastValue := uint64(0)
	minSegmentValue := uint64(math.MaxUint64)
	minSegmentEnd := uint64(math.MaxUint64)
	i := 0
loop:
	for ; i < n-2; i += utf8KanjiBytes {
		r := decodeUtf8Kanji(i, s)
		value, ok := ValueOf(r)
		if value > 0 && ok {
			if value < 10 { // 1 to 9
				// last number must not be < 10 as well
				if 0 < lastValue && lastValue < 10 {
					return 0, ErrInvalidSequence
				}
				lastValue = value
			} else if value < 10_000 { // 10, 100, 1000
				// check if we already encountered this number in the current segment
				if value >= minSegmentValue {
					return 0, ErrInvalidSequence
				}
				minSegmentValue = value
				// multiply with last number if allowed and possible
				if 0 < lastValue && lastValue < 10 {
					segment += lastValue * value
				} else {
					segment += value
				}
				lastValue = value
			} else { // >= 1_0000
				// check if we already encountered this number
				if value >= minSegmentEnd {
					return 0, ErrInvalidSequence
				}
				minSegmentEnd = value
				// create sum of current segment and add to sum
				if 0 < lastValue && lastValue < 10 {
					segment += lastValue
				}
				if segment == 0 {
					return 0, ErrInvalidSequence
				}

				var carry uint64
				overflow, segmentSum := bits.Mul64(segment, value)
				sum, carry = bits.Add64(sum, segmentSum, 0)
				if carry > 0 || overflow > 0 {
					return 0, ErrOverflow
				}
				// prepare for new segment
				minSegmentValue = uint64(math.MaxUint64)
				segment = 0
				lastValue = 0
			}
		} else {
			switch r {
			case '零', '〇':
				// zero is only valid if it is the only rune
				if i == 0 {
					i += utf8KanjiBytes
					break loop
				}
				return 0, ErrInvalidSequence
			// 10^20 - 10^68 overflows uint64
			// only the first kanji for multi kanji numbers
			case '垓', '秭', '穣', '溝',
				'澗', '正', '載', '極',
				'恒', '阿', '那', '不',
				'無':
				return 0, ErrOverflow
			default:
				return 0, checkUnexpectedRune(s[i:])
			}
		}
	}
	// are there still runes in the string after we are done?
	if i < n {
		return 0, checkUnexpectedRune(s[i:])
	}
	// add last segment to sum if there is one
	if 0 < lastValue && lastValue < 10 {
		segment += lastValue
	}
	if segment > 0 {
		var carry uint64
		sum, carry = bits.Add64(sum, segment, 0)
		if carry > 0 {
			return 0, ErrOverflow
		}
	}
	return sum, nil
}

// if our custom decoding fails, use the correct implementation from the standard library.
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
