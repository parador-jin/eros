package dao

import (
	"Eros/models"
	"github.com/go-xorm/xorm"
	"log"
)

type CodeDao struct {
	engine *xorm.Engine
}

func NewCodeDao(engine *xorm.Engine) *CodeDao {
	return &CodeDao{
		engine: engine,
	}
}

func (d *CodeDao) Get(id int) *models.LtCode {
	data := &models.LtCode{Id: id}
	ok, err := d.engine.Get(data)
	if err != nil || !ok {
		data.Id = 0
	}
	return data
}

func (d *CodeDao) GetAll(page, size int) []models.LtCode {
	offset := (page - 1) * size
	dataList := make([]models.LtCode, 0)
	err := d.engine.Desc("id").Limit(size, offset).Find(&dataList)
	if err != nil {
		log.Printf("code_dao.GetAll error, error info is %s\n", err)
	}
	return dataList
}

func (d *CodeDao) CountAll() int64 {
	num, err := d.engine.Count(&models.LtCode{})
	if err != nil {
		return 0
	}
	return num
}

func (d *CodeDao) Delete(id int) error {
	data := &models.LtCode{Id: id, SysStatus: 1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *CodeDao) Update(data *models.LtCode, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *CodeDao) UpdateByCode(data *models.LtCode, columns []string) error {
	_, err := d.engine.Where("code=?", data.Code).MustCols(columns...).Update(data)
	return err
}

func (d *CodeDao) Create(data *models.LtCode) error {
	_, err := d.engine.Insert(data)
	return err
}

func (d *CodeDao) Search(giftId int) []models.LtCode {
	datalist := make([]models.LtCode, 0)
	err := d.engine.Where("gift_id=?", giftId).Desc("id").Find(&datalist)
	if err != nil {
		return datalist
	} else {
		return datalist
	}
}

func (d *CodeDao) CountByGift(giftId int) int64 {
	num, err := d.engine.Count(&models.LtCode{GiftId: giftId})
	if err != nil {
		return 0
	}
	return num
}

func (d *CodeDao) NextUsingCode(giftId, codeId int) *models.LtCode {
	dataList := make([]models.LtCode, 0)
	err := d.engine.Where("gift_id=?", giftId).Where("sys_status=?", 0).Where("id>?", codeId).
		Asc("id").Limit(1).Find(&dataList)
	if err != nil || len(dataList) < 1 {
		return nil
	}
	return &dataList[0]
}
