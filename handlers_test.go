package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const (
	testdbpath      = "./test.db"
	testusername    = "test"
	cookievalueTest = "hellotestcookievalue"
	templatespath   = "templates"
	greetingAPI     = "http://172.18.0.20:8000/greetings"
)

func Init() (RuntimeState, error) {
	var state RuntimeState
	state.Config.Base.StorageURL = testdbpath
	state.Config.Base.GreetingAPIHost = greetingAPI
	state.Config.Base.TemplatesPath = templatespath
	err := initDB(&state)
	if err != nil {
		return state, err
	}
	state.authcookies = make(map[string]cookieInfo)
	expiresAt := time.Now().Add(time.Hour * cookieExpirationHours)
	usersession := cookieInfo{testusername, expiresAt}
	state.authcookies[cookievalueTest] = usersession
	return state, nil
}

func createCookie() http.Cookie {
	expiresAt := time.Now().Add(time.Hour * cookieExpirationHours)
	cookie := http.Cookie{Name: cookieName, Value: cookievalueTest, Path: indexPath, Expires: expiresAt, HttpOnly: true}
	return cookie
}

func TestRuntimeState_IndexHandler(t *testing.T) {
	state, err := Init()
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(getMethod, indexPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie := createCookie()
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(state.IndexHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRuntimeState_GreetingHandler(t *testing.T) {
	state, err := Init()
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest(getMethod, greetingPath, nil)
	if err != nil {
		t.Fatal(err)
	}
	cookie := createCookie()
	req.AddCookie(&cookie)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(state.IndexHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
