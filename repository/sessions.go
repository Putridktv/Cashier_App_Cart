package repository

import (
	"cashierAppCart/db"
	"cashierAppCart/model"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type SessionsRepository struct {
	db db.DB
}

func NewSessionsRepository(db db.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) ReadSessions() ([]model.Session, error) {
	records, err := u.db.Load("sessions")
	if err != nil {
		return nil, err
	}

	var listSessions []model.Session
	err = json.Unmarshal([]byte(records), &listSessions)
	if err != nil {
		return nil, err
	}

	return listSessions, nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}

	for i := 0; i < len(listSessions); i++ {
		if listSessions[i].Token == tokenTarget {
			listSessions[i] = model.Session{} //menghapus token yang sama dengan memasukkan nilai kosong
		}
	}

	jsonData, err := json.Marshal(listSessions)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	ses, err := u.ReadSessions() //buat variabel untuk nampung baca fungsi
	if err != nil {
		return err
	}

	ses = append(ses, session)

	data, err := json.Marshal(ses) //conv agar bia disave
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", data) //sava data ke db
	if err != nil {
		return err
	}
	return nil
}

func (u *SessionsRepository) CheckExpireToken(token string) (model.Session, error) {
	check, err := u.TokenExist(token) //membaca token dalam paramemter dan mengecek ada atau tdk dengan fungsi TokenExist
	if err != nil {
		return model.Session{}, err
	}

	checkToken := u.TokenExpired(check) //mengecek expired data
	if !checkToken {
		return check, nil
	}
	if checkToken {
		return model.Session{}, errors.New("Token is Expired!")
	}

	return model.Session{}, nil
}

func (u *SessionsRepository) ResetSessions() error {
	err := u.db.Reset("sessions", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) TokenExist(req string) (model.Session, error) {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return model.Session{}, err
	}
	for _, element := range listSessions {
		if element.Token == req {
			return element, nil
		}
	}
	return model.Session{}, fmt.Errorf("Token Not Found!")
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
