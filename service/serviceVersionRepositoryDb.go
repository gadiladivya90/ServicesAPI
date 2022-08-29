package service

import (
	"database/sql"
	"fmt"

	"github.com/divyag/services/entity"
	"github.com/divyag/services/errs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ServiceVersionRepositoryDb struct {
	client *sqlx.DB
}

func NewServiceVersionRepositoryDb(dbClient *sqlx.DB) ServiceVersionRepositoryDb {

	return ServiceVersionRepositoryDb{client: dbClient}
}

func (cdb ServiceVersionRepositoryDb) FindServiceVersionsByID(service_id string) ([]entity.ServiceVersionResponseDto, *errs.AppErr) {
	services := make([]entity.ServiceVersion, 0)

	findAllSql := fmt.Sprintf("select * from service_versions where service_package_id=%s", service_id)
	err := cdb.client.Select(&services, findAllSql)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.ServiceVersionResponseDto{}, nil
		} else {
			return nil, errs.InternalServerError(fmt.Sprintf("Internal Server Error :%+v", err.Error()))
		}
	}

	d := make([]entity.ServiceVersionResponseDto, 0)

	for _, c := range services {
		r := c.ToSVDto()
		d = append(d, r)
	}

	return d, nil
}

func (cdb ServiceVersionRepositoryDb) SaveServiceVersion(s *entity.ServiceVersion) (*entity.ServiceVersion, *errs.AppErr) {

	sqlInsert := fmt.Sprintf("INSERT INTO service_versions where (name, service_package_id ) values ('%s', '%s') RETURNING service_version_id", s.Name, s.ServiceId)

	var id string
	err := cdb.client.QueryRow(sqlInsert).Scan(&id)
	if err != nil {
		return nil, errs.InternalServerError(fmt.Sprintf("Error while inserting new service err:%+v\n", err.Error()))
	}

	s.Id = id
	return s, nil
}

func (cdb ServiceVersionRepositoryDb) DeleteServiceVersion(service_id string) *errs.AppErr {
	sqlDelete := "delete from service_versions where service_version_id=$1"

	_, err := cdb.client.Exec(sqlDelete, service_id)
	if err != nil {
		return errs.InternalServerError(fmt.Sprintf("Error while deleting service  version err:%+v\n", err.Error()))
	}

	return nil
}
