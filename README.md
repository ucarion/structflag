# structflag

[![PkgGoDev](https://pkg.go.dev/badge/github.com/ucarion/structflag)](https://pkg.go.dev/github.com/ucarion/structflag)

`structflag` is a Golang package that lets you conveniently use the standard
library's `flag` package with structs. In other words, `structflag` CLI-ifies
your config structs with a single line of code.

Just create a struct for your config, optionally put a `` `flag:"..."` `` and/or
`` `usage:"..."` `` tag on the fields you want to populate from flags, and then
call `structflag.Load` with a pointer to your struct.

The main benefits of this package are:

1. Putting all your config in a single struct. You can compose this struct on
   top of "sub-configs", and you can use these config structs in tests too.
2. A terser syntax for invoking the standard library's `flag` package.

## Installation

To use `structflag` in your program, run:

```bash
go get github.com/ucarion/structflag
```

## Basic Usage

In its simplest form, you can usually just invoke `structflag.Load` with an
existing config to CLI-ify it.

```go
type config struct {
	FirstName string
	LastName  string
}

func main() {
	var conf config
	structflag.Load(&conf)
	flag.Parse()

	fmt.Println(conf)
}
```

This is already in [`./examples/minimal/main.go`](./examples/minimal/main.go) in
this repo, so you can run it as:

```text
$ go run ./examples/minimal/... --help
Usage of /[snip]/minimal:
  -FirstName string

  -LastName string

exit status 2

$ go run ./examples/minimal/... -FirstName john -LastName doe
{john doe}
```

## Advanced Usage

For the common case of a relatively simple struct that you want to customize
flag and usage details for, and where you already have some default values, you
can do something like this:

```go
type name struct {
	FirstName string `flag:"first" usage:"first name"`
	LastName  string `flag:"last" usage:"last name"`
}

func main() {
	config := struct {
		Name  name          `flag:"name"`
		Count int           `flag:"count" usage:"how many times to say hello"`
		Wait  time.Duration `flag:"wait" usage:"how long to wait before greeting"`
	}{
		Count: 3,
		Wait:  1 * time.Second,
	}

	structflag.Load(&config)
	flag.Parse()

	time.Sleep(config.Wait)
	for i := 0; i < config.Count; i++ {
		fmt.Println("hello", config.Name)
	}
}
```

This is already in [`./examples/detailed/main.go`](./examples/detailed/main.go)
in this repo, so you can run it as:

```text
$ go run ./examples/detailed/... --help
Usage of /[snip]/detailed:
  -count int
    	how many times to say hello (default 3)
  -name-first string
    	first name
  -name-last string
    	last name
  -wait duration
    	how long to wait before greeting (default 1s)

$ go run ./examples/detailed/... --count 5 --name-first muhammad --name-last al-khwarizmi
hello {muhammad al-khwarizmi}
hello {muhammad al-khwarizmi}
hello {muhammad al-khwarizmi}
hello {muhammad al-khwarizmi}
hello {muhammad al-khwarizmi}
```
