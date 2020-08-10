package jnumber

import (
	"math/rand"
	"testing"
)

var uint64Kanjis = []rune{
	'零', '〇', '一', '二', '三',
	'四', '五', '六', '七', '八',
	'九', '十', '百', '千', '万',
	'億', '兆', '京',
	'壱', '弐', '参', '拾', '萬',
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
		return i一, true
	case '弐':
		return i二, true
	case '参':
		return i三, true
	case '伍':
		return i五, true
	case '拾':
		return i十, true
	case '萬':
		return i万, true
	default:
		return 0, false
	}
}
