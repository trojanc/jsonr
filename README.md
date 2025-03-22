# Type Wrapper

 A library that can marshal and unmarshal complex structures back into their original types.
 
This library is very new and definitely does not support all possible permutations of variables


# Supported Types

| Structure                                                                                  | Supported |
|--------------------------------------------------------------------------------------------|-----------|
| `MyStruct`                                                                                 | ✅         |
| `*MyStruct`                                                                                | ✅         |
| `[]MyStruct`, `*[]MyStruct`, `[]*MyStruct`, `*[]*MyStruct`                                 | ✅         |
| `map[string]MyStruct`, `map[string]*MyStruct`,`map[int]MyStruct`, `map[int]*MyStruct`,     | ✅         |
| `int`,`int8`,`int16`,`int32`,`int64`,`uint`,`uint8`,`uint16`,`uint32`,`uint64`             | ✅         |
| `*int`,`*int8`,`*int16`,`*int32`,`*int64`,  `*uint`,`*uint8`,`*uint16`,`*uint32`,`*uint64` | ✅         |
| `float32`, `float64`, `*float32`, `*float64`,                                              | ✅         |
| `bool`, `*bool`,                                                                           | ✅         |
| `byte`, `*byte`                                                                            | ✅         |
| `string`, `*string`                                                                        | ✅         |
| `map[string]any`, `map[MyStruct]any`, `map[MyStruct]MyStruct`                              | ❌         |