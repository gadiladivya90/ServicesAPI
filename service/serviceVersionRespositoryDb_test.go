package entity

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func Test_should_test_insert_into_ServiceVersion_failure(t *testing.T) {
	mockDB, mocksql, _ := sqlmock.New()
	defer mockDB.Close()

	//create a db client for serviceVersion repository
	client := sqlx.NewDb(mockDB, "sqlmock")
	svRepo := ServiceVersionRepositoryDb{client}

	// setting expectations and mock objects
	mocksql.ExpectBegin()

	// mock  instance
	sv := ServiceVersionRequestDto{"name", "service_id"}

	mocksql.ExpectExec(`INSERT INTO service_versions`).
		WithArgs(sv.name, sv.service_id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mocksql.ExpectExec(`INSERT INTO service_versions`).
		WithArgs(sv.name, sv.service_id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mocksql.ExpectExec(`UPDATE service_packages`).
		WithArgs(sv.name, sv.service_id).
		WillReturnError(fmt.Errorf("some error"))

	mocksql.ExpectRollback()

	// now we execute our method
	if _, appError := svRepo.SaveServiceVersion(sv); appError == nil {
		t.Errorf("was expecting an error, but there was none")
	}

	// we make sure that all expectations were met
	if err := mocksql.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
