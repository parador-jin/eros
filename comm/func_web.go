package comm

import (
	"Eros/conf"
	"Eros/models"
	"crypto/md5"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func ClientIP(request *http.Request) string {
	host, _, _ := net.SplitHostPort(request.RemoteAddr)
	return host
}

func Redirect(writer http.ResponseWriter, url string) {
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusFound)
}

func GetLoginUser(request *http.Request) *models.ObjLoginuser {
	c, err := request.Cookie("eros-loginUser")
	if err != nil {
		return nil
	}

	params, errValue := url.ParseQuery(c.Value)
	if errValue != nil {
		return nil
	}

	uid, errUid := strconv.Atoi(params.Get("uid"))
	if errUid != nil || uid < 1 {
		return nil
	}

	now, errNow := strconv.Atoi(params.Get("now"))
	if errNow != nil || NowUnix()-now > 86400*30 {
		return nil
	}

	loginUser := &models.ObjLoginuser{}
	loginUser.Uid = uid
	loginUser.Username = params.Get("username")
	loginUser.Now = now
	loginUser.Ip = ClientIP(request)
	loginUser.Sign = params.Get("sign")
	sign := createLoginUserSign(loginUser)
	if sign != loginUser.Sign {
		log.Println("func_web GetLoginUser createLoginUserSign not signed", sign, loginUser.Sign)
		return nil
	}
	return loginUser
}

func SetLoginUser(writer http.ResponseWriter, loginUser *models.ObjLoginuser) {
	if loginUser == nil || loginUser.Uid < 1 {
		c := &http.Cookie{
			Name:       "eros-loginUser",
			Value:      "",
			Path:       "/",
			Domain:     "",
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     -1,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		}
		http.SetCookie(writer, c)
		return
	}
	if loginUser.Sign == "" {
		loginUser.Sign = createLoginUserSign(loginUser)
	}
	params := url.Values{}
	params.Add("uid", strconv.Itoa(loginUser.Uid))
	params.Add("username", loginUser.Username)
	params.Add("now", strconv.Itoa(loginUser.Now))
	params.Add("ip", loginUser.Ip)
	params.Add("sign", loginUser.Sign)
	c := &http.Cookie{
		Name:       "eros-loginUser",
		Value:      params.Encode(),
		Path:       "/",
		Domain:     "",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(writer, c)
}

func createLoginUserSign(loginUser *models.ObjLoginuser) string {
	str := fmt.Sprintf("uid=%d&username=%s&secret=%s&now=%d",
		loginUser.Uid, loginUser.Username, conf.CookieSecret, loginUser.Now)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	return sign
}
