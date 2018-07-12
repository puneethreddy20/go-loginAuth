package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"reflect"
	"testing"
)

func TestRuntimeState_HashpasswordandInsertinDB(t *testing.T) {
	state, err := Init()
	if err != nil {
		log.Println(err)
	}
	useremail := "test20@gmail.com"
	userpassword := "test20"
	username := "test_test"
	hashvalue, err := bcrypt.GenerateFromPassword([]byte(userpassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	err = insertUserInDB(useremail, username, string(hashvalue), &state)
	if err != nil {
		log.Println(err)
	}

}

func TestRuntimeState_ValidateuserPassword(t *testing.T) {
	state, err := Init()
	if err != nil {
		log.Println(err)
	}
	userpassword := "test20"
	username := "test_test"
	getpassword, _, err := findpasswdofUserinDB(username, &state)
	if err != nil {
		log.Println(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(getpassword), []byte(userpassword))
	if err != nil {
		log.Println("password is wrong!")
	}
	_, errtest := state.ValidateuserPassword(username, userpassword)
	if errtest != nil {
		log.Println("Error in getting actual result")
	}
	if !reflect.DeepEqual(err, errtest) {
		t.Errorf("The Result obtained doesnt match with actual result")
	}
}
