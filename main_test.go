package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	"testing"

	"github.com/Golang-energetics-collection/models"
)

func TestGenerateUsers(t *testing.T) { //unit test

	capturedOutput := captureOutput(func() {
		GenerateUsers()
	})

	expectedOutput := "User"
	if !strings.Contains(capturedOutput, expectedOutput) { // if the output contains "User"
		t.Errorf("Expected output to contain '%s', got '%s'", expectedOutput, capturedOutput)
	}

	lines := strings.Split(capturedOutput, "\n")
	expectedLines := 101

	if len(lines) != expectedLines {
		t.Errorf("Expected %d lines, got %d lines", expectedLines, len(lines))

	}
}

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = old
	return string(out)
}

func TestGetNewUsers(t *testing.T) { // integration test
	req, err := http.NewRequest("GET", "/newusers", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	GetNewUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}

	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("Handler returned wrong content type: got %v, want %v", contentType, expectedContentType)
	}

	var receivedUsers []models.User
	err = json.Unmarshal(rr.Body.Bytes(), &receivedUsers) // check if there are users
	if err != nil {
		t.Errorf("Error decoding JSON response: %v", err)
	}

	expectedUsers := GenerateUsers()
	if len(receivedUsers) != len(expectedUsers) { // check users amount
		t.Errorf("Handler returned wrong number of users: got %v, want %v", len(receivedUsers), len(expectedUsers))
	}
}

// func TestLoginFront(t *testing.T) {
// 	const (
// 		seleniumPath     = "./tools/selenium-server-standalone-3.4.0.jar"
// 		chromeDriverPath = "./tools/chromedriver.exe"
// 		port             = 4444
// 	)
// 	opts := []selenium.ServiceOption{
// 		selenium.StartFrameBuffer(),
// 		selenium.ChromeDriver(chromeDriverPath),
// 	}
// 	selenium.SetDebug(true)
// 	service, err := selenium.NewSeleniumService(seleniumPath, port, opts...)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer service.Stop()

// 	caps := selenium.Capabilities{"browserName": "chrome"}
// 	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	// caps := selenium.Capabilities{"browserName": "chrome"}
// 	// chromeCaps := chrome.Capabilities{Path: "C:/Program Files (x86)/Google/Chrome/Application/chrome.exe"}
// 	// caps.AddChrome(chromeCaps)
// 	// wd, err := selenium.NewRemote(caps, "")
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	defer wd.Quit()
// 	wd.Get("http://localhost:8080/")
// 	loginButton, err := wd.FindElement(selenium.ByID, "nav-login")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	loginButton.Click()
// 	signinButton, err := wd.FindElement(selenium.ByID, "nav-signin")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	signinButton.Click()
// 	emailInpt, err := wd.FindElement(selenium.ByID, "logmail")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	emailInpt.SendKeys("letun_igor007@mail.ru")
// 	pwdInpt, err := wd.FindElement(selenium.ByID, "logpass")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	pwdInpt.SendKeys("321")
// 	submitButton, err := wd.FindElement(selenium.ByID, "logbtn")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	submitButton.Click()
// 	cookies, err := wd.GetCookies()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if len(cookies) == 0 {
// 		t.Error("No cookies found")
// 	}
// }
