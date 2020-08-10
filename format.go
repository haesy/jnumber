package jnumber

import (
	"math/big"
	"strings"
)

const (
	numberOfFastSmalls      = 101
	preferredFormatZero     = "〇"
	initialFormatBufferSize = 24
)

// FormatInt returns the given integer as a string of japanese numerals.
func FormatInt(i int64) string {
	if 0 <= i && i < numberOfFastSmalls {
		return formatSmall(int(i))
	}
	var result strings.Builder
	result.Grow(initialFormatBufferSize)
	var u uint64
	if i < 0 {
		u = uint64(-i)
		result.WriteRune('-')
	} else {
		u = uint64(i)
	}
	formatUnsigned(&result, u)
	return result.String()
}

// FormatUint returns the given unsigned integer as a string of japanese numerals.
func FormatUint(u uint64) string {
	if u < numberOfFastSmalls {
		return formatSmall(int(u))
	}
	var result strings.Builder
	result.Grow(initialFormatBufferSize)
	formatUnsigned(&result, u)
	return result.String()
}

func formatUnsigned(result *strings.Builder, u uint64) {
	if u >= i京 {
		u = formatAppend(result, u, '京', i京, u/i京)
	}
	if u >= i兆 {
		u = formatAppend(result, u, '兆', i兆, u/i兆)
	}
	if u >= i億 {
		u = formatAppend(result, u, '億', i億, u/i億)
	}
	if u >= i万 {
		u = formatAppend(result, u, '万', i万, u/i万)
	}
	if u >= i千 {
		u = formatAppend(result, u, '千', i千, u/i千)
	}
	if u > i百 {
		u = formatAppend(result, u, '百', i百, u/i百)
	}
	if 0 < u {
		result.WriteString(formatSmall(int(u)))
	}
}

func formatAppend(result *strings.Builder, u uint64, kanji rune, kanjiValue uint64, multiplier uint64) uint64 {
	totalValue := multiplier * kanjiValue
	if multiplier == 1 {
		if kanjiValue >= i万 {
			result.WriteRune('一')
		}
	} else if multiplier < numberOfFastSmalls {
		result.WriteString(formatSmall(int(multiplier)))
	} else {
		if multiplier >= i千 {
			multiplier = formatAppend(result, multiplier, '千', i千, multiplier/i千)
		}
		if multiplier > i百 {
			multiplier = formatAppend(result, multiplier, '百', i百, multiplier/i百)
		}
		if 0 < multiplier {
			result.WriteString(formatSmall(int(multiplier)))
		}
	}
	result.WriteRune(kanji)
	return u - totalValue
}

// FormatBigInt returns the given big integer as a string of japanese numerals.
// Supports only numbers |i| < 10^72.
func FormatBigInt(i *big.Int) string {
	if i.IsInt64() {
		return FormatInt(i.Int64())
	} else if i.IsUint64() {
		return FormatUint(i.Uint64())
	}
	var u big.Int
	u.Abs(i)
	var result strings.Builder
	result.Grow(initialFormatBufferSize)
	if i.Sign() < 0 {
		result.WriteRune('-')
	}
	initBigIntsOnce.Do(initBigInts)
	formatBigInt(&result, &u)
	return result.String()
}

func formatBigInt(result *strings.Builder, u *big.Int) {
	formatAppendBigInt(result, u, "無量大数", &b無量大数)
	formatAppendBigInt(result, u, "不可思議", &b不可思議)
	formatAppendBigInt(result, u, "那由他", &b那由他)
	formatAppendBigInt(result, u, "阿僧祇", &b阿僧祇)
	formatAppendBigInt(result, u, "恒河沙", &b恒河沙)
	formatAppendBigInt(result, u, "極", &b極)
	formatAppendBigInt(result, u, "載", &b載)
	formatAppendBigInt(result, u, "正", &b正)
	formatAppendBigInt(result, u, "澗", &b澗)
	formatAppendBigInt(result, u, "溝", &b溝)
	formatAppendBigInt(result, u, "穣", &b穣)
	formatAppendBigInt(result, u, "秭", &b秭)
	formatAppendBigInt(result, u, "垓", &b垓)
	formatAppendBigInt(result, u, "京", &b京)
	formatUnsigned(result, u.Uint64())
}

func formatAppendBigInt(result *strings.Builder, u *big.Int, kanji string, kanjiValue *big.Int) {
	if u.Cmp(kanjiValue) < 0 {
		return
	}
	var multiplier big.Int
	multiplier.Div(u, kanjiValue)
	var totalValue big.Int
	totalValue.Mul(&multiplier, kanjiValue)

	if multiplier.Cmp(&maxBigIntMultiplier) > 0 {
		multiplier.Set(&maxBigIntMultiplier)
	}

	if intMultiplier := multiplier.Uint64(); intMultiplier < numberOfFastSmalls {
		result.WriteString(formatSmall(int(intMultiplier)))
	} else {
		if intMultiplier >= i千 {
			intMultiplier = formatAppend(result, intMultiplier, '千', i千, intMultiplier/i千)
		}
		if intMultiplier > i百 {
			intMultiplier = formatAppend(result, intMultiplier, '百', i百, intMultiplier/i百)
		}
		if 0 < intMultiplier {
			result.WriteString(formatSmall(int(intMultiplier)))
		}
	}

	result.WriteString(kanji)
	u.Sub(u, &totalValue)
}

var smallInts = [...]string{
	preferredFormatZero, "一", "二", "三", "四", "五", "六", "七", "八", "九", "十",
	"十一", "十二", "十三", "十四", "十五", "十六", "十七", "十八", "十九", "二十",
	"二十一", "二十二", "二十三", "二十四", "二十五", "二十六", "二十七", "二十八", "二十九", "三十",
	"三十一", "三十二", "三十三", "三十四", "三十五", "三十六", "三十七", "三十八", "三十九", "四十",
	"四十一", "四十二", "四十三", "四十四", "四十五", "四十六", "四十七", "四十八", "四十九", "五十",
	"五十一", "五十二", "五十三", "五十四", "五十五", "五十六", "五十七", "五十八", "五十九", "六十",
	"六十一", "六十二", "六十三", "六十四", "六十五", "六十六", "六十七", "六十八", "六十九", "七十",
	"七十一", "七十二", "七十三", "七十四", "七十五", "七十六", "七十七", "七十八", "七十九", "八十",
	"八十一", "八十二", "八十三", "八十四", "八十五", "八十六", "八十七", "八十八", "八十九", "九十",
	"九十一", "九十二", "九十三", "九十四", "九十五", "九十六", "九十七", "九十八", "九十九", "百",
}

func formatSmall(i int) string {
	return smallInts[i]
}
