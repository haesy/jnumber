# JNumber - `strconv` for Japanese numerals

Go module that implements the conversion between UTF-8 encoded Japanese numerals and uint64/int64/big.Int.

## Features

- fast
- zero/low allocations
- zero external dependencies
- supports conversion from/to `int64`, `uint64` and `big.Int`
- supports numbers |x| < 10^72 (as long as they fit into the used datatype)
- supports daiji (大字), both current and obsolete ones
- supports serial numbers like 二〇二三 for 2023
- negative numbers use マイナス as a prefix

## Examples

```go
package main

import (
    "fmt"
    "math/big"
    
    "github.com/haesy/jnumber"
)

func main() {
    // int64/uint64/big.Int -> string
    fmt.Println(jnumber.FormatUint(299)) // "二百九十九"
    fmt.Println(jnumber.FormatInt(-299)) // "マイナス二百九十九"
    fmt.Println(jnumber.FormatBigInt(big.NewInt(299))) // "二百九十九"
    fmt.Println(jnumber.FormatSerialInt(2023)) // "二〇二三"

    // string -> int64/uint64/big.Int
    fmt.Println(jnumber.ParseUint("一千二百三十四")) // 1234
    fmt.Println(jnumber.ParseInt("マイナス二十三万四千五百六十七")) // -234567
    fmt.Println(jnumber.ParseInt("九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七")) // 9223372036854775807
    fmt.Println(jnumber.ParseBigInt("一無量大数")) // 10^68
    fmt.Println(jnumber.ParseSerialInt("二〇二三")) // 2023
    
    // support for daiji
    fmt.Println(jnumber.ParseInt("弐千")) // 2000
    fmt.Println(jnumber.ParseInt("壱万")) // 10000

    // numeric value of a single kanji
    fmt.Println(jnumber.ValueOf('零')) // 0
    fmt.Println(jnumber.ValueOf('〇')) // 0
    fmt.Println(jnumber.ValueOf('一')) // 1
    fmt.Println(jnumber.ValueOf('二')) // 2
    fmt.Println(jnumber.ValueOf('三')) // 3
    fmt.Println(jnumber.ValueOf('十')) // 10
    fmt.Println(jnumber.ValueOf('万')) // 10000
}
```

## Supported Numerals

Character | Value | Character | Value
--- | ---: | --- | ---:
零 / 〇 | 0 | 兆 | 10<sup>12</sup>
一 / 壱 * / 壹 * | 1 | 京 ** | 10<sup>16</sup>
二 / 弐 * / 貳 * | 2 | 垓 | 10<sup>20</sup>
三 / 参 * / 參 * | 3 | 秭 | 10<sup>24</sup>
四 / 肆 * | 4 | 穣 | 10<sup>28</sup>
五 / 伍 * | 5 | 溝 | 10<sup>32</sup>
六 / 陸 * | 6 | 澗 | 10<sup>36</sup>
七 / 柒 * / 漆 * | 7 | 正 | 10<sup>40</sup>
八 / 捌 * | 8 | 載 | 10<sup>44</sup>
九 / 玖 * | 9 | 極 | 10<sup>48</sup>
十 / 拾 * | 10 | 恒河沙 | 10<sup>52</sup>
百 / 佰 * | 100 | 阿僧祇 | 10<sup>56</sup>
千 / 阡 * / 仟 * | 1.000 | 那由他 | 10<sup>60</sup>
万 / 萬 * | 10<sup>4</sup> | 不可思議 | 10<sup>64</sup>
億 | 10<sup>8</sup> | 無量大数 | 10<sup>68</sup>

\* = Daiji / 大字

\*\* = Biggest numeral that fits into int64/uint64

## Contributing

If you find any bugs or want additional features please create an issue with details.

Merge requests out of the blue without any context or explanation are ignored.

## License

MIT
