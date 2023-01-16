package api

import (
	"cashierAppCart/model"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	// Read username and password request with FormValue.
	creds := model.Credentials{}
	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")

	if creds.Username == "" || creds.Password == "" {
		errorResponse := model.ErrorResponse{Error: "Username or Password empty"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(convErr)
		return
	} else {
		w.WriteHeader(200)
	}

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"

	err := api.usersRepo.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	filepath := path.Join("views", "status.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var data = map[string]string{"name": creds.Username, "message": "register success!"}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	// Read usernmae and password request with FormValue.
	creds := model.Credentials{}
	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")

	if creds.Username == "" || creds.Password == "" {
		errorResponse := model.ErrorResponse{Error: "Username or Password empty"}
		convErr, _ := json.Marshal(errorResponse)
		w.WriteHeader(400)
		w.Write(convErr)
		return
	}

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"

	err := api.usersRepo.LoginValid(creds)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	//set cookie
	sessionToken := uuid.New().String()
	// fmt.Println(sessionToken)

	expiresAt := time.Now().Add(5 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   sessionToken,
		Expires: expiresAt,
	})

	// Generate Cookie with Name "session_token", Path "/", Value "uuid generated with github.com/google/uuid", Expires time to 5 Hour.

	// session := model.Session{}
	session := model.Session{
		Username: creds.Username,
		Token:    sessionToken,
		Expiry:   expiresAt,
	}
	err = api.sessionsRepo.AddSessions(session)
	// if err != nil {
	// 	return err
	// }
	w.WriteHeader(http.StatusOK)

	api.dashboardView(w, r)
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	//Read session_token and get Value:
	w.WriteHeader(200)
	sessionToken := ""
	userName := fmt.Sprintf("%s", r.Context().Value("username")) //get username
	checkUser, _ := api.sessionsRepo.ReadSessions()
	for i := 0; i < len(checkUser); i++ {
		if userName == checkUser[i].Username {
			sessionToken = checkUser[i].Token
		}
	}
	w.WriteHeader(200)

	api.sessionsRepo.DeleteSessions(sessionToken)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Path:    "/",
		Value:   "",
		Expires: time.Now(),
	})
	//Set Cookie name session_token value to empty and set expires time to Now:

	filepath := path.Join("views", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}
