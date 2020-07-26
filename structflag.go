// Package structflag exposes configuration structs as CLI flags.
package structflag

import (
	"flag"
	"reflect"
	"time"
)

// Load creates a flag for each field of a struct. See LoadTo for details on how
// to configure flag naming, usage, and default values.
//
// These flags are created on flag.CommandLine, which is the default (global)
// FlagSet. The flag names are unprefixed.
func Load(v interface{}) {
	LoadTo(flag.CommandLine, "", v)
}

// LoadTo creates a flag for each field of a struct with the given FlagSet.
//
// Each created flag will be prefixed by the given prefix plus a dash ("-"),
// unless prefix is empty in which case no dash is prepended to the flag name.
//
// The created flags will be set up to update the fields of v; after calling
// fs.Parse, fields in v may be updated by the flag package.
//
// The values of the fields in v are used as the default values passed to the
// flag package.
//
// By default, flags will be named after the fields in the given struct. To set
// a custom name for a flag, use a tag named "flag". To disable a field from
// having any flags generated, use the name "-".
//
// By default, flags will not have any usage message. To set a usage message,
// use a tag named "usage".
//
// Examples of struct field tags and their meanings:
//
//  // Field appears as a flag named "Field" with no usage info.
//  Field int
//
//  // Field appears as a flag named "foo" with no usage info.
//  Field int `flag:"foo"`
//
//  // Field appears as a flag named "foo" with usage "bar".
//  Field int `flag:"foo" usage:"bar"`
//
//  // Field is ignored by this package.
//  Field int `flag:"-"`
//
// The following field types are supported by this package, and all other types
// are ignored:
//
//  bool
//  float64
//  int
//  uint
//  int64
//  uint64
//  time.Duration
//
// These correspond to the types supported natively by the flag package.
//
// If a field's value is a struct then that contained struct will be recursively
// loaded as well. Anonymous struct fields are loaded as though they were were
// named the same as their type's name, unless renamed by the "flag" tag.
//
// For example, given the following "config" struct:
//
//  type config struct {
//    Foo string `flag:"foo"`
//    Bar struct {
//      Baz string `flag:"baz"`
//    } `flag:"baaar"`
//    embedded `flag:"embezzled`
//  }
//
//  type embedded struct {
//    Quux string `flag:"quux"`
//  }
//
// If an instance of config were passed to LoadTo with an empty prefix, then it
// would generate the following flags:
//
//  foo
//  baaar-baz
//  embezzled-quux
//
// LoadTo is subject to the usual Go visibility rules. If a field is unexported,
// then no flag will be created for that field.
//
// Cyclic data structures will lead to a stack overflow.
//
// Panics if v is not a pointer to a struct.
func LoadTo(fs *flag.FlagSet, prefix string, v interface{}) {
	val := reflect.ValueOf(v).Elem()
	load(fs, prefix, val)
}

func load(fs *flag.FlagSet, prefix string, val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		usage := field.Tag.Get("usage")
		flagValue := field.Tag.Get("flag")

		// Skip struct fields that are marked as `flag:"-"`
		if flagValue == "-" {
			continue
		}

		// Name the flag after the value of the `flag:"xxx"` tag. If that's not
		// present, default to the field's name.
		//
		// This is similar to the default behavior of the encoding/json package.
		name := field.Name
		if flagValue != "" {
			name = flagValue
		}

		// Supposing the prefix is something like "prefix-", then the name of this
		// flag will be "prefix-name".
		//
		// However, if the prefix is "", then the name of this flag will just be
		// "name", without an additional dash.
		if prefix != "" {
			name = prefix + "-" + name
		}

		switch val.Field(i).Kind() {
		case reflect.Struct:
			load(fs, name, val.Field(i))
		case reflect.Bool, reflect.Int64, reflect.Float64, reflect.Int, reflect.Uint, reflect.Uint64, reflect.String:
			switch f := val.Field(i).Addr().Interface().(type) {
			case *bool:
				fs.BoolVar(f, name, *f, usage)
			case *time.Duration:
				fs.DurationVar(f, name, *f, usage)
			case *float64:
				fs.Float64Var(f, name, *f, usage)
			case *int:
				fs.IntVar(f, name, *f, usage)
			case *int64:
				fs.Int64Var(f, name, *f, usage)
			case *string:
				fs.StringVar(f, name, *f, usage)
			case *uint:
				fs.UintVar(f, name, *f, usage)
			case *uint64:
				fs.Uint64Var(f, name, *f, usage)
			}
		}
	}
}
