package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type UserClaims struct {
	ID        string    `json:"id"`
	TokenType TokenType `json:"typ"`
	jwt.RegisteredClaims
}

func NewUserClaims(id string, tokenType TokenType, duration time.Duration) (*UserClaims, error) {
	if tokenType != TokenTypeAccess && tokenType != TokenTypeRefresh {
		return nil, fmt.Errorf("invalid token type: %s", tokenType)
	}
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}
	return &UserClaims{
		ID:        id,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
