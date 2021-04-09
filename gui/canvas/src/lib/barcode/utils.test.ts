import { BitList } from "./utils";

describe("bitlist", () => {
	it("behaves as expected across a range of calls", async () => {
		let bl = new BitList();
		bl.AddBit([true, true, false, true]);
		expect(bl.GetBit(0)).toBe(true);
		expect(bl.GetBit(1)).toBe(true);
		expect(bl.GetBit(2)).toBe(false);
		expect(bl.GetBit(3)).toBe(true);
	})
});