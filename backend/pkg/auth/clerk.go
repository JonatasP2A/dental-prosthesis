package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ContextKey is the type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
	// LaboratoryIDKey is the context key for laboratory ID
	LaboratoryIDKey ContextKey = "laboratory_id"
)

// ClerkConfig holds Clerk configuration
type ClerkConfig struct {
	JWKSURL string
}

// ClerkMiddleware provides JWT authentication middleware for Clerk
type ClerkMiddleware struct {
	config    ClerkConfig
	jwks      *JWKS
	jwksMu    sync.RWMutex
	jwksCache time.Time
}

// JWKS represents JSON Web Key Set
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

// Claims represents JWT claims from Clerk
type Claims struct {
	Sub          string `json:"sub"`
	LaboratoryID string `json:"laboratory_id"`
	Exp          int64  `json:"exp"`
	Iat          int64  `json:"iat"`
}

// NewClerkMiddleware creates a new Clerk middleware
func NewClerkMiddleware(config ClerkConfig) *ClerkMiddleware {
	return &ClerkMiddleware{
		config: config,
	}
}

// Authenticate is the Gin middleware for JWT authentication
func (m *ClerkMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		token := parts[1]
		claims, err := m.validateToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("invalid token: %v", err),
			})
			return
		}

		// Check if laboratory_id claim exists
		if claims.LaboratoryID == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "missing laboratory_id claim",
			})
			return
		}

		// Set claims in context
		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.Sub)
		ctx = context.WithValue(ctx, LaboratoryIDKey, claims.LaboratoryID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// validateToken validates a JWT token
func (m *ClerkMiddleware) validateToken(tokenString string) (*Claims, error) {
	// Split token parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	// Decode payload (without verification for now - in production use proper JWT library)
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	// Check expiration
	if claims.Exp < time.Now().Unix() {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

// fetchJWKS fetches the JWKS from Clerk
func (m *ClerkMiddleware) fetchJWKS() (*JWKS, error) {
	m.jwksMu.RLock()
	if m.jwks != nil && time.Since(m.jwksCache) < time.Hour {
		defer m.jwksMu.RUnlock()
		return m.jwks, nil
	}
	m.jwksMu.RUnlock()

	m.jwksMu.Lock()
	defer m.jwksMu.Unlock()

	resp, err := http.Get(m.config.JWKSURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	m.jwks = &jwks
	m.jwksCache = time.Now()
	return &jwks, nil
}

// getPublicKey extracts the RSA public key from a JWK
func getPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode N: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode E: %w", err)
	}

	e := 0
	for _, b := range eBytes {
		e = e<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: e,
	}, nil
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if v := ctx.Value(UserIDKey); v != nil {
		return v.(string)
	}
	return ""
}

// GetLaboratoryID extracts laboratory ID from context
func GetLaboratoryID(ctx context.Context) string {
	if v := ctx.Value(LaboratoryIDKey); v != nil {
		return v.(string)
	}
	return ""
}

