import { BarcodeType } from "./barcode/barcode";
import { EncodingMode, ErrorCorrectionLevel } from "./barcode/qr";
import { Mutex, IMutex } from "@llkennedy/mutex.js";

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
	private ctx: CanvasRenderingContext2D;
	public PPI: number = 0;
	private mx: IMutex = new Mutex();
	constructor(ctx: CanvasRenderingContext2D) {
		if (ctx === undefined || ctx === null) {
			throw new Error("canvas ctxerence must not be null or undefined");
		}
		this.ctx = ctx;
	}
	public async SetUnderlyingImage(newImage: CanvasImageSource): Promise<void> {
		await this.mx.Run(() => {
			this.ctx.drawImage(newImage, 0, 0, this.ctx.canvas?.width ?? 0, this.ctx.canvas?.height ?? 0);
		})
	}
	public async GetUnderlyingImage(): Promise<ImageData> {
		return await this.mx.Run(() => {
			return this.ctx.getImageData(0, 0, this.ctx.canvas?.width ?? 0, this.ctx.canvas?.height ?? 0);
		})
	}
	public async GetWidth(): Promise<number> {
		return await this.mx.Run(() => {
			return this.ctx.canvas?.width ?? 0;
		})
	}
	public async GetHeight(): Promise<number> {
		return await this.mx.Run(() => {
			return this.ctx.canvas?.height ?? 0;
		})
	}
	public async Rectangle(topLeft: Point, width: number, height: number, colour: Colour): Promise<void> {
		await this.mx.Run(() => {
			this.ctx.fillStyle = colourToHex(this.ctx, colour);
			this.ctx.fillRect(topLeft.x, topLeft.y, width, height);
		})
	}
	public async Circle(centre: Point, radius: number, colour: Colour): Promise<void> {
		await this.mx.Run(() => {
			this.ctx.fillStyle = colourToHex(this.ctx, colour);
			this.ctx.beginPath();
			this.ctx.arc(centre.x, centre.y, radius, 0, 2 * Math.PI);
			this.ctx.fill();
		})
	}
	public async Text(text: string, start: Point, typeFace: string, colour: Colour, maxWidth: number): Promise<void> {
		await this.mx.Run(() => {
			this.ctx.font = typeFace;
			this.ctx.fillStyle = colourToHex(this.ctx, colour);
			this.ctx.fillText(text, start.x, start.y, maxWidth);
		})
	}
	/** TryText measures  */
	public async MeasureText(text: string, typeFace: string, maxWidth: number): Promise<number> {
		return await this.mx.Run(() => {
			this.ctx.font = typeFace;
			return this.ctx.measureText(text).width;
		})
	}
	/** Draw another image on top of this one */
	public async DrawImage(start: Point, subImage: CanvasImageSource): Promise<void> {
		await this.mx.Run(() => {
			this.ctx.drawImage(subImage, start.x, start.y);
		})
	}
	/** UNIMPLMENTED, DO NOT USE */
	public async Barcode(codeType: BarcodeType, content: Uint8Array, extra: BarcodeExtraData, start: Point, width: number, height: number, dataColour: Colour, bgColour: Colour): Promise<void> {
		await this.mx.Run(() => {
			throw new Error("barcodes are not implemented");
		})
	}
}

export interface ICanvas extends CanvasWrapper { }

function colourToHex(ctx: CanvasRenderingContext2D, colour: Colour): string {
	let [r, g, b, a] = colour.RGBA();
	if (r < 0 || g < 0 || b < 0 || a < 0 || r > 255 || g > 255 || b > 255 || a > 255) {
		throw new Error("R, G, B and A values must be between 0 and 255");
	}
	ctx.globalAlpha = a / 255;
	return `rgba(${r}, ${g}, ${b}, ${a})`;
}