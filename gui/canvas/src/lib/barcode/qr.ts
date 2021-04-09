export enum ErrorCorrectionLevel {
	L = 0,
	M = 1,
	Q = 2,
	H = 3
}

export enum EncodingMode {
	Invalid = 0,
	NumericMode = 1,
	AlphaNumericMode = 2,
	ByteMode = 4,
	KanjiMode = 8,
}

export class versionInfo {
	constructor(Version: number = 0,
		Level: ErrorCorrectionLevel = ErrorCorrectionLevel.L,
		ErrorCorrectionCodewordsPerBlock: number = 0,
		NumberOfBlocksInGroup1: number = 0,
		DataCodeWordsPerBlockInGroup1: number = 0,
		NumberOfBlocksInGroup2: number = 0,
		DataCodeWordsPerBlockInGroup2: number = 0
	) {
		this.Version = Version;
		this.Level = Level;
		this.ErrorCorrectionCodewordsPerBlock = ErrorCorrectionCodewordsPerBlock;
		this.NumberOfBlocksInGroup1 = NumberOfBlocksInGroup1;
		this.DataCodeWordsPerBlockInGroup1 = DataCodeWordsPerBlockInGroup1;
		this.NumberOfBlocksInGroup2 = NumberOfBlocksInGroup2;
		this.DataCodeWordsPerBlockInGroup2 = DataCodeWordsPerBlockInGroup2;
	}
	public Version: number = 0;
	public Level: ErrorCorrectionLevel = ErrorCorrectionLevel.L;
	public ErrorCorrectionCodewordsPerBlock: number = 0;
	public NumberOfBlocksInGroup1: number = 0;
	public DataCodeWordsPerBlockInGroup1: number = 0;
	public NumberOfBlocksInGroup2: number = 0;
	public DataCodeWordsPerBlockInGroup2: number = 0;
	public totalDataBytes(): number {
		let g1Data = this.NumberOfBlocksInGroup1 * this.DataCodeWordsPerBlockInGroup1;
		let g2Data = this.NumberOfBlocksInGroup2 * this.DataCodeWordsPerBlockInGroup2;
		return g1Data + g2Data;
	}
	public charCountBits(m: EncodingMode): number {
		switch (m) {
			case EncodingMode.NumericMode:
				if (this.Version < 10) {
					return 10;
				} else if (this.Version < 27) {
					return 12;
				}
				return 14;
			case EncodingMode.AlphaNumericMode:
				if (this.Version < 10) {
					return 9;
				} else if (this.Version < 27) {
					return 11;
				}
				return 13;
			case EncodingMode.ByteMode:
				if (this.Version < 10) {
					return 8;
				}
				return 16;
			case EncodingMode.KanjiMode:
				if (this.Version < 10) {
					return 8;
				} else if (this.Version < 27) {
					return 10;
				}
				return 12;
			default:
				return 0;
		}
	}
}

