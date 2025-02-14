package main

import (
	"os"

	"github.com/hpetrov29/resttemplate/app/services/api/v1/cmd"
)

func main() {
	if err := cmd.Main(cmd.Routes()); err != nil {
		os.Exit(1)
	}
}