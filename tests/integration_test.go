package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

const baseURL = "http://localhost:3000/api/v1"

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestEndToEndFlow(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// 1. Register Student
	studentUser := "student_" + randomString(5)
	studentPass := "password123"
	studentEmail := studentUser + "@example.com"
	
	t.Logf("Registering Student: %s", studentUser)
	registerPayload := map[string]string{
		"username": studentUser,
		"email":    studentEmail,
		"password": studentPass,
		"fullName": "Test Student",
		"roleName": "Mahasiswa",
	}
	resp := makeRequest(t, "POST", "/auth/register", nil, registerPayload)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to register student: %d", resp.StatusCode)
	}

	// 2. Login Student
	t.Log("Logging in Student")
	loginPayload := map[string]string{
		"username": studentUser,
		"password": studentPass,
	}
	resp = makeRequest(t, "POST", "/auth/login", nil, loginPayload)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to login student: %d", resp.StatusCode)
	}
	var loginData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&loginData)
	studentToken := loginData["data"].(map[string]interface{})["token"].(string)

	// 3. Create Achievement
	t.Log("Creating Achievement")
	achPayload := map[string]interface{}{
		"title":           "Test Achievement",
		"achievementType": "competition",
		"description":     "Won a test competition",
		"points":          100,
		"tags":            []string{"test", "competition"},
		"details": map[string]interface{}{
			"rank": 1,
		},
	}
	resp = makeRequest(t, "POST", "/achievements", &studentToken, achPayload)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to create achievement: %d", resp.StatusCode)
	}

	// 4. Get Student Achievements
	t.Log("Getting Student Achievements")
	resp = makeRequest(t, "GET", "/achievements", &studentToken, nil)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to get achievements: %d", resp.StatusCode)
	}
	var achList map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&achList)
	data := achList["data"].([]interface{})
	if len(data) == 0 {
		t.Fatal("No achievements found")
	}
	firstAch := data[0].(map[string]interface{})
	achID := firstAch["id"].(string)

	// 5. Submit Achievement
	t.Logf("Submitting Achievement: %s", achID)
	resp = makeRequest(t, "POST", "/achievements/"+achID+"/submit", &studentToken, nil)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to submit achievement: %d", resp.StatusCode)
	}

	// 6. Register Advisor
	advisorUser := "advisor_" + randomString(5)
	advisorPass := "password123"
	advisorEmail := advisorUser + "@example.com"
	
	t.Logf("Registering Advisor: %s", advisorUser)
	registerPayload["username"] = advisorUser
	registerPayload["email"] = advisorEmail
	registerPayload["password"] = advisorPass
	registerPayload["roleName"] = "Dosen Wali"
	resp = makeRequest(t, "POST", "/auth/register", nil, registerPayload)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to register advisor: %d", resp.StatusCode)
	}

	// 7. Login Advisor
	t.Log("Logging in Advisor")
	loginPayload["username"] = advisorUser
	loginPayload["password"] = advisorPass
	resp = makeRequest(t, "POST", "/auth/login", nil, loginPayload)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to login advisor: %d", resp.StatusCode)
	}
	var advisorLoginData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&advisorLoginData)
	advisorToken := advisorLoginData["data"].(map[string]interface{})["token"].(string)

	// 8. Verify Achievement
	t.Log("Verifying Achievement")
	resp = makeRequest(t, "POST", "/verification/"+achID+"/verify", &advisorToken, nil)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to verify achievement: %d", resp.StatusCode)
	}

	// 9. Check Statistics (Admin/Anyone)
	t.Log("Checking Statistics")
	resp = makeRequest(t, "GET", "/reports/statistics", &advisorToken, nil)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to get statistics: %d", resp.StatusCode)
	}
	var stats map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&stats)
	t.Logf("Stats: %v", stats["data"])

	// 10. Register Admin
	adminUser := "admin_" + randomString(5)
	adminPass := "password123"
	adminEmail := adminUser + "@example.com"

	t.Logf("Registering Admin: %s", adminUser)
	registerPayload["username"] = adminUser
	registerPayload["email"] = adminEmail
	registerPayload["password"] = adminPass
	registerPayload["roleName"] = "Admin"
	resp = makeRequest(t, "POST", "/auth/register", nil, registerPayload)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to register admin: %d", resp.StatusCode)
	}

	// 11. Login Admin
	t.Log("Logging in Admin")
	loginPayload["username"] = adminUser
	loginPayload["password"] = adminPass
	resp = makeRequest(t, "POST", "/auth/login", nil, loginPayload)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to login admin: %d", resp.StatusCode)
	}
	var adminLoginData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&adminLoginData)
	adminToken := adminLoginData["data"].(map[string]interface{})["token"].(string)

	// 12. Get All Users (Admin)
	t.Log("Getting All Users")
	resp = makeRequest(t, "GET", "/users", &adminToken, nil)
	if resp.StatusCode != 200 {
		t.Fatalf("Failed to get all users: %d", resp.StatusCode)
	}

	// 13. Update User (Admin updates Student)
	// We need the student's ID. We can get it from the login response or by querying users.
	// But we didn't save student ID from login.
	// Let's use the student token to get their own profile if there was a /me endpoint, but there isn't.
	// However, we can find the student in the list of all users.
	var usersList map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&usersList)
	usersData := usersList["data"].([]interface{})
	var targetUserID string
	for _, u := range usersData {
		userMap := u.(map[string]interface{})
		if userMap["username"] == studentUser {
			targetUserID = userMap["id"].(string)
			break
		}
	}

	if targetUserID != "" {
		t.Logf("Updating User: %s", targetUserID)
		updatePayload := map[string]interface{}{
			"fullName": "Updated Student Name",
			"roleName": "Mahasiswa",
		}
		resp = makeRequest(t, "PUT", "/users/"+targetUserID, &adminToken, updatePayload)
		if resp.StatusCode != 200 {
			t.Fatalf("Failed to update user: %d", resp.StatusCode)
		}
	} else {
		t.Log("Could not find student user to update")
	}
}

func makeRequest(t *testing.T, method, path string, token *string, body interface{}) *http.Response {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, baseURL+path, bodyReader)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", "Bearer "+*token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to execute request: %v", err)
	}
	return resp
}
