export interface Point {
	x: number;
	y: number;
}

export interface Colour {
	RGBA(): [number, number, number, number]
}

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

export enum QRErrorCorrectionLevel {
	L = 0,
	M = 1,
	Q = 2,
	H = 3
}

export enum QREncodingMode {
	invalid = 0,
	NumericMode = 1,
	AlphaNumericMode = 2,
	ByteMode = 4,
	KanjiMode = 8
}

export class BarcodeExtraData {
	/** AztecMinECCPercent is required for aztec barcodes*/
	public AztecMinECCPercent: number = 0;
	/** AztecUserSpecifiedLayers is required for aztec barcodes*/
	public AztecUserSpecifiedLayers: number = 0;
	/** Code39IncludeChecksum is required for code39 barcodes*/
	public Code39IncludeChecksum: boolean = false;
	/** Code39FullASCIIMode is required for code39 barcodes*/
	public Code39FullASCIIMode: boolean = false;
	/** Code93IncludeChecksum is required for code93 barcodes*/
	public Code93IncludeChecksum: boolean = false;
	/** Code93FullASCIIMode is required for code93 barcodes*/
	public Code93FullASCIIMode: boolean = false;
	/** PDFSecurityLevel is required for pdf417 barcodes*/
	public PDFSecurityLevel: number;
	/** QRLevel is required for qr barcodes*/
	public QRLevel: QRErrorCorrectionLevel = QRErrorCorrectionLevel.L;
	/** QRMode is required for qr barcodes*/
	public QRMode: QREncodingMode = QREncodingMode.invalid;
}

export class CanvasWrapper {
	private ref: CanvasRenderingContext2D;
	public PPI: number;
	constructor(ref: CanvasRenderingContext2D) {
		this.ref = ref;
		ref.font
	}
	public SetUnderlyingImage(newImage: CanvasImageSource): Promise<void> {
		throw new Error("unimplemented");
	}
	public GetUnderlyingImage(): Promise<ImageData> {
		throw new Error("unimplemented");
	}
	public GetWidth(): Promise<number> {
		throw new Error("unimplemented");
	}
	public GetHeight(): Promise<number> {
		throw new Error("unimplemented");
	}
	public Rectangle(topLeft: Point, width: number, height: number, colour: Colour): Promise<void> {
		throw new Error("unimplemented");
	}
	public Circle(centre: Point, radius: number, colour: Colour): Promise<void> {
		throw new Error("unimplemented");
	}
	public Text(text: string, start: Point, typeFace: string, colour: Colour, maxWidth: number): Promise<void> {
		throw new Error("unimplemented");
	}
	public TryText(text: string, start: Point, typeFace: string, colour: Colour, maxWidth: number): Promise<boolean> {
		throw new Error("unimplemented");
	}
	public DrawImage(start: Point, subImage: CanvasImageSource): Promise<void> {
		throw new Error("unimplemented");
	}
	public Barcode(codeType: BarcodeType, content: Uint8Array, extra: BarcodeExtraData, start: Point, width: number, height: number, dataColour: Colour, bgColour: Colour): Promise<void> {
		throw new Error("unimplemented");
	}
}

export interface ICanvas extends CanvasWrapper { }