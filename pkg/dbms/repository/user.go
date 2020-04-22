package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"
)

//UserDAOInterface UserDAOInterface
type UserDAOInterface interface {
	ListAllUsers() ([]model.User, error)
	FindByEmail(email string) (model.User, error)
}

//UserDAOImpl UserDAOImpl
type UserDAOImpl struct {
	Db *gorm.DB
}

//ListAllUsers - List all users
func (dao UserDAOImpl) ListAllUsers() ([]model.User, error) {
	users := make([]model.User, 0)
	if err := dao.Db.Preload("Environments").Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

//FindByEmail FindByEmail
func (dao UserDAOImpl) FindByEmail(email string) (model.User, error) {
	var user model.User
	if err := dao.Db.Where(&model.User{Email: email}).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
