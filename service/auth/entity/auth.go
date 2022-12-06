package entity

// Token represents token.
type Token struct {
	AccessToken           string
	AccessTokenExpiresIn  uint32
	RefreshToken          string
	RefreshTokenExpiresIn uint32
	TokenType             string
}
