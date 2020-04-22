package model

import (
	"github.com/jinzhu/gorm"
)

//Environment - Environment Model
type Environment struct {
	gorm.Model
	Group          string `json:"group"`
	Name           string `json:"name"`
	ClusterURI     string `json:"cluster_uri"`
	CACertificate  string `json:"ca_certificate"`
	Token          string `json:"token"`
	Namespace      string `json:"namespace"`
	Gateway        string `json:"gateway"`
	ProductVersion string `json:"productVersion"`
	CurrentRelease string `json:"currentRelease"`
	EnvType        string `json:"envType"`
	Host           string `json:"host"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

//User struct
type User struct {
	gorm.Model
	Email                string        `json:"email"`
	DefaultEnvironmentID int           `json:"defaultEnvironmentID"`
	Environments         []Environment `gorm:"many2many:user_environment;"`
}

//UserEnvironmentRole UserEnvironmentRole
type UserEnvironmentRole struct {
	gorm.Model
	UserID              uint `json:"userId"`
	EnvironmentID       uint `json:"environmentId"`
	SecurityOperationID uint `json:"securityOperationId"`
}

//Variable Structure
type Variable struct {
	gorm.Model
	Scope         string `json:"scope" gorm:"index:var_scope"`
	ChartVersion  string `gorm:"-" json:"chartVersion"`
	Name          string `json:"name" gorm:"index:var_name"`
	Value         string `json:"value"`
	Secret        bool   `json:"secret"`
	Description   string `json:"description"`
	EnvironmentID int    `json:"environmentId"`
}
