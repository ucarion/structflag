package structflag_test

import (
	"flag"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ucarion/structflag"
)

func TestLoadTo(t *testing.T) {
	// structflag will not work when you pass it an opaque interface{}. For that
	// reason, all of these tests need to use a concrete type when they call
	// structflag.LoadTo.
	//
	// For that reason, testCase uses this config type instead of an arbitrary
	// interface{} type.

	type configEmbedBar struct {
		Baz string `flag:"baz" usage:"baz usage"`
	}

	type configEmbed struct {
		Foo string         `flag:"foo" usage:"foo usage"`
		Bar configEmbedBar `flag:"bar"`
	}

	type config struct {
		Bool           bool          `flag:"bool" usage:"bool usage"`
		Duration       time.Duration `flag:"duration" usage:"duration usage"`
		Float          float64       `flag:"float" usage:"float usage"`
		Int            int           `flag:"int" usage:"int usage"`
		Int64          int64         `flag:"int64" usage:"int64 usage"`
		String         string        `flag:"string" usage:"string usage"`
		Uint           uint          `flag:"uint" usage:"uint usage"`
		Uint64         uint64        `flag:"uint64" usage:"uint64 usage"`
		EmbeddedStruct configEmbed   `flag:"embed"`
	}

	type testCase struct {
		Name string
		Args []string
		In   config
		Out  config
	}

	testCases := []testCase{
		testCase{
			Name: "bool",
			Args: []string{"--bool"},
			In:   config{},
			Out:  config{Bool: true},
		},
		testCase{
			Name: "duration",
			Args: []string{"--duration=5m"},
			In:   config{},
			Out:  config{Duration: 5 * time.Minute},
		},
		testCase{
			Name: "float",
			Args: []string{"--float=3.14"},
			In:   config{},
			Out:  config{Float: 3.14},
		},
		testCase{
			Name: "int",
			Args: []string{"--int=42"},
			In:   config{},
			Out:  config{Int: 42},
		},
		testCase{
			Name: "int64",
			Args: []string{"--int64=42"},
			In:   config{},
			Out:  config{Int64: 42},
		},
		testCase{
			Name: "uint",
			Args: []string{"--uint=42"},
			In:   config{},
			Out:  config{Uint: 42},
		},
		testCase{
			Name: "uint64",
			Args: []string{"--uint64=42"},
			In:   config{},
			Out:  config{Uint64: 42},
		},
		testCase{
			Name: "embedded struct",
			Args: []string{"--embed-foo=xxx", "--embed-bar-baz=yyy"},
			In:   config{},
			Out:  config{EmbeddedStruct: configEmbed{Foo: "xxx", Bar: configEmbedBar{Baz: "yyy"}}},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			fs := flag.NewFlagSet(tt.Name, flag.PanicOnError)
			structflag.LoadTo(fs, &tt.In)
			fs.Parse(tt.Args)

			assert.Equal(t, tt.Out, tt.In)
		})
	}

	t.Run("usage and default", func(t *testing.T) {
		fs := flag.NewFlagSet("test", flag.PanicOnError)
		structflag.LoadTo(fs, &config{
			Bool:     true,
			Duration: 5 * time.Minute,
			Float:    3.14,
			Int:      42,
			Int64:    42,
			String:   "default string",
			Uint:     42,
			Uint64:   42,
			EmbeddedStruct: configEmbed{
				Foo: "default foo",
				Bar: configEmbedBar{
					Baz: "default baz",
				},
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
			{Name: "embed-foo", Usage: "foo usage", DefValue: "default foo"}:     {},
			{Name: "embed-bar-baz", Usage: "baz usage", DefValue: "default baz"}: {},
		}, flags))
	})
}
