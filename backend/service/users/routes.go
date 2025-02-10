package users

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexhool2/TimeCard/config"
	"github.com/alexhool2/TimeCard/service/auth"
	"github.com/alexhool2/TimeCard/types"
	"github.com/alexhool2/TimeCard/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	store      types.UserStore
	eventStore types.EventStore
}

func NewHandler(store types.UserStore, eventStore types.EventStore) *Handler {
	return &Handler{store: store, eventStore: eventStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.Handle("/me", auth.AuthMiddleWare(http.HandlerFunc(h.handleGetUser))).Methods("GET")
	router.Handle("/logout", auth.AuthMiddleWare(http.HandlerFunc(h.handleLogout))).Methods("POST")
	router.Handle("/create", auth.AuthMiddleWare(auth.AdminOnly(http.HandlerFunc(h.handleRegister)))).Methods("POST")
	router.Handle("/{id}/reset-password", auth.AuthMiddleWare(auth.AdminOnly(http.HandlerFunc(h.handleResetPassword)))).Methods("PUT")
	router.Handle("", auth.AuthMiddleWare(auth.AdminOnly(http.HandlerFunc(h.handleGetAllUsers)))).Methods("GET")

}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user types.LoginUserPayload
	if err := utils.ParseJson(r, &user); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate the user

	if err := utils.Validate.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user %v", errors))
		return
	}

	// check if user exists

	u, err := h.store.GetUserByUserName(user.UserName)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid username or password"))
		return
	}

	if !auth.ComparePasswords(u.Password, []byte(user.Password)) {
		log.Printf("Password mismatch for username '%s'", user.UserName)

		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found invalid username or password"))
		return
	}
	secret := []byte(config.Envs.JWTSecret)
	token, err := auth.CreateJWT(secret, u.ID, string(u.Role))
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	// Create a cookie to store the JWT token

	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	// Set the cookie in the response
	http.SetCookie(w, &cookie)

	response := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       u.ID,
			"userName": u.UserName,
			"role":     u.Role,
		},
	}

	utils.WriteJSON(w, http.StatusOK, response)

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	// get json payload
	var payload types.RegisterUserPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	//validate role

	if !types.IsValidRole(payload.Role) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid role %s", payload.Role))
		return
	}

	// validate the payload

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	//check if user exists
	_, err := h.store.GetUserByUserName(payload.UserName)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with username %s already exists", payload.UserName))
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	//creating a new user

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		UserName:  payload.UserName,
		Password:  hashedPassword,
		Email:     payload.Email,
		Role:      payload.Role,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleResetPassword(w http.ResponseWriter, r *http.Request) {

	//get user id
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id: %v", err))
	}

	// get the new password from payload
	var payload types.ResetPassword
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate new password
	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid password %v", errors))
		return
	}

	// generate new hash
	hashedPassword, err := auth.HashPassword(payload.NewPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to hash password: %v", err))
		return
	}

	// update password
	err = h.store.UpdatePassword(id, hashedPassword)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to reset password: %v", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, map[string]interface{}{"message": "password reset successfully"})

}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	// get user id
	userID, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized %v", err))
		return
	}
	user, err := h.store.GetUserById(userID.ID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to get user: %v", err))
		return
	}

	// go does not accept nil values so we are turning event to -1, if wont change means that we dont have a new event or its nil
	eventID := -1

	// Tenta buscar o evento do usuário
	retrievedEventID, err := h.eventStore.GetTodayEvent(user.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// if its not possible to find the event id, will keep the same
			eventID = -1
		} else if eventID != -1 {
			// checking others errors
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to get event: %w", err))
			return
		}
	} else {
		eventID = retrievedEventID // find the new event
	}
	response := map[string]interface{}{
		"user":    user,
		"eventID": eventID,
	}
	utils.WriteJSON(w, http.StatusOK, response)

}

func (h *Handler) handleGetAllUsers(w http.ResponseWriter, r *http.Request) {
	searchQuery := r.URL.Query().Get("search")
	var users []types.User
	var err error

	if searchQuery != "" {

		users, err = h.store.GetDynamicUsers(searchQuery)
	} else {
		users, err = h.store.GetAllUsers()
	}

	if err != nil {

		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching user: %v", err))
		return
	}
	utils.WriteJSON(w, http.StatusOK, users)
}
func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
	w.Header().Set("Cache-Control", "no-store")
	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Logout successful"})
}
