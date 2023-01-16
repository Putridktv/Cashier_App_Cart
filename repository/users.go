package repository

import (
	"cashierAppCart/db"
	"cashierAppCart/model"
	"encoding/json"
	"fmt"
)

type UserRepository struct {
	db db.DB
}

func NewUserRepository(db db.DB) UserRepository {
	return UserRepository{db}
}

func (u *UserRepository) ReadUser() ([]model.Credentials, error) {
	records, err := u.db.Load("users")
	if err != nil {
		return nil, err
	}

	var listUser []model.Credentials
	err = json.Unmarshal([]byte(records), &listUser)
	if err != nil {
		return nil, err
	}

	return listUser, nil
}

func (u *UserRepository) AddUser(creds model.Credentials) error {
	checkUser, err := u.ReadUser() //read data user
	if err != nil {
		return err
	}

	checkUser = append(checkUser, creds) //add user ke variabel

	data, err := json.Marshal(checkUser) // conv agar bisa add ke db
	if err != nil {
		return err
	}

	err = u.db.Save("users", data) //save data ke db dengan func save json.go
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) ResetUser() error {
	err := u.db.Reset("users", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) LoginValid(req model.Credentials) error {
	listUser, err := u.ReadUser()
	if err != nil {
		return err
	}

	for _, element := range listUser {
		if element.Username == req.Username && element.Password == req.Password {
			return nil
		}
	}

	return fmt.Errorf("Wrong User or Password!")
}
