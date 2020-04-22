//+build !test

package dbms

import (
	"github.com/jinzhu/gorm"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"

	//postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
	//sqllite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//Database Structure
type Database struct {
	Db *gorm.DB
}

//Connect - Connect to a database
func (database *Database) Connect(dbmsURI string, local bool) {
	var err error

	if local {
		database.Db, err = gorm.Open("sqlite3", "/tmp/tekai.db")
	} else {
		database.Db, err = gorm.Open("postgres", dbmsURI)
	}

	if err != nil {
		panic("failed to connect database")
	}

	database.Db.AutoMigrate(&model.Environment{})
	database.Db.AutoMigrate(&model.Variable{})
	database.Db.AutoMigrate(&model.User{})
	database.Db.AutoMigrate(&model.DockerRepo{})
	database.Db.AutoMigrate(&model.SecurityOperation{})
	database.Db.AutoMigrate(&model.UserEnvironmentRole{})
}
