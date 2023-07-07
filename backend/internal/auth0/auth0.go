package auth0

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
  "strings"

	ginadapter "github.com/gwatts/gin-adapter"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type CustomClaims struct {
	Scope string `json:"scope"`
}

func (c CustomClaims) Validate(context.Context) error {
	return nil
}

func extractDiscordIdFromSub(sub string) string {
  parts := strings.Split(sub, "|")
  return parts[len(parts) - 1]
}

func EnsureValidClaims() gin.HandlerFunc {
  return func (c *gin.Context) {
      claims, ok := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
      if !ok {
        c.AbortWithStatusJSON(
          http.StatusInternalServerError,
          map[string]string{"message": "Failed to get validated JWT claims."},
        )
      }

      claimID := extractDiscordIdFromSub(claims.RegisteredClaims.Subject)
      targetID := c.Param("userId")

      if claimID != targetID {
        c.AbortWithStatusJSON(
          http.StatusUnauthorized,
          map[string]string{"message": "Not permitted to access resource"},
        )
      }
			return
  }
}

// EnsureValidToken is a middleware that will check the validity of our JWT.
func EnsureValidToken() gin.HandlerFunc {
	issuerURL, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	if err != nil {
		log.Fatalf("Failed to parse the issuer url: %v", err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 5*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		[]string{os.Getenv("AUTH0_AUDIENCE")},
		validator.WithCustomClaims(
			func() validator.CustomClaims {
				return &CustomClaims{}
			},
		),
		validator.WithAllowedClockSkew(time.Minute),
	)
	if err != nil {
		log.Fatalf("Failed to set up the jwt validator")
	}


	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Encountered error while validating JWT: %v", err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"headers": {"Content-Type": "application/json"}, "statusCode": "401", "isBase64Encoded": "false", "body":"Failed to validate JWT."}`))
	}

	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
	)

  handler := func(next http.Handler) http.Handler {
		return middleware.CheckJWT(next)
	}

	return ginadapter.Wrap(handler)
}
