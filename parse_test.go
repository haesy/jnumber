package jnumber

import (
	"errors"
	"math"
	"strings"
	"testing"
)

func TestParseInt(t *testing.T) {
	testParse(t, commonTestCases, ParseInt)
	testParse(t, boundaryTestCases, ParseInt)
	testParse(t, bankNoteTestCases, ParseInt)
	testParse(t, parseTestCases, ParseInt)
}

func TestParseIntError(t *testing.T) {
	testParseError(t, commonErrorCases, ParseInt)
	testParseError(t, intOverflowTestCases, ParseInt)
}

func TestParseUint(t *testing.T) {
	testParse(t, uintTestCases, ParseUint)
}

func TestParseUintError(t *testing.T) {
	testParseError(t, uintOverflowTestCases, ParseUint)
}

func TestParseSerialInt(t *testing.T) {
	testParse(t, serialTestCases, ParseSerialInt)
}

func testParse[T comparable](t *testing.T, tcs []testCase[T], fn func(string) (T, error)) {
	for _, tc := range tcs {
		t.Run(tc.String, func(st *testing.T) {
			actual, err := fn(tc.String)
			expectEqual(st, tc.Value, actual)
			expectErrNil(st, err)
		})
	}
}

func testParseError[T comparable](t *testing.T, tcs []parseErrorTestCase, fn func(string) (T, error)) {
	for _, tc := range tcs {
		t.Run(tc.Text, func(st *testing.T) {
			var zero T
			actualValue, actualErr := fn(tc.Text)
			expectEqual(st, zero, actualValue)
			expectErrIs(st, tc.Expected, actualErr)
		})
	}
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

func TestParseSerialIntError(t *testing.T) {
	testParseError(t, parseSerialErrorCases, ParseSerialInt)
}

func BenchmarkParseInt(b *testing.B) {
	for _, tc := range commonTestCases {
		b.Run(tc.String, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseInt(tc.String)
			}
		})
	}
}

func BenchmarkParseUint(b *testing.B) {
	for _, tc := range commonErrorCases {
		if strings.HasPrefix(tc.Text, negativePrefix) {
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
	for _, tc := range []testCase[int64]{
		{"〇", 0},
		{"一", 1},
		{"百二十三", 123},
		{"一万二千三百四十五", 12_345},
		{"二十三万四千五百六十七", 234_567},
		{"三百四十五万六千七百八十九", 3_456_789},
		{"九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七", math.MaxInt64},
	} {
		b.Run(tc.String, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ParseUint(tc.String)
			}
		})
	}
}
