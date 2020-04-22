package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"
	"github.com/stretchr/testify/assert"
)

func getUser() model.User {
	item := model.User{}
	item.Email = "musk@mars.com"
	item.DefaultEnvironmentID = 999
	return item
}

func TestListAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	mock.MatchExpectationsInOrder(false)
	assert.Nil(t, err)

	gormDB, err := gorm.Open("postgres", db)
	defer gormDB.Close()

	userDAO := UserDAOImpl{}
	userDAO.Db = gormDB

	payload := getUser()
	payload.ID = 888

	env := getEnvironmentTestData()
	env.ID = 999

	payload.Environments = append(payload.Environments, env)

	row1 := sqlmock.NewRows([]string{"id", "email", "default_environment_id"}).AddRow(payload.ID, payload.Email, payload.DefaultEnvironmentID)
	mock.ExpectQuery(`SELECT (.*) FROM "users"
		WHERE "users"."deleted_at" IS NULL
	`).WillReturnRows(row1)

	row2 := sqlmock.NewRows([]string{"id"}).AddRow(env.ID)
	mock.ExpectQuery(`SELECT (.*) FROM "environments" 
		INNER JOIN "user_environment" ON "user_environment"."environment_id" = "environments"."id" 
		WHERE "environments"."deleted_at" IS NULL AND \(\("user_environment"."user_id" IN (.*)\)\)
	`).WithArgs(888).WillReturnRows(row2)

	u, e := userDAO.ListAllUsers()
	assert.NoError(t, e)
	assert.NotNil(t, u)

	mock.ExpectationsWereMet()
}

func TestFindByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	mock.MatchExpectationsInOrder(false)
	assert.Nil(t, err)

	gormDB, err := gorm.Open("postgres", db)
	defer gormDB.Close()

	userDAO := UserDAOImpl{}
	userDAO.Db = gormDB

	payload := getUser()
	payload.ID = 888

	row1 := sqlmock.NewRows([]string{"id", "email", "default_environment_id"}).
		AddRow(payload.ID, payload.Email, payload.DefaultEnvironmentID)

	mock.ExpectQuery(`SELECT .+ FROM "users" WHERE "users"."deleted_at" IS NULL AND \(\("users"."email" =`).
		WithArgs("musk@mars.com").
		WillReturnRows(row1)

	u, e := userDAO.FindByEmail("musk@mars.com")
	assert.NoError(t, e)
	assert.NotNil(t, u)

	mock.ExpectationsWereMet()
}
