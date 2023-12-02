package jnumber

import (
	"math/big"
	"unsafe"
)

const (
	numberOfFastSmalls      = 101
	initialFormatBufferSize = 8 * utf8KanjiBytes
)

func AppendInt(dst []byte, i int64) []byte {
	if 0 <= i && i < numberOfFastSmalls {
		return append(dst, formatSmall(int(i))...)
	}
	var u uint64
	if i < 0 {
		u = uint64(-i)
		dst = append(dst, negativePrefix...)
	} else {
		u = uint64(i)
	}
	return formatUnsigned(dst, u)
}

func AppendUint(dst []byte, u uint64) []byte {
	if u < numberOfFastSmalls {
		return append(dst, formatSmall(int(u))...)
	}
	return formatUnsigned(dst, u)
}

func AppendSerialInt(dst []byte, i int64) []byte {
	var u uint64
	if i < 0 {
		u = uint64(-i)
		dst = append(dst, negativePrefix...)
	} else {
		u = uint64(i)
	}
	return AppendSerialUint(dst, u)
}

func AppendSerialUint(dst []byte, u uint64) []byte {
	if u < 10 {
		return append(dst, serialInts[int(u)]...)
	}
	var buffer [19]string
	b := buffer[:0]
	for u > 0 {
		r := u % 10
		u = u / 10
		c := serialInts[int(r)]
		b = append(b, c)
	}
	for i := len(b) - 1; i >= 0; i -= 1 {
		dst = append(dst, b[i]...)
	}
	return dst
}

// FormatInt returns the given integer as a string of Japanese numerals.
func FormatInt(i int64) string {
	if 0 <= i && i < numberOfFastSmalls {
		return formatSmall(int(i))
	}
	dst := make([]byte, 0, initialFormatBufferSize)
	var u uint64
	if i < 0 {
		u = uint64(-i)
		dst = append(dst, negativePrefix...)
	} else {
		u = uint64(i)
	}
	dst = formatUnsigned(dst, u)
	return unsafe.String(unsafe.SliceData(dst), len(dst))
}

// FormatUint returns the given unsigned integer as a string of Japanese numerals.
func FormatUint(u uint64) string {
	if u < numberOfFastSmalls {
		return formatSmall(int(u))
	}
	dst := make([]byte, 0, initialFormatBufferSize)
	dst = formatUnsigned(dst, u)
	return unsafe.String(unsafe.SliceData(dst), len(dst))
}

// FormatSerialInt returns the given integer as a string of Japanese numerals
// where the decimal digits 0 to 9 are replaced by the kanjis 〇 to 九.
func FormatSerialInt(i int64) string {
	if 0 <= i && i < 10 {
		return serialInts[i]
	}
	dst := make([]byte, 0, initialFormatBufferSize)
	var u uint64
	if i < 0 {
		u = uint64(-i)
		dst = append(dst, negativePrefix...)
	} else {
		u = uint64(i)
	}
	dst = AppendSerialUint(dst, u)
	return unsafe.String(unsafe.SliceData(dst), len(dst))
}

// FormatSerialInt returns the given unsigned integer as a string of Japanese
// numerals where the decimal digits 0 to 9 are replaced by the kanjis 〇 to 九.
func FormatSerialUint(u uint64) string {
	if u < 10 {
		return serialInts[u]
	}
	dst := make([]byte, 0, initialFormatBufferSize)
	dst = AppendSerialUint(dst, u)
	return unsafe.String(unsafe.SliceData(dst), len(dst))
}

func formatUnsigned(dst []byte, u uint64) []byte {
	if u >= i京 {
		dst, u = formatAppend(dst, u, "京", i京, u/i京)
	}
	if u >= i兆 {
		dst, u = formatAppend(dst, u, "兆", i兆, u/i兆)
	}
	if u >= i億 {
		dst, u = formatAppend(dst, u, "億", i億, u/i億)
	}
	if u >= i万 {
		dst, u = formatAppend(dst, u, "万", i万, u/i万)
	}
	if u >= i千 {
		dst, u = formatAppend(dst, u, "千", i千, u/i千)
	}
	if u > i百 {
		dst, u = formatAppend(dst, u, "百", i百, u/i百)
	}
	if 0 < u {
		dst = append(dst, formatSmall(int(u))...)
	}
	return dst
}