const versionInfos: readonly Readonly<versionInfo>[] = [
	new versionInfo(1, ErrorCorrectionLevel.L, 7, 1, 19, 0, 0),
	new versionInfo(1, ErrorCorrectionLevel.M, 10, 1, 16, 0, 0),
	new versionInfo(1, ErrorCorrectionLevel.Q, 13, 1, 13, 0, 0),
	new versionInfo(1, ErrorCorrectionLevel.H, 17, 1, 9, 0, 0),
	new versionInfo(2, ErrorCorrectionLevel.L, 10, 1, 34, 0, 0),
	new versionInfo(2, ErrorCorrectionLevel.M, 16, 1, 28, 0, 0),
	new versionInfo(2, ErrorCorrectionLevel.Q, 22, 1, 22, 0, 0),
	new versionInfo(2, ErrorCorrectionLevel.H, 28, 1, 16, 0, 0),
	new versionInfo(3, ErrorCorrectionLevel.L, 15, 1, 55, 0, 0),
	new versionInfo(3, ErrorCorrectionLevel.M, 26, 1, 44, 0, 0),
	new versionInfo(3, ErrorCorrectionLevel.Q, 18, 2, 17, 0, 0),
	new versionInfo(3, ErrorCorrectionLevel.H, 22, 2, 13, 0, 0),
	new versionInfo(4, ErrorCorrectionLevel.L, 20, 1, 80, 0, 0),
	new versionInfo(4, ErrorCorrectionLevel.M, 18, 2, 32, 0, 0),
	new versionInfo(4, ErrorCorrectionLevel.Q, 26, 2, 24, 0, 0),
	new versionInfo(4, ErrorCorrectionLevel.H, 16, 4, 9, 0, 0),
	new versionInfo(5, ErrorCorrectionLevel.L, 26, 1, 108, 0, 0),
	new versionInfo(5, ErrorCorrectionLevel.M, 24, 2, 43, 0, 0),
	new versionInfo(5, ErrorCorrectionLevel.Q, 18, 2, 15, 2, 16),
	new versionInfo(5, ErrorCorrectionLevel.H, 22, 2, 11, 2, 12),
	new versionInfo(6, ErrorCorrectionLevel.L, 18, 2, 68, 0, 0),
	new versionInfo(6, ErrorCorrectionLevel.M, 16, 4, 27, 0, 0),
	new versionInfo(6, ErrorCorrectionLevel.Q, 24, 4, 19, 0, 0),
	new versionInfo(6, ErrorCorrectionLevel.H, 28, 4, 15, 0, 0),
	new versionInfo(7, ErrorCorrectionLevel.L, 20, 2, 78, 0, 0),
	new versionInfo(7, ErrorCorrectionLevel.M, 18, 4, 31, 0, 0),
	new versionInfo(7, ErrorCorrectionLevel.Q, 18, 2, 14, 4, 15),
	new versionInfo(7, ErrorCorrectionLevel.H, 26, 4, 13, 1, 14),
	new versionInfo(8, ErrorCorrectionLevel.L, 24, 2, 97, 0, 0),
	new versionInfo(8, ErrorCorrectionLevel.M, 22, 2, 38, 2, 39),
	new versionInfo(8, ErrorCorrectionLevel.Q, 22, 4, 18, 2, 19),
	new versionInfo(8, ErrorCorrectionLevel.H, 26, 4, 14, 2, 15),
	new versionInfo(9, ErrorCorrectionLevel.L, 30, 2, 116, 0, 0),
	new versionInfo(9, ErrorCorrectionLevel.M, 22, 3, 36, 2, 37),
	new versionInfo(9, ErrorCorrectionLevel.Q, 20, 4, 16, 4, 17),
	new versionInfo(9, ErrorCorrectionLevel.H, 24, 4, 12, 4, 13),
	new versionInfo(10, ErrorCorrectionLevel.L, 18, 2, 68, 2, 69),
	new versionInfo(10, ErrorCorrectionLevel.M, 26, 4, 43, 1, 44),
	new versionInfo(10, ErrorCorrectionLevel.Q, 24, 6, 19, 2, 20),
	new versionInfo(10, ErrorCorrectionLevel.H, 28, 6, 15, 2, 16),
	new versionInfo(11, ErrorCorrectionLevel.L, 20, 4, 81, 0, 0),
	new versionInfo(11, ErrorCorrectionLevel.M, 30, 1, 50, 4, 51),
	new versionInfo(11, ErrorCorrectionLevel.Q, 28, 4, 22, 4, 23),
	new versionInfo(11, ErrorCorrectionLevel.H, 24, 3, 12, 8, 13),
	new versionInfo(12, ErrorCorrectionLevel.L, 24, 2, 92, 2, 93),
	new versionInfo(12, ErrorCorrectionLevel.M, 22, 6, 36, 2, 37),
	new versionInfo(12, ErrorCorrectionLevel.Q, 26, 4, 20, 6, 21),
	new versionInfo(12, ErrorCorrectionLevel.H, 28, 7, 14, 4, 15),
	new versionInfo(13, ErrorCorrectionLevel.L, 26, 4, 107, 0, 0),
	new versionInfo(13, ErrorCorrectionLevel.M, 22, 8, 37, 1, 38),
	new versionInfo(13, ErrorCorrectionLevel.Q, 24, 8, 20, 4, 21),
	new versionInfo(13, ErrorCorrectionLevel.H, 22, 12, 11, 4, 12),
	new versionInfo(14, ErrorCorrectionLevel.L, 30, 3, 115, 1, 116),
	new versionInfo(14, ErrorCorrectionLevel.M, 24, 4, 40, 5, 41),
	new versionInfo(14, ErrorCorrectionLevel.Q, 20, 11, 16, 5, 17),
	new versionInfo(14, ErrorCorrectionLevel.H, 24, 11, 12, 5, 13),
	new versionInfo(15, ErrorCorrectionLevel.L, 22, 5, 87, 1, 88),
	new versionInfo(15, ErrorCorrectionLevel.M, 24, 5, 41, 5, 42),
	new versionInfo(15, ErrorCorrectionLevel.Q, 30, 5, 24, 7, 25),
	new versionInfo(15, ErrorCorrectionLevel.H, 24, 11, 12, 7, 13),
	new versionInfo(16, ErrorCorrectionLevel.L, 24, 5, 98, 1, 99),
	new versionInfo(16, ErrorCorrectionLevel.M, 28, 7, 45, 3, 46),
	new versionInfo(16, ErrorCorrectionLevel.Q, 24, 15, 19, 2, 20),
	new versionInfo(16, ErrorCorrectionLevel.H, 30, 3, 15, 13, 16),
	new versionInfo(17, ErrorCorrectionLevel.L, 28, 1, 107, 5, 108),
	new versionInfo(17, ErrorCorrectionLevel.M, 28, 10, 46, 1, 47),
	new versionInfo(17, ErrorCorrectionLevel.Q, 28, 1, 22, 15, 23),
	new versionInfo(17, ErrorCorrectionLevel.H, 28, 2, 14, 17, 15),
	new versionInfo(18, ErrorCorrectionLevel.L, 30, 5, 120, 1, 121),
	new versionInfo(18, ErrorCorrectionLevel.M, 26, 9, 43, 4, 44),
	new versionInfo(18, ErrorCorrectionLevel.Q, 28, 17, 22, 1, 23),
	new versionInfo(18, ErrorCorrectionLevel.H, 28, 2, 14, 19, 15),
	new versionInfo(19, ErrorCorrectionLevel.L, 28, 3, 113, 4, 114),
	new versionInfo(19, ErrorCorrectionLevel.M, 26, 3, 44, 11, 45),
	new versionInfo(19, ErrorCorrectionLevel.Q, 26, 17, 21, 4, 22),
	new versionInfo(19, ErrorCorrectionLevel.H, 26, 9, 13, 16, 14),
	new versionInfo(20, ErrorCorrectionLevel.L, 28, 3, 107, 5, 108),
	new versionInfo(20, ErrorCorrectionLevel.M, 26, 3, 41, 13, 42),
	new versionInfo(20, ErrorCorrectionLevel.Q, 30, 15, 24, 5, 25),
	new versionInfo(20, ErrorCorrectionLevel.H, 28, 15, 15, 10, 16),
	new versionInfo(21, ErrorCorrectionLevel.L, 28, 4, 116, 4, 117),
	new versionInfo(21, ErrorCorrectionLevel.M, 26, 17, 42, 0, 0),
	new versionInfo(21, ErrorCorrectionLevel.Q, 28, 17, 22, 6, 23),
	new versionInfo(21, ErrorCorrectionLevel.H, 30, 19, 16, 6, 17),
	new versionInfo(22, ErrorCorrectionLevel.L, 28, 2, 111, 7, 112),
	new versionInfo(22, ErrorCorrectionLevel.M, 28, 17, 46, 0, 0),
	new versionInfo(22, ErrorCorrectionLevel.Q, 30, 7, 24, 16, 25),
	new versionInfo(22, ErrorCorrectionLevel.H, 24, 34, 13, 0, 0),
	new versionInfo(23, ErrorCorrectionLevel.L, 30, 4, 121, 5, 122),
	new versionInfo(23, ErrorCorrectionLevel.M, 28, 4, 47, 14, 48),
	new versionInfo(23, ErrorCorrectionLevel.Q, 30, 11, 24, 14, 25),
	new versionInfo(23, ErrorCorrectionLevel.H, 30, 16, 15, 14, 16),
	new versionInfo(24, ErrorCorrectionLevel.L, 30, 6, 117, 4, 118),
	new versionInfo(24, ErrorCorrectionLevel.M, 28, 6, 45, 14, 46),
	new versionInfo(24, ErrorCorrectionLevel.Q, 30, 11, 24, 16, 25),
	new versionInfo(24, ErrorCorrectionLevel.H, 30, 30, 16, 2, 17),
	new versionInfo(25, ErrorCorrectionLevel.L, 26, 8, 106, 4, 107),
	new versionInfo(25, ErrorCorrectionLevel.M, 28, 8, 47, 13, 48),
	new versionInfo(25, ErrorCorrectionLevel.Q, 30, 7, 24, 22, 25),
	new versionInfo(25, ErrorCorrectionLevel.H, 30, 22, 15, 13, 16),
	new versionInfo(26, ErrorCorrectionLevel.L, 28, 10, 114, 2, 115),
	new versionInfo(26, ErrorCorrectionLevel.M, 28, 19, 46, 4, 47),
	new versionInfo(26, ErrorCorrectionLevel.Q, 28, 28, 22, 6, 23),
	new versionInfo(26, ErrorCorrectionLevel.H, 30, 33, 16, 4, 17),
	new versionInfo(27, ErrorCorrectionLevel.L, 30, 8, 122, 4, 123),
	new versionInfo(27, ErrorCorrectionLevel.M, 28, 22, 45, 3, 46),
	new versionInfo(27, ErrorCorrectionLevel.Q, 30, 8, 23, 26, 24),
	new versionInfo(27, ErrorCorrectionLevel.H, 30, 12, 15, 28, 16),
	new versionInfo(28, ErrorCorrectionLevel.L, 30, 3, 117, 10, 118),
	new versionInfo(28, ErrorCorrectionLevel.M, 28, 3, 45, 23, 46),
	new versionInfo(28, ErrorCorrectionLevel.Q, 30, 4, 24, 31, 25),
	new versionInfo(28, ErrorCorrectionLevel.H, 30, 11, 15, 31, 16),
	new versionInfo(29, ErrorCorrectionLevel.L, 30, 7, 116, 7, 117),
	new versionInfo(29, ErrorCorrectionLevel.M, 28, 21, 45, 7, 46),
	new versionInfo(29, ErrorCorrectionLevel.Q, 30, 1, 23, 37, 24),
	new versionInfo(29, ErrorCorrectionLevel.H, 30, 19, 15, 26, 16),
	new versionInfo(30, ErrorCorrectionLevel.L, 30, 5, 115, 10, 116),
	new versionInfo(30, ErrorCorrectionLevel.M, 28, 19, 47, 10, 48),
	new versionInfo(30, ErrorCorrectionLevel.Q, 30, 15, 24, 25, 25),
	new versionInfo(30, ErrorCorrectionLevel.H, 30, 23, 15, 25, 16),
	new versionInfo(31, ErrorCorrectionLevel.L, 30, 13, 115, 3, 116),
	new versionInfo(31, ErrorCorrectionLevel.M, 28, 2, 46, 29, 47),
	new versionInfo(31, ErrorCorrectionLevel.Q, 30, 42, 24, 1, 25),
	new versionInfo(31, ErrorCorrectionLevel.H, 30, 23, 15, 28, 16),
	new versionInfo(32, ErrorCorrectionLevel.L, 30, 17, 115, 0, 0),
	new versionInfo(32, ErrorCorrectionLevel.M, 28, 10, 46, 23, 47),
	new versionInfo(32, ErrorCorrectionLevel.Q, 30, 10, 24, 35, 25),
	new versionInfo(32, ErrorCorrectionLevel.H, 30, 19, 15, 35, 16),
	new versionInfo(33, ErrorCorrectionLevel.L, 30, 17, 115, 1, 116),
	new versionInfo(33, ErrorCorrectionLevel.M, 28, 14, 46, 21, 47),
	new versionInfo(33, ErrorCorrectionLevel.Q, 30, 29, 24, 19, 25),
	new versionInfo(33, ErrorCorrectionLevel.H, 30, 11, 15, 46, 16),
	new versionInfo(34, ErrorCorrectionLevel.L, 30, 13, 115, 6, 116),
	new versionInfo(34, ErrorCorrectionLevel.M, 28, 14, 46, 23, 47),
	new versionInfo(34, ErrorCorrectionLevel.Q, 30, 44, 24, 7, 25),
	new versionInfo(34, ErrorCorrectionLevel.H, 30, 59, 16, 1, 17),
	new versionInfo(35, ErrorCorrectionLevel.L, 30, 12, 121, 7, 122),
	new versionInfo(35, ErrorCorrectionLevel.M, 28, 12, 47, 26, 48),
	new versionInfo(35, ErrorCorrectionLevel.Q, 30, 39, 24, 14, 25),
	new versionInfo(35, ErrorCorrectionLevel.H, 30, 22, 15, 41, 16),
	new versionInfo(36, ErrorCorrectionLevel.L, 30, 6, 121, 14, 122),
	new versionInfo(36, ErrorCorrectionLevel.M, 28, 6, 47, 34, 48),
	new versionInfo(36, ErrorCorrectionLevel.Q, 30, 46, 24, 10, 25),
	new versionInfo(36, ErrorCorrectionLevel.H, 30, 2, 15, 64, 16),
	new versionInfo(37, ErrorCorrectionLevel.L, 30, 17, 122, 4, 123),
	new versionInfo(37, ErrorCorrectionLevel.M, 28, 29, 46, 14, 47),
	new versionInfo(37, ErrorCorrectionLevel.Q, 30, 49, 24, 10, 25),
	new versionInfo(37, ErrorCorrectionLevel.H, 30, 24, 15, 46, 16),
	new versionInfo(38, ErrorCorrectionLevel.L, 30, 4, 122, 18, 123),
	new versionInfo(38, ErrorCorrectionLevel.M, 28, 13, 46, 32, 47),
	new versionInfo(38, ErrorCorrectionLevel.Q, 30, 48, 24, 14, 25),
	new versionInfo(38, ErrorCorrectionLevel.H, 30, 42, 15, 32, 16),
	new versionInfo(39, ErrorCorrectionLevel.L, 30, 20, 117, 4, 118),
	new versionInfo(39, ErrorCorrectionLevel.M, 28, 40, 47, 7, 48),
	new versionInfo(39, ErrorCorrectionLevel.Q, 30, 43, 24, 22, 25),
	new versionInfo(39, ErrorCorrectionLevel.H, 30, 10, 15, 67, 16),
	new versionInfo(40, ErrorCorrectionLevel.L, 30, 19, 118, 6, 119),
	new versionInfo(40, ErrorCorrectionLevel.M, 28, 18, 47, 31, 48),
	new versionInfo(40, ErrorCorrectionLevel.Q, 30, 34, 24, 34, 25),
	new versionInfo(40, ErrorCorrectionLevel.H, 30, 20, 15, 61, 16),
]

function encodeAlphaNumeric(content: string, ecl: ErrorCorrectionLevel) {
	const contentLenIsOdd = content.length % 2 == 1;
	let contentBitCount = Math.floor(content.length / 2) * 11;
	if (contentLenIsOdd) {
		contentBitCount += 6;
	}
	let vi = findSmallestVersionInfo(ecl, EncodingMode.AlphaNumericMode, contentBitCount);
}

function findSmallestVersionInfo(ecl: ErrorCorrectionLevel, mode: EncodingMode, dataBits: number): Readonly<versionInfo> | undefined {
	dataBits += 4;
	for (let vi of versionInfos) {
		if (vi.Level === ecl) {
			if ((vi.totalDataBytes() * 8) >= (dataBits + vi.charCountBits(mode))) {
				return vi;
			}
		}
	}
	return undefined;
}