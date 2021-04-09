import { BarcodeType } from "./barcode/barcode";
import { EncodingMode, ErrorCorrectionLevel } from "./barcode/qr";

export interface Point {
	x: number;
	y: number;
}

export class RGBA implements Colour {
	constructor(r: number = 0, g: number = 0, b: number = 0, a: number = 0) {
		this.R = r;
		this.G = g;
		this.B = b;
		this.A = a;
	}
	public R: number = 0;
	public G: number = 0;
	public B: number = 0;
	public A: number = 0;
	RGBA(): [number, number, number, number] {
		return [this.R, this.G, this.B, this.A];
	}
}

export interface Colour {
	RGBA(): [number, number, number, number];
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
	public PDFSecurityLevel: number = 0;
	/** QRLevel is required for qr barcodes*/
	public QRLevel: ErrorCorrectionLevel = ErrorCorrectionLevel.L;
	/** QRMode is required for qr barcodes*/
	public QRMode: EncodingMode = EncodingMode.Invalid;
}

export class CanvasWrapper {
	private ref: CanvasRenderingContext2D;
	public PPI: number = 0;
	constructor(ref: CanvasRenderingContext2D) {
		if (ref === undefined || ref === null) {
			throw new Error("canvas reference must not be null or undefined");
		}
		this.ref = ref;
	}
	public async SetUnderlyingImage(newImage: CanvasImageSource): Promise<void> {
		this.ref.drawImage(newImage, 0, 0, this.ref.canvas?.width ?? 0, this.ref.canvas?.height ?? 0);
	}
	public async GetUnderlyingImage(): Promise<ImageData> {
		return this.ref.getImageData(0, 0, this.ref.canvas?.width ?? 0, this.ref.canvas?.height ?? 0);
	}
	public async GetWidth(): Promise<number> {
		return this.ref.canvas?.width ?? 0;
	}
	public async GetHeight(): Promise<number> {
		return this.ref.canvas?.height ?? 0;
	}
	public async Rectangle(topLeft: Point, width: number, height: number, colour: Colour): Promise<void> {
		this.ref.fillStyle = colourToHex(this.ref, colour);
		this.ref.fillRect(topLeft.x, topLeft.y, width, height);
	}
	public async Circle(centre: Point, radius: number, colour: Colour): Promise<void> {
		this.ref.fillStyle = colourToHex(this.ref, colour);
		throw new Error("unimplemented");
	}
	public async Text(text: string, start: Point, typeFace: string, colour: Colour, maxWidth: number): Promise<void> {
		this.ref.font = typeFace;
		this.ref.fillText(text, start.x, start.y, maxWidth);
	}
	public async TryText(text: string, start: Point, typeFace: string, colour: Colour, maxWidth: number): Promise<boolean> {
		this.ref.font = typeFace;
		let measured = this.ref.measureText(text);
		return measured.width <= maxWidth;
	}
	public async DrawImage(start: Point, subImage: CanvasImageSource): Promise<void> {
		this.ref.drawImage(subImage, start.x, start.y);
	}
	public async Barcode(codeType: BarcodeType, content: Uint8Array, extra: BarcodeExtraData, start: Point, width: number, height: number, dataColour: Colour, bgColour: Colour): Promise<void> {
		throw new Error("unimplemented");
	}
}

export interface ICanvas extends CanvasWrapper { }

function colourToHex(ref: CanvasRenderingContext2D, colour: Colour): string {
	let [r, g, b, a] = colour.RGBA();
	if (r < 0 || g < 0 || b < 0 || a < 0 || r > 255 || g > 255 || b > 255 || a > 255) {
		throw new Error("R, G, B and A values must be between 0 and 255");
	}
	ref.globalAlpha = a / 255;
	return `rgba(${r}, ${g}, ${b}, ${a})`;
}