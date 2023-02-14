package fetcher

type TaskModel struct {
	Property
	Root  string      `json:"root_script"`
	Rules []RuleModel `json:"rules"`
}

type RuleModel struct {
	Name      string `json:"name"`
	ParseFunc string `json:"parse_script"`
}
