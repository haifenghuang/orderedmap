# OrderedMap

Chinese version: [中文](README_cn.md)

# Summary

OrderedMap is a map which preserving the key order as they are added to the map.

# Feature

* Simple
* Fast(Use `json.Decoder` for unmarshaling)
* Powerful API
* Well documented
* Marshaling & UnMarshaling support.


# Usage

```go
package main

import (
    "fmt"
    "encoding/json"
    "github.com/haifenghuang/orderedmap"
)

func main() {
    om := orderedmap.New()
    om.Set("Name", "HuangHaiFeng")
    om.Set("Sex", "Male")
    om.Set("Hobby", "Programming")
    om.Set("Country", "China")

    hobby, _ := om.Get("Hobby")
    fmt.Printf("Hobby = %s\n", hobby)

    sex, _ := om.GetAt(1)
    fmt.Printf("sex = %v\n", sex)

    om.SetAt(2, "Married", true)
    married, _ := om.GetAt(2)
    fmt.Printf("married = %t\n", married)

    fmt.Printf("=============================\n")
    fmt.Printf("keys = %v\n", om.Keys())
    fmt.Printf("values = %v\n", om.Values())
    fmt.Printf("mapLen = %d\n", om.Len())

    fmt.Printf("=============================\n")
    om.DeleteAt(2)
    fmt.Printf("OrderedMap = %s\n", om)

    fmt.Printf("=============================\n")
    has := om.Exists("Married")
    fmt.Printf("Married? - %t\n", has)
    has = om.Exists("Country")
    fmt.Printf("Country? - %t\n", has)

    fmt.Printf("=============================\n")
    idx := om.Index("Hobby")
    fmt.Printf("Hobby key's index = %d\n", idx)

    fmt.Printf("=============================\n")
    b, _ := json.MarshalIndent(om, "", "    ")
    fmt.Printf("Marshal result = %s\n", string(b))

    fmt.Printf("=============================\n")
    jsonStream := `{
    "Name": "HuangHaiFeng",
    "Sex": "Male",
    "Hobby": "Programming",
    "Country": "China"
}`
    om2 := orderedmap.New()
    _ = json.Unmarshal([]byte(jsonStream), om2)
    fmt.Printf("om2 = %v\n", om2)
}
```

# Limitation

* OrderedMap only takes strings for the key.
* OrderedMap is not thread-saft for concurrent use.(It is simple enough to add one using sync.RWMutex)

# Alternatives

* [cevaris/ordered_map](https://github.com/cevaris/ordered_map)
Not support marshaling/unmarshaling.

* [iancoleman/orderedmap](https://github.com/iancoleman/orderedmap)
Unmarshaling is self-written, not using `json.Decoder`.

# License

MIT

# Remarks

If you like this repo，please star & fork. Thanks!