package jnumber

import (
	"errors"
	"math/rand"
	"testing"
)

var uint64Kanjis = [...]rune{
	'零', '〇', '一', '二', '三',
	'四', '五', '六', '七', '八',
	'九', '十', '百', '千', '万',
	'億', '兆', '京',
	'壱', '弐', '参', '拾', '萬',
	'壹', '貳', '參', '肆', '伍',
	'陸', '柒', '漆', '捌', '玖',
	'佰', '阡', '仟',
}

func TestFormatParseIntRandom(t *testing.T) {
	for i := 0; i < 1_000_000; i++ {
		expected := rand.Int63()
		if i%2 == 0 {
			expected = -expected
		}
		str := FormatInt(expected)
		actual, err := ParseInt(str)
		if err != nil {
			t.Errorf("err: %v", err)
			t.FailNow()
		}
		if actual != expected {
			t.Errorf("expected: %d, actual: %d, str: %s", expected, actual, str)
			t.FailNow()
		}
	}
}

func TestFormatParseUintRandom(t *testing.T) {
	for i := 0; i < 1_000_000; i++ {
		expected := rand.Uint64()
		str := FormatUint(expected)
		actual, err := ParseUint(str)
		if err != nil {
			t.Errorf("err: %v", err)
			t.FailNow()
		}
		if actual != expected {
			t.Errorf("expected: %d, actual: %d, str: %s", expected, actual, str)
			t.FailNow()
		}
	}
}

func TestValueOf(t *testing.T) {
	for _, k := range uint64Kanjis {
		t.Run(string(k), func(t *testing.T) {
			expectedValue, expectedOk := valueOfSwitch(k)
			actualValue, actualOk := ValueOf(k)
			if actualOk != expectedOk {
				t.Errorf("ok expected: %v, actual: %v", expectedOk, actualOk)
			}
			if actualValue != expectedValue {
				t.Errorf("value expected: %d, actual: %d", expectedValue, actualValue)
			}
		})
	}
}

type findTestCase struct {
	Text     string
	Expected []*SearchResult
}

var findCases = []findTestCase{
	{"", []*SearchResult{}},
	{"零", []*SearchResult{
		{Start: 0, End: 3, Str: "零", Value: 0},
	}},
	{"〇", []*SearchResult{
		{Start: 0, End: 3, Str: "〇", Value: 0},
	}},
	{"一", []*SearchResult{
		{Start: 0, End: 3, Str: "一", Value: 1},
	}},
	{"一 二 三", []*SearchResult{
		{Start: 0, End: 3, Str: "一", Value: 1},
		{Start: 4, End: 7, Str: "二", Value: 2},
		{Start: 8, End: 11, Str: "三", Value: 3},
	}},
	{"一一 二二二", []*SearchResult{
		{Start: 0, End: 6, Str: "一一", Value: 0, Err: ErrInvalidSequence},
		{Start: 7, End: 16, Str: "二二二", Value: 0, Err: ErrInvalidSequence},
	}},
}

func TestFind(t *testing.T) {
	for _, k := range findCases {
		t.Run(string(k.Text), func(t *testing.T) {
			actual := Find(k.Text)
			if len(actual) != len(k.Expected) {
				t.Errorf("len expected: %v, actual: %v", len(k.Expected), len(actual))
			}
			for i, actualMatch := range actual {
				expectedMatch := k.Expected[i]
				if actualMatch.Start != expectedMatch.Start {
					t.Errorf("start expected: %v, actual: %v", expectedMatch.Start, actualMatch.Start)
				}
				if actualMatch.End != expectedMatch.End {
					t.Errorf("end expected: %v, actual: %v", expectedMatch.End, actualMatch.End)
				}
				if actualMatch.Str != expectedMatch.Str {
					t.Errorf("str expected: %v, actual: %v", expectedMatch.Str, actualMatch.Str)
				}
				if actualMatch.Value != expectedMatch.Value {
					t.Errorf("value expected: %v, actual: %v", expectedMatch.Value, actualMatch.Value)
				}
				if expectedMatch.Err != nil && !errors.Is(actualMatch.Err, expectedMatch.Err) {
					t.Errorf("err expected: %v, actual: %v", expectedMatch.Err, actualMatch.Err)
				}
			}
		})
	}
}

func BenchmarkValueOf(b *testing.B) {
	for _, k := range uint64Kanjis {
		b.Run(string(k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				valueOfSwitch(k)
			}
		})
	}
}

func BenchmarkValueOfPerfectHash(b *testing.B) {
	for _, k := range uint64Kanjis {
		b.Run(string(k), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ValueOf(k)
			}
		})
	}
}

// valueOfSwitch returns the numeric value of a single kanji, if it has some.
// Old version of ValueOf that is slower than the current version that uses
// a perfect hash. Keep in case we want to add new kanji.
func valueOfSwitch(r rune) (value uint64, ok bool) {
	switch r {
	case '零':
		return i零, true
	case '〇':
		return i零, true
	case '一':
		return i一, true
	case '二':
		return i二, true
	case '三':
		return i三, true
	case '四':
		return i四, true
	case '五':
		return i五, true
	case '六':
		return i六, true
	case '七':
		return i七, true
	case '八':
		return i八, true
	case '九':
		return i九, true
	case '十':
		return i十, true
	case '百':
		return i百, true
	case '千':
		return i千, true
	case '万':
		return i万, true
	case '億':
		return i億, true
	case '兆':
		return i兆, true
	case '京':
		return i京, true
	// formal numbers / daiji / 大字
	case '壱':
		return i壱, true
	case '弐':
		return i弐, true
	case '参':
		return i参, true
	case '拾':
		return i拾, true
	case '萬':
		return i萬, true
	case '壹':
		return i壹, true
	case '貳':
		return i貳, true
	case '參':
		return i參, true
	case '肆':
		return i肆, true
	case '伍':
		return i伍, true
	case '陸':
		return i陸, true
	case '柒':
		return i柒, true
	case '漆':
		return i漆, true
	case '捌':
		return i捌, true
	case '玖':
		return i玖, true
	case '佰':
		return i佰, true
	case '阡':
		return i阡, true
	case '仟':
		return i仟, true
	default:
		return 0, false
	}
}
