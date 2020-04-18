package structflag

import (
	"flag"
	"reflect"
	"time"
)

func Load(v interface{}) {
	LoadTo(flag.CommandLine, v)
}

func LoadTo(fs *flag.FlagSet, v interface{}) {
	val := reflect.ValueOf(v).Elem()
	load(fs, "", val)
}

func load(fs *flag.FlagSet, prefix string, val reflect.Value) {
	for i := 0; i < val.NumField(); i++ {
		tags := val.Type().Field(i).Tag
		name := tags.Get("flag")

		if name == "" {
			continue
		}

		usage := tags.Get("usage")

		switch val.Field(i).Kind() {
		case reflect.Struct:
			load(fs, prefix+name+"-", val.Field(i))
		case reflect.Bool, reflect.Int64, reflect.Float64, reflect.Int, reflect.Uint, reflect.Uint64, reflect.String:
			switch f := val.Field(i).Addr().Interface().(type) {
			case *bool:
				fs.BoolVar(f, prefix+name, *f, usage)
			case *time.Duration:
				fs.DurationVar(f, prefix+name, *f, usage)
			case *float64:
				fs.Float64Var(f, prefix+name, *f, usage)
			case *int:
				fs.IntVar(f, prefix+name, *f, usage)
			case *int64:
				fs.Int64Var(f, prefix+name, *f, usage)
			case *string:
				fs.StringVar(f, prefix+name, *f, usage)
			case *uint:
				fs.UintVar(f, prefix+name, *f, usage)
			case *uint64:
				fs.Uint64Var(f, prefix+name, *f, usage)
			}
		}
	}
}
