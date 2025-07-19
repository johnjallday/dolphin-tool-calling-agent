package user


type UserManager interface {
    Users() []string                  // List available user IDs (usernames)
    CurrentUser() *User               // Get the currently loaded user
    LoadUser(userID string) error     // Load a user by ID
    UnloadUser()                      // Unload current user
}

type SimpleUserManager struct{

}
