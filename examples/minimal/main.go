package main

import (
	"flag"
	"fmt"

	"github.com/ucarion/structflag"
)

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
