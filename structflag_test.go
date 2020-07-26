package structflag_test

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ucarion/structflag"
)

func ExampleLoad() {
	config := struct {
		FirstName string `flag:"first" usage:"first name"`
		LastName  string `flag:"last" usage:"last name"`
	}{
		FirstName: "default first name",
		LastName:  "default last name",
	}

	// Just as a hack for example purposes, we'll modify os.Args ourselves here.
	os.Args = []string{"my-cool-program", "-first=Foo", "-last=Bar"}

	structflag.Load(&config)
	flag.Parse()

	fmt.Println(config)
	// Output:
	//
	// {Foo Bar}
}

// config is the type used in test calls to LoadTo.
type config struct {
	Bool           bool          `flag:"bool" usage:"bool usage"`
	Duration       time.Duration `flag:"duration" usage:"duration usage"`
	Float          float64       `flag:"float" usage:"float usage"`
	Int            int           `flag:"int" usage:"int usage"`
	Int64          int64         `flag:"int64" usage:"int64 usage"`
	String         string        `flag:"string" usage:"string usage"`
	Uint           uint          `flag:"uint" usage:"uint usage"`
	Uint64         uint64        `flag:"uint64" usage:"uint64 usage"`
	SkipBool       bool          `flag:"-"`
	PlainBool      bool
	EmbeddedStruct configEmbed `flag:"embed"`
	anonEmbed      `flag:"anon"`
}

type configEmbed struct {
	Foo string         `flag:"foo" usage:"foo usage"`
	Bar configEmbedBar `flag:"bar"`
}

type configEmbedBar struct {
	Baz string `flag:"baz" usage:"baz usage"`
}

type anonEmbed struct {
	Quux string `flag:"quux" usage:"quux usage"`
}

