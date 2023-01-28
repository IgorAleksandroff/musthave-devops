package main

import "os"

func main() {
	os.Exit(0) // want "function main should not have os exit"
}
