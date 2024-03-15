package app

import (
	"encoding/json"
	"fmt"
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

func (h AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// converting from Query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		appErr := h.service.Verify(urlParams)
		fmt.Println(appErr, urlParams)
		if appErr != nil {
			WriteJson(w, notAuthorizedResponse(appErr.Message), appErr.Code)
		} else {
			WriteJson(w, authorizedResponse(), http.StatusOK)
		}
	} else {
		WriteJson(w, notAuthorizedResponse("missing token"), http.StatusForbidden)
	}
}

func notAuthorizedResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"isAuthorized": false,
		"message":      msg,
	}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}
}
