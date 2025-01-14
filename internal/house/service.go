package house

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
)

type HouseService interface {
	CreateHouse(ctx context.Context, req dto.PostHouseCreateJSONRequestBody) (dto.House, error)
	Subscribe(ctx context.Context, houseId int, req dto.Email) error
}

type houseService struct {
	houseRepository HouseRepository
}

func NewHouseService(houseRepository HouseRepository) HouseService {
	return &houseService{
		houseRepository: houseRepository,
	}
}

func (s *houseService) CreateHouse(ctx context.Context, req dto.PostHouseCreateJSONRequestBody) (dto.House, error) {
	house, err := s.houseRepository.CreateHouse(ctx, req)
	if err != nil {
		return dto.House{}, err
	}

	return house, nil
}

func (s *houseService) Subscribe(ctx context.Context, houseId int, req dto.Email) error {
	err := s.houseRepository.Subscribe(ctx, houseId, req)
	if err != nil {
		return err
	}
	return nil
}
