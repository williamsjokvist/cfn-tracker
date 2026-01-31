package model

type SF6BattleStatsMetric struct {
	Key     string  `json:"key"`
	Name    string  `json:"name"`
	Unit    string  `json:"unit"`
	Kind    string  `json:"kind"`
	Current float64 `json:"current"`
	TopAvg  float64 `json:"topAvg"`
	TopMin  float64 `json:"topMin"`
	TopMax  float64 `json:"topMax"`
}

type SF6BattleStatsComparisonReport struct {
	UserCode          string                 `json:"userCode"`
	CharacterName     string                 `json:"characterName"`
	CharacterToolName string                 `json:"characterToolName"`
	TopN              int                    `json:"topN"`
	Metrics           []SF6BattleStatsMetric `json:"metrics"`
}
