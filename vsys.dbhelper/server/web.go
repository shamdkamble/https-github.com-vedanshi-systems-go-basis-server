package server

import (
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	m "vsys.commons/model"
	dao "vsys.dbhelper/dao"

	"github.com/gorilla/mux"
)

func Web() {
	r := mux.NewRouter()

	// TODO- apply middleware

	r.HandleFunc("/create-user", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/get-user-by-mobile", GetUserByMobileNoHandler).Methods("GET")
	r.HandleFunc("/get-user-by-email", GetUserByEmailHandler).Methods("GET")
	r.HandleFunc("/get-otp-count", GetOtpCountHandler).Methods("GET")
	r.HandleFunc("/save-otp-details", SaveOtpDetailsHandler).Methods("POST")
	r.HandleFunc("/find-otp-details", FindOtpDetailsHandler).Methods("GET") // Assuming mobileNo is passed as query param

	port := os.Getenv("DBHELPER_PORT")
	if port == "" {
		slog.Info("DBHELPER_PORT environment variable not set, defaulting to :7200")
		port = "7200"
	}

	slog.Info("Server listening on port: " + port)

	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var endUser m.EndUser
	err := json.NewDecoder(r.Body).Decode(&endUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userExistsErr := dao.CheckIfEndUserExists(&endUser, endUser.MobileNo)

	if userExistsErr != nil && strings.Contains(userExistsErr.Error(), "user already exists") {
		log.Printf("Error creating end user, already exists: %v", err)
		http.Error(w, "Failed to create end user (101)", http.StatusForbidden)
		return
	}

	if err := dao.CreateNewOrUpdateExistingEndUser(&endUser); err != nil {
		log.Printf("Error creating end user: %v", err)
		http.Error(w, "Failed to create end user (102)", http.StatusInternalServerError)
		return
	}

	token, expire, err := dao.GenerateToken(endUser) // with expiration time
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := m.ApiResp{
		Code:     http.StatusOK,
		Token:    token,
		Expire:   expire.Format(time.RFC3339),
		EndUsers: endUser,
		// Assuming you handle EndUserRoles and EndUsers appropriately here
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Printf("LoginHandler request received: in dbhandler")

	var credentials m.LoginReq
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// log.Println("LoginHandler request params: username- ", credentials.Username, "  password- ", credentials.Password)
	// Validate login credentials
	endUser, endUserRoles, validateErr := dao.CheckUserCredentials(credentials.Username, credentials.Password)
	if validateErr != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, expire, err := dao.GenerateToken(endUser) // with expiration time
	if err != nil {
		log.Printf("Error generating token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	resp := m.ApiResp{
		Code:         http.StatusOK,
		Token:        token,
		Expire:       expire.Format(time.RFC3339),
		EndUsers:     endUser,
		EndUserRoles: endUserRoles,
		// Assuming you handle EndUserRoles and EndUsers appropriately here
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetUserByMobileNoHandler handles the API request to get user info by mobile number
func GetUserByMobileNoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// log.Printf("GetUserByMobileNoHandler request received: in dbhandler")

	// Extract the mobile number from query params
	mobileNo := r.URL.Query().Get("mobile_no")
	if mobileNo == "" {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	// log.Println("GetUserByMobileNoHandler request received: in dbhandler, mobileNo- ", mobileNo)

	// Fetch user information based on mobile number
	endUser, err := dao.GetUserWithMobileNo(mobileNo)
	if err != nil {
		log.Printf("Error fetching user with mobile number %s: %v", mobileNo, err)
		http.Error(w, "Failed to get user information", http.StatusNotFound)
		return
	}

	// Assuming (m.EndUser{}) is the zero value when no user is found; this may need adjustment
	if (m.EndUser{}) == endUser {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// log.Println("GetUserByMobileNoHandler request received: in dbhandler, endUser- ", endUser.ID)

	// Construct the response with user details, token, and its expiry
	resp := m.ApiResp{
		Code: http.StatusOK,
		// Token:    token,
		// Expire:   expiry.Format(time.RFC3339),
		EndUsers: endUser,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetUserByEmailHandler handles the API request to get user info by email
func GetUserByEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the email from query params
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Fetch user information based on email
	endUser, err := dao.GetUserWithEmail(email)
	if err != nil {
		log.Printf("Error fetching user with email %s: %v", email, err)
		http.Error(w, "Failed to get user information", http.StatusNotFound)
		return
	}

	// Assuming (m.EndUser{}) is the zero value when no user is found; this may need adjustment
	if (m.EndUser{}) == endUser {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Construct the response with user details, token, and its expiry
	resp := m.ApiResp{
		Code: http.StatusOK,
		// Token:    token,
		// Expire:   expiry.Format(time.RFC3339),
		EndUsers: endUser,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func GetOtpCountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	count, err := dao.GetOtpCount()
	if err != nil {
		http.Error(w, "Failed to get OTP count", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]int64{"count": count}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func SaveOtpDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rec m.OtpDetail
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := dao.SaveOtpDetails(&rec); err != nil {
		log.Printf("Error saving OTP details: %v", err)
		http.Error(w, "Failed to save OTP details", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP details saved successfully"})
}

func FindOtpDetailsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mobileNo, ok := r.URL.Query()["mobileNo"]
	if !ok || len(mobileNo[0]) < 1 {
		http.Error(w, "Mobile number is required", http.StatusBadRequest)
		return
	}

	mobileNoFromReq, err := strconv.ParseUint(mobileNo[0], 10, 64)
	if err != nil {
		http.Error(w, "Invalid mobile number format", http.StatusBadRequest)
		return
	}

	otpDetail, err := dao.FindOtpDetailsForMobileNo(mobileNoFromReq)
	if err != nil {
		log.Printf("Error finding OTP details: %v", err)
		http.Error(w, "Failed to find OTP details", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(otpDetail); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
