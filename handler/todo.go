package handler

import (
	"context"
	"encoding/json"
	"log"
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

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var cReq *model.CreateTODORequest

		if err := json.NewDecoder(r.Body).Decode(&cReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if cReq.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		res, err := h.Create(r.Context(), cReq)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}

		return
	case http.MethodPut:
		var uReq *model.UpdateTODORequest

		if err := json.NewDecoder(r.Body).Decode(&uReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if uReq.Subject == "" || uReq.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		res, err := h.Update(r.Context(), uReq)

		if err != nil {
			switch err.(type) {
			case model.ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}

		return
	case http.MethodGet:
		prevID := r.URL.Query().Get("prev_id")

		p := 0

		if prevID != "" {
			p, _ = strconv.Atoi(prevID)
		}

		size := r.URL.Query().Get("size")

		s := 0

		if size != "" {
			s, _ = strconv.Atoi(size)
		}

		req := &model.ReadTODORequest{
			PrevID: p,
			Size:   s,
		}

		res, err := h.Read(r.Context(), req)

		if err != nil {
			switch err.(type) {
			case model.ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}

		return
	case http.MethodDelete:
		var dReq *model.DeleteTODORequest

		if err := json.NewDecoder(r.Body).Decode(&dReq); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if dReq.IDs == nil || len(dReq.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		res, err := h.Delete(r.Context(), dReq)

		if err != nil {
			switch err.(type) {
			case model.ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Println(err)
		}

		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)

	if err != nil {
		return nil, err
	}

	return &model.CreateTODOResponse{
		TODO: *todo,
	}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	res, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)

	if err != nil {
		return nil, err
	}

	return &model.ReadTODOResponse{
		TODOs: res,
	}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	res, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)

	if err != nil {
		return nil, err
	}

	return &model.UpdateTODOResponse{
		TODO: *res,
	}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	if err := h.svc.DeleteTODO(ctx, req.IDs); err != nil {
		return &model.DeleteTODOResponse{}, err
	}

	return &model.DeleteTODOResponse{}, nil
}
