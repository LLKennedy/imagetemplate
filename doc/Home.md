# Image Template Documentation

## Introduction
This documentation is intended to cover the design and usage of the packages in this image templating library. This home page serves as a high level overview of the concepts and terminology.

The main readme file for the repository provides code for a basic usage example, but full self-contained examples can be found in the top level [examples folder](../examples). A compilation of the JSON template files used in those example applications, along with other sample template files, is located in the [Examples](Examples.md) file within this documentation.

## Contents

### 1. [Terminology](#Terminology)
### 2. [Template Files](#Template-Files)
### 3. [Builder](#Builder)
### 4. [Default Components](#Default-Components)
### 5. [Example Results](#Example-Results)

## Terminology

### Application
As a utility library, the assumption is made that the functions provided will always be called by some other application with its own purpose(s). The term "the application" refers to this external caller, and no assumptions are made about the custom variables this application is aware of or the usage of the data provided to it.

### Canvas
A Canvas is a target for rendering. It holds the cumulative result of each component rendering onto it in the specified order, and is responsible for producing the resultant image for use by the application.

### Component
A Component is an individual piece of the image, such as a rectangle, a text string, a photo or similar image embedded in the main image, etc. Components consist of a JSON representation which will appear in the JSON template file, an in-memory struct holding their properties, variables and state of processing, and the final rendered pixels in the resultant image.

## Template Files
See the main [Template File](TemplateFile.md) page for full detail.

## Builder
See the main [Builder](Builder.md) page for full detail.

## Default Components
The following components are built into the package and are always available to be used in template files.

### Barcode
Barcodes can be rendered with full RGBA control over both data and background channels in a variety of formats. See the main [Barcode](Barcode.md) page for full detail.

### Circle
Primitive circles of custom colour. See the main [Circle](Circle.md) page for full detail.

### DateTime
Timestamps of any granularity can be rendered as custom-formatted text. See the main [DateTime](DateTime.md) page for full detail.

### Image
Photos and other pre-rendered images implementing golang's image.Image interface can be scaled, transformed and cropped onto the canvas. See the main [Image](Image.md) page for full detail.

### Rectangle
Primitive rectangles of custom colour. See the main [Rectangle](Rectangle.md) page for full detail.

### Text
Single-line text can be rendered with any TrueType font, in custom colour, with automatic scaling down to a hard-set maximum width to prevent overrun. See the full [Text](Text.md) page for full detail.


## Example Results