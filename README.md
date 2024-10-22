# OrderedMap

[![CI Status](https://github.com/schoren/orderedmap/actions/workflows/ci-validate.yaml/badge.svg)](https://github.com/schoren/orderedmap/actions/workflows/ci-validate.yaml)

[![Go Report](https://goreportcard.com/badge/github.com/schoren/orderedmap)](https://goreportcard.com/report/github.com/schoren/orderedmap)

[![Go Coverage](https://github.com/schoren/orderedmap/wiki/coverage.svg)](https://raw.githack.com/wiki/schoren/orderedmap/coverage.html)


`OrderedMap` is a Go package that provides a map-like data structure which maintains the order of insertion. It supports JSON marshaling and unmarshaling, and provides various utility methods.

`OrderedMap`s are immutable, meaning each time you set or delete an item, a new instance is created.

## Example

```go
package main

import (
    "fmt"
    "log"
    "github.com/schoren/orderedmap"
)

func main() {
    om := orderedmap.New[string, string]().
        MustSet("first", "1").
        MustSet("second", "2")

    fmt.Println(om.Get("first")) // Output: 1

    om.ForEach(func(key, value string) error {
        fmt.Printf("Key: %s, Value: %s\n", key, value)
        return nil
    })

    data, err := om.MarshalJSON()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(data)) // Output: [{"Key":"first","Value":"1"},{"Key":"second","Value":"2"}]

    var newOM orderedmap.OrderedMap[string, string]
    err = newOM.UnmarshalJSON(data)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(newOM.Get("second")) // Output: 2
}
```

## Installation

To install the package, run:

```sh
go get github.com/schoren/orderedmap
```

## Usage

### Creating an OrderedMap

```go
import "github.com/schoren/orderedmap"

om := orderedmap.New[string, string]()
```

### Setting and Getting Values

```go
om = om.MustSet("key1", "value1")
value := om.Get("key1") // "value1"
```


### Checking for Key Existence

```go
exists := om.Contains("key1") // true
```

### Iterating Over the Map

```go
err := om.ForEach(func(key, value string) error {
    fmt.Printf("Key: %s, Value: %s\n", key, value)
    return nil
})
if err != nil {
    log.Fatal(err)
}
```

### JSON Marshaling and Unmarshaling

```go
data, err := om.MarshalJSON()
if err != nil {
    log.Fatal(err)
}

var newOM orderedmap.OrderedMap[string, string]
err = newOM.UnmarshalJSON(data)
if err != nil {
    log.Fatal(err)
}
```

## Testing 

To run the tests, use:

```sh
go test ./...
```


## License

This project is licensed under the MIT License.