package data

import (
	"context"
	"testing"

	"doe/src/mocks"
	"doe/src/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
)

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	logMock := mocks.NewMockLogger(ctrl)
	repo := &repoImpl{
		log: logMock,
		tx:  &sqlImpl{db},
	}
	sqlRows := sqlmock.NewRows([]string{
		"port_id",
		"name",
		"city",
		"province",
		"country",
		"regions",
		"coordinates",
		"timezone",
		"unlocs",
		"code",
		"alias",
	}).AddRow(
		"a", "b", "c", "d", "u", nil, nil, "k", nil, "g", nil,
	)
	mock.ExpectQuery("SELECT").WillReturnRows(sqlRows)
	_, err = repo.GetByID(ctx, "foobar")
	if err != nil {
		t.Fatal(err)
	}
}
func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	logMock := mocks.NewMockLogger(ctrl)
	repo := &repoImpl{
		log: logMock,
		tx:  &sqlImpl{db},
	}
	sqlRows := sqlmock.NewRows([]string{
		"port_id",
		"name",
		"city",
		"province",
		"country",
		"regions",
		"coordinates",
		"timezone",
		"unlocs",
		"code",
		"alias",
	}).AddRow(
		"a", "b", "c", "d", "u", nil, nil, "k", nil, "g", nil,
	)
	mock.ExpectQuery("SELECT").WillReturnRows(sqlRows)
	_, err = repo.GetAll(ctx, 100)
	if err != nil {
		t.Fatal(err)
	}
}

func TestInsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.TODO()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}

	logMock := mocks.NewMockLogger(ctrl)
	repo := &repoImpl{
		log: logMock,
		tx:  &sqlImpl{db},
	}
	logMock.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO").
		WithArgs("", "", "", "", nil, nil, "", "", nil, "").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("DELETE FROM").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = repo.Insert(ctx, []*models.Port{{}})
	if err != nil {
		t.Fatal(err)
	}
}
