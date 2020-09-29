package services

import (
	"Eros/comm"
	"Eros/dao"
	"Eros/datasource"
	"Eros/models"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
)

type UserService interface {
	GetAll(page, size int) []models.LtUser
	CountAll() int64
	Get(id int) *models.LtUser
	Update(data *models.LtUser, columns []string) error
	Create(data *models.LtUser) error
}

type userService struct {
	dao *dao.UserDao
}

func NewUserService() UserService {
	return &userService{
		dao: dao.NewUserDao(datasource.InstanceDbMaster()),
	}
}

func (s *userService) GetAll(page, size int) []models.LtUser {
	return s.dao.GetAll(page, size)
}

func (s *userService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *userService) Get(id int) *models.LtUser {
	data := s.getByCache(id)
	if data == nil || data.Id < 0 {
		data = s.dao.Get(id)
		if data == nil || data.Id < 0 {
			data = &models.LtUser{Id: id}
		}
		s.setByCache(data)
	}
	return data
}

func (s *userService) Update(data *models.LtUser, columns []string) error {
	s.updateByCache(data)
	return s.dao.Update(data, columns)
}

func (s *userService) Create(data *models.LtUser) error {
	return s.dao.Create(data)
}

func (s *userService) getByCache(id int) *models.LtUser {
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("user_service.getByCache HGETALL key is ", key, ", error is ", err)
		return nil
	}
	dataId := comm.GetInt64FromStringMap(dataMap, "Id", 0)
	if dataId <= 0 {
		return nil
	}
	data := &models.LtUser{
		Id:         int(dataId),
		Username:   comm.GetStringFromStringMap(dataMap, "Username", ""),
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		Realname:   comm.GetStringFromStringMap(dataMap, "Realname", ""),
		Mobile:     comm.GetStringFromStringMap(dataMap, "Mobile", ""),
		Address:    comm.GetStringFromStringMap(dataMap, "Address", ""),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
		SysIp:      comm.GetStringFromStringMap(dataMap, "SysIp", ""),
	}
	return data
}

func (s *userService) setByCache(data *models.LtUser) {
	if data == nil || data.Id <= 0 {
		return
	}
	id := data.Id
	key := fmt.Sprintf("info_user_%d", id)
	rds := datasource.InstanceCache()
	params := redis.Args{key}
	params = params.Add(id)
	if data.Username != "" {
		params = params.Add(params, "Username", data.Username)
		params = params.Add(params, "Blacktime", data.Blacktime)
		params = params.Add(params, "Realname", data.Realname)
		params = params.Add(params, "Mobile", data.Mobile)
		params = params.Add(params, "Address", data.Address)
		params = params.Add(params, "SysUpdated", data.SysUpdated)
		params = params.Add(params, "SysCreated", data.SysCreated)
		params = params.Add(params, "SysIp", data.SysIp)
	}
	_, err := rds.Do("HMSET", params)
	if err != nil {
		log.Println("user_service.setByCache HMSET params is ", params, ", error is ", err)
	}
}

func (s *userService) updateByCache(data *models.LtUser) {
	if data == nil || data.Id <= 0 {
		return
	}
	key := fmt.Sprintf("info_user_%d", data.Id)
	rds := datasource.InstanceCache()
	rds.Do("DEL", key)
}
