import { BitList } from "./utils";

describe("bitlist", () => {
	it("behaves as expected across a range of calls", async () => {
		let bl = new BitList();
		bl.AddBit([true, true, false, true]);
		expect(bl.GetBit(0)).toBe(true);
		expect(bl.GetBit(1)).toBe(true);
		expect(bl.GetBit(2)).toBe(false);
		expect(bl.GetBit(3)).toBe(true);
		bl.AddByte(1);
		expect(bl.GetBit(10)).toBe(false);
		expect(bl.GetBit(11)).toBe(true);
		expect(bl.GetBit(12)).toBe(false);
		expect(bl.GetBytes()).toStrictEqual(new Uint8Array([0b11010000, 0b00010000]));
		expect(bl.Len()).toBe(12);
		bl.AddBits(0b11111101, 4)
		expect(bl.GetBytes()).toStrictEqual(new Uint8Array([0b11010000, 0b00011101]));
	})
});