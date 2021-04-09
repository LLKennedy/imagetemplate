export interface Component {
	Write(canvas: CanvasRenderingContext2D): Promise<void>;
	SetNamedProperties(properties: ReadonlyMap<string, any>): Promise<void>;
	GetJSONFormat(): Promise<Object>;
	VerifyAndSetJSONData(data: Object): Promise<void>;
}

export enum ConditionalOperator {
	invalid = "",
	Equals = "equals",
	Contains = "contains",
	StartsWith = "startswith",
	EndsWith = "endswith",
	CIEquals = "ci_equals",
	CIContains = "ci_contains",
	CIStartsWith = "ci_startswith",
	CIEndsWith = "ci_endswith",
	CINumericEquals = "==",
	CILessThan = "<",
	CIGreaterThan = ">",
	CILessOrEqual = "<=",
	CIGreaterOrEqual = ">="
}

export enum GroupOperator {
	invalid = "",
	OR = "or",
	AND = "and",
	NOR = "nor",
	NAND = "nand",
	XOR = "xor"
}

export class ConditionalGroup {
	public Operator: GroupOperator = GroupOperator.invalid;
	public Conditionals: ComponentConditional[] = [];
}

export class ComponentConditional {
	public Name: string = "";
	public Not: boolean = false;
	public Operator: ConditionalOperator = ConditionalOperator.invalid;
	public Value: string = "";
	public Group?: ConditionalGroup;
	private valueSet: boolean = false;
	private validated: boolean = false;
}

export class ToggleableComponent {
	// TODO: the rest of these defs
}