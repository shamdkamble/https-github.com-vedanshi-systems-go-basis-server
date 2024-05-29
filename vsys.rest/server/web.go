// server\web.go

package server

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"

	// "time"
	m "vsys.commons/model"
	websecure "vsys.commons/websecure"
	restUtils "vsys.rest/services"

	"github.com/gorilla/mux"
)

const DefaultDBHelperHost string = "0.0.0.0" // Default port if not set in env
const DefaultDBHelperPort string = "7200"    // Default port if not set in env

var dbHelperHost string
var dbHelperPort string // Global variable to hold the DB helper port

func init() {
	// Initialize dbHelperHost and dbHelperPort with value from environment variable or fallback to default
	dbHelperHost = os.Getenv("DBHELPER_HOST")
	if strings.TrimSpace(dbHelperHost) == "" {
		dbHelperHost = DefaultDBHelperHost
	}

	dbHelperPort = os.Getenv("DBHELPER_PORT")
	if strings.TrimSpace(dbHelperPort) == "" {
		dbHelperPort = DefaultDBHelperPort
	}
}

func Web() {
	r := mux.NewRouter()

	r.Use(websecure.CommonMiddleware)

	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/get-user-by-mobile", GetUserByMobileNoHandler).Methods("GET")
	r.HandleFunc("/get-user-by-email", GetUserByEmailHandler).Methods("GET")

	port := os.Getenv("RESTSRV_PORT")
	if port == "" {
		slog.Info("RESTSRV_PORT environment variable not set, defaulting to :7100")
		port = "7100"
	}

	slog.Info("Server listening on port: " + port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// TODO- call db helper APIs using tokens

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var endUser m.EndUser
	err := json.NewDecoder(r.Body).Decode(&endUser)
	if err != nil {
		log.Println("RegisterHandler: error- ", err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(endUser)
	if err != nil {
		log.Printf("Error marshaling user data: %v", err)
		http.Error(w, "Failed to marshal user data", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://"+dbHelperHost+":"+dbHelperPort+"/create-user", "application/json", bytes.NewBuffer(jsonValue)) // Changed port to 3100
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to create user", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	restUtils.RespondFailure(resp, w)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Println("LoginHandler: request received")

	var credentials m.LoginReq
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jsonValue, err := json.Marshal(credentials)
	if err != nil {
		log.Printf("Error marshaling login credentials: %v", err)
		http.Error(w, "Failed to marshal login credentials", http.StatusBadRequest)
		return
	}

	// log.Println("LoginHandler: API http://"+dbHelperHost+":"+dbHelperPort+"/login called")

	resp, err := http.Post("http://"+dbHelperHost+":"+dbHelperPort+"/login", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to login", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	// log.Println("LoginHandler: http://"+dbHelperHost+":"+dbHelperPort+"/login response received")

	restUtils.RespondFailure(resp, w)

	// Proceed with handling a successful response (200 OK or 201 Created)
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	// Write the successful response back to the client
	w.Write(responseBody)
}

func GetUserByMobileNoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Println("GetUserByMobileNoHandler: request received")

	mobileNo := r.URL.Query().Get("mobile_no")
	if len(strings.TrimSpace(mobileNo)) <= 0 {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	// log.Println("GetUserByMobileNoHandler: mobileNo", mobileNo)

	// log.Println("GetUserByMobileNoHandler: http://"+dbHelperHost+":"+dbHelperPort+"/get-user-by-mobile request sent", mobileNo)

	// Use net/url to build the query parameters
	queryParams := url.Values{}
	queryParams.Add("mobile_no", mobileNo)

	// Make HTTP POST request
	resp, err := http.Get("http://" + dbHelperHost + ":" + dbHelperPort + "/get-user-by-mobile?" + queryParams.Encode())
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	restUtils.RespondFailure(resp, w)

	// Process response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}

func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Println("GetUserByEmailHandler: request received")

	email := r.URL.Query().Get("email")
	if len(strings.TrimSpace(email)) <= 0 {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	// log.Println("GetUserByMobileNoHandler: email", email)

	// log.Println("GetUserByMobileNoHandler: http://"+dbHelperHost+":"+dbHelperPort+"/get-user-by-email request sent", email)

	// Use net/url to build the query parameters
	queryParams := url.Values{}
	queryParams.Add("email", email)

	// Make HTTP Get request
	resp, err := http.Get("http://" + dbHelperHost + ":" + dbHelperPort + "/get-user-by-email?" + queryParams.Encode())
	if err != nil {
		log.Printf("Error making POST request: %v", err)
		http.Error(w, "Failed to get user information", http.StatusForbidden)
		return
	}
	defer resp.Body.Close()

	restUtils.RespondFailure(resp, w)

	// log.Println("GetUserByMobileNoHandler: http://"+dbHelperHost+":"+dbHelperPort+"/get-user-by-mobile response", resp.Body)

	// Process response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	w.Write(responseBody)
}
