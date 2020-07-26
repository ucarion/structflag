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
