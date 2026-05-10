package model

type Observation struct {
	Path         string  `json:"path"`
	P95MS        int     `json:"p95_ms"`
	P99MS        int     `json:"p99_ms"`
	ErrorRate    float64 `json:"error_rate"`
	DependencyMS int     `json:"dependency_ms"`
	Region       string  `json:"region"`
}

type BudgetRequest struct {
	SystemName    string        `json:"system_name"`
	ServiceOwner  string        `json:"service_owner"`
	P95BudgetMS   int           `json:"p95_budget_ms"`
	P99BudgetMS   int           `json:"p99_budget_ms"`
	MaxErrorRate  float64       `json:"max_error_rate"`
	Observations  []Observation `json:"observations"`
}

type Finding struct {
	Code                  string   `json:"code"`
	Severity              string   `json:"severity"`
	Score                 int      `json:"score"`
	Summary               string   `json:"summary"`
	Evidence              []string `json:"evidence"`
	RecommendedNextAction string   `json:"recommended_next_action"`
}

type Report struct {
	SystemName           string    `json:"system_name"`
	ServiceOwner         string    `json:"service_owner"`
	ObservationsAnalyzed int       `json:"observations_analyzed"`
	OverallScore         int       `json:"overall_score"`
	Findings             []Finding `json:"findings"`
}

type HealthResponse struct {
	Status        string `json:"status"`
	Service       string `json:"service"`
	Docs          string `json:"docs"`
	SampleBudget  string `json:"sample_budget"`
}
