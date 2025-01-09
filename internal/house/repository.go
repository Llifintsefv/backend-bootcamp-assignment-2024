package house

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
	"database/sql"
)

type HouseRepository interface {
	CreateHouse(ctx context.Context, req dto.PostHouseCreateJSONRequestBody) (dto.House, error)
}

type houseRepository struct {
	db *sql.DB
}

func NewHouseRepository(db *sql.DB) HouseRepository {
	return &houseRepository{
		db: db,
	}
}

func (r *houseRepository) CreateHouse(ctx context.Context, req dto.PostHouseCreateJSONRequestBody) (dto.House, error) {
	var house dto.House
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO houses (address,year_built,developer) VALUES ($1,$2,$3) RETURNING id,created_at,update_at")
	if err != nil {
		return house, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, req.Address, req.Year, req.Developer).Scan(&house.Id, &house.CreatedAt, &house.UpdateAt)
	if err != nil {
		return house, err
	}

	house = dto.House{
		Id:        house.Id,
		CreatedAt: house.CreatedAt,
		UpdateAt:  house.UpdateAt,
		Developer: req.Developer,
		Address:   req.Address,
		Year:      req.Year,
	}

	return house, nil
}

