# JNumber

Go module that implements the conversion between UTF-8 encoded Japenese numerals and uint64/int64.

Support for big.Int is limited to FormatBigInt, ParseBigInt is not implemented yet.

Negative numbers use the char '-' (0x45) as prefix and are mainly implemented to allow serialization/deserialization without losing information.

## Examples

```go
package main

import (
    "fmt"
    "github.com/haesy/jnumber"
)

func main() {
    // int64/uint64 -> string
    fmt.Println(FormatUint(299)) // "二百九十九"
    fmt.Println(FormatInt(-299)) // "-二百九十九"

    // string -> int64/uint64
    fmt.Println(ParseUint("一千二百三十四")) // 1234
    fmt.Println(ParseInt("-二十三万四千五百六十七")) // -234567
    fmt.Println(ParseInt("九百二十二京三千三百七十二兆三百六十八億五千四百七十七万五千八百七") // 9223372036854775807
    
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

## License

MIT
