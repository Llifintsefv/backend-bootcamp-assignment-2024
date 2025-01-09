package house

import (
	"backend-bootcamp-assignment-2024/dto"
	"encoding/json"
	"net/http"
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
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if req.Address == "" || req.Year <= 0 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	houseResponse, err := h.houseService.CreateHouse(ctx, req)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(houseResponse)

}
