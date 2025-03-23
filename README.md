# jsonr – JSON Parsing with Reflection for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/trojanc/jsonr.svg)](https://pkg.go.dev/github.com/trojanc/jsonr)  
A Go library for marshalling to json while maintaining custom types. When JSON is unmarshalled the original types
will be recreated.


Why would you need this? Go's native `json` package does support marshalling back into a struct...yes...but...
As soon as you are using `any/interface{}` the type being used in that field is being lost, and it is impossible
to recreate without explicitly knowing the type when unmarshalling.

For example given the following:

```go
 data := map[string]any {
	"key1": MyStruct{ Name: "John" },
	"key2": OtherStruct { Score: 12},
 }
```
After this data has been marshalled to json, it is impossible to restore it back to a `map[string]any` with arbitrary
structs in it.

This library embeds type information into the JSON, that assists when Unmarshalling to recreate the correct type of
structs.
---

## 🔹 Features
- **Dynamic JSON parsing** using reflection.
- **Type inference** for automatic struct instantiation.
- **Supports custom types** dynamically.
- **Works with arbitrary and deeply nested JSON**.

---

## 📦 Installation
Requires **Go 1.23+**.

```sh
go get github.com/trojanc/jsonr
```
