package flat

import (
	"backend-bootcamp-assignment-2024/dto"
	"encoding/json"
	"net/http"
)

type FlatHandler struct {
	flatService FlatService
}

func NewFlatHandler(flatService FlatService) *FlatHandler {
	return &FlatHandler{flatService: flatService}
}

func (h *FlatHandler) CreateFlat(w http.ResponseWriter, r *http.Request){
	var reg dto.PostFlatCreateJSONRequestBody

	err := json.NewDecoder(r.Body).Decode(&reg)
	if err != nil {

	}

	if reg.Price <= 0 || *reg.Rooms <= 0 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	flatCreateResponse,err := h.flatService.CreateFlat(ctx,reg)
	if err != nil {

	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(flatCreateResponse)
}