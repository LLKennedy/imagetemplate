# <a name="imagetemplatefile"></a>Image Template File
An Image Template is a JSON-formatted file, comprised of a individual components in the order they should be rendered, and conditions which must be met for each component to render.

Please note the standard below is the intended file format and implementation for the 1.0.0 release, features which may not yet be fully functional have been *italicised*.

## Sample
The following is a stripped-down sample showing only the most basic elements of a template file. For specific components/conditionals please refer to the individual component references, and for full examples of files and the rendered output, please refer to the [Examples](Examples.md) page.

```
{
	"baseImage": {
		"data": "$appPhoto$",
		"width": "1920",
		"height": "1080"
	}
	"components": [
		{
			"type": "circle",
			"properties": {
				
			}
		},
		{
			"type": "barcode",
			"properties" {

			},
			conditional: {
				
			}
		}
	]
}
```

## <a name="terminology"></a>Terminology
### <a name="templatefile"></a>Template File
The JSON file created to render images for a specific application

### <a name="tendererandorparser"></a>Renderer and/or Parser
The software which parses a Template File and presents rendering methods to an Application. This repository contains a golang implementation of an Image Template Renderer.

### <a name="application"></a>Application
A script, executable or other program using the API of the renderer to load and render template files.

## <a name="databalues"></a>Data Values
In version 1.0.0 and earlier, all raw values in the template file must be strings. This allows the parser to easily identify named variables to present to the application

## <a name="namedproperties"></a>Named Properties / Named Variables / Custom Variables
In almost all primitive properties (that is, non-JSON properties) either a raw primitive value or a variable can be declared. Variables are wrapped by dollar symbols in the form `$variableName$`. Any string which does not contain dollar symbols is a valid variable name.

Custom variables declared in the template file will bubble through to the application, available through the `GetNamedPropertiesList()` method.

## <a name="structure"></a>Structure
The basic structure of the file is as follows:

### <a name="baseimage"></a>1. Base Image
- [`baseImage`](#baseimage): JSON structure
The base image upon which to render all other components. This can be a filename, a byte array or a rectangle of a pure colour. Mandatory components are exactly one of [`data`](#data), [`baseColour`](#basecolour) and [`fileName`](#filename). Optional components are [`width`](#widthandheight), [`height`](#widthandheight) and `components`(#Components). A value or variable declared in more than one of the mandatory type properties is invalid and will not render any components.

#### <a name="widthandheight"></a>Width and Height
- [`width`](#widthandheight): string-encoded integer in pixels
- [`height`](#widthandheight): string-encoded integer in pixels

The width and height of the base image.

For a filename or byte array base image, a width and height are optional, and will result in scaling the image to fit the desired parameters. Zero for width or height is an invalid value, but *negative values in one or the other will lock the positively-specified dimension and automatically scale the other to maintain aspect ratio*. For example, `width: "200", height: "-1"` will result in an image which is exactly 200 pixels wide but will scale height to maintain the aspect ratio of the input image.

For a coloured rectange base image, both width and height must be specified exactly with positive integers.

#### <a name="filename"></a>File Name
- [`fileName`](#filename): string representing a name of or path to the file

The path or name should be either absolute, or relative to the application.

#### <a name="data"></a>Data
- [`data`](#data): base64-encoded raw image bytes

The raw image data, which must be formatted as one of png, jpg, bmp, tiff. Please submit an issue if further image format support is desired.

*If specifying this value as a variable rather than a base64-encoded string, a byte array may be passed instead of base64-encoded data.*

#### <a name="basecolour"></a>Base Colour
- [`baseColour`](#basecolour): JSON structure

The NRGBA colour value of the base image rectangle.

##### <a name="rgba"></a>Red/Green/Blue/Alpha
- `R`: string-encoded uint8
- `G`: string-encoded uint8
- `B`: string-encoded uint8
- `A`: string-encoded uint8

The non-premultiplied RGBA values of the colour.

### <a name="components"></a>2. Components
- `component`: Ordered array of JSON structures

Each component object has the mandatory components `type` and `properties`, and the optional component `conditional`

#### <a name="type"></a>Type
- `type`: String matching the a known component type.

Valid options are [`circle`](Rectangle.md), [`text`](Text.md), [`image`](Image.md), [`barcode`](Barcode.md), and [*`dateTime`*](DateTime.md). See each relevant page for further detail.

#### <a name="properties"></a>Properties
- `properties`: JSON structure

Properties used to render the specific component.

The contents of `properties` varies wildly with component types, see each component's document (linked in [Type](#type)) for specifics.

#### <a name="conditional"></a>Conditional
- `conditional`: JSON structure

The conditions under which the component will render. See the main [Conditional](Conditional.md) page for more information.