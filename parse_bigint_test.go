package jnumber

import (
	"errors"
	"math/big"
	"testing"
)

type parseBigIntTestCase struct {
	Text     string
	Expected *big.Int
}

var parseBigIntCases = []parseBigIntTestCase{
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

func TestParseBigIntPrimitiveTypeCases(t *testing.T) {
	for _, tc := range commonTestCases {
		t.Run(tc.String, func(t *testing.T) {
			actual, err := ParseBigInt(tc.String)
			if err != nil {
				t.Errorf("err: %v", err)
			}
			if actual == nil || !actual.IsInt64() {
				t.Errorf("expected: not nil and IsInt64, actual: %s", actual)
				return
			}
			if actualInt64 := actual.Int64(); actualInt64 != tc.Value {
				t.Errorf("expected: %d, actual: %d", tc.Value, actualInt64)
			}
		})
	}
}

func TestParseBigInt(t *testing.T) {
	for _, tc := range parseBigIntCases {
		t.Run(tc.Text, func(t *testing.T) {
			actual, err := ParseBigInt(tc.Text)
			if err != nil {
				t.Errorf("err: %v", err)
			}
			if actual == nil || tc.Expected.Cmp(actual) != 0 {
				t.Errorf("expected: %s, actual: %s", tc.Expected, actual)
			}
		})
	}
}

func TestParseBigIntError(t *testing.T) {
	for _, tc := range commonErrorCases {
		testParseBigIntError(t, tc.Text, tc.Expected)
	}
	// 恒河沙
	testParseBigIntError(t, "一恒", ErrEOF)
	testParseBigIntError(t, "一恒a", ErrUnexpectedRune)
	testParseBigIntError(t, "一恒一", ErrUnexpectedRune)
	testParseBigIntError(t, "一恒河", ErrEOF)
	testParseBigIntError(t, "一恒河一", ErrUnexpectedRune)
	// 阿僧祇
	testParseBigIntError(t, "一阿", ErrEOF)
	testParseBigIntError(t, "一阿一", ErrUnexpectedRune)
	testParseBigIntError(t, "一阿僧", ErrEOF)
	testParseBigIntError(t, "一阿僧一", ErrUnexpectedRune)
	// 那由他
	testParseBigIntError(t, "一那", ErrEOF)
	testParseBigIntError(t, "一那一", ErrUnexpectedRune)
	testParseBigIntError(t, "一那由", ErrEOF)
	testParseBigIntError(t, "一那由一", ErrUnexpectedRune)
	// 不可思議
	testParseBigIntError(t, "一不", ErrEOF)
	testParseBigIntError(t, "一不一", ErrUnexpectedRune)
	testParseBigIntError(t, "一不可", ErrEOF)
	testParseBigIntError(t, "一不可一", ErrUnexpectedRune)
	testParseBigIntError(t, "一不可思", ErrEOF)
	testParseBigIntError(t, "一不可思一", ErrUnexpectedRune)
	// 無量大数
	testParseBigIntError(t, "一無", ErrEOF)
	testParseBigIntError(t, "一無一", ErrUnexpectedRune)
	testParseBigIntError(t, "一無量", ErrEOF)
	testParseBigIntError(t, "一無量一", ErrUnexpectedRune)
	testParseBigIntError(t, "一無量大", ErrEOF)
	testParseBigIntError(t, "一無量大一", ErrUnexpectedRune)
}

func testParseBigIntError(t *testing.T, str string, expectedErr error) {
	t.Run(str, func(t *testing.T) {
		actualValue, actualErr := ParseBigInt(str)
		if !errors.Is(actualErr, expectedErr) {
			t.Errorf("expected: %v, actual: %v, result: %d", expectedErr, actualErr, actualValue)
		}
	})
}
