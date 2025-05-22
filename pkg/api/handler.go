package api

import (
	"context"
	"encoding/json"
	"net/http"

	"infra-scheduler/pkg/models"
	"infra-scheduler/pkg/scheduler"
	pb "infra-scheduler/pkg/schedulerpb"

	"github.com/gorilla/mux"
)

type Handler struct {
	sched  *scheduler.Scheduler
	router *mux.Router
}

// NewHandler creates a new API handler with routes
func NewHandler(s *scheduler.Scheduler) *Handler {
	h := &Handler{sched: s, router: mux.NewRouter()}
	h.routes()
	return h
}

// routes registers HTTP endpoints
func (h *Handler) routes() {
	h.router.HandleFunc("/schedule", h.handleSchedule).Methods(http.MethodPost)
}

// Router returns the HTTP router
func (h *Handler) Router() http.Handler {
	return h.router
}

// handleSchedule processes VM scheduling requests
func (h *Handler) handleSchedule(w http.ResponseWriter, r *http.Request) {
	var req models.VMRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Convert to scheduler request
	schedReq := &pb.ScheduleRequest{
		Name:   req.UserID,
		Cpu:    int32(req.Cores),
		Memory: int32(req.Memory),
		Disk:   int32(req.Disk),
	}

	result, err := h.sched.ScheduleVM(context.Background(), schedReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
