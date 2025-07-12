package app

import (
  "github.com/johnjallday/dolphin-tool-calling-agent/internal/user"
)

type UserLoader interface {
    Load(userID string) (*user.User, error)
}
