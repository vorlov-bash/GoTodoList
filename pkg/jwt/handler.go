package jwt

import (
	"net/http"
	"strings"
)

const DefaultExp = 3600
const DefaultAlg = "HS256"
const DefaultTyp = "JWT"

type JWT struct {
	Alg    string
	Typ    string
	Secret string
}

type Controller interface {
	NewToken(target string, exp int64) string
	VerifyToken(token string) bool
}

func NewJWT(secret, alg string) *JWT {
	if alg == "" {
		alg = DefaultAlg
	}

	return &JWT{
		Alg:    alg,
		Typ:    DefaultTyp,
		Secret: secret,
	}
}

func (j *JWT) CreateToken(target string, exp int64) string {
	if exp == 0 {
		exp = DefaultExp
	}

	header := generateHeader(j.Alg, j.Typ)
	payload := generatePayload(target, exp)
	signature := generateSignature(header, payload, j.Secret)
	return header + "." + payload + "." + signature
}

func (j *JWT) VerifyToken(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	header := parts[0]
	payload := parts[1]
	signature := parts[2]

	expectedSignature := generateSignature(header, payload, j.Secret)
	return signature == expectedSignature
}

func (j *JWT) WrapHTTPHandler(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !j.VerifyToken(token) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (j *JWT) GetSubFromToken(token string) string {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return ""
	}

	payload := parts[1]
	pMap, err := getPayload(payload)
	if err != nil {
		return ""
	}

	return pMap["sub"].(string)
}
