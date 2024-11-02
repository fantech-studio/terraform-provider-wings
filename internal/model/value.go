package model

type Value struct {
	ID             string            `json:"id"`
	Enabled        bool              `json:"enabled"`
	Description    string            `json:"description"`
	DefaultVariant string            `json:"defaultVariant"`
	Variants       Variants          `json:"variants"`
	Targeting      Targeting         `json:"targeting"`
	Tests          []*EvaluationTest `json:"tests,omitempty"`
}

type (
	Variants map[string]ValueEvaluation

	ValueEvaluation struct {
		Bool   *Bool   `json:"bool"`
		String *String `json:"string"`
		Object *Object `json:"object"`
		Int    *Int    `json:"int"`
	}

	Bool struct {
		Value bool `json:"value"`
	}

	Int struct {
		Value int64 `json:"value"`
	}

	String struct {
		Value string `json:"value"`
	}

	Object struct {
		Value      map[string]any    `json:"value"`
		Transforms []*ValueTransform `json:"transforms,omitempty"`
	}
)

type EvaluationTest struct {
	Variables map[string]any `json:"variables"`
	Expected  string         `json:"expected"`
}

type Targeting struct {
	Rules []ValueTargetingRule `json:"rules"`
}

type ValueTargetingRule struct {
	Variant string `json:"variant"`
	Expr    string `json:"expr"`
}

type ValueTransform struct {
	Expr string `json:"expr"`
}
