package main

import "github.com/bernielomax/dagger-hello/internal/client"

func main() {
	client.Get("https://google.com")
}
