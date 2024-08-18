package main

import (
	"fmt"

	"github.com/msalemor/gorag/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Println(err)
	}
}
