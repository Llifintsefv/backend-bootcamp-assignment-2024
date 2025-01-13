package flat

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
	"fmt"
)

type FlatService interface {
	CreateFlat(ctx context.Context, req dto.PostFlatCreateJSONRequestBody) (dto.Flat, error)
	GetFlats(ctx context.Context, houseId string) ([]dto.Flat, error)
	UpdateFlatStatus(ctx context.Context, req dto.PostFlatUpdateJSONRequestBody) (dto.Flat, error)
}

type flatService struct {
	flatRepository FlatRepository
}

func NewFlatService(flatRepository FlatRepository) FlatService {
	return &flatService{flatRepository: flatRepository}
}

func (s *flatService) CreateFlat(ctx context.Context, req dto.PostFlatCreateJSONRequestBody) (dto.Flat, error) {
	flat, err := s.flatRepository.CreateFlat(ctx, req)
	if err != nil {
		fmt.Errorf("error creating flat: %w", err)
		return dto.Flat{}, nil
	}

	return flat, nil
}

func (s *flatService) GetFlats(ctx context.Context, houseId string) ([]dto.Flat, error) {
	flats, err := s.flatRepository.GetFlats(ctx, houseId)
	if err != nil {
		fmt.Errorf("error getting flats: %w", err)
		return []dto.Flat{}, err
	}

	return flats, nil
}


func (s *flatService) UpdateFlatStatus(ctx context.Context, req dto.PostFlatUpdateJSONRequestBody) (dto.Flat, error) {
	flat, err := s.flatRepository.UpdateFlatStatus(ctx, req)
	if err != nil {
		fmt.Errorf("error updating flat status: %w", err)
		return dto.Flat{}, nil
	}

	return flat, nil
}