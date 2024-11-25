package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ViniNepo/secretfriend/domain"
	"github.com/ViniNepo/secretfriend/services"
	"github.com/gorilla/mux"
)

type FriendHandlers struct {
	friendService services.FriendService
}

func NewFriendHandlers(service services.FriendService) *FriendHandlers {
	return &FriendHandlers{
		friendService: service,
	}
}

func (h *FriendHandlers) CreateFriendHandlers(router *mux.Router) {
	router.HandleFunc("/friend", h.create).Methods("POST")
	router.HandleFunc("/friend/reminder", h.reminder).Methods("GET")
	router.HandleFunc("/friend/shuffle", h.shuffle).Methods("GET")
	router.HandleFunc("/friend/validate", h.validate).Methods("PATCH")
}

func (h *FriendHandlers) create(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)

	var friend domain.Friend
	if err := json.NewDecoder(r.Body).Decode(&friend); err != nil {
		HandleError(ErrRequestBodyIsInvalid, w)
		return
	}

	code, err := generateRandomCode(6)
	if err != nil {
		HandleError(err, w)
	}

	friend.ValidateCode = code

	id, err := h.friendService.Create(friend)
	if err != nil {
		HandleError(err, w)
		return
	}

	// Construindo a resposta com o ID
	response := map[string]interface{}{
		"id": id,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		HandleError(fmt.Errorf("failed to encode response: %w", err), w)
	}
}

func (h *FriendHandlers) reminder(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)

	err := h.friendService.Reminder()
	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reminder send"))
}

func (h *FriendHandlers) shuffle(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)

	err := h.friendService.Shuffle()
	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shuffle send"))
}

func (h *FriendHandlers) validate(w http.ResponseWriter, r *http.Request) {
	defer handleError(w)

	var request domain.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		HandleError(ErrRequestBodyIsInvalid, w)
		return
	}

	err := h.friendService.Validate(request)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Friend validated"))
}

func generateRandomCode(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		randomInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[randomInt.Int64()]
	}
	return string(result), nil
}
