package ctxKey

type contextKey int

const (
	UID contextKey = iota
	OrgID
	DbTransaction
)
