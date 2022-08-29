package service

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/divyag/services/entity"
	"github.com/divyag/services/errs"
	"github.com/divyag/services/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ServiceRepositoryDb struct {
	client *sqlx.DB
}

func NewServiceRepositoryDb(dbClient *sqlx.DB) ServiceRepositoryDb {

	return ServiceRepositoryDb{client: dbClient}
}

func (cdb ServiceRepositoryDb) FindAll(filters entity.FilterParams) (*entity.PaginationResponseDto, *errs.AppErr) {
	services := make([]entity.ServicePackage, 0)

	//sql
	filterQuery := getFilterSql(filters)
	findAllSql := fmt.Sprintf("select * from service_packages %s", filterQuery)
	err := cdb.client.Select(&services, findAllSql)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NotFoundError("Service Not Found")
		} else {
			return nil, errs.InternalServerError(fmt.Sprintf("Internal Server Error :%+v", err.Error()))
		}
	}

	//getcount for pagination
	findAllCountSql := fmt.Sprintf("select count(*) from service_packages %s", filterQuery)
	row := cdb.client.QueryRow(findAllCountSql)
	var total int
	err = row.Scan(&total)
	if err != nil {
		total = len(services)
	}

	d := make([]entity.ServiceResponseDto, 0)

	for _, c := range services {
		r := c.ToDto()
		d = append(d, r)
	}

	return entity.Paginate(d, total, filters), nil
}

func (cdb ServiceRepositoryDb) FindServiceByID(id string) (*entity.ServicePackage, *errs.AppErr) {
	var service entity.ServicePackage
	findServiceSql := "select * from service_packages where service_id = $1"

	err := cdb.client.Get(&service, findServiceSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NotFoundError("Service Not Found")
		} else {
			return nil, errs.InternalServerError(fmt.Sprintf("Internal Server Error: %v", err.Error()))
		}
	}

	return &service, nil
}

func (cdb ServiceRepositoryDb) SaveService(s *entity.ServicePackage) (*entity.ServicePackage, *errs.AppErr) {

	sqlInsert := fmt.Sprintf("INSERT INTO service_packages (name, description, created_at, updated_at) values ('%s', '%s', '%s', '%s') RETURNING service_id", s.Name, s.Description, s.CreatedAt, s.UpdatedAt)

	var id string
	err := cdb.client.QueryRow(sqlInsert).Scan(&id)
	if err != nil {
		return nil, errs.InternalServerError(fmt.Sprintf("Error while inserting new service err:%+v\n", err.Error()))
	}

	s.Id = id
	return s, nil
}

func (cdb ServiceRepositoryDb) UpdateService(s *entity.ServicePackage) (*entity.ServicePackage, *errs.AppErr) {

	s.UpdatedAt = time.Now().Format(time.RFC3339)
	sqlUpdate := "UPDATE service_packages SET name=:name, description=:description, created_at=:created_at, updated_at=:updated_at where service_id=:service_id"
	_, err := cdb.client.NamedQuery(sqlUpdate, s)
	if err != nil {
		return nil, errs.InternalServerError(fmt.Sprintf("Error while updating existing service err:%+v\n", err.Error()))
	}

	return s, nil
}

func (cdb ServiceRepositoryDb) DeleteService(service_id string) *errs.AppErr {
	sqlDelete := "delete from service_packages where service_id=$1"

	_, err := cdb.client.Exec(sqlDelete, service_id)
	if err != nil {
		return errs.InternalServerError(fmt.Sprintf("Error while deleting service err:%+v\n", err.Error()))
	}

	return nil
}

func getFilterSql(filters entity.FilterParams) string {
	var filterQuery string

	if filters.Filter != "" {
		//TODO: This can be better. Need to understand usage better
		s := strings.Split(filters.Filter, " ")
		if len(s) != 3 {
			logger.Warn("Invalid fitler")
		} else {
			filterQuery = fmt.Sprintf("where %s like '%%%s%%' ", s[0], s[2])
		}
	}

	filterQuery = fmt.Sprintf("%s order by %s offset %d", filterQuery, filters.Sort, filters.Offset)
	if filters.Limit > 0 {
		filterQuery = fmt.Sprintf("%s limit %d", filterQuery, filters.Limit)
	}

	return filterQuery
}
