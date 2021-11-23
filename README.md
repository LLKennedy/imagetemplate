# Image Template

[![GoDoc](https://godoc.org/github.com/LLKennedy/imagetemplate?status.svg)](https://godoc.org/github.com/LLKennedy/imagetemplate)
[![Build Status](https://travis-ci.org/disintegration/imaging.svg?branch=master)](https://travis-ci.org/LLKennedy/imagetemplate)
![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/LLKennedy/imagetemplate.svg)
[![Coverage Status](https://coveralls.io/repos/github/LLKennedy/imagetemplate/badge.svg?branch=master)](https://coveralls.io/github/LLKennedy/imagetemplate?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/LLKennedy/imagetemplate)](https://goreportcard.com/report/github.com/LLKennedy/imagetemplate)
[![Maintainability](https://api.codeclimate.com/v1/badges/22d24397a4cccf8471d4/maintainability)](https://codeclimate.com/github/LLKennedy/imagetemplate/maintainability)
[![GitHub](https://img.shields.io/github/license/LLKennedy/imagetemplate.svg)](https://github.com/LLKennedy/imagetemplate/blob/master/LICENSE)

An image templating library for golang. Builder provides the templating engine, render provides the canvas, individual components provide the elements which can be templated and rendered. Component registration follows the pattern of built-in package "image", any package which implements the Component interface (found under render) and uses the RegisterComponent function during initialisation will be available for use.

Several default components are included, see their documentation as well as detailed usage guidelines for the template files and Builder in [the documentation pages](/doc/Home.md).

## Installation
`go get "github.com/LLKennedy/imagetemplate/v3"`

## Basic Usage
```
loader, props, err := imagetemplate.New().Load().FromFile("template.json")

// Check props here, set any discovered variables with real values.

data, err := loader.Write().ToBMP(props)
err = ioutil.WriteFile("output.bmp", data, os.ModeExclusive)
```

## Testing
On windows, the simplest way to test is to use the powershell script.

`./test.ps1`

To emulate the testing which occurs in build pipelines for linux and mac, run the following:

`go test ./... -race`

## Lessons Learned

This was one of my earliest go projects, and I learned a lot both while developing it and afterwards. I intend to come back to this project eventually and improve it, but in the meantime there are several aspects of it which are worth noting, since I would do these things differently if I was to complete this project now.

### Misuse of Interfaces

I was enamoured with the concept of interaces in Go, and used them in many situations in this project. While interfaces are very valuable for abstraction and inversion of control, the general rule I have learned with Go interfaces can be summarised as follows: Use interfaces for your dependencies, but export and return concrete types where possible. Exporting interfaces wrapping your concrete types, instead of the concrete types directly, makes your code brittle. Any change at all, including changes that would be backwards-compatible on a concrete type, will be breaking in an exported interface.

This principle follows the advice in official guides to idiomatic Go, but it's one that I misunderstood in my early Go development, so this library exports and returns a number of interfaces which should be concrete types.

### Naming Conventions

I wrote the core of this project over a long weekend, and didn't take much time to plan out a structure for the different components with clear and distinct names. As a result, many types and functions use words like "properties" and "component" and "value" in ways that overlap and become confusing. A good refactor of this project would address these naming issues and define some terms more clearly in this readme.

### Documentation of Integration

This library supports the definition and usage of third party "components" (extensions on a basic JSON definition and associated logic for rendering onto a canvas), but no documentation on how to write or use one. The internal components do use this system and can be used as an example, but this is not intuitive or pointed to end users beyond this paragraph. Given this system requires the use of init() functions for type registration, this should be comprehensively documented.

### Overly TDD

Test-Driven Design is a good concept, but if you refactor large chunks of code to be less readable and less functionally clear purely in order to reach 100% test coverage, this is a net negative to the codebase. In this project I was striving for 100% test coverage as a novelty, and this impacted the quality of the code.