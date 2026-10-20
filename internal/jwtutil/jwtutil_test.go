package jwtutil_test

import (
	"authservice/internal/jwtutil"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSecret  = "test-secret-bot-for-production"
	otherSecret = "some-other-test-secret"
)

func TestGenerateAndValidate_RoundTrip(t *testing.T) {
	userID := uuid.New()
	token, err := jwtutil.GenerateToken(userID, testSecret)
	require.NoError(t, err)
	gotID, err := jwtutil.ValidateToken(token, testSecret)
	require.NoError(t, err)
	assert.Equal(t, userID, gotID)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	userID := uuid.New()
	token, err := jwtutil.GenerateToken(userID, otherSecret)
	require.NoError(t, err)
	gotID, err := jwtutil.ValidateToken(token, testSecret)
	assert.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
	assert.Equal(t, uuid.Nil, gotID)
}

func TestValidateToken_Expired(t *testing.T) {
	userID := uuid.New()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Date(2005, time.October, 7, 11, 37, 15, 1, time.Local).Unix(),
	}
	signedToken := signToken(t, claims, testSecret)
	gotID, err := jwtutil.ValidateToken(signedToken, testSecret)
	require.ErrorIs(t, err, jwt.ErrTokenExpired)
	assert.Equal(t, uuid.Nil, gotID)
}

func signToken(t *testing.T, claims jwt.MapClaims, secret string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	require.NoError(t, err)
	return tokenStr
}

func TestValidateToken_Malformed(t *testing.T) {
	tests := []struct {
		name  string
		token string
	}{
		{"empty token", ""},
		{"one word", "ehfivnrogfregreey5rsw"},
		{"two segments", "rhioegipakjepor.fhriehgiiwehfgpiwehg"},
		{"three segments", "hwefuhwofhjewpifh.fahfewhfiowhefo.werhfjherwitohrewgh"},
		{"four segments", "jgfkrfjgjriegpwerwpgreig.grfeujgjrgirejg.gjreijgperjgper.jgierjgrejgpeorgjqogf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := jwtutil.ValidateToken(tt.token, testSecret)
			require.ErrorIs(t, err, jwt.ErrTokenMalformed)
			assert.Equal(t, uuid.Nil, userID)
		})
	}
}

func TestValidateToken_InvalidClaims(t *testing.T) {
	tests := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{"no user_id", jwt.MapClaims{
			"exp": time.Now().Add(24 * time.Hour).Unix(),
		}},
		{"user_id is int", jwt.MapClaims{
			"user_id": 12345,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}},
		{"user_id is not uuid", jwt.MapClaims{
			"user_id": "not-a-uuid",
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}},
		{"user_id is empty", jwt.MapClaims{
			"user_id": "",
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			signedToken := signToken(t, tt.claims, testSecret)
			userID, err := jwtutil.ValidateToken(signedToken, testSecret)
			require.ErrorIs(t, err, jwtutil.ErrInvalidClaims)
			assert.Equal(t, uuid.Nil, userID)
		})
	}
}

func TestValidateToken_UnexpectedAlgorithm(t *testing.T) {
	tests := []struct {
		name   string
		method jwt.SigningMethod
		key    any
	}{
		{"none", jwt.SigningMethodNone, jwt.UnsafeAllowNoneSignatureType},
		{"HS384", jwt.SigningMethodHS384, []byte(testSecret)},
		{"HS512", jwt.SigningMethodHS512, []byte(testSecret)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := jwt.MapClaims{
				"user_id": uuid.New(),
				"exp":     time.Now().Add(time.Hour).Unix(),
			}
			token, err := jwt.NewWithClaims(tt.method, claims).SignedString(tt.key)
			require.NoError(t, err)
			userID, err := jwtutil.ValidateToken(token, testSecret)
			require.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
			assert.Equal(t, uuid.Nil, userID)
		})
	}
}

func TestValidateToken_MissingExp(t *testing.T) {
	claims := jwt.MapClaims{
		"user_id": uuid.New(),
	}
	token := signToken(t, claims, testSecret)
	userID, err := jwtutil.ValidateToken(token, testSecret)
	require.ErrorIs(t, err, jwt.ErrTokenRequiredClaimMissing)
	assert.Equal(t, uuid.Nil, userID)
}

func TestGenerateToken_Claims(t *testing.T) {
	userID := uuid.New()
	tokenStr, err := jwtutil.GenerateToken(userID, testSecret)
	require.NoError(t, err)
	token, err := jwt.Parse(tokenStr, func(tt *jwt.Token) (interface{}, error) {
		if _, ok := tt.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tt.Header["alg"])
		}
		return []byte(testSecret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	require.NoError(t, err)
	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	require.Equal(t, userID.String(), claims["user_id"])
	exp, ok := claims["exp"].(float64)
	require.True(t, ok, "exp must be a number")
	assert.InDelta(t, time.Now().Add(24*time.Hour).Unix(), exp, 5)
}
