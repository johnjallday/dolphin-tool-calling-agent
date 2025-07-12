package app

import (
    "fmt"
    "os"
)

type OSUserAuth struct{}

func (a *OSUserAuth) Login() (string, error) {
    u := os.Getenv("USER")
    if u == "" {
        return "", fmt.Errorf("$USER not set, cannot log in")
    }
    return u, nil
}
