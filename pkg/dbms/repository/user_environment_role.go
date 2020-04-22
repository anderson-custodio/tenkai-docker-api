package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"
)

//UserEnvironmentRoleDAOInterface - UserEnvironmentRoleDAOInterface
type UserEnvironmentRoleDAOInterface interface {
	CreateOrUpdate(so model.UserEnvironmentRole) error
	GetRoleByUserAndEnvironment(user model.User, envID uint) (*model.SecurityOperation, error)
}

//UserEnvironmentRoleDAOImpl UserEnvironmentRoleDAOImpl
type UserEnvironmentRoleDAOImpl struct {
	Db *gorm.DB
}

//CreateOrUpdate - Create or update a security operation
func (dao UserEnvironmentRoleDAOImpl) CreateOrUpdate(so model.UserEnvironmentRole) error {
	loadSO, err := dao.isEdit(so)
	if err != nil {
		return err
	}
	if loadSO != nil {
		return dao.edit(so, loadSO)
	}
	return dao.create(so)
}

func (dao UserEnvironmentRoleDAOImpl) isEdit(so model.UserEnvironmentRole) (*model.UserEnvironmentRole, error) {
	var loadSO model.UserEnvironmentRole
	if err := dao.Db.Where(model.UserEnvironmentRole{UserID: so.UserID, EnvironmentID: so.EnvironmentID}).First(&loadSO).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
		return nil, nil
	}
	return &loadSO, nil
}

func (dao UserEnvironmentRoleDAOImpl) edit(so model.UserEnvironmentRole, loadSo *model.UserEnvironmentRole) error {
	loadSo.SecurityOperationID = so.SecurityOperationID
	if err := dao.Db.Save(&so).Error; err != nil {
		return err
	}
	return nil
}

func (dao UserEnvironmentRoleDAOImpl) create(so model.UserEnvironmentRole) error {
	if err := dao.Db.Create(&so).Error; err != nil {
		return err
	}
	return nil
}

//GetRoleByUserAndEnvironment - GetRoleByUserAndEnvironment
func (dao UserEnvironmentRoleDAOImpl) GetRoleByUserAndEnvironment(user model.User,
	envID uint) (*model.SecurityOperation, error) {
	var userEnvironmentRole model.UserEnvironmentRole
	if err := dao.Db.Where(model.UserEnvironmentRole{UserID: user.ID, EnvironmentID: envID}).Find(&userEnvironmentRole).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
		return nil, nil
	}
	var result model.SecurityOperation
	if err := dao.Db.First(&result, userEnvironmentRole.SecurityOperationID).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			return nil, err
		}
		return nil, nil
	}
	return &result, nil
}
