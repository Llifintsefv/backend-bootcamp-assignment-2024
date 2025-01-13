package flat

import (
	"backend-bootcamp-assignment-2024/dto"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type FlatHandler struct {
	flatService FlatService
}

func NewFlatHandler(flatService FlatService) *FlatHandler {
	return &FlatHandler{flatService: flatService}
}

func (h *FlatHandler) CreateFlat(w http.ResponseWriter, r *http.Request) {
	var req dto.PostFlatCreateJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

	}

	if req.Price <= 0 || *req.Rooms <= 0 {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	flatCreateResponse, err := h.flatService.CreateFlat(ctx, req)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(flatCreateResponse)
}

func (h *FlatHandler) GetFlats(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	houseId := params["id"]

	ctx := r.Context()

	flatGetResponse, err := h.flatService.GetFlats(ctx, houseId)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flatGetResponse)
}

func (h *FlatHandler) UpdateFlatStatus(w http.ResponseWriter, r *http.Request) {
	var req dto.PostFlatUpdateJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if *req.Status == "" && *req.Status != "created" && *req.Status != "updated" && *req.Status != "declined" && *req.Status != "on moderation" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	flatUpdateResponse, err := h.flatService.UpdateFlatStatus(ctx, req)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(flatUpdateResponse)
}
