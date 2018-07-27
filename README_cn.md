# OrderedMap

英文版: [English](README.md)

# 概要

OrderedMap是一个能够保持键(key)顺序的map。

# 特点

* 简单
* 快速(使用`json.Decoder`来反序列化)
* 强大的API
* 完备的文档
* 序列化/反序列化支持

# 使用

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

# 局限性

* OrderedMap仅支持字符串类型的键(key)。
* OrderedMap对于并发访问不是线程安全的。(使用sync.RWMutex来增加并发访问应该比较简单)

# 其它可选项

* [cevaris/ordered_map](https://github.com/cevaris/ordered_map)
不支持序列化/反序列化

* [iancoleman/orderedmap](https://github.com/iancoleman/orderedmap)
反序列化是作者自己写的，没有使用`json.Decoder`

# 许可证

MIT

# 备注

如果你喜欢此项目，请多多star，fork。谢谢！

