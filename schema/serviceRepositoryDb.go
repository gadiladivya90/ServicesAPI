package schema

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/divyag/services/dto"
	"github.com/divyag/services/errs"
	"github.com/divyag/services/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type ServiceRepositoryDb struct {
	client *sqlx.DB
}

func NewServiceRepositoryDb() ServiceRepositoryDb {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	dbClient, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic("unable to establish connection")
	}

	dbClient.SetConnMaxLifetime(time.Minute * 3)
	dbClient.SetMaxOpenConns(10)
	dbClient.SetMaxIdleConns(10)

	return ServiceRepositoryDb{client: dbClient}
}

func (cdb ServiceRepositoryDb) FindAll(filters dto.FilterParams) (*dto.PaginationResponseDto, *errs.AppErr) {
	services := make([]ServicePackage, 0)

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

	findAllCountSql := fmt.Sprintf("select count(*) from service_packages %s", filterQuery)
	row := cdb.client.QueryRow(findAllCountSql)
	var total int
	err = row.Scan(&total)
	if err != nil {
		total = len(services)
	}

	d := make([]dto.ServiceResponseDto, 0)

	for _, c := range services {
		r := c.ToDto()
		d = append(d, r)
	}

	return dto.Paginate(d, total, filters), nil
}

func (cdb ServiceRepositoryDb) FindServiceByID(id string) (*ServicePackage, *errs.AppErr) {
	var service ServicePackage
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

func (cdb ServiceRepositoryDb) SaveService(s *ServicePackage) (*ServicePackage, *errs.AppErr) {

	sqlInsert := fmt.Sprintf("INSERT INTO service_packages (name, description, created_at, updated_at) values ('%s', '%s', %s, %s) RETURNING service_id", s.Name, s.Description, s.CreatedAt, s.UpdatedAt)

	var id int64
	err := cdb.client.QueryRow(sqlInsert).Scan(&id)
	if err != nil {
		return nil, errs.InternalServerError(fmt.Sprintf("Error while inserting new service err:%+v\n", err.Error()))
	}

	s.Id = strconv.FormatInt(id, 10)
	return s, nil
}

func (cdb ServiceRepositoryDb) UpdateService(service_id string, s *ServicePackage) (*ServicePackage, *errs.AppErr) {

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

func getFilterSql(filters dto.FilterParams) string {
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
