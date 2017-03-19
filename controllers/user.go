package controllers

import (

	"github.com/astaxie/beego"
	"cloud-platform-ua/models"
	"time"
)

// Operations about Users
type UserController struct {
	BaseController
}
// @Description register
// @Param phone formData string true "注册用户的手机号"
// @Param name formData string true "注册用户的名称"
// @Param password fromData string true "注册用户的密码"
// @router /register [post]
func (this *UserController) Register(){
	form := models.RegisterForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseRegsiterForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	beego.Debug("ParseRegsiterForm:", &form)

	if err := this.VerifyForm(&form); err != nil {
		beego.Debug("ValidRegsiterForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}

	regDate := time.Now()
	user, err := models.NewUser(&form, regDate)
	if err != nil {
		beego.Error("NewUser:", err)
		this.Data["json"] = models.NewErrorInfo(ErrSystem)
		this.ServeJSON()
		return
	}
	beego.Debug("NewUser:", user)

	if code, err := user.Insert(); err != nil {
		beego.Error("InsertUser:", err)
		if code == models.ErrDupRows {
			this.Data["json"] = models.NewErrorInfo(ErrDupUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJSON()
		return
	}

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJSON()
}
// @Description User login
// @Param phone formData string false "用户手机号"
// @Param name formData string false "用户名"
// @Param password formData string true "密码"
// @router /login [post]
func (this *UserController) Login(){

}

