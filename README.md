# Image Template

[![GoDoc](https://godoc.org/github.com/LLKennedy/imagetemplate?status.svg)](https://godoc.org/github.com/LLKennedy/imagetemplate)
[![Build Status](https://travis-ci.org/disintegration/imaging.svg?branch=master)](https://travis-ci.org/LLKennedy/imagetemplate)
[![Coverage Status](https://coveralls.io/repos/github/LLKennedy/imagetemplate/badge.svg?branch=master)](https://coveralls.io/github/LLKennedy/imagetemplate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/LLKennedy/imagetemplate)](https://goreportcard.com/report/github.com/LLKennedy/imagetemplate)

An image templating library for golang. Builder provides the templating engine, render provides the canvas, individual components provide the elements which can be templated and rendered. Component registration follows the pattern of built-in package "image", any package which implements the Component interface (found under render) and uses the RegisterComponent function during initialisation will be available for use.

Several default components are included, see their documentation as well as detailed usage guidelines for the template files and Builder in [the documentation pages](/doc/Home.md).

## Installation
`go get "github.com/LLKennedy/imagetemplate/v3"`

## Basic Usage
```
loader, props, err := imagetemplate.New().Load().FromFile("template.json")

// Check props here, set any discovered variables with real values

data, err := loader.Write().ToBMP(props)
err = ioutil.WriteFile("output.bmp", data, os.ModeExclusive)
```

## Testing
On windows, the simplest way to test is to use the powershell script.

`./test.ps1`

To emulate the testing which occurs in build pipelines for linux and mac, run the following:

`go test . ./components/... ./render/... ./scaffold/... ./internal/filesystem -race -coverprofile=coverage.out;`