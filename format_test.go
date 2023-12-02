package jnumber

import (
	"math/big"
	"testing"
)

func TestAppendInt(t *testing.T) {
	testAppend(t, commonTestCases, AppendInt)
	testAppend(t, boundaryTestCases, AppendInt)
}

func TestAppendUint(t *testing.T) {
	testAppend(t, uintTestCases, AppendUint)
}

func TestAppendSerialInt(t *testing.T) {
	testAppend(t, serialTestCases, AppendSerialInt)
}

func TestFormatInt(t *testing.T) {
	testFormat(t, commonTestCases, FormatInt)
	testFormat(t, boundaryTestCases, FormatInt)
}

func TestFormatUint(t *testing.T) {
	testFormat(t, uintTestCases, FormatUint)
}

func TestFormatSerialInt(t *testing.T) {
	testFormat(t, serialTestCases, FormatSerialInt)
}

func testAppend[T comparable](t *testing.T, tcs []testCase[T], fn func([]byte, T) []byte) {
	const prefix = "prefix "
	for _, tc := range tcs {
		t.Run(tc.String, func(st *testing.T) {
			dst := make([]byte, 0)
			dst = append(dst, prefix...)
			dst = fn(dst, tc.Value)
			actual := string(dst)
			expected := prefix + tc.String
			expectEqual(st, expected, actual)
		})
	}
}

func testFormat[T comparable](t *testing.T, tcs []testCase[T], fn func(T) string) {
	for _, tc := range tcs {
		t.Run(tc.String, func(st *testing.T) {
			actual := fn(tc.Value)
			expectEqual(st, tc.String, actual)
		})
	}
}

func BenchmarkFormatInt(b *testing.B) {
	for _, tc := range commonTestCases {
		b.Run(tc.String, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				FormatInt(tc.Value)
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
	for _, tc := range commonTestCases {
		t.Run(tc.String, func(st *testing.T) {
			i := big.NewInt(tc.Value)
			actual := FormatBigInt(i)
			expectEqual(st, tc.String, actual)
		})
	}
}

func TestFormatBigInt(t *testing.T) {
	for _, tc := range formatBigIntCases {
		t.Run(tc.Expected, func(st *testing.T) {
			actual := FormatBigInt(tc.Number)
			expectEqual(st, tc.Expected, actual)
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
