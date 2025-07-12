package app

import (
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

type TOMLUserLoader struct{}

func (l *TOMLUserLoader) Load(userID string) (*user.User, error) {
    return user.LoadUser(userID)
}
