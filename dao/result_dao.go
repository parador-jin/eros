package dao

import (
	"Eros/models"
	"github.com/go-xorm/xorm"
	"log"
)

type ResultDao struct {
	engine *xorm.Engine
}

func NewResultDao(engine *xorm.Engine) *ResultDao {
	return &ResultDao{
		engine: engine,
	}
}

func (d *ResultDao) Get(id int) *models.LtResult {
	data := &models.LtResult{Id: id}
	ok, err := d.engine.Get(data)
	if err != nil || !ok {
		data.Id = 0
	}
	return data
}

func (d *ResultDao) GetAll(page, size int) []models.LtResult {
	offset := (page - 1) * size
	dataList := make([]models.LtResult, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&dataList)
	if err != nil {
		log.Printf("result_dao.GetAll error, error info is %s\n", err)
	}
	return dataList
}

func (d *ResultDao) CountAll() int64 {
	num, err := d.engine.Count(&models.LtResult{})
	if err != nil {
		return 0
	}
	return num
}

func (d *ResultDao) CountByGift(giftId int) int64 {
	num, err := d.engine.Where("gift_id=?", giftId).Count(&models.LtResult{})
	if err != nil {
		num = 0
	}
	return num
}

func (d *ResultDao) CountByUid(uid int) int64 {
	num, err := d.engine.Where("uid=?", uid).Count(&models.LtResult{})
	if err != nil {
		num = 0
	}
	return num
}

func (d *ResultDao) Delete(id int) error {
	data := &models.LtResult{Id: id, SysStatus: 1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *ResultDao) Update(data *models.LtResult, column []string) error {
	_, err := d.engine.Id(data.Id).MustCols(column...).Update(data)
	return err
}

func (d *ResultDao) Create(data *models.LtResult) error {
	_, err := d.engine.Insert(data)
	return err
}

func (d *ResultDao) SearchByGift(giftId, page, size int) []models.LtResult {
	offset := (page - 1) * size
	dataList := make([]models.LtResult, 0)
	err := d.engine.Where("gift_id=?", giftId).Desc("id").Limit(size, offset).Find(&dataList)
	if err != nil {
		log.Printf("result_dao.SearchByGift error, error info is %s\n", err)
	}
	return dataList
}

func (d *ResultDao) SearchByUid(uid, page, size int) []models.LtResult {
	offset := (page - 1) * size
	dataList := make([]models.LtResult, 0)
	err := d.engine.Where("uid=?", uid).Desc("id").Limit(size, offset).Find(&dataList)
	if err != nil {
		log.Printf("result_dao.SearchByUid error, error info is %s\n", err)
	}
	return dataList
}

func (d *ResultDao) GetNewPrize(size int, giftIds []int) []models.LtResult {
	dataList := make([]models.LtResult, 0)
	err := d.engine.In("gift_id", giftIds).Limit(size).Find(&dataList)
	if err != nil {
		log.Printf("result_dao.GetNewPrize error, error info is %s\n", err)
	}
	return dataList
}
