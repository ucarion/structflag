# structflag

`structflag` is a Golang package that lets you conveniently use the standard
library's `flag` package with structs. Just create a struct for your config, put
a `flag:"..."` tag on the fields you want to populate from flags, and then call
`structflag.Load` with a pointer to your struct.

The main benefits of this package are:

1. Putting all your config in a single struct. You can compose this struct on
   top of "sub-configs", and you can use these config structs in tests too.
2. A terser syntax for invoking the standard library's `flag` package.

## Example

Here's a simple program that uses `structflag` to avoid having to do a bunch of
`flag.String`, `flag.Int`, `flag.Duration` calls:

```go
package main

import (
  "flag"
  "fmt"
  "time"

  "github.com/ucarion/structflag"
)

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

Assuming you put this in a file called `./examples/simple/main.go` ([like the
one you can find in this repo](./examples/simple/main.go)), you can invoke it as
so:

```text
$ go run ./examples/simple/... --help
Usage of simple:
  -count int
    	how many times to say hello (default 3)
  -name-first string
    	first name
  -name-last string
    	last name
  -wait duration
    	how long to wait before greeting (default 1s)
exit status 2

$ go run ./examples/simple/... --name-first=john --name-last=doe
hello {john doe}
hello {john doe}
hello {john doe}
```

By default, `Load` writes to `flag.CommandLine`, the "default" `flag.FlagSet`.
If you prefer to use a different `flag.FlagSet`, use `LoadTo` instead:

```go
flagSet := flag.NewFlagSet("my-cool-flagset", flag.PanicOnError)
structflag.LoadTo(&flagSet, ...)
```
