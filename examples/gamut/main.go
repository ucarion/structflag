package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ucarion/structflag"
)

func main() {
	config := struct {
		Bool     bool
		Float64  float64
		Int      int
		Uint     uint
		Int64    int64
		Uint64   uint64
		Duration time.Duration
	}{
		Bool:     true,
		Float64:  42,
		Int:      42,
		Uint:     42,
		Int64:    42,
		Uint64:   42,
		Duration: time.Minute,
	}

	structflag.Load(&config)
	flag.Parse()

	fmt.Println(config)
}
