package websecure

import (
	"log"
	"net/http"
	"strings"

	u "vsys.commons/utils"
)

func CommonMiddleware(next http.Handler) http.Handler {
	// Array of URLs allowed without a token
	allowedURLs := []string{"/login", "/register"}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request path is in the list of allowed URLs
		path := r.URL.Path
		for _, url := range allowedURLs {
			if path == url {
				next.ServeHTTP(w, r)
				return
			}
		}

		// For all other paths, check for the presence of a valid Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// Extract and validate the JWT token from the Authorization header
			if checkApiToken(r) {
				next.ServeHTTP(w, r)
				return
			}

			// If the token is not valid, respond with Unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		// If there's no Authorization header, respond with Unauthorized status
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
	})
}

// checkApiToken extracts and validates JWT token from the Authorization header.
func checkApiToken(r *http.Request) bool {

	// log.Println("checkApiToken: request received")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		log.Println("checkApiToken: request rejected, no Authorization in header")
		return false
	}

	// Typically, Authorization header is in the format "Bearer <token>",
	// so we need to split by space and get the second part
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		log.Println("checkApiToken: request rejected, no bearer in header")
		return false
	}

	token := parts[1]

	// log.Println("checkApiToken: token: ", token)

	return u.ValidateJwtToken(token)
}
