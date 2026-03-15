package types

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

type RefreshToken struct {
	Value string `json:"value"`
}
