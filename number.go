package jnumber

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
)

var (
	// ErrEmpty is returned if a function does not allow empty strings as a parameter.
	ErrEmpty = errors.New("empty string")
	// ErrEOF is returned if a function expects a rune and detects an EOF instead.
	ErrEOF = errors.New("unexpected eof")
	// ErrOverflow is returned if the number is too big for the function.
	ErrOverflow = errors.New("number overflows datatype")
	// ErrEncoding ist returned if a function recognizes an invalid UTF-8 encoding.
	ErrEncoding = errors.New("invalid utf-8 encoding")
	// ErrInvalidSequence is returned if a string contains an invalid sequence of digits.
	// Examples: "一一" or "十百"
	ErrInvalidSequence = errors.New("invalid sequence of digits")
	// ErrUnexpectedRune is returned if a functions finds a rune that it does not expect.
	ErrUnexpectedRune = errors.New("unexpected rune")
)

// UnexpectedRuneError is returned if a functions finds a rune that it does not expect.
type UnexpectedRuneError struct {
	Actual, Expected rune
}

func (e *UnexpectedRuneError) Error() string {
	if e.Expected != 0 {
		return fmt.Sprintf("unexpected rune: expected %s, actual %s", string(e.Expected), string(e.Actual))
	}
	return fmt.Sprintf("unexpected rune: %s", string(e.Actual))
}

func (e *UnexpectedRuneError) Unwrap() error {
	return ErrUnexpectedRune
}

const (
	i零                     = 0
	i一                     = 1
	i二                     = 2
	i三                     = 3
	i四                     = 4
	i五                     = 5
	i六                     = 6
	i七                     = 7
	i八                     = 8
	i九                     = 9
	i十                     = 10
	i百                     = 100
	i千                     = 1_000
	i万                     = 10_000
	i億                     = 100_000_000
	i兆                     = 1_000_000_000_000
	i京                     = 10_000_000_000_000_000
	i壱                     = 1
	i弐                     = 2
	i参                     = 3
	i拾                     = 10
	i萬                     = 10_000
	maxParseUintMultiplier = math.MaxUint64 / i京
)

var (
	b零                  big.Int
	b一                  big.Int
	b二                  big.Int
	b三                  big.Int
	b四                  big.Int
	b五                  big.Int
	b六                  big.Int
	b七                  big.Int
	b八                  big.Int
	b九                  big.Int
	b十                  big.Int // 10^1
	b百                  big.Int // 10^2
	b千                  big.Int // 10^3
	b万                  big.Int // 10^4
	b億                  big.Int // 10^8
	b兆                  big.Int // 10^12
	b京                  big.Int // 10^16
	b垓                  big.Int // 10^20
	b秭                  big.Int // 10^24
	b穣                  big.Int // 10^28
	b溝                  big.Int // 10^32
	b澗                  big.Int // 10^36
	b正                  big.Int // 10^40
	b載                  big.Int // 10^44
	b極                  big.Int // 10^48
	b恒河沙                big.Int // 10^52
	b阿僧祇                big.Int // 10^56
	b那由他                big.Int // 10^60
	b不可思議               big.Int // 10^64
	b無量大数               big.Int // 10^68
	maxBigIntMultiplier big.Int
	initBigIntsOnce     sync.Once
)

func initBigInts() {
	b零.SetUint64(i零)
	b一.SetUint64(i一)
	b二.SetUint64(i二)
	b三.SetUint64(i三)
	b四.SetUint64(i四)
	b五.SetUint64(i五)
	b六.SetUint64(i六)
	b七.SetUint64(i七)
	b八.SetUint64(i八)
	b九.SetUint64(i九)
	b十.SetUint64(i十)
	b百.SetUint64(i百)
	b千.SetUint64(i千)
	b万.SetUint64(i万)
	b億.SetUint64(i億)
	b兆.SetUint64(i兆)
	b京.SetUint64(i京)
	var ten big.Int
	ten.SetUint64(10)
	b垓.Exp(&ten, big.NewInt(20), nil)
	b秭.Exp(&ten, big.NewInt(24), nil)
	b穣.Exp(&ten, big.NewInt(28), nil)
	b溝.Exp(&ten, big.NewInt(32), nil)
	b澗.Exp(&ten, big.NewInt(36), nil)
	b正.Exp(&ten, big.NewInt(40), nil)
	b載.Exp(&ten, big.NewInt(44), nil)
	b極.Exp(&ten, big.NewInt(48), nil)
	b恒河沙.Exp(&ten, big.NewInt(52), nil)
	b阿僧祇.Exp(&ten, big.NewInt(56), nil)
	b那由他.Exp(&ten, big.NewInt(60), nil)
	b不可思議.Exp(&ten, big.NewInt(64), nil)
	b無量大数.Exp(&ten, big.NewInt(68), nil)
	maxBigIntMultiplier.SetUint64(9999)
}

var (
	toDaijiReplacer = strings.NewReplacer(
		"一", "壱",
		"二", "弐",
		"三", "参",
		"五", "伍",
		"十", "拾",
		"万", "萬",
	)
	fromDaijiReplacer = strings.NewReplacer(
		"壱", "一",
		"弐", "二",
		"参", "三",
		"伍", "五",
		"拾", "十",
		"萬", "万",
	)
)

// ValueOf returns the numeric value of a single kanji, if it has one.
func ValueOf(r rune) (value uint64, ok bool) {
	hash := runeValuesPerfectHash(r)
	val := runeValues[hash]
	return val.Value, val.Rune == r
}

type runeValue struct {
	Rune  rune
	Value uint64
}

// runeValues contains the values of all kanjis that fit into uint64.
// Index must be calculated with runeValuesPerfectHash. Roughly 10 times
// faster than a switch.
var runeValues = [...]runeValue{
	{'参', i参}, {0, 0}, {0, 0}, {'一', i一}, {'京', i京}, {'七', i七},
	{'〇', 0}, {'万', i万}, {'三', i三}, {0, 0}, {0, 0}, {'壱', i壱},
	{'百', i百}, {0, 0}, {0, 0}, {'八', i八}, {'六', i六}, {'零', i零},
	{'二', i二}, {'拾', i拾}, {0, 0}, {'億', i億}, {'五', i五}, {'九', i九},
	{'十', i十}, {'千', i千}, {'兆', i兆}, {0, 0}, {0, 0}, {'弐', i弐},
	{'萬', i萬}, {'四', i四},
}

// runeValuesPerfectHash is a perfect hash function for all kanji with numeric values.
func runeValuesPerfectHash(r rune) int {
	return int(((75520486 * r) >> 27) & 0b11111)
}

// ToDaiji replaces some kanji with current daiji (大字).
func ToDaiji() *strings.Replacer {
	return toDaijiReplacer
}

// FromDaiji replaces current daiji with regular kanji.
func FromDaiji() *strings.Replacer {
	return fromDaijiReplacer
}