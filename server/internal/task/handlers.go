package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Handlers holds the dependencies HTTP handlers need and exposes plain
// http.HandlerFunc values — no framework-specific handler signature, so
// they register directly on the stdlib http.ServeMux in main.go.
type Handlers struct {
	store *Store
}

func NewHandlers(store *Store) *Handlers {
	return &Handlers{store: store}
}

type taskInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.List(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var in taskInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	if in.Title == "" {
		writeError(w, http.StatusBadRequest, errors.New("title is required"))
		return
	}

	t, err := h.store.Create(r.Context(), in.Title, in.Description)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handlers) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid task id"))
		return
	}

	var in taskInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}
	if in.Title == "" {
		writeError(w, http.StatusBadRequest, errors.New("title is required"))
		return
	}

	t, err := h.store.Update(r.Context(), id, in.Title, in.Description, in.Done)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, t)
}

func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid task id"))
		return
	}

	if err := h.store.Delete(r.Context(), id); errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
