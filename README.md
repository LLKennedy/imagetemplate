# Image Template
This project defines a template file for drawing custom images from pre-defined components. The intended application is smartcard printing, and some assumptions may be made with that in mind, but this format should be appropriate for general use.

## Testing
`go test . -covermode=count -coverprofile="coverage.out"; go tool cover -html="coverage.out"`