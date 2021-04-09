export class BitList {
	private count: number = 0;
	private data: Uint32Array = new Uint32Array();
	public Len(): number {
		return this.count;
	}
	private grow() {
		let growBy = this.data.length;
		if (growBy < 128) {
			growBy = 128;
		} else if (growBy > 1024) {
			growBy = 1024;
		}
		let newVals: number[] = [];
		for (let i = 0; i < growBy; i++) {
			newVals.push(0);
		}
		let newData = new Uint32Array(this.data.length + growBy);
		newData.set(this.data, 0);
		this.data = newData;
	}
	public AddBit(bits: boolean[]) {
		for (let bit of bits) {
			let itmIndex = Math.floor(this.count / 32);
			while (itmIndex >= this.data.length) {
				this.grow();
			}
			this.SetBit(this.count, bit);
			this.count++;
		}
	}
	public SetBit(index: number, value: boolean) {
		let itmIndex = Math.floor(index / 32);
		let itmBitShift = 31 - (index % 32);
		if (value) {
			this.data[itmIndex] = this.data[itmIndex] | (1 << itmBitShift);
		} else {
			this.data[itmIndex] = this.data[itmIndex] & (0xFFFFFFFF ^ (1 << itmBitShift))
		}
	}
	public GetBit(index: number): boolean {
		let itmIndex = Math.floor(index / 32);
		let itmBitShift = 31 - (index % 32);
		return ((this.data[itmIndex] >>> itmBitShift) & 1) === 1;
	}
	public AddByte(b: number) {
		// Enforce byte
		b = 0x000000FF & b;
		for (let i = 7; i >= 0; i--) {
			this.AddBit([((b >>> i) & 1) === 1]);
		}
	}
	public AddBits(b: number, count: number) {
		for (let i = count; i >= 0; i--) {
			this.AddBit([((b >>> i) & 1) === 1]);
		}
	}
	public GetBytes(): Uint8Array {
		let len = this.count >>> 3;
		if ((this.count % 8) !== 0) {
			len++;
		}
		let result = new Uint8Array(len);
		for (let i = 0; i < len; i++) {
			let shift = (3 - (i % 4)) * 8;
			result[i] = this.data[Math.floor(i / 4)] >>> shift & 0xFF;
		}
		return result
	}
}