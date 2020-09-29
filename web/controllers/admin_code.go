package controllers

import (
	"Eros/comm"
	"Eros/conf"
	"Eros/models"
	"Eros/services"
	"Eros/web/utils"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"strings"
)

type AdminCodeController struct {
	Ctx         iris.Context
	ServiceCode services.CodeService
	ServiceGift services.GiftService
}

// http://localhost:8080/admin/code
func (c *AdminCodeController) Get() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	pagePrev := ""
	pageNext := ""
	// 数据列表
	var dataList []models.LtCode
	var total, num, cacheNum int
	if giftId > 0 {
		dataList = c.ServiceCode.Search(giftId)
		num, cacheNum = utils.GetCacheCodeNum(giftId, c.ServiceCode)
	} else {
		dataList = c.ServiceCode.GetAll(page, size)
	}
	total = (page-1)*size + len(dataList)
	if len(dataList) >= size {
		if giftId > 0 {
			total = int(c.ServiceCode.CountByGift(giftId))
		} else {
			total = int(c.ServiceCode.CountAll())
		}
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/code.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "code",
			"GiftId":   giftId,
			"Datalist": dataList,
			"Total":    total,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
			"CodeNum":  num,
			"CacheNum": cacheNum,
		},
		Layout: "admin/layout.html",
	}
}

// http://localhost:8080/admin/code/import
func (c *AdminCodeController) PostImport() {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	if giftId < 1 {
		c.Ctx.Text("没有指定奖品ID，无法进行导入，<a href='' onclick='history.go(-1);return false;'>返回</a>")
		return
	}
	gift := c.ServiceGift.Get(giftId, true)
	if gift == nil || gift.Id < 1 || gift.Gtype != conf.GtypeCodeDiff {
		c.Ctx.Text("奖品信息不存在或者奖品类型不是差异化优惠券，无法进行导入，<a href='' onclick='history.go(-1);return false;'>返回</a>")
		return
	}
	codes := c.Ctx.PostValue("codes")
	now := comm.NowUnix()
	list := strings.Split(codes, "\n")
	sucNum := 0
	errNum := 0
	for _, code := range list {
		code := strings.TrimSpace(code)
		if code != "" {
			data := &models.LtCode{
				GiftId:     giftId,
				Code:       code,
				SysCreated: now,
			}
			err := c.ServiceCode.Create(data)
			if err != nil {
				errNum++
			} else {
				// 成功导入数据库，下一步还需要导入缓存
				ok := utils.ImportCacheCodes(giftId, code)
				if ok {
					sucNum++
				} else {
					errNum++
				}
			}
		}
	}
	c.Ctx.HTML(fmt.Sprintf("成功导入%d条，导入失败%d条，<a href='/admin/code?gift_id=%d'>返回</a>",
		sucNum, errNum, giftId))
}

// http://localhost:8080/admin/code/delete
func (c *AdminCodeController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceCode.Delete(id)
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	return mvc.Response{
		Path: refer,
	}
}

// http://localhost:8080/admin/code/reset
func (c *AdminCodeController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceCode.Update(&models.LtCode{Id: id, SysStatus: 0}, []string{"sys_status"})
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	return mvc.Response{
		Path: refer,
	}
}

// http://localhost:8080/admin/code/recache
func (c *AdminCodeController) GetRecache() {
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	id, err := c.Ctx.URLParamInt("id")
	if id < 1 || err != nil {
		rs := fmt.Sprintf("没有指定优惠券所属的奖品id, <a href='%s'>饭呼救</a>", refer)
		c.Ctx.HTML(rs)
		return
	}
	sucNum, errNum := utils.ReCacheCodes(id, c.ServiceCode)
	rs := fmt.Sprintf("sucNum is %d, errNum is %d, <a href='%s'></a>", sucNum, errNum, refer)
	c.Ctx.HTML(rs)
}
