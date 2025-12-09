package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwks"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	"github.com/gin-gonic/gin"
)

// ContextKey is the type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "user_id"
)

// ClerkConfig holds Clerk configuration
type ClerkConfig struct {
	SecretKey string
}

// ClerkMiddleware provides JWT authentication middleware for Clerk
type ClerkMiddleware struct {
	config    ClerkConfig
	jwksClient *jwks.Client
	jwkStore   *JWKStore
	jwkMu      sync.RWMutex
}

// JWKStore stores the cached JSON Web Key
type JWKStore struct {
	jwk *clerk.JSONWebKey
}

// Claims represents JWT claims from Clerk
type Claims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
}

// NewClerkMiddleware creates a new Clerk middleware
func NewClerkMiddleware(config ClerkConfig) *ClerkMiddleware {
	// Set the Clerk secret key for SDK initialization
	clerk.SetKey(config.SecretKey)

	// Create JWKS client
	clientConfig := &clerk.ClientConfig{}
	clientConfig.Key = clerk.String(config.SecretKey)
	jwksClient := jwks.NewClient(clientConfig)

	return &ClerkMiddleware{
		config:     config,
		jwksClient: jwksClient,
		jwkStore:   &JWKStore{},
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
		claims, err := m.validateToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": fmt.Sprintf("invalid token: %v", err),
			})
			return
		}

		// Set user ID in context
		ctx := context.WithValue(c.Request.Context(), UserIDKey, claims.Sub)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// validateToken validates a JWT token using Clerk SDK
func (m *ClerkMiddleware) validateToken(ctx context.Context, tokenString string) (*Claims, error) {
	// Try to get cached JWK
	jwk := m.getJWK()
	
	if jwk == nil {
		// Decode token to get key ID (without verification)
		unsafeClaims, err := jwt.Decode(ctx, &jwt.DecodeParams{
			Token: tokenString,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to decode token: %w", err)
		}

		// Fetch the JSON Web Key using the key ID
		jwk, err = jwt.GetJSONWebKey(ctx, &jwt.GetJSONWebKeyParams{
			KeyID:      unsafeClaims.KeyID,
			JWKSClient: m.jwksClient,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to fetch JWK: %w", err)
		}

		// Cache the JWK
		m.setJWK(jwk)
	}

	// Verify the token signature and claims
	verifiedClaims, err := jwt.Verify(ctx, &jwt.VerifyParams{
		Token: tokenString,
		JWK:   jwk,
	})
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}

	// Convert Clerk claims to our Claims struct
	claims := &Claims{
		Sub: verifiedClaims.Subject,
	}

	// Extract expiration time (Expiry field from RegisteredClaims)
	if verifiedClaims.Expiry != nil {
		claims.Exp = *verifiedClaims.Expiry
	}

	// Extract issued at time
	if verifiedClaims.IssuedAt != nil {
		claims.Iat = *verifiedClaims.IssuedAt
	}

	return claims, nil
}

// getJWK returns the cached JWK
func (m *ClerkMiddleware) getJWK() *clerk.JSONWebKey {
	m.jwkMu.RLock()
	defer m.jwkMu.RUnlock()
	return m.jwkStore.jwk
}

// setJWK caches the JWK
func (m *ClerkMiddleware) setJWK(jwk *clerk.JSONWebKey) {
	m.jwkMu.Lock()
	defer m.jwkMu.Unlock()
	m.jwkStore.jwk = jwk
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
	if v := ctx.Value(UserIDKey); v != nil {
		return v.(string)
	}
	return ""
}

// GetLaboratoryID returns empty string (laboratory_id is no longer extracted from JWT)
// Handlers should extract laboratory_id from query parameters instead
func GetLaboratoryID(ctx context.Context) string {
	return ""
}

