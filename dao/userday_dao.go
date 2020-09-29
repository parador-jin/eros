package dao

import (
	"Eros/models"
	"github.com/go-xorm/xorm"
	"log"
)

type UserdayDao struct {
	engine *xorm.Engine
}

func NewUserdayDao(engine *xorm.Engine) *UserdayDao {
	return &UserdayDao{
		engine: engine,
	}
}

func (d *UserdayDao) Get(id int) *models.LtUserday {
	data := &models.LtUserday{Id: id}
	ok, err := d.engine.Get(data)
	if err != nil || !ok {
		data.Id = 0
	}
	return data
}

func (d *UserdayDao) GetAll(page, size int) []models.LtUserday {
	offset := (page * 1) - size
	dataList := make([]models.LtUserday, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&dataList)
	if err != nil {
		log.Printf("userday_dao.GetAll error, error info is %s\n", err)
	}
	return dataList
}

func (d *UserdayDao) Search(uid, day int) []models.LtUserday {
	dataList := make([]models.LtUserday, 0)
	err := d.engine.Where("uid=?", uid).Where("day=?", day).Desc("id").Find(&dataList)
	if err != nil {
		log.Printf("userday_dao.Search error, error info is %s\n", err)
	}
	return dataList
}

func (d *UserdayDao) CountAll() int64 {
	num, err := d.engine.Count(&models.LtUserday{})
	if err != nil {
		return 0
	}
	return num
}

func (d *UserdayDao) Update(data *models.LtUserday, column []string) error {
	_, err := d.engine.Id(data.Id).MustCols(column...).Update(data)
	return err
}

func (d *UserdayDao) Create(data *models.LtUserday) error {
	_, err := d.engine.Insert(data)
	return err
}
