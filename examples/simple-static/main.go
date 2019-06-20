package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/LLKennedy/imagetemplate/v3"
	"golang.org/x/tools/godoc/vfs"
)

func main() {
	loader, props, err := imagetemplate.NewUsing(vfs.OS("./")).Load().FromFile("template.json")
	if err != nil {
		fmt.Printf("failed to load file: %v\n", err)
		os.Exit(1)
	}
	data, err := loader.Write().ToBMP(props)
	err = ioutil.WriteFile("simple-static.bmp", data, os.ModeExclusive)
	if err != nil {
		fmt.Printf("failed to write file: %v\n", err)
		os.Exit(1)
	}
}
