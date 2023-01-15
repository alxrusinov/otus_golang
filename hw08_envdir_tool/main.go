package main

import (
	"errors"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal(errors.New("must be more than 2 arguments"))
	}

	dir := os.Args[1]
	cmd := os.Args[2:]

	environment, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(cmd, environment)

	os.Exit(returnCode)
}
