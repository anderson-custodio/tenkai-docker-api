package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"
)

//mockPrincipal injects a http header with the specified role to be used only for testing.
func mockPrincipal(req *http.Request) {
	var roles []string
	roles = append(roles, "tenkai-admin")
	principal := model.Principal{Name: "alfa", Email: "beta@alfa.com", Roles: roles}
	pSe, _ := json.Marshal(principal)
	req.Header.Set("principal", string(pSe))
}
