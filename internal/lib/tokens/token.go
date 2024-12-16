package tokens

import (
	"auth_rest/internal/models"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

type TokenManager interface {
	NewAccess(user models.User, duration time.Duration) (string, error)
	NewRefresh(user models.User, duration time.Duration) (string, error)
	ParseAccess(accessToken string) (TokenClaims, error)
	ParseRefresh(refreshToken string) (TokenClaims, error)
}

type TokenClaims struct {
	jwt.StandardClaims
	GUID string `json:"guid"`
	IP   string `json:"ip"`
}

type Manager struct {
	signingKey string
}

func NewManager(signingKey string) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}

	return &Manager{signingKey: signingKey}, nil
}

func (m *Manager) NewAccess(user models.User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		GUID: user.GUID,
		IP:   user.IP,
	})

	tokenString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) NewRefresh(user models.User, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		GUID: user.GUID,
		IP:   user.IP,
	})

	tokenString, err := token.SignedString([]byte(m.signingKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *Manager) ParseAccess(accessToken string) (TokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return TokenClaims{}, err
	}

	if !token.Valid {
		return TokenClaims{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return TokenClaims{}, fmt.Errorf("token claims are not of type *tokenClaims")
	}

	return *claims, nil
}

func (m *Manager) ParseRefresh(refreshToken string) (TokenClaims, error) {
	token, err := jwt.ParseWithClaims(refreshToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return TokenClaims{}, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return TokenClaims{}, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return TokenClaims{}, fmt.Errorf("token claims are not of type *tokenClaims")
	}

	return *claims, nil
}
