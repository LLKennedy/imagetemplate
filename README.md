# Image Template
This project defines a template file for drawing custom images from pre-defined components. The intended application is smartcard printing, and some assumptions may be made with that in mind, but this format should be appropriate for general use.

[![GoDoc](https://godoc.org/github.com/LLKennedy/imagetemplate?status.svg)](https://godoc.org/github.com/LLKennedy/imagetemplate)
[![Build Status](https://travis-ci.org/disintegration/imaging.svg?branch=master)](https://travis-ci.org/LLKennedy/imagetemplate)
[![Coverage Status](https://coveralls.io/repos/github/LLKennedy/imagetemplate/badge.svg?branch=master)](https://coveralls.io/github/LLKennedy/imagetemplate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/LLKennedy/imagetemplate)](https://goreportcard.com/report/github.com/LLKennedy/imagetemplate)

## Testing
`go test . -covermode=count -coverprofile="coverage.out"; go tool cover -html="coverage.out"`