func TestLoadTo(t *testing.T) {
	type testCase struct {
		Name string
		Args []string
		In   config
		Out  config
	}

	testCases := []testCase{
		testCase{
			Name: "bool",
			Args: []string{"--prefix-bool"},
			In:   config{},
			Out:  config{Bool: true},
		},
		testCase{
			Name: "duration",
			Args: []string{"--prefix-duration=5m"},
			In:   config{},
			Out:  config{Duration: 5 * time.Minute},
		},
		testCase{
			Name: "float",
			Args: []string{"--prefix-float=3.14"},
			In:   config{},
			Out:  config{Float: 3.14},
		},
		testCase{
			Name: "int",
			Args: []string{"--prefix-int=42"},
			In:   config{},
			Out:  config{Int: 42},
		},
		testCase{
			Name: "int64",
			Args: []string{"--prefix-int64=42"},
			In:   config{},
			Out:  config{Int64: 42},
		},
		testCase{
			Name: "uint",
			Args: []string{"--prefix-uint=42"},
			In:   config{},
			Out:  config{Uint: 42},
		},
		testCase{
			Name: "uint64",
			Args: []string{"--prefix-uint64=42"},
			In:   config{},
			Out:  config{Uint64: 42},
		},
		testCase{
			Name: "embedded struct",
			Args: []string{"--prefix-embed-foo=xxx", "--prefix-embed-bar-baz=yyy"},
			In:   config{},
			Out:  config{EmbeddedStruct: configEmbed{Foo: "xxx", Bar: configEmbedBar{Baz: "yyy"}}},
		},
		testCase{
			Name: "anonymous embedded struct",
			Args: []string{"--prefix-anon-quux=xxx"},
			In:   config{},
			Out:  config{anonEmbed: anonEmbed{Quux: "xxx"}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			fs := flag.NewFlagSet(tt.Name, flag.PanicOnError)
			structflag.LoadTo(fs, "prefix", &tt.In)
			fs.Parse(tt.Args)

			assert.Equal(t, tt.Out, tt.In)
		})
	}

	t.Run("usage and default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.PanicOnError)
		structflag.LoadTo(fs, "prefix", &config{
			Bool:      true,
			Duration:  5 * time.Minute,
			Float:     3.14,
			Int:       42,
			Int64:     42,
			String:    "default string",
			Uint:      42,
			Uint64:    42,
			SkipBool:  true,
			PlainBool: true,
			EmbeddedStruct: configEmbed{
				Foo: "default foo",
				Bar: configEmbedBar{
					Baz: "default baz",
				},
			},
			anonEmbed: anonEmbed{
				Quux: "default quux",
			},
		})

		flags := map[flag.Flag]struct{}{}
		fs.VisitAll(func(f *flag.Flag) {
			flag := *f
			flag.Value = nil
			flags[flag] = struct{}{}
		})

		assert.True(t, reflect.DeepEqual(map[flag.Flag]struct{}{
			{Name: "prefix-bool", Usage: "bool usage", DefValue: "true"}:                {},
			{Name: "prefix-duration", Usage: "duration usage", DefValue: "5m0s"}:        {},
			{Name: "prefix-float", Usage: "float usage", DefValue: "3.14"}:              {},
			{Name: "prefix-int", Usage: "int usage", DefValue: "42"}:                    {},
			{Name: "prefix-int64", Usage: "int64 usage", DefValue: "42"}:                {},
			{Name: "prefix-string", Usage: "string usage", DefValue: "default string"}:  {},
			{Name: "prefix-uint", Usage: "uint usage", DefValue: "42"}:                  {},
			{Name: "prefix-uint64", Usage: "uint64 usage", DefValue: "42"}:              {},
			{Name: "prefix-PlainBool", DefValue: "true"}:                                {},
			{Name: "prefix-embed-foo", Usage: "foo usage", DefValue: "default foo"}:     {},
			{Name: "prefix-embed-bar-baz", Usage: "baz usage", DefValue: "default baz"}: {},
			{Name: "prefix-anon-quux", Usage: "quux usage", DefValue: "default quux"}:   {},
		}, flags))
	})

	t.Run("usage and default with empty prefix", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.PanicOnError)
		structflag.LoadTo(fs, "", &config{
			Bool:      true,
			Duration:  5 * time.Minute,
			Float:     3.14,
			Int:       42,
			Int64:     42,
			String:    "default string",
			Uint:      42,
			Uint64:    42,
			SkipBool:  true,
			PlainBool: true,
			EmbeddedStruct: configEmbed{
				Foo: "default foo",
				Bar: configEmbedBar{
					Baz: "default baz",
				},
			},
			anonEmbed: anonEmbed{
				Quux: "default quux",
			},
		})

		flags := map[flag.Flag]struct{}{}
		fs.VisitAll(func(f *flag.Flag) {
			flag := *f
			flag.Value = nil
			flags[flag] = struct{}{}
		})

		assert.True(t, reflect.DeepEqual(map[flag.Flag]struct{}{
			{Name: "bool", Usage: "bool usage", DefValue: "true"}:                {},
			{Name: "duration", Usage: "duration usage", DefValue: "5m0s"}:        {},
			{Name: "float", Usage: "float usage", DefValue: "3.14"}:              {},
			{Name: "int", Usage: "int usage", DefValue: "42"}:                    {},
			{Name: "int64", Usage: "int64 usage", DefValue: "42"}:                {},
			{Name: "string", Usage: "string usage", DefValue: "default string"}:  {},
			{Name: "uint", Usage: "uint usage", DefValue: "42"}:                  {},
			{Name: "uint64", Usage: "uint64 usage", DefValue: "42"}:              {},
			{Name: "PlainBool", DefValue: "true"}:                                {},
			{Name: "embed-foo", Usage: "foo usage", DefValue: "default foo"}:     {},
			{Name: "embed-bar-baz", Usage: "baz usage", DefValue: "default baz"}: {},
			{Name: "anon-quux", Usage: "quux usage", DefValue: "default quux"}:   {},
		}, flags))
	})
}
