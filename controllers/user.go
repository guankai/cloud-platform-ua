package controllers

import (
	"github.com/astaxie/beego"
	"cloud-platform-ua/models"
	"time"
	"github.com/astaxie/beego/httplib"
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

	// 在gogs上创建用户
	statusCode, gitErr := CreateGitUser(&form)
	if gitErr != nil {
		beego.Error("Git register user error:", gitErr)
		this.Data["json"] = models.NewErrorInfo(ErrGitReg)
		this.ServeJSON()
		return
	}
	if statusCode != 201 {
		this.Data["json"] = models.NewErrorInfo(ErrGitReg)
		this.ServeJSON()
		return
	}

	// 创建镜像仓库的repository
	hubStatusCode, hubErr := CreateHub(&form)
	if hubErr != nil {
		beego.Error("Hub repository create error:", hubErr)
		this.Data["json"] = models.NewErrorInfo(ErrHubReg)
		this.ServeJSON()
		return
	}
	beego.Debug("the hub status code is ", hubStatusCode)
	if hubStatusCode != 201 {
		this.Data["json"] = models.NewErrorInfo(ErrHubReg)
		this.ServeJSON()
		return
	}

	//创建k8s的namespace
	k8sStatusCode, k8sErr := CreateK8sNamespace(&form)
	if k8sErr != nil {
		beego.Error("k8s namespace create error:", hubErr)
		this.Data["json"] = models.NewErrorInfo(ErrK8sReg)
		this.ServeJSON()
		return
	}
	if k8sStatusCode != 200 {
		this.Data["json"] = models.NewErrorInfo(ErrK8sReg)
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

	this.Data["json"] = models.NewNormalInfo("Success")
	this.ServeJSON()
}
// @Description User login
// @Param phone formData string false "用户手机号"
// @Param name formData string false "用户名"
// @Param password formData string true "密码"
// @router /login [post]
func (this *UserController) Login() {
	name := this.GetString("name")
	password := this.GetString("password")
	// 验证输入信息
	if name == "" {
		beego.Error("请输入用户名!")
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	// 验证用户是否存在
	user := models.User{}

	if code, err := user.FindByName(name); err != nil {
		beego.Error("通过用户名查找用户失败", err)
		if code == models.ErrNotFound {
			this.Data["json"] = models.NewErrorInfo(ErrNoUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJSON()
		return
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

	this.SetSession(SessId + user.Name, user.Name)

	this.Data["json"] = &models.LoginInfo{Code: 0, UserInfo: &user}
	this.ServeJSON()

}
// @Description user logout
// @Param name formData string true "用户名"
// @router /logout [post]
func (this *UserController) Logout() {
	form := models.LogoutForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseLogoutForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	beego.Debug("ParseLogoutForm:", &form)

	if err := this.VerifyForm(&form); err != nil {
		beego.Debug("ValidLogoutForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}

	if this.GetSession(SessId + form.Name) != form.Name {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJSON()
		return
	}

	this.DelSession(SessId + form.Name)

	this.Data["json"] = models.NewNormalInfo("Success")
	this.ServeJSON()
}
// @Description User update user information
// @Param phone formData string true "用户手机号"
// @Param name formData string true "用户名"
// @Param email formData string tru "用户邮箱"
// @router /update [post]
func (this *UserController) UserUpdate() {
	updateForm := models.UpdateForm{}
	if err := this.ParseForm(&updateForm); err != nil {
		beego.Debug("ParseLogoutForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	beego.Debug("UserUpateForm:", &updateForm)

	if err := this.VerifyForm(&updateForm); err != nil {
		beego.Debug("ValidLogoutForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}

	if this.GetSession(SessId + updateForm.Phone) != updateForm.Phone {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJSON()
		return
	}

	user := models.User{}
	if code, err := user.FindByName(updateForm.Name); err != nil {
		beego.Error("通过用户名查找用户失败", err)
		if code == models.ErrNotFound {
			this.Data["json"] = models.NewErrorInfo(ErrNoUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJSON()
		return
	}

	user.Name = updateForm.Name
	user.Email = updateForm.Email

	if err := user.UpdateUser(); err != nil {
		beego.Error("更新用户信息失败", err)
		this.Data["json"] = models.NewErrorInfo(ErrDatabase)
	}

	this.Data["json"] = models.NewNormalInfo("Success")
	this.ServeJSON()

}
// @Description get user information
// @Param name path string true "用户名"
// @router /:name [get]
func (this *UserController) GetUserInfo() {
	name := this.GetString(":name")
	user := models.User{}
	if code, err := user.FindByName(name); err != nil {
		beego.Error("通过手机号查找用户失败", err)
		if code == models.ErrNotFound {
			this.Data["json"] = models.NewErrorInfo(ErrNoUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJSON()
		return
	}
	this.Data["json"] = &models.LoginInfo{Code: 0, UserInfo: &user}
	this.ServeJSON()
}
// 创建镜像仓库的repository
func CreateHub(form *models.RegisterForm) (code int, err error) {
	req := httplib.Post(beego.AppConfig.String("hub::url"))
	req.SetBasicAuth(beego.AppConfig.String("hub::user"), beego.AppConfig.String("hub::password"))
	hub := models.Hub{ProjectName:form.Name, Public:1}
	req.JSONBody(hub)
	resp, err := req.Response()
	return resp.StatusCode, err
}

// 在gogs上创建git用户
func CreateGitUser(form *models.RegisterForm) (code int, err error) {
	req := httplib.Post(beego.AppConfig.String("gogs::url") + CreateUser)
	req.SetBasicAuth(beego.AppConfig.String("gogs::admin"), beego.AppConfig.String("gogs::password"))
	req.Header("Content-Type", "application/json")
	req.Param("source_id", "0")
	req.Param("login_name", form.Name)
	req.Param("username", form.Name)
	req.Param("email", form.Email)
	req.Param("password", form.Password)
	resp, err := req.Response()
	return resp.StatusCode, err
}

// 创建kubernetes的namespace
func CreateK8sNamespace(form *models.RegisterForm) (code int, err error) {
	req := httplib.Post(beego.AppConfig.String("k8s::url"))
	req.Param("namespace", form.Name)
	resp, err := req.Response()
	return resp.StatusCode, err
}



