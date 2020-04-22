package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/olivere/elastic"
	"github.com/softplan/tenkai-docker-api/pkg/audit"
	"github.com/softplan/tenkai-docker-api/pkg/configs"
	"github.com/softplan/tenkai-docker-api/pkg/dbms"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/repository"
	"github.com/softplan/tenkai-docker-api/pkg/global"
	dockerapi "github.com/softplan/tenkai-docker-api/pkg/service/docker"
)

//Repositories  Repositories
type Repositories struct {
	DockerDAO repository.DockerDAOInterface
}

//AppContext AppContext
type AppContext struct {
	DockerServiceAPI dockerapi.DockerServiceInterface
	Auditing         audit.AuditingInterface
	Configuration    *configs.Configuration
	Repositories     Repositories
	Database         dbms.Database
	Elk              *elastic.Client
	Mutex            sync.Mutex
	DockerTagsCache  sync.Map
}

func defineRotes(r *mux.Router, appContext *AppContext) {
	r.HandleFunc("/dockerRepo", appContext.listDockerRepositories).Methods("GET")
	r.HandleFunc("/dockerRepo", appContext.newDockerRepository).Methods("POST")
	r.HandleFunc("/dockerRepo/{id}", appContext.deleteDockerRepository).Methods("DELETE")
	r.HandleFunc("/listDockerTags", appContext.listDockerTags).Methods("POST")
	r.HandleFunc("/listDockerVariables", appContext.listDockerVariables).Methods("POST")
	r.HandleFunc("/", appContext.rootHandler)
}

//StartHTTPServer StartHTTPServer
func StartHTTPServer(appContext *AppContext) {

	port := appContext.Configuration.Server.Port
	global.Logger.Info(global.AppFields{global.Function: "startHTTPServer", "port": port}, "online - listen and server")

	r := mux.NewRouter()

	defineRotes(r, appContext)

	log.Fatal(http.ListenAndServe(":"+port, commonHandler(r)))

}

func extractToken(reqToken string) *model.Principal {
	var principal model.Principal

	splitToken := strings.Split(reqToken, "Bearer ")
	reqToken = splitToken[1]
	token, _, earl := new(jwt.Parser).ParseUnverified(reqToken, jwt.MapClaims{})
	if earl == nil {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			in := claims["realm_access"]
			if in != nil {
				realmAccessMap := in.(map[string]interface{})
				roles := realmAccessMap["roles"]
				if roles != nil {
					elements := roles.([]interface{})
					for _, element := range elements {
						principal.Roles = append(principal.Roles, element.(string))
					}
				}
			}
			principal.Name = fmt.Sprintf("%v", claims["name"])
			principal.Email = fmt.Sprintf("%v", claims["email"])
			return &principal
		}
	}
	return nil
}

func commonHandler(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		reqToken := r.Header.Get("Authorization")
		if len(reqToken) > 0 {
			principal := extractToken(reqToken)
			if principal != nil {
				data, _ := json.Marshal(*principal)
				r.Header.Set("principal", string(data))
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (appContext *AppContext) rootHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"service": "TENKAI",
		"status":  "ready",
	}

	json, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set(global.ContentType, "application/json")
	w.Write(json)
}
