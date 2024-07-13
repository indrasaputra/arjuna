package entity

// Token represents token.
type Token struct {
	AccessToken           string
	TokenType             string
	RefreshToken          string
	AccessTokenExpiresIn  uint32
	RefreshTokenExpiresIn uint32
}
