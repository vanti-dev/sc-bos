package token

type Claims struct {
	Roles     []string `json:"roles"`      // The names of the roles that the subject has been granted
	Scopes    []string `json:"scopes"`     // The scopes that this authorization is limited to
	Zones     []string `json:"zones"`      // The zones that this token is authorized for, for tenant tokens
	IsService bool     `json:"is_service"` // True if the subject is an application acting on its own behalf, false if it's a user
}
