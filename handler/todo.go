package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// ServeHTTP handles an HTTP request for a TODO resource.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req model.CreateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			http.Error(w, "Bad Request: subject is empty", http.StatusBadRequest)
			return
		}

		createdTODO, err := h.svc.CreateTODO(r.Context(), req.Subject, req.Description)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		resp := model.CreateTODOResponse{
			TODO: *createdTODO,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)

	case http.MethodPut:
		var req model.UpdateTODORequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if req.ID == 0 || req.Subject == "" {
			http.Error(w, "Bad Request: id or subject is empty", http.StatusBadRequest)
			return
		}

		updatedTODO, err := h.svc.UpdateTODO(r.Context(), req.ID, req.Subject, req.Description)
		if err != nil {
			if errors.Is(err, &model.ErrNotFound{}) {
				http.Error(w, "Not Found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		resp := model.UpdateTODOResponse{
			TODO: *updatedTODO,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)

	case http.MethodGet:
		query := r.URL.Query()
		prevIDStr := query.Get("prev_id")
		sizeStr := query.Get("size")

		var prevID, size int64
		var err error

		if prevIDStr != "" {
			prevID, err = strconv.ParseInt(prevIDStr, 10, 64)
			if err != nil {
				http.Error(w, "Bad Request: invalid prev_id", http.StatusBadRequest)
				return
			}
		}
		if sizeStr != "" {
			size, err = strconv.ParseInt(sizeStr, 10, 64)
			if err != nil {
				http.Error(w, "Bad Request: invalid size", http.StatusBadRequest)
				return
			}
		}

		const defaultSize = 10
		if size <= 0 {
			size = defaultSize
		}

		todos, err := h.svc.ReadTODO(r.Context(), prevID, size)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		resp := model.ReadTODOResponse{
			TODOs: todos,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	createdTODO, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{
		TODO: *createdTODO,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	_, _ = h.svc.UpdateTODO(ctx, 0, "", "")
	return &model.UpdateTODOResponse{}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}
