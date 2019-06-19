package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/LLKennedy/imagetemplate/v3"
)

func main() {
	props, callback, err := imagetemplate.LoadTemplate("template.json")
	if err != nil {
		fmt.Printf("failed to load file: %v\n", err)
		os.Exit(1)
	}
	data, err := callback(props)
	err = ioutil.WriteFile("simple-static.bmp", data, os.ModeExclusive)
	if err != nil {
		fmt.Printf("failed to write file: %v\n", err)
		os.Exit(1)
	}
}
