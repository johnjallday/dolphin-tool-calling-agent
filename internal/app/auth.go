package app

// Authenticator is responsible for "logging in" a user
// and returning a canonical userID (e.g. "alice").
type Authenticator interface {
    Login() (userID string, err error)
}
