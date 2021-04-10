export enum BarcodeType {
	Aztec = "Aztec",
	Codabar = "Codabar",
	Code128 = "Code 128",
	Code39 = "Code 39",
	Code93 = "Code 93",
	DataMatrix = "DataMatrix",
	EAN8 = "EAN 8",
	EAN13 = "EAN 13",
	PDF = "PDF417",
	QR = "QR Code",
	TwoOfFive = "2 of 5",
	TwoOfFiveInterleaved = "2 of 5 (interleaved)"
}

export interface Metadata {
	CodeKind: string;
	Dimensions: number;
}

export interface Barcode {
	Metadata(): Metadata;
	Content(): string;
	Draw(ref: CanvasRenderingContext2D): Promise<void>;
}