package main

import (
	img "github.com/LLKennedy/imagetemplate"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	log.Println("Starting imagetemplate demo")
	var builder img.Builder
	builder = img.NewBuilder(100, 100)
	bytes, err := builder.WriteToBMP()
	if err != nil {
		log.Fatalf("Failed to write canvas to bitmap: %v", err)
	}
	err = ioutil.WriteFile("demo.bmp", bytes, os.ModeExclusive)
	if err != nil {
		log.Fatalf("Failed to write bitmap to file: %v", err)
	}
}
