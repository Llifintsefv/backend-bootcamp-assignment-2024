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