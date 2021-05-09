package proj_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMethod(t *testing.T) {
	// Switch to test mode so you don't get such noisy output
	// gin.SetMode(gin.TestMode)

	// Create the mock request you'd like to test. Make sure the second argument
	// here is the same as one of the routes you defined in the router setup
	// block!
	body := bytes.NewBuffer([]byte("{\"ApiTest\":\"GetReq\"}"))
	req, err := http.NewRequest("GET", "/api/v1/PersonId/Id456", body)
	if err != nil {
		t.Fatalf("Couldn't create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "adminpass")

	// Setup your router, just like you did in your main function, and
	// register your routes
	// router := gin.Default()
	// router.GET("/api/v1/PersonId/:IdValue", GetMethod)

	// Create a response recorder so you can inspect the response
	w := httptest.NewRecorder()

	// Perform the request
	//router.ServeHTTP(w, req)
	fmt.Println(w.Body)
}
