package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	SECRET_KEY = "vsys_jwt_token"
)

// concatenate all provided strings and return single resultant string
func JoinStr(sentences ...string) (result string) {
	var final strings.Builder
	for _, sentence := range sentences {
		final.WriteString(sentence)
	}
	result = final.String()
	return
}

// compares the data in both interface and returns true if both interface content is identical else false
func CompareData(first interface{}, second interface{}) (res bool) {
	var fir []byte
	var sec []byte

	fir, firOk := first.([]byte)
	if !firOk {
		byteArray, firErr := json.Marshal(first)
		if firErr != nil {
			res = false
		}
		fir = byteArray
	}

	sec, secOk := second.([]byte)
	if !secOk {
		byteArray, secErr := json.Marshal(second)
		if secErr != nil {
			res = false
		}
		sec = byteArray
	}
	res = bytes.Contains(fir, sec)

	return
}

// returns date corresponding to yyyy-mm-dd hh-mm-ss + nsec local(+530 IST)
func DateAndTime(year int, month time.Month, day int, hour int, min int, sec int, nsec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, time.Local)
}

// Create Jwt token using mobile number and hard coded secret key
// Token expiry timing is kept as 1 hr
func CreateJwtToken(mobileNo uint64) (string, time.Time, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = mobileNo
	expire := time.Now().UTC().Add(time.Hour * 1)
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().UTC().Unix()
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", time.Now(), err
	}
	return tokenString, expire, nil
}

// ValidateJwtToken checks the validity of the JWT token
func ValidateJwtToken(tokenString string) bool {
	secretKey := []byte(SECRET_KEY)

	// log.Println("ValidateJwtToken: tokenString: ", tokenString)

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("unexpected signing method: %v", token.Header["alg"])
			return nil, nil // Should return an error here instead of nil, nil if the signing method is not as expected
		}
		return secretKey, nil
	})

	if err != nil {
		log.Println("Token parsing error, token- "+tokenString+" error- ", err)
		return false
	}

	// Validate token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if the token is expired
		if exp, ok := claims["exp"].(float64); ok {
			// Convert to int64, then compare with current time
			if int64(exp) > time.Now().UTC().Unix() {
				return true // Token is valid and not expired
			}
		}
	}

	return false // Token is either invalid or expired
}

func GenerateOtp(max int) (int, string, error) {
	var digitsTable = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	var b = make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return 0, "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = digitsTable[int(b[i])%len(digitsTable)]
	}
	byteToInt, err := strconv.Atoi(string(b))
	if err != nil {
		return 0, "", err
	}
	// hard coding, need to replace this with near to optimal logic
	if max == 6 && byteToInt < 100000 {
		return GenerateOtp(max)
	}
	return byteToInt, string(b), nil
}
