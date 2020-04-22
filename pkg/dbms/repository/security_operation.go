package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/softplan/tenkai-docker-api/pkg/dbms/model"
)

//SecurityOperationDAOInterface - SecurityOperationDAOInterface
type SecurityOperationDAOInterface interface {
	List() ([]model.SecurityOperation, error)
}

//SecurityOperationDAOImpl SecurityOperationDAOImpl
type SecurityOperationDAOImpl struct {
	Db *gorm.DB
}

//List - List
func (dao SecurityOperationDAOImpl) List() ([]model.SecurityOperation, error) {
	oss := make([]model.SecurityOperation, 0)
	if err := dao.Db.Find(&oss).Error; err != nil {
		return nil, err
	}
	return oss, nil
}
