package evaluations

type EvaluationHistory struct {
	Evaluation []Token
	Result     *string `json:"Result,omitempty"`
}

type ListResponse struct {
	EvaluationHistory []EvaluationHistory
}
