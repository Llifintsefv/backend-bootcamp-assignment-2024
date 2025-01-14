package house

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
	"database/sql"
)

type HouseRepository interface {
	CreateHouse(ctx context.Context, req dto.PostHouseCreateJSONRequestBody) (dto.House, error)
	Subscribe(ctx context.Context, houseId int, req dto.Email) error
	GetEmailsByHouseID(ctx context.Context, houseId int) ([]dto.Email, error)
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

func (r *houseRepository) Subscribe(ctx context.Context, houseId int, req dto.Email) error {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO subscriptions (house_id,email) VALUES ($1,$2)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, houseId, req)
	if err != nil {
		return err
	}

	return nil
}

func (r *houseRepository) GetEmailsByHouseID(ctx context.Context, houseId int) ([]dto.Email, error) {
	var emails []dto.Email
	stmt, err := r.db.PrepareContext(ctx, "SELECT email FROM subscriptions WHERE house_id = $1")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, houseId)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var email dto.Email
		err = rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}

	return emails, nil
}
