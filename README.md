# gjs

`gjs` is a Go package for generating and validating JSON Schemas from Go structs using generics. It leverages `github.com/invopop/jsonschema` for schema generation and `github.com/xeipuuv/gojsonschema` for schema validation.

## Features

- Generate JSON Schema from Go structs using reflection.
- Validate data against generated schemas.
- Support for generic types.
- Export schema as JSON string with customizable indentation.
- Save schema to a file with configurable file permissions and write modes.

## Installation

```bash
go get github.com/twoojoo/gjs
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/twoojoo/gjs"
)

type User struct {
    Name  		string `json:"name" jsonschema:"required,minLength=3"`
    Email 		string `json:"email" jsonschema:"required,format=email"`
    Age   		int    `json:"age" jsonschema:"required,minimum=0"`
    Description	 string `json:"description,omitempty" jsonschema:"maxLength=500"`
}

func main() {
    schema := gjs.NewSchema[User]()
    
    // Generate JSON schema string
    jsonSchema, err := schema.String(gjs.WithIndent("  "))
    if err != nil {
        panic(err)
    }
    fmt.Println(jsonSchema)

    // Validate data
    user := User{Name: "John Doe", Email: "john@example.com"}
    result, err := schema.Validate(&user)
    if err != nil {
        panic(err)
    }

    if result.Valid() {
        fmt.Println("User is valid")
    } else {
        fmt.Println("User validation errors:", result.Errors())
    }
}
```

## API

`NewSchema[T any](data ...T) *Schema[T]`

Creates a new schema generator for the given type T. Optionally accepts an instance of T to generate schema from.

`(s *Schema[T]) Validate(data *T) (*gojsonschema.Result, error)`

Validates the given data against the schema.

`(s *Schema[T]) ValidateAny(data any) (*gojsonschema.Result, error)`

Validates any data type against the schema.

`(s *Schema[T]) String(options ...Option) (string, error)`

Returns the JSON schema as a string. Supports options such as indentation.

`(s *Schema[T]) Store(filename string, options ...Option) error`

Stores the schema JSON to a file with customizable options.

## Options

- `WithIndent(string)` Option: Format JSON schema output with the specified indent string.
- `WithAppend()` Option: Open the file in append mode when storing.
- `WithTruncate()` Option: Open the file truncating existing content (default).
- `WithPermissions(os.FileMode)` Option: Set file permissions when storing schema.

Dependencies

This package uses the following libraries:

- [github.com/invopop/jsonschema](https://github.com/invopop/jsonschema)
For generating JSON Schema definitions from Go structs via reflection.

- [github.com/xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema)
For validating JSON data against JSON Schema.

## License

[MIT License](LICENSE.md)
