package middleware

type contextKey string

const (
	UserIDKey    contextKey = "userId"
	SessionIDKey contextKey = "sessionId"
	UserRoleKey  contextKey = "userRole"
)
