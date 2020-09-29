package dao

import (
	"Eros/models"
	"github.com/go-xorm/xorm"
	"log"
)

type BlackipDao struct {
	engine *xorm.Engine
}

func NewBlackipDao(engine *xorm.Engine) *BlackipDao {
	return &BlackipDao{
		engine: engine,
	}
}

func (d *BlackipDao) Get(id int) *models.LtBlackip {
	data := &models.LtBlackip{Id: id}
	ok, err := d.engine.Get(data)
	if err != nil || !ok {
		data.Id = 0
	}
	return data
}

func (d *BlackipDao) GetByIP(ip string) *models.LtBlackip {
	dataList := make([]models.LtBlackip, 0)
	err := d.engine.Where("ip=?", ip).Desc("id").Limit(1).Find(&dataList)
	if err != nil {
		return nil
	}
	return &dataList[0]
}

func (d *BlackipDao) GetAll(page, size int) []models.LtBlackip {
	offset := (page - 1) * size
	dataList := make([]models.LtBlackip, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&dataList)
	if err != nil {
		log.Printf("blackip_dao.GetAll error, error info is %s\n", err)
	}
	return dataList
}

func (d *BlackipDao) CountAll() int64 {
	num, err := d.engine.Count(&models.LtBlackip{})
	if err != nil {
		return 0
	}
	return num
}

func (d *BlackipDao) Delete(id int) error {
	data := &models.LtBlackip{Id: id}
	_, err := d.engine.Id(data.Id).Delete(data)
	return err
}

func (d *BlackipDao) Update(data *models.LtBlackip, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *BlackipDao) Create(data *models.LtBlackip) error {
	_, err := d.engine.Insert(data)
	return err
}
