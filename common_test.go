package jnumber

import (
	"errors"
	"math"
	"testing"
	"unicode/utf8"
)

type testCase[T comparable] struct {
	String string
	Value  T
}

type parseErrorTestCase struct {
	Text     string
	Expected error
}

var commonTestCases = []testCase[int64]{
	{"零", 0},
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
	{"千四", 1_004},
	{"千三十四", 1_034},
	{"千二百三十四", 1_234},
	{"二千", 2_000},
	{"二万", 20_000},
	{"一万二千三百四十五", 12_345},
	{"二十三万四千五百六十七", 234_567},
	{"三百四十五万六千七百八十九", 3_456_789},
	{"二億", 2 * i億},
	{"二兆", 2 * i兆},
	{"二京", 2 * i京},
}

var bankNoteTestCases = []testCase[int64]{
	{"千", 1_000},
	{"弐千", 2_000},
	{"五千", 5_000},
	{"壱万", 10_000},
}

var boundaryTestCases = []testCase[int64]{
	{"九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七", math.MaxInt64},
	{negativePrefix + "九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百八", math.MinInt64},
}

var uintTestCases = []testCase[uint64]{
	{"零", 0},
	{"十", 10},
	{"百", 100},
	{"千", 1_000},
	{"一万", 10_000},
	{"一億", i億},
	{"一兆", i兆},
	{"一京", i京},
	{"九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七", math.MaxInt64},
	{"九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百八", math.MaxInt64 + 1},
	{"千八百四十四京六千七百四十四兆七百三十七億九百五十五万千六百十五", math.MaxUint64},
}

var parseTestCases = []testCase[int64]{
	{"〇", 0},
}

var serialTestCases = []testCase[int64]{
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
	{"一〇", 10},
	{"二〇〇", 200},
	{"三〇〇〇", 3000},
	{"一二三四五六七八九〇", 1234567890},
	{negativePrefix + "一二三四五六七八九〇", -1234567890},
	{"九二二三三七二〇三六八五四七七五八〇七", math.MaxInt64},
	{negativePrefix + "九二二三三七二〇三六八五四七七五八〇八", math.MinInt64},
}

var commonErrorCases = []parseErrorTestCase{
	{"", ErrEmpty},
	{"京", ErrInvalidSequence},
	{"一一", ErrInvalidSequence},
	{"一二", ErrInvalidSequence},
	{"二一", ErrInvalidSequence},
	{"二一十", ErrInvalidSequence},
	{"一二十", ErrInvalidSequence},
	{"十二一", ErrInvalidSequence},
	{"十一二", ErrInvalidSequence},
	{"十百", ErrInvalidSequence},
	{"十千", ErrInvalidSequence},
	{"十千", ErrInvalidSequence},
	{"十千", ErrInvalidSequence},
	{"一〇", ErrInvalidSequence},
	{"一零", ErrInvalidSequence},
	{"〇一", &UnexpectedRuneError{'一', 0}},
	{"零一", &UnexpectedRuneError{'一', 0}},
	{"二十一十", ErrInvalidSequence},
	{"一十二十", ErrInvalidSequence},
	{"一万二万", ErrInvalidSequence},
	{"二万一万", ErrInvalidSequence},
	{string(utf8.RuneError), ErrEncoding},
	{string(utf8.RuneError) + "一", ErrEncoding},
	{"一" + string(utf8.RuneError), ErrEncoding},
}

var intOverflowTestCases = []parseErrorTestCase{
	{"九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百八", ErrOverflow},
	{negativePrefix + "九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百九", ErrOverflow},
}

var uintOverflowTestCases = []parseErrorTestCase{
	{"千八百四十四京六千七百四十四兆七百三十七億九百五十五万千六百十六", ErrOverflow},
	{"二千八百四十四京六千七百四十四兆七百三十七億九百五十五万千六百十五", ErrOverflow},
	{"一垓", ErrOverflow},
}

var parseSerialErrorCases = []parseErrorTestCase{
	{"", ErrEmpty},
	{"百二十三", ErrInvalidSequence},
}

func expectEqual[T comparable](t *testing.T, expected, actual T) {
	if expected != actual {
		t.Errorf("expected: %v, actual: %v", expected, actual)
	}
}

func expectErrNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func expectErrIs(t *testing.T, expected, actual error) {
	if !errors.Is(actual, expected) {
		t.Errorf("expected error: %v, actual error: %v", expected, actual)
	}
}
