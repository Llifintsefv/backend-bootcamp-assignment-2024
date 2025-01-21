package flat

import (
	"backend-bootcamp-assignment-2024/dto"
	"context"
	"fmt"
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


func Test_GetFlat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db,mock,err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}

		defer db.Close()

		repo := NewFlatRepository(db)

		houseId := "1"
		expectedFlats := []dto.Flat{
		{Id: 1, HouseId: 1, Price: 100000, Rooms: 2, Status: "created"},
		{Id: 2, HouseId: 1, Price: 150000, Rooms: 3, Status: "booked"},
		}

	
		rows := sqlmock.NewRows([]string{"id", "house_id", "price", "rooms", "status"}).
		AddRow(expectedFlats[0].Id, expectedFlats[0].HouseId, expectedFlats[0].Price, expectedFlats[0].Rooms, expectedFlats[0].Status).
		AddRow(expectedFlats[1].Id, expectedFlats[1].HouseId, expectedFlats[1].Price, expectedFlats[1].Rooms, expectedFlats[1].Status)

		mock.ExpectPrepare("SELECT id,house_id,price,rooms,status FROM flats WHERE house_id = \\$1").
		ExpectQuery().
		WithArgs(houseId).
		WillReturnRows(rows)

		actualFlats, err := repo.GetFlats(context.Background(), houseId)
		
		assert.NoError(t, err)
		assert.Equal(t, expectedFlats, actualFlats)
		assert.NoError(t, mock.ExpectationsWereMet())

	})

	t.Run("DBerror", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		repo := NewFlatRepository(db)

		houseId := "1"


		mock.ExpectPrepare("SELECT id,house_id,price,rooms,status FROM flats WHERE house_id = \\$1").
			ExpectQuery().
			WithArgs(houseId).
			WillReturnError(fmt.Errorf("some database error"))

		_, err = repo.GetFlats(context.Background(), houseId)

		assert.Error(t, err)
		assert.EqualError(t, err, "error executing query: some database error")

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}

	})
	
}