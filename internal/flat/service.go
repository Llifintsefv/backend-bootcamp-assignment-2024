package flat

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
)

type FlatService interface {
	CreateFlat(ctx context.Context,req dto.PostFlatCreateJSONRequestBody) (dto.Flat,error) 
}

type flatService struct {
	flatRepository FlatRepository
}

func NewFlatService(flatRepository FlatRepository) FlatService {
	return &flatService{flatRepository: flatRepository}
}

func (s *flatService) CreateFlat(ctx context.Context,req dto.PostFlatCreateJSONRequestBody) (dto.Flat,error)  {
		flat,err := s.flatRepository.CreateFlat(ctx,req)
		if err != nil {

		}

		return flat,nil
}

