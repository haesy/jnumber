package jnumber

import (
	"math"
	"math/big"
	"testing"
)

type formatIntTestCase struct {
	Expected string
	Number   int64
}

var formatIntCases = []formatIntTestCase{
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
	{"九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七", math.MaxInt64},
}

func TestFormatInt(t *testing.T) {
	for _, tc := range formatIntCases {
		t.Run(tc.Expected, func(t *testing.T) {
			actual := FormatInt(tc.Number)
			if actual != tc.Expected {
				t.Errorf("expected: %s, actual: %s", tc.Expected, actual)
			}
		})
	}
}

func BenchmarkFormatInt(b *testing.B) {
	for _, tc := range formatIntCases {
		b.Run(tc.Expected, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FormatInt(tc.Number)
			}
		})
	}
}

type formatBigIntTestCase struct {
	Expected string
	Number   *big.Int
}

func newTestBigInt(m int64, e int64, a int64) *big.Int {
	result := &big.Int{}
	ten := big.NewInt(10)
	mul := big.NewInt(m)
	exp := big.NewInt(e)
	add := big.NewInt(a)
	return result.Exp(ten, exp, nil).Mul(result, mul).Add(result, add)
}

var formatBigIntCases = []formatBigIntTestCase{
	{"一垓", newTestBigInt(1, 20, 0)},
	{"二垓", newTestBigInt(2, 20, 0)},
	{"一垓一", newTestBigInt(1, 20, 1)},
	{"二垓二", newTestBigInt(2, 20, 2)},

	{"一秭", newTestBigInt(1, 24, 0)},
	{"二秭", newTestBigInt(2, 24, 0)},
	{"一秭一", newTestBigInt(1, 24, 1)},
	{"二秭二", newTestBigInt(2, 24, 2)},

	{"一穣", newTestBigInt(1, 28, 0)},
	{"二穣", newTestBigInt(2, 28, 0)},
	{"一穣一", newTestBigInt(1, 28, 1)},
	{"二穣二", newTestBigInt(2, 28, 2)},

	{"一溝", newTestBigInt(1, 32, 0)},
	{"二溝", newTestBigInt(2, 32, 0)},
	{"一溝一", newTestBigInt(1, 32, 1)},
	{"二溝二", newTestBigInt(2, 32, 2)},

	{"一澗", newTestBigInt(1, 36, 0)},
	{"二澗", newTestBigInt(2, 36, 0)},
	{"一澗一", newTestBigInt(1, 36, 1)},
	{"二澗二", newTestBigInt(2, 36, 2)},

	{"一正", newTestBigInt(1, 40, 0)},
	{"二正", newTestBigInt(2, 40, 0)},
	{"一正一", newTestBigInt(1, 40, 1)},
	{"二正二", newTestBigInt(2, 40, 2)},

	{"一載", newTestBigInt(1, 44, 0)},
	{"二載", newTestBigInt(2, 44, 0)},
	{"一載一", newTestBigInt(1, 44, 1)},
	{"二載二", newTestBigInt(2, 44, 2)},

	{"一極", newTestBigInt(1, 48, 0)},
	{"二極", newTestBigInt(2, 48, 0)},
	{"一極一", newTestBigInt(1, 48, 1)},
	{"二極二", newTestBigInt(2, 48, 2)},

	{"一恒河沙", newTestBigInt(1, 52, 0)},
	{"二恒河沙", newTestBigInt(2, 52, 0)},
	{"一恒河沙一", newTestBigInt(1, 52, 1)},
	{"二恒河沙二", newTestBigInt(2, 52, 2)},

	{"一阿僧祇", newTestBigInt(1, 56, 0)},
	{"二阿僧祇", newTestBigInt(2, 56, 0)},
	{"一阿僧祇一", newTestBigInt(1, 56, 1)},
	{"二阿僧祇二", newTestBigInt(2, 56, 2)},

	{"一那由他", newTestBigInt(1, 60, 0)},
	{"二那由他", newTestBigInt(2, 60, 0)},
	{"一那由他一", newTestBigInt(1, 60, 1)},
	{"二那由他二", newTestBigInt(2, 60, 2)},

	{"一不可思議", newTestBigInt(1, 64, 0)},
	{"二不可思議", newTestBigInt(2, 64, 0)},
	{"一不可思議一", newTestBigInt(1, 64, 1)},
	{"二不可思議二", newTestBigInt(2, 64, 2)},

	{"一無量大数", newTestBigInt(1, 68, 0)},
	{"二無量大数", newTestBigInt(2, 68, 0)},
	{"一無量大数一", newTestBigInt(1, 68, 1)},
	{"二無量大数二", newTestBigInt(2, 68, 2)},

	{"二無量大数一京二", newTestBigInt(2, 68, i京+2)},
	{"二無量大数一万二千三百四十五", newTestBigInt(2, 68, 12_345)},
	{"九千九百九十九無量大数", newTestBigInt(9999, 68, 0)},
}

func TestFormatBigIntSmall(t *testing.T) {
	for _, tc := range formatIntCases {
		t.Run(tc.Expected, func(t *testing.T) {
			i := big.NewInt(tc.Number)
			actual := FormatBigInt(i)
			if actual != tc.Expected {
				t.Errorf("expected: %s, actual: %s", tc.Expected, actual)
			}
		})
	}
}

func TestFormatBigInt(t *testing.T) {
	for _, tc := range formatBigIntCases {
		t.Run(tc.Expected, func(t *testing.T) {
			actual := FormatBigInt(tc.Number)
			if actual != tc.Expected {
				t.Errorf("expected: %s, actual: %s", tc.Expected, actual)
			}
		})
	}
}

func BenchmarkFormatBigInt(b *testing.B) {
	for _, tc := range formatBigIntCases {
		b.Run(tc.Expected, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FormatBigInt(tc.Number)
			}
		})
	}
}
