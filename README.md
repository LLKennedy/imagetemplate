# Image Template

[![GoDoc](https://godoc.org/github.com/LLKennedy/imagetemplate?status.svg)](https://godoc.org/github.com/LLKennedy/imagetemplate)
[![Build Status](https://travis-ci.org/disintegration/imaging.svg?branch=master)](https://travis-ci.org/LLKennedy/imagetemplate)
[![Coverage Status](https://coveralls.io/repos/github/LLKennedy/imagetemplate/badge.svg?branch=master)](https://coveralls.io/github/LLKennedy/imagetemplate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/LLKennedy/imagetemplate)](https://goreportcard.com/report/github.com/LLKennedy/imagetemplate)

An image templating library for golang. Builder provides the templating engine, render provides the canvas, individual components provide the elements which can be templated and rendered. Component registration follows the pattern of built-in package "image", any package which implements the Component interface (found under render) and uses the RegisterComponent function during initialisation will be available for use.

Several default components are included, see their documentation as well as detailed usage guidelines for the template files and Builder in [the documentation pages](/doc/Home.md).

## Installation
`go get "github.com/LLKennedy/imagetemplate/v2"`

## Basic Usage
```
// Load your custom template file
newBuilder := imagetemplate.NewBuilder()
newBuilder, err := newBuilder.LoadComponentsFile("myTemplate.json")

// Handle file parsing errors here

// Extract the "named properties" (custom variables) from the file then insert corresponding values according to your application logic

customVariables := newBuilder.GetNamedPropertiesList()

// Process your custom variables here
// eg customVariables["username"] = "John Smith"

newBuilder, err = newBuilder.SetNamedProperties(customVariables)

// Handle value parsing errors here

// Write the loaded components to the canvas
newBuilder, err = newBuilder.ApplyComponents()

// Handle rendering errors here

// Export the canvas to a BMP image. 
// The underlying image.Image object can also be extracted for export to any other image format using newBuilder.GetCanvas().GetUnderlyingImage()
imgBytes, err := newBuilder.WriteToBMP()
// imgBytes now contains the rendered BMP image
```

## Testing
On windows, the simplest way to test is to use the powershell script.

`./test.ps1`

To emulate the testing which occurs in build pipelines for linux and mac, run the following:

`go test . ./components/... ./internal/filesystem ./render/... -race -coverprofile=coverage.out;`