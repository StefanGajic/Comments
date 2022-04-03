package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Comments/internal/comment"
	"github.com/gorilla/mux"
)

// Handler stores pointer to comments service
type Handler struct {
	Router  *mux.Router
	Service *comment.Service
}

// Response is an object to store responses from API
type Response struct {
	Message string
	Error   string
}

// NewHandler returns a pointer to a Handler
func NewHandler(service *comment.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

// SetupRoutes sets up all the routes
func (h *Handler) SetupRoutes() {
	fmt.Println("Setting up Routers")
	h.Router = mux.NewRouter()

	h.Router.HandleFunc("/api/comment/{id}", h.GetAllComments).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.PostComment).Methods("POST")
	h.Router.HandleFunc("/api/comment/{id}", h.GetComment).Methods("GET")
	h.Router.HandleFunc("/api/comment/{id}", h.UpdateComment).Methods("PUT")
	h.Router.HandleFunc("/api/comment/{id}", h.DeleteComment).Methods("DELETE")

	h.Router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		if err := sendOkResponse(w, Response{Message: "Working fine!"}); err != nil {
			panic(err)
		}
	})

}

// GetAllComments retrieves all comments
func (h *Handler) GetAllComments(w http.ResponseWriter, r *http.Request) {
	comments, err := h.Service.GetAllComments()
	if err != nil {
		sendErrorResponse(w, "Failed to retrieve all comments", err)
		return
	}
	if sendOkResponse(w, comments); err != nil {
		panic(err)
	}
}

// PostComment posts new comment to the database
func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Failed to decode JSON body", err)
		return
	}

	comment, err := h.Service.PostComment(comment)
	if err != nil {
		sendErrorResponse(w, "Failed to post new comment", err)
		return
	}

	if sendOkResponse(w, comment); err != nil {
		panic(err)
	}
}

// GetComment retrieve comment by ID
func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	i, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UNIT for ID", err)
		return
	}

	comment, err := h.Service.GetComment(uint(i))
	if err != nil {
		sendErrorResponse(w, "Error retrieving comment by ID", err)
		return
	}

	if sendOkResponse(w, comment); err != nil {
		panic(err)
	}
}

// UpdateComment updates comment
func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	var comment comment.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		sendErrorResponse(w, "Failed to decode JSON body", err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UNIT for ID", err)
		return
	}

	comment, err = h.Service.UpdateComment(uint(commentID), comment)
	if err != nil {
		sendErrorResponse(w, "Failed to update comment", err)
		return
	}

	if sendOkResponse(w, comment); err != nil {
		panic(err)
	}
}

// DeleteComment deletes comment by ID
func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	commentID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		sendErrorResponse(w, "Unable to parse UNIT for ID", err)
		return
	}

	err = h.Service.DeleteComment(uint(commentID))
	if err != nil {
		sendErrorResponse(w, "Failed to delete comment by ID", err)
		return
	}

	if err = sendOkResponse(w, Response{Message: "Successfully deleted comment"}); err != nil {
		panic(err)
	}
}

func sendOkResponse(w http.ResponseWriter, resp interface{}) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(resp)
}

func sendErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	if err := json.NewEncoder(w).Encode(Response{
		Message: message,
		Error:   err.Error(),
	}); err != nil {
		panic(err)
	}
}
