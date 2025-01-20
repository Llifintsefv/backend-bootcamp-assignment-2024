package flat

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
)



func Test_CreateFlat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db,mock,err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()

		repo := NewFlatRepository(db)

		
		request := dto.PostFlatCreateJSONRequestBody{
			Price: 100,
			Rooms: new(dto.Rooms),
			HouseId: 1,
		}
		*request.Rooms = 2
		
		mock.ExpectPrepare("INSERT INTO flats").ExpectQuery().
		WithArgs(request.HouseId, request.Price, request.Rooms, "created").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		createdFlat, err := repo.CreateFlat(context.Background(), request)
		
		assert.NoError(t, err)
		assert.Equal(t, 1, createdFlat.Id)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("error", func(t *testing.T) {
		db,mock,err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()

		repo := NewFlatRepository(db)

		
		request := dto.PostFlatCreateJSONRequestBody{
			Price: 100,
			Rooms: new(dto.Rooms),
			HouseId: 1,
		}
		*request.Rooms = 2
		
		mock.ExpectPrepare("INSERT INTO flats").ExpectQuery().
		WithArgs(request.HouseId, request.Price, request.Rooms, "created").
		WillReturnError(err)

		_, err = repo.CreateFlat(context.Background(), request)
		
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("ForeignKeyViolation", func(t *testing.T) {
		db,mock,err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()

		repo := NewFlatRepository(db)

		
		request := dto.PostFlatCreateJSONRequestBody{
			Price: 100,
			Rooms: new(dto.Rooms),
			HouseId: 1,
		}
		*request.Rooms = 2	
		
		mock.ExpectPrepare("INSERT INTO flats").ExpectQuery().
		WithArgs(request.HouseId, request.Price, request.Rooms, "created").
		WillReturnError(&pgconn.PgError{Code: "23503"})

		_, err = repo.CreateFlat(context.Background(), request)
		
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "house with id 1 does not exist")
		assert.NoError(t, mock.ExpectationsWereMet())

	})
}