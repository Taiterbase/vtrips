package auth

import (
    "net/http"
    "os"
    "time"

    "github.com/Taiterbase/vtrips/apps/users/internal/storage"
    "github.com/golang-jwt/jwt"
)

const (
	// AuthCookieName is the name of the cookie that holds the JWT
	AuthCookieName = "auth_token"
)

func GenerateJWT(userID, username, tenant string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"name":     username,
        "iat":      time.Now().Unix(),
	}
	if tenant != "" {
		claims["tenant"] = tenant
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func InvalidateCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     AuthCookieName,
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	if os.Getenv("ENVIRONMENT") == "development" {
		cookie.Secure = false
	} else {
		cookie.Secure = true
	}
	http.SetCookie(w, cookie)
}

func SetCookieWithJWT(w http.ResponseWriter, userID, username, tenant string) error {
	tokenString, err := GenerateJWT(userID, username, tenant)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     AuthCookieName,
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 14),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	}
	if os.Getenv("ENVIRONMENT") == "development" {
		cookie.Secure = false
	} else {
		cookie.Secure = true
	}

	http.SetCookie(w, cookie)
	return err
}

func Validate(tokenString string) (jwt.Claims, error) {
	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok { // validate signing method
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid { // validate token
		return nil, jwt.ErrSignatureInvalid
	}
	if token.Claims.Valid() != nil { // validate claims
		return nil, jwt.ErrInvalidKey
	}
    claims := token.Claims.(jwt.MapClaims)
    // Check revocation by user ID and issued-at
    if sub, ok := claims["sub"].(string); ok {
        if iatf, ok := claims["iat"].(float64); ok {
            revokedBefore, _ := storage.GetRevokedBefore(sub)
            if revokedBefore > 0 && int64(iatf) <= revokedBefore {
                return nil, jwt.ErrSignatureInvalid
            }
        }
    }
	return claims, nil
}
