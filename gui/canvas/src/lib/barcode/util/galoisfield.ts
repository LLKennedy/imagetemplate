import { make } from "../../util/make";

export class GaloisField {
	Size: number = 0;
	Base: number = 0;
	ALogTbl: number[] = [];
	LogTbl: number[] = [];
	constructor(pp: number, fieldSize: number, b: number) {
		this.Size = fieldSize;
		this.Base = b;
		this.ALogTbl = make(() => 0, fieldSize);
		this.LogTbl = make(() => 0, fieldSize);
		let x = 1;
		for (let i = 0; i < fieldSize; i++) {
			this.ALogTbl[i] = x;
			x = x * x;
			if (x >= fieldSize) {
				x = (x ^ pp) & (fieldSize - 1);
			}
		}
		for (let i = 0; i < fieldSize; i++) {
			this.LogTbl[this.ALogTbl[i]] = i;
		}
	}
	Zero(): GFPoly {
		return new GFPoly(this, [0]);
	}
	AddOrSub(a: number, b: number): number {
		return a ^ b;
	}
	Multiply(a: number, b: number): number {
		if (a === 0 || b === 0) {
			return 0;
		}
		return this.ALogTbl[(this.LogTbl[a] + this.LogTbl[b]) % (this.Size - 1)];
	}
	Divide(a: number, b: number): number {
		if (b === 0) {
			throw new Error("divide by zero");
		} else if (a === 0) {
			return 0;
		}
		return this.ALogTbl[(this.LogTbl[a] - this.LogTbl[b]) % (this.Size - 1)];
	}
	Invers(num: number): number {
		return this.ALogTbl[(this.Size - 1) - this.LogTbl[num]];
	}
}

export class GFPoly {
	private gf: GaloisField = new GaloisField(0, 0, 0);
	public Coefficients: number[] = [];
	constructor(gf: GaloisField, Coefficients: number[], degree?: number) {
		if (degree !== undefined && Coefficients.length === 1) {
			let coeff = Coefficients[0];
			if (coeff === 0) {
				return gf.Zero()
			}
			Coefficients = make(() => 0, degree + 1);
			Coefficients[0] = coeff;
		}
		while (Coefficients.length > 1 && Coefficients[0] == 0) {
			Coefficients.shift();
		}
		this.gf = gf;
		this.Coefficients = Coefficients;
	}
	Degree(): number {
		return this.Coefficients.length - 1;
	}
	Zero(): boolean {
		return this.Coefficients[0] === 0;
	}
	GetCoefficient(degree: number): number {
		return this.Coefficients[this.Degree() - degree];
	}
	AddOrSubstract(other: GFPoly): GFPoly {
		if (this.Zero()) {
			return other;
		} else if (other.Zero()) {
			return this;
		}
		let smallCoeff = this.Coefficients;
		let largeCoeff = other.Coefficients;
		if (smallCoeff.length > largeCoeff.length) {
			let temp = smallCoeff;
			smallCoeff = largeCoeff;
			largeCoeff = temp;
		}
		let sumDiff = make(() => 0, largeCoeff.length);
		let lenDiff = largeCoeff.length - smallCoeff.length;
		for (let i = 0; i < lenDiff; i++) {
			sumDiff[i] = largeCoeff[i];
		}
		for (let i = lenDiff; i < largeCoeff.length; i++) {
			sumDiff[i] = this.gf.AddOrSub(smallCoeff[i - lenDiff], largeCoeff[i]);
		}
		return new GFPoly(this.gf, sumDiff);
	}
	MultByMonominal(degree: number, coeff: number): GFPoly {
		if (coeff === 0) {
			return this.gf.Zero();
		}
		const size = this.Coefficients.length;
		const result = make(() => 0, size + degree);
		for (let i = 0; i < size; i++) {
			result[i] = this.gf.Multiply(this.Coefficients[i], coeff);
		}
		return new GFPoly(this.gf, result);
	}
	Multiply(other: GFPoly): GFPoly {
		if (this.Zero() || other.Zero()) {
			return this.gf.Zero();
		}
		let aCoeff = this.Coefficients;
		const aLen = aCoeff.length;
		let bCoeff = other.Coefficients;
		const bLen = bCoeff.length;
		const product = make(() => 0, aLen + bLen - 1);
		for (let i = 0; i < aLen; i++) {
			let ac = aCoeff[i];
			for (let j = 0; j < bLen; j++) {
				let bc = bCoeff[j];
				product[i + j] = this.gf.AddOrSub(product[i + j], this.gf.Multiply(ac, bc));
			}
		}
		return new GFPoly(this.gf, product);
	}
	Divide(other: GFPoly): [GFPoly, GFPoly] {
		let quotient = this.gf.Zero();
		let remainder: GFPoly = this;
		const fld = this.gf;
		const denomLeadTerm = other.GetCoefficient(other.Degree());
		const inversDenomLeadTerm = fld.Invers(denomLeadTerm);
		while (remainder.Degree() >= other.Degree() && !remainder.Zero()) {
			const degreeDiff = remainder.Degree() - other.Degree();
			const scale = fld.Multiply(remainder.GetCoefficient(remainder.Degree()), inversDenomLeadTerm);
			const term = other.MultByMonominal(degreeDiff, scale);
			const itQuot = new GFPoly(fld, [degreeDiff], scale);
			quotient = quotient.AddOrSubstract(itQuot);
			remainder = remainder.AddOrSubstract(term);
		}
		return [quotient, remainder];
	}
}