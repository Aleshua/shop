package services

import (
	"context"
	"crypto/rsa"
	"time"

	"auth/internal/config"
	d "auth/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	ttlSeconds int64
}

func NewJWT(params config.JWT, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *JWTService {
	return &JWTService{
		privateKey: privateKey,
		publicKey:  publicKey,
		ttlSeconds: params.AccessTokenTTLSeconds,
	}
}

func (j *JWTService) GenerateToken(ctx context.Context, userId int32) (string, error) {
	now := time.Now().Unix()

	issuedAt := now
	expiresAt := now + j.ttlSeconds

	token := jwt.NewWithClaims(
		jwt.SigningMethodRS256,
		jwt.MapClaims{"user_Id": userId, "issued_at": issuedAt, "expires_at": expiresAt},
	)
	return token.SignedString(j.privateKey)
}

func (j *JWTService) ValidateToken(ctx context.Context, token string) (int32, error) {
	parsedToken, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, d.ErrJWTInvalidMethodSignature
			}
			return j.publicKey, nil
		},
	)
	if err != nil {
		return 0, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		if time.Now().Unix() > int64(claims["expires_at"].(float64)) {
			return 0, d.ErrJWTTokenExpired
		}

		return int32(claims["user_Id"].(float64)), nil
	}

	return 0, d.ErrJWTInvalidToken
}

func (j *JWTService) ExtractClaims(ctx context.Context, token string) (int32, error) {
	parsedToken, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			return j.publicKey, nil
		},
		jwt.WithoutClaimsValidation(),
	)
	if err != nil {
		return 0, err
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
		return int32(claims["user_Id"].(float64)), nil
	}

	return 0, d.ErrJWTInvalidToken
}

func (j *JWTService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	userId, err := j.ValidateToken(ctx, oldToken)
	if err != nil {
		return "", err
	}

	newToken, err := j.GenerateToken(ctx, userId)
	if err != nil {
		return "", err
	}

	return newToken, nil
}
