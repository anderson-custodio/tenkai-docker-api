package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"
)

//EnvironmentDAOInterface EnvironmentDAOInterface
type EnvironmentDAOInterface interface {
	GetAllEnvironments(principal string) ([]model.Environment, error)
	GetByID(envID int) (*model.Environment, error)
}

//EnvironmentDAOImpl EnvironmentDAOImpl
type EnvironmentDAOImpl struct {
	Db *gorm.DB
}

//GetAllEnvironments - Retrieve all environments
func (dao EnvironmentDAOImpl) GetAllEnvironments(principal string) ([]model.Environment, error) {
	envs := make([]model.Environment, 0)
	if len(principal) > 0 {
		var user model.User
		if err := dao.Db.Where(model.User{Email: principal}).First(&user).Error; err != nil {
			return checkNotFound(err)
		}

		if err := dao.Db.Model(&user).Related(&envs, "Environments").Error; err != nil {
			return checkNotFound(err)
		}
	} else {
		if err := dao.Db.Find(&envs).Error; err != nil {
			return checkNotFound(err)
		}
	}
	return envs, nil
}

//GetByID - Get Environment By Id
func (dao EnvironmentDAOImpl) GetByID(envID int) (*model.Environment, error) {
	var result model.Environment
	if err := dao.Db.First(&result, envID).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func checkNotFound(err error) ([]model.Environment, error) {
	if err == gorm.ErrRecordNotFound {
		return make([]model.Environment, 0), nil
	}
	return nil, err

}
