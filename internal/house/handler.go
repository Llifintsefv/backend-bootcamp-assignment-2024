package house

import (
	"backend-bootcamp-assignment-2024/dto"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type HouseHandler struct {
	houseService HouseService
}

func NewHouseHandler(houseService HouseService) *HouseHandler {
	return &HouseHandler{
		houseService: houseService,
	}
}

func (h *HouseHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	var req dto.PostHouseCreateJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Address == "" || req.Year <= 0 {
		http.Error(w, "Address and year are required", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	houseResponse, err := h.houseService.CreateHouse(ctx, req)

	if err != nil {
		http.Error(w, "Failed to create house", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(houseResponse)

}

func (h *HouseHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var req dto.Email

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)

	houseId, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid house ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	err = h.houseService.Subscribe(ctx, houseId, req)
	if err != nil {
		http.Error(w, "Failed to subscribe", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
