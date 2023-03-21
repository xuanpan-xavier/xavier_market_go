package controller

// 会员的注册与登录
import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"xmarket_gin/common"
	"xmarket_gin/dto"
	"xmarket_gin/model"
	"xmarket_gin/response"
)

func Register(c *gin.Context) {
	db := common.GetDB()
	//使用结构体获取请求的参数
	var requestUser = model.User{}
	//json.NewDecoder(c.Request.Body).Decode(&requestUser)
	c.Bind(&requestUser)
	//获取参数
	name := requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password
	//数据验证
	if len(telephone) != 11 {
		log.Println(telephone)
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位数")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	if len(password) > 15 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能多于15位")
		return
	}
	if len(name) < 3 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户名不能少于3位")
		return
	}
	if len(name) > 10 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户名不能多于10位")
		return
	}
	//log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExist(db, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已存在")
		return
	}

	//创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
		VipCard:   model.VipCard{Telephone: telephone},
	}
	db.Create(&newUser)
	//发放token
	token, err := common.ReleaseToken(newUser)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token genneratr error : %v", err)
		return
	}
	response.Success(c, gin.H{"token": token}, "注册成功！")

}
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if phone, _ := strconv.Atoi(user.Telephone); phone != 0 {
		return true
	}
	return false
}

func Login(c *gin.Context) {
	db := common.GetDB()
	var requestUser = model.User{}
	//json.NewDecoder(c.Request.Body).Decode(&requestUser)
	c.Bind(&requestUser)
	//获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password
	//数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位数")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能少于6位")
		return
	}
	if len(password) > 15 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码不能多于15位")
		return
	}
	//手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if phone, _ := strconv.Atoi(user.Telephone); phone == 0 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}
	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(c, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "系统异常")
		log.Printf("token genneratr error : %v", err)
		return
	}
	//返回结果
	response.Success(c, gin.H{"token": token}, "登录成功！")
}

func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": dto.ToUserDto(user.(model.User))}})
}
