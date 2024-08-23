package main

import (
	"fmt"

	"github.com/bernielomax/dagger-hello/internal/client"
)

func main() {
	out, err := client.Get("https://icanhazip.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}
