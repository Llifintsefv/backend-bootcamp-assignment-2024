package flat

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)


type FlatRepository interface {
	CreateFlat(ctx context.Context,req dto.PostFlatCreateJSONRequestBody) (dto.Flat,error) 
	GetFlats(ctx context.Context, houseId string) ([]dto.Flat,error)
}

type flatRepository struct{
	db *sql.DB
}

func NewFlatRepository(db *sql.DB) FlatRepository{
	return &flatRepository{db: db}
}

func (r *flatRepository) CreateFlat(ctx context.Context,req dto.PostFlatCreateJSONRequestBody) (dto.Flat,error) {
	var flat dto.Flat
	stmt,err := r.db.PrepareContext(ctx,"INSERT INTO flats(house_id,price,rooms,status) VALUES($1,$2,$3,$4) RETURNING id")
	if err != nil {

	}

	defer stmt.Close()

	err = stmt.QueryRowContext(ctx,req.HouseId,req.Price,req.Rooms,"created").Scan(&flat.Id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" { // foreign_key_violation
				return dto.Flat{}, fmt.Errorf("house with id %d does not exist", req.HouseId)
			}
		}
		return dto.Flat{}, fmt.Errorf("error executing query: %w", err)
	}

	flat = dto.Flat{
		Id: flat.Id,
		HouseId: req.HouseId,
		Price: req.Price,
		Rooms: *req.Rooms,
		Status: "created",
	}

	return flat, nil

}

func (r *flatRepository) GetFlats(ctx context.Context, houseId string) ([]dto.Flat,error) {
	var flats []dto.Flat
	stmt,err := r.db.PrepareContext(ctx,"SELECT id,house_id,price,rooms,status FROM flats WHERE house_id = $1")
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %w", err)
	}

	defer stmt.Close()

	rows,err := stmt.QueryContext(ctx,houseId)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	for rows.Next() {
		var flat dto.Flat
		err = rows.Scan(&flat.Id,&flat.HouseId,&flat.Price,&flat.Rooms,&flat.Status)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		flats = append(flats, flat)
	}

	return flats, nil
}