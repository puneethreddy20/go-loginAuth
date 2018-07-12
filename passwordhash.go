package main

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log"
)

//hash user password using bcrypt and insert in DB. (while registering)
func (state *RuntimeState) HashpasswordandInsertinDB(useremail string, username string, userpassword string) error {

	hashvalue, err := bcrypt.GenerateFromPassword([]byte(userpassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = insertUserInDB(useremail, username, string(hashvalue), state)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (state *RuntimeState) ValidateuserPassword(username string, userpassword string) (bool, error) {
	getpassword, _, err := findpasswdofUserinDB(username, state)
	if err != nil {
		log.Println(err)
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(getpassword), []byte(userpassword))
	if err != nil {
		log.Println("password is wrong!")
		return false, errors.New("password is wrong!")
	}
	return true, nil
}
