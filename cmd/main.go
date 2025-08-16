package main

import (
	"go-boilerplate/cmd/cli"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cli.Execute()
}
