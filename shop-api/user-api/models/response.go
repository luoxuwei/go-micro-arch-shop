package models


import (
	"fmt"
	"time"
)

//定制时间转字符串的格式，不能在time.Time上加方法，但是可以在它之上定义一个新类型，可以绕过限制
type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error){
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id int32 `json:"id"`
	NickName string `json:"name"`
	//Birthday string `json:"birthday"`
	Birthday JsonTime `json:"birthday"`
	Gender string `json:"gender"`
	Mobile string `json:"mobile"`
}