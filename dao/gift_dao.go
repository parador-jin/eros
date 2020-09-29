package dao

import (
	"Eros/comm"
	"Eros/models"
	"github.com/go-xorm/xorm"
	"log"
)

type GiftDao struct {
	engine *xorm.Engine
}

func NewGiftDao(engine *xorm.Engine) *GiftDao {
	return &GiftDao{
		engine: engine,
	}
}

func (d *GiftDao) Get(id int) *models.LtGift {
	data := &models.LtGift{Id: id}
	ok, err := d.engine.Get(data)
	if err != nil || !ok {
		data.Id = 0
	}
	return data
}

func (d *GiftDao) GetAll() []models.LtGift {
	dataList := make([]models.LtGift, 0)
	err := d.engine.Asc("sys_status").Asc("displayorder").Find(&dataList)
	if err != nil {
		log.Printf("gift_dao.GetAll error, error info is %s\n", err)
	}
	return dataList
}

// 获取到当前可以获取的奖品列表
// 有奖品限定，状态正常，时间期间内
// gtype倒序，displayorder正序
func (d GiftDao) GetAllUse() []models.LtGift {
	now := comm.NowUnix()
	dataList := make([]models.LtGift, 0)
	err := d.engine.Cols("id", "title", "prize_num", "left_num", "prize_code", "prize_time", "img",
		"displayorder", "gtype", "gdata").Desc("gtype").Asc("displayorder").
		Where("prize_num>=?", 0).Where("sys_status=?", 0).
		Where("time_begin<=?", now).Where("time_end>=?", now).Find(&dataList)
	if err != nil {
		log.Printf("gift_dao.GetAllUse error, error info is %s\n", err)
	}
	return dataList
}

func (d *GiftDao) CountAll() int64 {
	num, err := d.engine.Count(&models.LtGift{})
	if err != nil {
		return 0
	}
	return num
}

func (d *GiftDao) Delete(id int) error {
	data := &models.LtGift{Id: id, SysStatus: 1}
	_, err := d.engine.Id(data.Id).Update(data)
	return err
}

func (d *GiftDao) Update(data *models.LtGift, columns []string) error {
	_, err := d.engine.Id(data.Id).MustCols(columns...).Update(data)
	return err
}

func (d *GiftDao) Create(data *models.LtGift) error {
	_, err := d.engine.Insert(data)
	return err
}

func (d *GiftDao) DecrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.Id(id).Decr("left_num", num).Where("left_num>=?", num).Update(&models.LtGift{Id: id})
	return r, err
}

func (d *GiftDao) IncrLeftNum(id, num int) (int64, error) {
	r, err := d.engine.Id(id).Incr("left_num", num).Update(&models.LtGift{Id: id})
	return r, err
}
