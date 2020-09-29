package services

import (
	"Eros/comm"
	"Eros/dao"
	"Eros/datasource"
	"Eros/models"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

type GiftService interface {
	GetAll(useCache bool) []models.LtGift
	CountAll() int64
	Get(id int, useCache bool) *models.LtGift
	Delete(id int) error
	Update(data *models.LtGift, columns []string) error
	Create(data *models.LtGift) error
	GetAllUse(useCache bool) []models.ObjGiftPrize
	DecrLeftNum(id, num int) (int64, error)
	IncrLeftNum(id, num int) (int64, error)
}

type giftService struct {
	dao *dao.GiftDao
}

func NewGiftService() GiftService {
	return &giftService{
		dao: dao.NewGiftDao(datasource.InstanceDbMaster()),
	}
}

func (s *giftService) GetAll(useCache bool) []models.LtGift {
	if !useCache {
		return s.dao.GetAll()
	}
	gifts := s.getAllByCache()
	if len(gifts) < 1 {
		gifts := s.dao.GetAll()
		s.setAllByCache(gifts)
	}
	return gifts
}

func (s *giftService) CountAll() int64 {
	gifts := s.GetAll(true)
	return int64(len(gifts))
}

func (s *giftService) Get(id int, useCache bool) *models.LtGift {
	if !useCache {
		return s.dao.Get(id)
	}
	gifts := s.GetAll(true)
	for _, gift := range gifts {
		if gift.Id == id {
			return &gift
		}
	}
	return nil
}

func (s *giftService) Delete(id int) error {
	data := &models.LtGift{Id: id}
	s.updateByCache(data)
	return s.dao.Delete(id)
}

func (s *giftService) Update(data *models.LtGift, columns []string) error {
	s.updateByCache(data)
	return s.dao.Update(data, columns)
}

func (s *giftService) Create(data *models.LtGift) error {
	s.updateByCache(data)
	return s.dao.Create(data)
}

// 获取到当前可以获取的奖品列表
// 有奖品限定，状态正常，时间期间内
// gtype倒序，displayorder正序
func (s *giftService) GetAllUse(useCache bool) []models.ObjGiftPrize {
	list := make([]models.LtGift, 0)
	if !useCache {
		list = s.dao.GetAllUse()
	} else {
		now := comm.NowUnix()
		gifts := s.GetAll(true)
		for _, gift := range gifts {
			if gift.Id > 0 && gift.SysStatus == 0 && gift.PrizeNum >= 0 && gift.TimeBegin <= now && gift.TimeEnd >= now {
				list = append(list, gift)
			}
		}
	}
	if list != nil {
		gifts := make([]models.ObjGiftPrize, 0)
		for _, gift := range list {
			codes := strings.Split(gift.PrizeCode, "-")
			if len(codes) == 2 {
				// 设置了获奖编码范围 a-b 才可以进行抽奖
				codeA := codes[0]
				codeB := codes[1]
				a, e1 := strconv.Atoi(codeA)
				b, e2 := strconv.Atoi(codeB)
				if e1 == nil && e2 == nil && b >= a && a >= 0 && b < 10000 {
					data := models.ObjGiftPrize{
						Id:           gift.Id,
						Title:        gift.Title,
						PrizeNum:     gift.PrizeNum,
						LeftNum:      gift.LeftNum,
						PrizeCodeA:   a,
						PrizeCodeB:   b,
						Img:          gift.Img,
						Displayorder: gift.Displayorder,
						Gtype:        gift.Gtype,
						Gdata:        gift.Gdata,
					}
					gifts = append(gifts, data)
				}
			}
		}
		return gifts
	}
	return []models.ObjGiftPrize{}
}

func (s *giftService) DecrLeftNum(id, num int) (int64, error) {
	return s.dao.DecrLeftNum(id, num)
}

func (s *giftService) IncrLeftNum(id, num int) (int64, error) {
	return s.dao.IncrLeftNum(id, num)
}

func (s *giftService) getAllByCache() []models.LtGift {
	key := "allgift"
	rds := datasource.InstanceCache()
	result, err := rds.Do("GET", key)
	if err != nil {
		log.Println("gift_service.getAllByCache GET key is ", key, ", error is ", err)
		return nil
	}
	str := comm.GetString(result, "")
	if str == "" {
		return nil
	}
	var dataList []map[string]interface{}
	if errJson := json.Unmarshal([]byte(str), &dataList); errJson != nil {
		log.Println("gift_service.getAllByCache json.Unmarshal error is ", errJson)
		return nil
	}
	gifts := make([]models.LtGift, len(dataList))
	for i := 0; i < len(dataList); i++ {
		data := dataList[i]
		id := comm.GetInt64FromMap(data, "Id", 0)
		if id <= 0 {
			gifts[i] = models.LtGift{}
		} else {
			gift := models.LtGift{
				Id:           int(id),
				Title:        comm.GetStringFromMap(data, "Title", ""),
				PrizeNum:     int(comm.GetInt64FromMap(data, "PrizeNum", 0)),
				LeftNum:      int(comm.GetInt64FromMap(data, "LeftNum", 0)),
				PrizeCode:    comm.GetStringFromMap(data, "PrizeCode", ""),
				PrizeTime:    int(comm.GetInt64FromMap(data, "PrizeTime", 0)),
				Img:          comm.GetStringFromMap(data, "Img", ""),
				Displayorder: int(comm.GetInt64FromMap(data, "Displayorder", 0)),
				Gtype:        int(comm.GetInt64FromMap(data, "Gtype", 0)),
				Gdata:        comm.GetStringFromMap(data, "Gdata", ""),
				TimeBegin:    int(comm.GetInt64FromMap(data, "TimeBegin", 0)),
				TimeEnd:      int(comm.GetInt64FromMap(data, "TimeEnd", 0)),
				PrizeBegin:   int(comm.GetInt64FromMap(data, "PrizeBegin", 0)),
				PrizeEnd:     int(comm.GetInt64FromMap(data, "PrizeEnd", 0)),
				SysStatus:    int(comm.GetInt64FromMap(data, "SysStatus", 0)),
				SysCreated:   int(comm.GetInt64FromMap(data, "SysCreated", 0)),
				SysUpdated:   int(comm.GetInt64FromMap(data, "SysUpdated", 0)),
				SysIp:        comm.GetStringFromMap(data, "SysIp", ""),
			}
			gifts[i] = gift
		}
	}
	return gifts
}

func (s *giftService) setAllByCache(gifts []models.LtGift) {
	strValue := ""
	if len(gifts) > 0 {
		dataList := make([]map[string]interface{}, len(gifts))
		// 格式转换
		for i := 0; i < len(gifts); i++ {
			gift := gifts[i]
			data := make(map[string]interface{})
			data["Id"] = gift.Id
			data["Title"] = gift.Title
			data["PrizeNum"] = gift.PrizeNum
			data["LeftNum"] = gift.LeftNum
			data["PrizeCode"] = gift.PrizeCode
			data["PrizeTime"] = gift.PrizeTime
			data["Img"] = gift.Img
			data["Displayorder"] = gift.Displayorder
			data["Gtype"] = gift.Gtype
			data["Gdata"] = gift.Gdata
			data["TimeBegin"] = gift.TimeBegin
			data["TimeEnd"] = gift.TimeEnd
			data["PrizeBegin"] = gift.PrizeBegin
			data["PrizeEnd"] = gift.PrizeEnd
			data["SysStatus"] = gift.SysStatus
			data["SysCreated"] = gift.SysCreated
			data["SysUpdated"] = gift.SysUpdated
			data["SysIp"] = gift.SysIp
			dataList[i] = data
		}
		str, errMar := json.Marshal(dataList)
		if errMar != nil {
			log.Println("gift_service.setAllByCache json.Marshal error is ", errMar)
		}
		strValue = string(str)
	}
	key := "allgift"
	rds := datasource.InstanceCache()
	if _, err := rds.Do("SET", key, strValue); err != nil {
		log.Println("gift_service.setAllByCache redis SET key is ", key, ", error is ", err)
	}
}

// 数据更新，需要更新缓存，直接清空缓存数据
func (s *giftService) updateByCache(data *models.LtGift) {
	if data == nil || data.Id <= 0 {
		return
	}
	// 集群模式，redis缓存
	key := "allgift"
	rds := datasource.InstanceCache()
	// 删除redis中的缓存
	rds.Do("DEL", key)
}
