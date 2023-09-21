package main

import "github.com/alecthomas/kong"

var CLI struct {
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	default:
		panic(ctx.Command())
	}
}