func formatAppend(dst []byte, u uint64, kanji string, kanjiValue uint64, multiplier uint64) ([]byte, uint64) {
	totalValue := multiplier * kanjiValue
	if multiplier == 1 {
		if kanjiValue >= i万 {
			dst = append(dst, "一"...)
		}
	} else if multiplier < numberOfFastSmalls {
		dst = append(dst, formatSmall(int(multiplier))...)
	} else {
		if multiplier >= i千 {
			dst, multiplier = formatAppend(dst, multiplier, "千", i千, multiplier/i千)
		}
		if multiplier > i百 {
			dst, multiplier = formatAppend(dst, multiplier, "百", i百, multiplier/i百)
		}
		if 0 < multiplier {
			dst = append(dst, formatSmall(int(multiplier))...)
		}
	}
	dst = append(dst, kanji...)
	return dst, u - totalValue
}

// FormatBigInt returns the given big integer as a string of Japanese numerals.
// Supports only numbers |i| < 10^72.
func FormatBigInt(i *big.Int) string {
	if i.IsInt64() {
		return FormatInt(i.Int64())
	} else if i.IsUint64() {
		return FormatUint(i.Uint64())
	}
	var u big.Int
	u.Abs(i)
	dst := make([]byte, 0, initialFormatBufferSize)
	if i.Sign() < 0 {
		dst = append(dst, negativePrefix...)
	}
	initBigIntsOnce.Do(initBigInts)
	dst = formatBigInt(dst, &u)
	return unsafe.String(unsafe.SliceData(dst), len(dst))
}

func formatBigInt(dst []byte, u *big.Int) []byte {
	dst = formatAppendBigInt(dst, u, "無量大数", &b無量大数)
	dst = formatAppendBigInt(dst, u, "不可思議", &b不可思議)
	dst = formatAppendBigInt(dst, u, "那由他", &b那由他)
	dst = formatAppendBigInt(dst, u, "阿僧祇", &b阿僧祇)
	dst = formatAppendBigInt(dst, u, "恒河沙", &b恒河沙)
	dst = formatAppendBigInt(dst, u, "極", &b極)
	dst = formatAppendBigInt(dst, u, "載", &b載)
	dst = formatAppendBigInt(dst, u, "正", &b正)
	dst = formatAppendBigInt(dst, u, "澗", &b澗)
	dst = formatAppendBigInt(dst, u, "溝", &b溝)
	dst = formatAppendBigInt(dst, u, "穣", &b穣)
	dst = formatAppendBigInt(dst, u, "秭", &b秭)
	dst = formatAppendBigInt(dst, u, "垓", &b垓)
	dst = formatAppendBigInt(dst, u, "京", &b京)
	return formatUnsigned(dst, u.Uint64())
}

func formatAppendBigInt(dst []byte, u *big.Int, kanji string, kanjiValue *big.Int) []byte {
	if u.Cmp(kanjiValue) < 0 {
		return dst
	}
	var multiplier big.Int
	multiplier.Div(u, kanjiValue)
	var totalValue big.Int
	totalValue.Mul(&multiplier, kanjiValue)

	if multiplier.Cmp(&maxBigIntMultiplier) > 0 {
		multiplier.Set(&maxBigIntMultiplier)
	}

	if intMultiplier := multiplier.Uint64(); intMultiplier < numberOfFastSmalls {
		dst = append(dst, formatSmall(int(intMultiplier))...)
	} else {
		if intMultiplier >= i千 {
			dst, intMultiplier = formatAppend(dst, intMultiplier, "千", i千, intMultiplier/i千)
		}
		if intMultiplier > i百 {
			dst, intMultiplier = formatAppend(dst, intMultiplier, "百", i百, intMultiplier/i百)
		}
		if 0 < intMultiplier {
			dst = append(dst, formatSmall(int(intMultiplier))...)
		}
	}
	dst = append(dst, kanji...)
	u.Sub(u, &totalValue)
	return dst
}

var smallInts = [...]string{
	"零", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十",
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

var serialInts = [...]string{
	"〇", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十",
}
