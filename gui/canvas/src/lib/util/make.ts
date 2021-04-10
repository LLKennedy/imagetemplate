export function make<T>(n: () => T, len: number): T[] {
	let res: T[] = [];
	for (let i = 0; i < len; i++) {
		res.push(n());
	}
	return res;
}