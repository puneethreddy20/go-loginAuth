package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

const postMethod = "POST"
const getMethod = "GET"

//Checks CSRF
func checkCSRF(w http.ResponseWriter, r *http.Request) (bool, error) {
	if r.Method != getMethod {
		referer := r.Referer()
		if len(referer) > 0 && len(r.Host) > 0 {
			log.Println(3, "ref =%s, host=%s", referer, r.Host)
			refererURL, err := url.Parse(referer)
			if err != nil {
				log.Println(err)
				return false, err
			}
			log.Println(3, "refHost =%s, host=%s", refererURL.Host, r.Host)
			if refererURL.Host != r.Host {
				log.Printf("CSRF detected.... rejecting with a 400")
				http.Error(w, "you are not authorized", http.StatusUnauthorized)
				err := errors.New("CSRF detected... rejecting")
				return false, err

			}
		}
	}
	return true, nil
}

//generate 32bit random string which will be used as cookie value.
func randomStringGeneration() (string, error) {
	const size = 32
	bytesval := make([]byte, size)
	_, err := rand.Read(bytesval)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytesval), nil
}

//Authentication using cookie and returns the username.
func (state *RuntimeState) GetRemoteUserName(w http.ResponseWriter, r *http.Request) (string, error) {
	_, err := checkCSRF(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return "", err
	}
	remoteCookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, loginPath, http.StatusFound)
		return "", err
	}
	state.cookiemutex.Lock()
	cookieInfo, ok := state.authcookies[remoteCookie.Value]
	state.cookiemutex.Unlock()

	if !ok {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return "", nil
	}
	if cookieInfo.ExpiresAt.Before(time.Now()) {
		http.Redirect(w, r, loginPath, http.StatusFound)
		return "", nil
	}
	return cookieInfo.Username, nil
}

//Initial login handler and checks username, password if authenticated then creates cookies for future authentication.
func (state *RuntimeState) Checklogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != postMethod {
		http.Error(w, "you are not authorized", http.StatusMethodNotAllowed)
		return
	}
	_, err := checkCSRF(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return
	}
	if r.Method != postMethod {
		http.Error(w, "Unaurthorized", http.StatusMethodNotAllowed)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("Error while Parsing Form")
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	checkpasswd, err := state.ValidateuserPassword(username, password)
	if !checkpasswd || err != nil {
		log.Println("Password is wrong")
		http.Error(w, "password is wrong", http.StatusUnauthorized)
		return
	}
	randomString, err := randomStringGeneration()
	if err != nil {
		log.Println(err)
		http.Error(w, "cannot generate random string", http.StatusInternalServerError)
		return
	}

	expires := time.Now().Add(time.Hour * cookieExpirationHours)

	usercookie := http.Cookie{Name: cookieName, Value: randomString, Path: indexPath, Expires: expires, HttpOnly: true}

	http.SetCookie(w, &usercookie)

	Cookieinfo := cookieInfo{username, usercookie.Expires}

	state.cookiemutex.Lock()
	state.authcookies[usercookie.Value] = Cookieinfo
	state.cookiemutex.Unlock()

	http.Redirect(w, r, indexPath, http.StatusFound)
}

//Handler to register user.
func (state *RuntimeState) RegisterUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != postMethod {
		http.Error(w, "you are not authorized", http.StatusMethodNotAllowed)
		return
	}
	_, err := checkCSRF(w, r)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusUnauthorized)
		return
	}
	if r.Method != postMethod {
		http.Error(w, "Unaurthorized", http.StatusMethodNotAllowed)
		return
	}
	err = r.ParseForm()
	if err != nil {
		log.Println("Error while Parsing Form")
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	confirmpassword := r.PostFormValue("confirm-password")
	email := r.PostFormValue("email")
	if password != confirmpassword {
		log.Println("Bad Request")
		http.Error(w, "Passwords didnot match", http.StatusBadRequest)
		return
	}
	err = state.HashpasswordandInsertinDB(email, username, password)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	http.Error(w, "User successfully Created! please login now.", http.StatusOK)
}

//logout handler
func (state *RuntimeState) logoutHanler(w http.ResponseWriter, r *http.Request){
	cookie := http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
		Expires:time.Now(),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w,r,loginPath,http.StatusFound)

}
//Login page (or) Sign Up Web page.
func (state *RuntimeState) loginPageHandler(w http.ResponseWriter, r *http.Request) {

	generateHTML(w, nil, state.Config.Base.TemplatesPath, "newlogin")

}

//Main Page
func (state *RuntimeState) IndexHandler(w http.ResponseWriter, r *http.Request) {
	_, err := state.GetRemoteUserName(w, r)
	if err != nil {
		return
	}
	generateHTML(w, nil, state.Config.Base.TemplatesPath, "Index")

}

// Greeting page, which makes API call to greetingAPI with JWT token
func (state *RuntimeState) GreetingHandler(w http.ResponseWriter, r *http.Request) {
	username, err := state.GetRemoteUserName(w, r)
	if err != nil {
		return
	}

	req, err := http.NewRequest("GET", state.Config.Base.GreetingAPIHost, nil)

	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	state.tokenmutex.Lock()
	tokeninfo, ok := state.jwtTokenmap[username]
	state.tokenmutex.Unlock()

	if !ok || tokeninfo.ExpiresAt.Before(time.Now()) {

		JWTtoken, err := state.createandSetToken(username)

		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprint("hey", err), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", JWTtoken)

	} else {
		req.Header.Set("Authorization", tokeninfo.Token)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}

	if res.StatusCode == 400 {
		log.Println("Invalid token. Please Contact Administrator.")
		delete(state.jwtTokenmap, username)
		http.Error(w, "Invalid token. Please Contact Administrator.", http.StatusInternalServerError)
		return
	}
	if res.StatusCode == 401 {
		log.Println("Signature expired. Please try again/ reload the page")
		delete(state.jwtTokenmap, username)
		http.Error(w, "Signature expired. Please try again/reload the page", http.StatusInternalServerError)
		return
	}

	// Read the response body
	buf := new(bytes.Buffer)
	io.Copy(buf, res.Body)
	res.Body.Close()
	fmt.Println(buf.String())
	//response:=Response{buf.String()}
	io.WriteString(w, buf.String())
	//generateHTML(w, response, state.Config.Base.TemplatesPath, "Index")

}
