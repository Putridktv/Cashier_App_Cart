package api

import (
	"cashierAppCart/model"
	"context"
	"encoding/json"
	"net/http"
)

func (api *API) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			errorResponse := model.ErrorResponse{Error: "http: named cookie not present"}
			convErr, _ := json.Marshal(errorResponse)
			w.WriteHeader(401)
			w.Write(convErr)
			return
		}

		sessionToken := cookie.Value

		sessionFound, err := api.sessionsRepo.CheckExpireToken(sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), "username", sessionFound.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (api *API) Get(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
		} else {
			errorResponse := model.ErrorResponse{Error: "Method is not allowed!"}
			convErr, _ := json.Marshal(errorResponse)
			w.WriteHeader(405)
			w.Write(convErr)
			return
		}
	})
}

func (api *API) Post(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
		} else {
			errorResponse := model.ErrorResponse{Error: "Method is not allowed!"}
			convErr, _ := json.Marshal(errorResponse)
			w.WriteHeader(405)
			w.Write(convErr)
			return
		}

	})
}
