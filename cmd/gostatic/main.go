package main

import (
	"context"
	"fmt"

	"github.com/ericselin/gostatic"
)

func main() {
	ctx := context.Background()
	err := gostatic.Build(
		"https://github.com/ericselin/ericselin",
		"example",
		"denoland/deno",
		[]string{"run", "-A", "bob.ts"},
		ctx,
	)
	if err != nil {
		panic(err)
	}
	fmt.Println("Built site successfully")
}
