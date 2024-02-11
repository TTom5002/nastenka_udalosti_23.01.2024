package helpers

import (
	"errors"
	"fmt"
	"nastenka_udalosti/internal/config"
	"nastenka_udalosti/internal/models"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers nastaví app config pro helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError vypíše error ze strany uživatele
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError vypíše error ze strany serveru
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// IsAuthenticated ověří, že uživatel je příhlášen
func IsAuthenticated(r *http.Request) (bool, error) {
	if !app.Session.Exists(r.Context(), "userInfo") {
		return false, nil
	}

	userInfo, ok := app.Session.Get(r.Context(), "userInfo").(models.User)
	if !ok {
		return false, errors.New("nepodařilo se přetypovat hodnotu na User")
	}
	id := userInfo.ID
	if id != 0 {
		return true, nil
	}
	return false, errors.New("nepodařilo se získat id uživatele")
}

// IsVerified ověří, že uživatel je ověřený (user_verified = true, false)
func IsVerified(r *http.Request) (bool, error) {
	userInfo, ok := app.Session.Get(r.Context(), "userInfo").(models.User)
	if !ok {
		return false, errors.New("nepodařilo se přetypovat hodnotu na User")
	}
	ver := userInfo.Verified
	return ver, nil
}

// GetAdminAccessLevel ověří platnost admina (user_access_level = 3)
func GetAdminAccessLevel(r *http.Request) (bool, error) {
	userInfo, ok := app.Session.Get(r.Context(), "userInfo").(models.User)
	if !ok {
		return false, errors.New("nepodařilo se přetypovat hodnotu na User")
	}
	accessLevel := userInfo.AccessLevel
	if accessLevel != 3 {
		return false, errors.New("nepovedlo se ověřit uživatele")
	}
	return true, nil
}

// GetUserInfo získá potřebné informace od uživatele
func GetUserInfo(r *http.Request) (models.User, error) {
	userInfo, ok := app.Session.Get(r.Context(), "userInfo").(models.User)
	if !ok {
		return models.User{}, errors.New("nepodařilo se přetypovat hodnotu na User")
	}
	return userInfo, nil
}
