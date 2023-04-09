package request

type GrantScoreMapping struct {
	DpMapping  BalanceMapping `json:"dp_mapping"`
	ExpMapping BalanceMapping `json:"exp_mapping"`
}

type BalanceMapping struct {
	BalanceAmount uint                     `json:"balance_amount"`
	ActionRewards map[string]ActionMapping `json:"action_rewards"`
}

type ActionMapping struct {
	Amount uint `json:"amount"`
	Limit  uint `json:"limit"`
}

type UpdateUserScore struct {
	Scores []UserScore `json:"scores"`
}

type UserScore struct {
	UserId     uint                `json:"user_id"`
	DpMapping  []UserActionMapping `json:"dp_mapping"`
	ExpMapping []UserActionMapping `json:"exp_mapping"`
}

type UserActionMapping struct {
	RuleId         uint `json:"rule_id"`
	RewardedAmount uint `json:"amount"`
}
