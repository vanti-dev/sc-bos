package token

type Permission string

const (
	TraitWriteAll Permission = "trait:write"
	TraitReadAll  Permission = "trait:read"
)
