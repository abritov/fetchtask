package main

import (
	"fetchtask"
)

func main() {
	generateId := fetchtask.NewIdGeneratorMock()
	server := fetchtask.NewServer(generateId)
	server.Listen()
}
