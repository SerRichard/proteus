package main

import (
  "fmt"
  "os"

  "github.com/SerRichard/proteus/cli/cmd"
)

func main() {

  if err := cmd.ProteusCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}