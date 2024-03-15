package app

import (
	"encoding/json"
	"github.com/Dubjay18/gobank_auth/dto"
	"github.com/Dubjay18/gobank_auth/logger"
	"github.com/Dubjay18/gobank_auth/service"
	"net/http"
)

func WriteJson(w http.ResponseWriter, i interface{}, code ...int) {
	if code == nil {

		code = append(code, http.StatusOK)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code[0])
	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		return
	}
}

type AuthHandler struct {
	service service.AuthService
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Error("Error while decoding login request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Login(loginRequest)
		if appErr != nil {
			WriteJson(w, appErr.AsMessage(), appErr.Code)
		} else {
			WriteJson(w, *token, http.StatusOK)
		}
	}
}
