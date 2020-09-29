package services

import (
	"Eros/dao"
	"Eros/datasource"
	"Eros/models"
)

type ResultService interface {
	GetAll(page, size int) []models.LtResult
	CountAll() int64
	CountByGift(giftId int) int64
	CountByUid(Uid int) int64
	Get(id int) *models.LtResult
	Delete(id int) error
	Update(data *models.LtResult, columns []string) error
	Create(data *models.LtResult) error
	SearchByGift(giftId, page, size int) []models.LtResult
	SearchByUid(uid, page, size int) []models.LtResult
	GetNewPrize(size int, giftIds []int) []models.LtResult
}

type resultService struct {
	dao *dao.ResultDao
}

func NewResultService() ResultService {
	return &resultService{
		dao: dao.NewResultDao(datasource.InstanceDbMaster()),
	}
}

func (s *resultService) GetAll(page, size int) []models.LtResult {
	return s.dao.GetAll(page, size)
}

func (s *resultService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *resultService) CountByGift(giftId int) int64 {
	return s.dao.CountByGift(giftId)
}

func (s *resultService) CountByUid(uid int) int64 {
	return s.dao.CountByUid(uid)
}

func (s *resultService) Get(id int) *models.LtResult {
	return s.dao.Get(id)
}

func (s *resultService) Delete(id int) error {
	return s.dao.Delete(id)
}

func (s *resultService) Update(data *models.LtResult, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *resultService) Create(data *models.LtResult) error {
	return s.dao.Create(data)
}

func (s *resultService) SearchByGift(giftId, page, size int) []models.LtResult {
	return s.dao.SearchByGift(giftId, page, size)
}

func (s *resultService) SearchByUid(uid, page, size int) []models.LtResult {
	return s.dao.SearchByGift(uid, page, size)
}

func (s *resultService) GetNewPrize(size int, giftIds []int) []models.LtResult {
	return s.dao.GetNewPrize(size, giftIds)
}
