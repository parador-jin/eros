package dao

import (
	"Eros/models"
	"github.com/go-xorm/xorm"
	"log"
)

type UserDao struct {
	engine *xorm.Engine
}

func NewUserDao(engine *xorm.Engine) *UserDao {
	return &UserDao{
		engine: engine,
	}
}

func (d *UserDao) Get(id int) *models.LtUser {
	data := &models.LtUser{Id: id}
	ok, err := d.engine.Get(data)
	if err != nil || !ok {
		data.Id = 0
	}
	return data
}

func (d *UserDao) GetAll(page, size int) []models.LtUser {
	offset := (page - 1) * size
	dataList := make([]models.LtUser, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&dataList)
	if err != nil {
		log.Printf("user_dao.GetAll error, error info is %s\n", err)
	}
	return dataList
}

func (d *UserDao) CountAll() int64 {
	num, err := d.engine.Count(&models.LtUser{})
	if err != nil {
		return 0
	}
	return num
}

func (d *UserDao) Update(data *models.LtUser, column []string) error {
	_, err := d.engine.Id(data.Id).MustCols(column...).Update(data)
	return err
}

func (d *UserDao) Create(data *models.LtUser) error {
	_, err := d.engine.Insert(data)
	return err
}
