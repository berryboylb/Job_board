package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-gonic/gin"
)

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

// holds fields related to the JSON Web Key Set for the API. These keys contain the public keys, which will be used to verify JWTs
type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// declaring variable outside function
var GetToken map[string]interface{}

func Middleware() (*jwtmiddleware.JWTMiddleware, map[string]interface{}) {
	// jwtMiddleware is a handler that will verify access tokens
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// get token for database
			// assigning value to earlier declared variable through pointer
			p := &GetToken
			*p = token.Claims.(jwt.MapClaims)

			// Verify 'aud' claim
			// 'aud' = audience (where you deploy the backend, either locally or Heroku)
			aud := "https://" + os.Getenv("AUTH0_DOMAIN") + "/api/v2/"

			// convert audience in the JWT token to []interface{} if multiple audiences
			convAud, ok := token.Claims.(jwt.MapClaims)["aud"].([]interface{})
			if !ok {
				// convert audience in the JWT token to string if only 1 audience
				strAud, ok := token.Claims.(jwt.MapClaims)["aud"].(string)
				// return error if can't convert to string
				if !ok {
					return token, errors.New("invalid audience")
				}
				// return error if audience doesn't match
				if strAud != aud {
					return token, errors.New("invalid audience")
				}
			} else {
				for _, v := range convAud {
					// verify if audience in JWT is the one you've set
					if v == aud {
						break
					} else {
						return token, errors.New("invalid audience")
					}
				}
			}

			// Verify 'iss' claim
			// 'iss' = issuer
			iss := "https://" + os.Getenv("AUTH0_DOMAIN") + "/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return jwtMiddleware, GetToken
}

// function to grab JSON Web Key Set and return the certificate with the public key
func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://" + os.Getenv("AUTH0_DOMAIN") + "/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

func GinJWTMiddleware() gin.HandlerFunc {
	jwtMiddleware, _ := Middleware()

	return func(c *gin.Context) {
		// Convert the Gin context to an HTTP handler context
		h := jwtMiddleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Continue to the next handler if the JWT validation passes
			c.Next()
		}))

		// Create a ResponseWriter that writes to the Gin context
		rw := &responseWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}

		// Use the ResponseWriter to serve the HTTP handler
		h.ServeHTTP(rw, c.Request)

		// If the response status is not  200, then the JWT validation failed
		if rw.status != http.StatusOK {
			// Write the response body and status to the Gin context
			c.Writer.WriteHeader(rw.status)
			c.Writer.Write(rw.body.Bytes())
			// Abort the request
			c.Abort()
			return
		}
	}
}

// responseWriter is a ResponseWriter that writes to a bytes.Buffer
type responseWriter struct {
	body *bytes.Buffer
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.body.Write(b)
}
