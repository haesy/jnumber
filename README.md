# JNumber

Go module that implements the conversion between UTF-8 encoded Japenese numerals and uint64/int64/big.Int.

Support for big.Int is limited to numbers |x| < 10^72.

Negative numbers use the char '-' (0x45) as prefix and are mainly implemented to allow serialization/deserialization without losing information.

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
    fmt.Println(FormatUint(299)) // "二百九十九"
    fmt.Println(FormatInt(-299)) // "-二百九十九"
    fmt.Println(FormatBigInt(big.NewInt(299))) // "二百九十九"

    // string -> int64/uint64/big.Int
    fmt.Println(ParseUint("一千二百三十四")) // 1234
    fmt.Println(ParseInt("-二十三万四千五百六十七")) // -234567
    fmt.Println(ParseInt("九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七")) // 9223372036854775807
    fmt.Println(ParseBigInt("一無量大数")) // 10^68
    
    // support for daiji
    fmt.Println(ParseInt("弐千")) // 2000
    fmt.Println(ParseInt("壱万")) // 10000

    // numeric value of a single kanji
    fmt.Println(ValueOf('零')) // 0
    fmt.Println(ValueOf('〇')) // 0
    fmt.Println(ValueOf('一')) // 1
    fmt.Println(ValueOf('二')) // 2
    fmt.Println(ValueOf('三')) // 3
    fmt.Println(ValueOf('十')) // 10
    fmt.Println(ValueOf('万')) // 10000
}
```

## Supported Numerals

Character | Value
--- | ---:
零 / 〇 | 0
一 / 壱 * | 1
二 / 弐 * | 2
三 / 参 * | 3
四 | 4
五 | 5
六 | 6
七 | 7
八 | 8
九 | 9
十 / 拾 * | 10
百 | 100
千 | 1.000
万 / 萬 * | 10<sup>4</sup>
億 | 10<sup>8</sup>
兆 | 10<sup>12</sup>
京 ** | 10<sup>16</sup>
垓 | 10<sup>20</sup>
秭 | 10<sup>24</sup>
穣 | 10<sup>28</sup>
溝 | 10<sup>32</sup>
澗 | 10<sup>36</sup>
正 | 10<sup>40</sup>
載 | 10<sup>44</sup>
極 | 10<sup>48</sup>
恒河沙 | 10<sup>52</sup>
阿僧祇 | 10<sup>56</sup>
那由他 | 10<sup>60</sup>
不可思議 | 10<sup>64</sup>
無量大数 | 10<sup>68</sup>

\* = Daiji / 大字

\*\* = Biggest numeral that fits into int64/uint64

## License

MIT
