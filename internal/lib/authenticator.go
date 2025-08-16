package lib

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nrednav/cuid2"
)

// --- JWT Authenticator ---
type JWTAuthenticator struct {
	secret string
	aud    string
	iss    string
}

func NewJWTAuthenticator(secret, aud, iss string) *JWTAuthenticator {
	return &JWTAuthenticator{secret: secret, aud: aud, iss: iss}
}

type AccessTokenClaims struct {
	// No custom claims needed for access token if Casbin handles roles server-side
	jwt.RegisteredClaims
	IsSuperuser bool `json:"isSuperuser,omitempty"`
}

// VerificationTokenClaims for email verification or password reset tokens
type VerificationTokenClaims struct {
	Purpose string `json:"pur,omitempty"` // e.g., "email_verification", "password_reset"
	jwt.RegisteredClaims
}

const (
	PurposeEmailVerification = "email_verification"
	PurposePasswordReset     = "password_reset"
)

// GenerateAccessToken generates a standard access token for a user.
func (a *JWTAuthenticator) GenerateAccessToken(userID string, isSuperuser bool, lifetime time.Duration) (string, error) {
	expirationTime := time.Now().Add(lifetime)
	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    a.iss,
			Audience:  jwt.ClaimStrings{a.aud},
			ID:        cuid2.Generate(),
		},
		IsSuperuser: isSuperuser,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}

// ValidateAndParseAccessToken validates a token string and parses it into AccessTokenClaims.
func (a *JWTAuthenticator) ValidateAndParseAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	},
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, jwt.ErrTokenSignatureInvalid
	}
	return claims, nil
}

// generatePurposeToken generates a JWT for a specific purpose (email verification, password reset).
func (a *JWTAuthenticator) generatePurposeToken(userID string, purpose string, lifetime time.Duration) (string, error) {
	expirationTime := time.Now().Add(lifetime)
	purposeAudience := a.aud + "-" + purpose // e.g., "ethnocopia-app-email_verification"
	purposeIssuer := a.iss                   // Can use the same issuer or a purpose-specific one

	claims := VerificationTokenClaims{
		Purpose: purpose,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    purposeIssuer,
			Audience:  jwt.ClaimStrings{purposeAudience},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.secret))
}

// validatePurposeToken validates a purpose-specific JWT and returns the userID (subject).
func (a *JWTAuthenticator) validatePurposeToken(tokenStr string, expectedPurpose string) (userID string, err error) {
	claims := &VerificationTokenClaims{}
	purposeAudience := a.aud + "-" + expectedPurpose
	purposeIssuer := a.iss

	parsedToken, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.secret), nil
	},
		jwt.WithAudience(purposeAudience),
		jwt.WithIssuer(purposeIssuer),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return "", err
	}
	if !parsedToken.Valid {
		return "", jwt.ErrTokenSignatureInvalid
	}
	if claims.Purpose != expectedPurpose {
		return "", fmt.Errorf("token purpose mismatch: expected '%s', got '%s'", expectedPurpose, claims.Purpose)
	}

	return claims.Subject, nil
}

// GenerateVerificationToken delegates to generatePurposeToken
func (a *JWTAuthenticator) GenerateVerificationTokenJWT(userID string, lifetime time.Duration) (string, error) {
	return a.generatePurposeToken(userID, PurposeEmailVerification, lifetime)
}

// ValidateAndParseVerificationTokenJWT delegates to validatePurposeToken
func (a *JWTAuthenticator) ValidateAndParseVerificationTokenJWT(tokenStr string) (userID string, err error) {
	return a.validatePurposeToken(tokenStr, PurposeEmailVerification)
}

// GeneratePasswordResetToken delegates to generatePurposeToken
func (a *JWTAuthenticator) GeneratePasswordResetTokenJWT(userID string, lifetime time.Duration) (string, error) {
	return a.generatePurposeToken(userID, PurposePasswordReset, lifetime)
}

// ValidatePasswordResetToken delegates to validatePurposeToken
func (a *JWTAuthenticator) ValidatePasswordResetTokenJWT(tokenStr string) (userID string, err error) {
	return a.validatePurposeToken(tokenStr, PurposePasswordReset)
}
