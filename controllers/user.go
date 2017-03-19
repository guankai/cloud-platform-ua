package controllers

import (
	"github.com/astaxie/beego"
	"cloud-platform-ua/models"
	"time"
	"golang.org/x/tools/go/gcimporter15/testdata"
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
func (this *UserController) Register() {
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
func (this *UserController) Login() {
	phone := this.GetString("phone")
	name := this.GetString("name")
	password := this.GetString("password")
	// 验证输入信息
	if phone == "" && name == "" {
		beego.Error("至少输入phone number或者name中的一个")
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	// 验证用户是否存在
	user := models.User{}
	if phone != nil || phone != "" {
		//通过手机号查找
		if code, err := user.FindByID(phone); err != nil {
			beego.Error("通过手机号查找用户失败" + err)
			if code == models.ErrNotFound {
				this.Data["json"] = models.NewErrorInfo(ErrNoUser)
			} else {
				this.Data["json"] = models.NewErrorInfo(ErrDatabase)
			}
			this.ServeJSON()
			return
		}
	} else {
		//通过用户名查找
		if code, err := user.FindByName(name); err != nil {
			beego.Error("通过用户名查找用户失败" + err)
			if code == models.ErrNotFound {
				this.Data["json"] = models.NewErrorInfo(ErrNoUser)
			} else {
				this.Data["json"] = models.NewErrorInfo(ErrDatabase)
			}
			this.ServeJSON()
			return
		}
	}
	beego.Debug("UserInfo:", &user)
	// 验证用户密码
	if ok, err := user.CheckPass(password); err != nil {
		beego.Error("验证用户密码失败:", err)
		this.Data["json"] = models.NewErrorInfo(ErrSystem)
		this.ServeJSON()
		return
	} else if !ok {
		this.Data["json"] = models.NewErrorInfo(ErrPass)
		this.ServeJSON()
		return
	}
	user.ClearPass()

	this.SetSession("user_id", user.ID)

	this.Data["json"] = &models.LoginInfo{Code: 0, UserInfo: &user}
	this.ServeJSON()

}

