package gjs

import (
	"encoding/json"
	"os"

	js "github.com/invopop/jsonschema"
	gjs "github.com/xeipuuv/gojsonschema"
)

type Schema[T any] struct {
	schema *js.Schema
}

func NewSchema[T any](data ...T) *Schema[T] {
	var mock T

	if len(data) > 0 {
		mock = data[0]
	} else {
		mock = *new(T)
	}

	schema := js.Reflect(&mock)
	return &Schema[T]{schema: schema}
}

func (v *Schema[T]) Validate(data *T) (*gjs.Result, error) {
	schemaLoader := gjs.NewGoLoader(v.schema)
	dataLoader := gjs.NewGoLoader(data)
	return gjs.Validate(schemaLoader, dataLoader)
}

func (v *Schema[T]) ValidateAny(data any) (*gjs.Result, error) {
	schemaLoader := gjs.NewGoLoader(v.schema)
	dataLoader := gjs.NewGoLoader(&data)
	return gjs.Validate(schemaLoader, dataLoader)
}

// Struct returns the struct schema as struct.
func (v *Schema[T]) Struct() *js.Schema {
	return v.schema
}

// String returns the JSON version of the struct schema.
func (v *Schema[T]) String(options ...Option) (string, error) {
	opts := &Options{
		flag:   os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
		indent: "",
		perm:   0644,
	}

	// Apply passed options
	for _, opt := range options {
		opt(opts)
	}

	if opts.indent != "" {
		schemaBytes, err := json.MarshalIndent(v.schema, "", opts.indent)
		return string(schemaBytes), err
	}

	schemaBytes, err := json.Marshal(v.schema)
	return string(schemaBytes), err
}

func (v *Schema[T]) Store(filename string, options ...Option) error {
	// Default options: create/truncate, write only, permission 0644
	opts := &Options{
		flag: os.O_CREATE | os.O_WRONLY | os.O_TRUNC,
		perm: 0644,
	}

	// Apply passed options
	for _, opt := range options {
		opt(opts)
	}

	f, err := os.OpenFile(filename, opts.flag, opts.perm)
	if err != nil {
		return err
	}
	defer f.Close()

	stringSchema, err := v.String(options...)
	if err != nil {
		return err
	}

	_, err = f.WriteString(stringSchema)
	return err
}

type Option func(*Options)

type Options struct {
	flag   int
	indent string
	perm   os.FileMode
}

func WithAppend() Option {
	return func(opts *Options) {
		opts.flag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	}
}

func WithTruncate() Option {
	return func(opts *Options) {
		opts.flag = os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	}
}

func WithPermissions(perm os.FileMode) Option {
	return func(opts *Options) {
		opts.perm = perm
	}
}

func WithIndent(indent ...string) Option {
	if len(indent) == 0 {
		return func(opts *Options) {
			opts.indent = ""
		}
	}

	return func(opts *Options) {
		opts.indent = indent[0]
	}
}
