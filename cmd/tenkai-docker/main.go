package main

import (
	"log"
	"sync"

	"github.com/softplan/tenkai-docker-api/pkg/audit"
	"github.com/softplan/tenkai-docker-api/pkg/configs"
	"github.com/softplan/tenkai-docker-api/pkg/dbms"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/repository"
	"github.com/softplan/tenkai-docker-api/pkg/global"
	"github.com/softplan/tenkai-docker-api/pkg/handlers"
	dockerapi "github.com/softplan/tenkai-docker-api/pkg/service/docker"
)

const (
	configFileName = "tenkai-docker-api"
)

func main() {
	logFields := global.AppFields{global.Function: "main"}

	global.Logger.Info(logFields, "loading config properties")

	config, error := configs.ReadConfig(configFileName)
	checkFatalError(error)

	appContext := &handlers.AppContext{Configuration: config}

	dbmsURI := config.App.Dbms.URI

	initCache(appContext)
	initAPIs(appContext)

	//Dbms connection
	appContext.Database.Connect(dbmsURI, dbmsURI == "")
	defer appContext.Database.Db.Close()

	appContext.Repositories = initRepository(&appContext.Database)

	//Elk setup
	appContext.Elk, _ = appContext.Auditing.ElkClient(config.App.Elastic.URL, config.App.Elastic.Username, config.App.Elastic.Password)

	global.Logger.Info(logFields, "http server started")
	handlers.StartHTTPServer(appContext)
}

func initCache(appContext *handlers.AppContext) {
	appContext.DockerTagsCache = sync.Map{}
}

func initAPIs(appContext *handlers.AppContext) {
	appContext.DockerServiceAPI = dockerapi.DockerServiceBuilder()
	appContext.Auditing = audit.AuditingBuilder()
}

func initRepository(database *dbms.Database) handlers.Repositories {
	repositories := handlers.Repositories{}
	repositories.DockerDAO = &repository.DockerDAOImpl{Db: database.Db}
	repositories.EnvironmentDAO = &repository.EnvironmentDAOImpl{Db: database.Db}
	repositories.UserDAO = &repository.UserDAOImpl{Db: database.Db}
	repositories.VariableDAO = &repository.VariableDAOImpl{Db: database.Db}
	repositories.SecurityOperationDAO = &repository.SecurityOperationDAOImpl{Db: database.Db}
	repositories.UserEnvironmentRoleDAO = &repository.UserEnvironmentRoleDAOImpl{Db: database.Db}

	return repositories
}

func checkFatalError(err error) {
	if err != nil {
		global.Logger.Error(global.AppFields{global.Function: "upload", "error": err}, "erro fatal")
		log.Fatal(err)
	}
}
