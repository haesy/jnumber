package jnumber

import (
	"errors"
	"math"
	"strings"
	"testing"
	"unicode/utf8"
)

type parseTestCase struct {
	Text     string
	Expected int64
}

var parseCases = []parseTestCase{
	{"零", 0},
	{"〇", 0},
	{"一", 1},
	{"二", 2},
	{"三", 3},
	{"四", 4},
	{"五", 5},
	{"六", 6},
	{"七", 7},
	{"八", 8},
	{"九", 9},
	{"十", 10},
	{"百", 100},
	{"千", 1_000},
	{"一千", 1_000},
	{"一万", 10_000},
	{"一億", i億},
	{"一兆", i兆},
	{"一京", i京},
	{"十一", 11},
	{"十二", 12},
	{"十三", 13},
	{"十四", 14},
	{"十五", 15},
	{"十六", 16},
	{"十七", 17},
	{"十八", 18},
	{"十九", 19},
	{"二十", 20},
	{"二十一", 21},
	{"二十二", 22},
	{"二十三", 23},
	{"二十四", 24},
	{"二十五", 25},
	{"二十六", 26},
	{"二十七", 27},
	{"二十八", 28},
	{"二十九", 29},
	{"三十", 30},
	{"三十一", 31},
	{"三十二", 32},
	{"三十三", 33},
	{"三十四", 34},
	{"三十五", 35},
	{"三十六", 36},
	{"三十七", 37},
	{"三十八", 38},
	{"三十九", 39},
	{"九十九", 99},
	{"百一", 101},
	{"百十", 110},
	{"百十一", 111},
	{"百二十一", 121},
	{"百二十二", 122},
	{"百二十三", 123},
	{"百三十三", 133},
	{"百九十九", 199},
	{"二百", 200},
	{"二百一", 201},
	{"二百十", 210},
	{"二百十一", 211},
	{"二百九十九", 299},
	{"三百", 300},
	{"一千四", 1_004},
	{"一千三十四", 1_034},
	{"一千二百三十四", 1_234},
	{"二千", 2_000},
	{"三千", 3_000},
	{"二万", 20_000},
	{"三万", 30_000},
	{"一万二千三百四十五", 12_345},
	{"二十三万四千五百六十七", 234_567},
	{"三百四十五万六千七百八十九", 3_456_789},
	{"二億", 2 * i億},
	{"二兆", 2 * i兆},
	{"二京", 2 * i京},
	{"九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七", math.MaxInt64},
	{"-九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百八", math.MinInt64},
	// bank notes
	{"千", 1_000},
	{"弐千", 2_000},
	{"五千", 5_000},
	{"壱万", 10_000},
}

func TestParseInt(t *testing.T) {
	for _, tc := range parseCases {
		t.Run(tc.Text, func(t *testing.T) {
			actual, err := ParseInt(tc.Text)
			if err != nil {
				t.Errorf("err: %v", err)
			}
			if actual != tc.Expected {
				t.Errorf("expected: %d, actual: %d", tc.Expected, actual)
			}
		})
	}
}

func TestParseUintMax(t *testing.T) {
	expected := uint64(math.MaxUint64)
	actual, err := ParseUint("千八百四十四京六千七百四十四兆七百三十七億九百五十五万千六百十五")
	if err != nil {
		t.Errorf("err: %v", err)
	}
	if actual != expected {
		t.Errorf("expected: %d, actual: %d", expected, actual)
	}
}

func TestParseUintOverflowOne(t *testing.T) {
	testParseUintError(t, "千八百四十四京六千七百四十四兆七百三十七億九百五十五万千六百十六", ErrOverflow)
}

func TestParseUintOverflowTwoTimes(t *testing.T) {
	testParseUintError(t, "二千八百四十四京六千七百四十四兆七百三十七億九百五十五万千六百十五", ErrOverflow)
}

func TestParseUintOverflowKanji(t *testing.T) {
	testParseUintError(t, "一垓", ErrOverflow)
}

func TestParseUintUnexpectedRune(t *testing.T) {
	expectedRune := 'a'
	value, err := ParseUint(string(expectedRune))
	if !errors.Is(err, ErrUnexpectedRune) {
		t.Errorf("expected: %v, actual: %v, result: %d", ErrUnexpectedRune, err, value)
	}
	result, ok := err.(*UnexpectedRuneError)
	if ok && expectedRune != result.Actual {
		t.Errorf("expected: %s, actual: %s, result: %d", string(expectedRune), string(result.Actual), value)
	}
}

func TestParseUintUnexpectedRuneAfterValid(t *testing.T) {
	expectedRune := 'a'
	value, err := ParseUint("一a")
	if !errors.Is(err, ErrUnexpectedRune) {
		t.Errorf("expected: %v, actual: %v, result: %d", ErrUnexpectedRune, err, value)
	}
	result, ok := err.(*UnexpectedRuneError)
	if ok && expectedRune != result.Actual {
		t.Errorf("expected: %s, actual: %s, result: %d", string(expectedRune), string(result.Actual), value)
	}
}

func TestParseUintError(t *testing.T) {
	testParseUintError(t, "", ErrEmpty)
	testParseUintError(t, "京", ErrInvalidSequence)
	testParseUintError(t, "一一", ErrInvalidSequence)
	testParseUintError(t, "一二", ErrInvalidSequence)
	testParseUintError(t, "二一", ErrInvalidSequence)
	testParseUintError(t, "二一十", ErrInvalidSequence)
	testParseUintError(t, "一二十", ErrInvalidSequence)
	testParseUintError(t, "十二一", ErrInvalidSequence)
	testParseUintError(t, "十一二", ErrInvalidSequence)
	testParseUintError(t, "十百", ErrInvalidSequence)
	testParseUintError(t, "十千", ErrInvalidSequence)
	testParseUintError(t, "十千", ErrInvalidSequence)
	testParseUintError(t, "十千", ErrInvalidSequence)
	testParseUintError(t, "一〇", ErrInvalidSequence)
	testParseUintError(t, "一零", ErrInvalidSequence)
	testParseUintError(t, "〇一", ErrUnexpectedRune)
	testParseUintError(t, "零一", ErrUnexpectedRune)
	testParseUintError(t, string(utf8.RuneError), ErrEncoding)
	testParseUintError(t, string(utf8.RuneError)+"一", ErrEncoding)
	testParseUintError(t, "一"+string(utf8.RuneError), ErrEncoding)
}

func testParseUintError(t *testing.T, str string, expectedErr error) {
	t.Run(str, func(t *testing.T) {
		actualValue, actualErr := ParseUint(str)
		if !errors.Is(actualErr, expectedErr) {
			t.Errorf("expected: %v, actual: %v, result: %d", expectedErr, actualErr, actualValue)
		}
	})
}

func BenchmarkParseInt(b *testing.B) {
	for _, tc := range parseCases {
		b.Run(tc.Text, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseInt(tc.Text)
			}
		})
	}
}

func BenchmarkParseUint(b *testing.B) {
	for _, tc := range parseCases {
		if strings.HasPrefix(tc.Text, "-") {
			continue
		}
		b.Run(tc.Text, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseUint(tc.Text)
			}
		})
	}
}

func BenchmarkQuickParseUint(b *testing.B) {
	for _, tc := range []parseTestCase{
		{"〇", 0},
		{"一", 1},
		{"百二十三", 123},
		{"一万二千三百四十五", 12_345},
		{"二十三万四千五百六十七", 234_567},
		{"三百四十五万六千七百八十九", 3_456_789},
		{"九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七", math.MaxInt64},
	} {
		b.Run(tc.Text, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseUint(tc.Text)
			}
		})
	}
}
