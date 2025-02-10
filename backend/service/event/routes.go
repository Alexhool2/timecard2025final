package event

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/alexhool2/TimeCard/service/auth"
	"github.com/alexhool2/TimeCard/types"
	"github.com/alexhool2/TimeCard/utils"
	"github.com/gorilla/mux"
)

type Handler struct {
	store     types.EventStore
	userStore types.UserStore
}

func NewHandler(store types.EventStore, userStore types.UserStore) *Handler {

	return &Handler{store: store, userStore: userStore}

}

func (h *Handler) RegisterRoutes(router *mux.Router) {

	authRouter := router.PathPrefix("").Subrouter()
	authRouter.Use(auth.AuthMiddleWare)
	authRouter.HandleFunc("/create", h.handleCreateEvent).Methods("POST")

	authRouter.HandleFunc("/start-lunch/{id}", h.UpdateStartLunchEvent).Methods("PUT")
	authRouter.HandleFunc("/end-lunch/{id}", h.UpdateEndLunchEvent).Methods("PUT")
	authRouter.HandleFunc("/end-time/{id}", h.UpdateEndTimeEvent).Methods("PUT")
	router.Handle("/user/{id}/date", auth.AuthMiddleWare(auth.AdminOnly(http.HandlerFunc(h.handleGetEvent)))).Methods("GET")
	router.Handle("/user/{id}/period", auth.AuthMiddleWare(auth.AdminOnly(http.HandlerFunc(h.handlerGetEventBySelectedDates)))).Methods("GET")

}

func (h *Handler) handleCreateEvent(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateEventPayload

	// json to body
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid JSON: %v", err))
		return
	}

	payload.Date = time.Now() // converting to show only date
	payload.Start_time = time.Now()

	// validate user
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return
	}

	// check if user id authentication is valid
	if payload.UserID != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("forbidden: cannot create an event for another user"))
		return
	}

	//validate fields
	if err := utils.Validate.Struct(payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", err))
	}

	// check user on db
	_, err = h.userStore.GetUserById(payload.UserID)
	if err != nil {
		if err.Error() == "user not found" {

			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with ID %d does not exist", payload.UserID))
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching user: %v", err))
		return
	}

	// Insert event on db
	eventID, err := h.store.CreateEvent(payload)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to insert event: %v", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.WriteJSON(w, http.StatusCreated, map[string]interface{}{
		"eventId":   eventID,
		"startTime": payload.Start_time.Format("2006-01-02 15:04:05"),
	})
}

func (h *Handler) UpdateStartLunchEvent(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid event id"))
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized %v", err))
		return
	}

	// check if is there start lunch event
	event, err := h.store.GetStartLunch(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// check if already this event exists
	if event.UserID != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("forbidden: cannot update an event for another user"))
		return
	}

	// update start lunch
	startLunch := time.Now()
	err = h.store.UpdateStartLunchEvent(event.ID, startLunch)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update start_lunch: %v", err))
		return
	}

	response := map[string]interface{}{
		"message":    "start_lunch updated successfully",
		"eventId":    event.ID,
		"startLunch": startLunch.Format("2006-01-02 15:04:05"),
	}
	utils.WriteJSON(w, http.StatusOK, response)
}
func (h *Handler) UpdateEndLunchEvent(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid event id"))
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized %v", err))
		return
	}

	// check if is there already created field
	event, err := h.store.GetEndLunch(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if event.UserID != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("forbidden: cannot create an event for another user"))
		return
	}

	// update end lunch
	endLunch := time.Now()
	err = h.store.UpdateEndLunchEvent(event.ID, endLunch)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update end_lunch: %v", err))
		return
	}

	response := map[string]interface{}{
		"message":  "end_lunch updated successfully",
		"eventId":  event.ID,
		"endLunch": endLunch.Format("2006-01-02 15:04:05"),
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) UpdateEndTimeEvent(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid event id"))
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized %v", err))
		return
	}

	// check if event exists
	event, err := h.store.GetEndTime(id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}
	if event.UserID != user.ID {
		utils.WriteError(w, http.StatusForbidden, fmt.Errorf("forbidden: cannot create an event for another user"))
		return
	}

	// update field
	endTime := time.Now()
	err = h.store.UpdateEndTimeEvent(event.ID, endTime)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to update end_Time: %v", err))
		return
	}

	response := map[string]interface{}{
		"message": "end_time updated successfully",
		"eventId": event.ID,
		"endTime": endTime.Format("2006-01-02 15:04:05"),
	}
	utils.WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	query := r.URL.Query()
	dateStr := query.Get("date")
	if dateStr == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("date query parameter is missing"))
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid date format, use YYYY-MM-DD"))
		return
	}

	_, err = auth.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return
	}

	event, err := h.store.GetEventByDate(date, userId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching event: %v", err))
		return
	}
	if event == nil {
		utils.WriteError(w, http.StatusNotFound, fmt.Errorf("event not found"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, event)
}

func (h *Handler) handlerGetEventBySelectedDates(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		log.Printf("Error converting user id: %v", err)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	query := r.URL.Query()
	startDateStr := query.Get("start_date")
	endDateStr := query.Get("end_date")
	if startDateStr == "" || endDateStr == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("date query parameter is missing"))
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid date format, use YYYY-MM-DD"))
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid date format, use YYYY-MM-DD"))
		return
	}

	_, err = auth.GetUserFromContext(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return
	}

	if startDate.After(time.Now()) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("cannot check future dates"))
		return
	}

	event, err := h.store.GetEventBySelectedDates(userId, startDate, endDate)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error fetching event: %v", err))
		return
	}

	utils.WriteJSON(w, http.StatusOK, event)
}
