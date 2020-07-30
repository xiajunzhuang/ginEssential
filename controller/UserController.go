package controller

import (
	"ginEssential/common"
	"ginEssential/dto"
	"ginEssential/model"
	"ginEssential/response"
	"ginEssential/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func Register(context *gin.Context) {
	db := common.GetDB()
	//使用map获取参数
	//var requestMap = make(map[string]string)
	//json.NewDecoder(context.Request.Body).Decode(&requestMap)

	//var requestUser = model.User{}
	//json.NewDecoder(context.Request.Body).Decode(&requestUser)

	var requestUser = model.User{}
	context.Bind(&requestUser)
	//获取参数
	name :=requestUser.Name
	telephone := requestUser.Telephone
	password := requestUser.Password

	if len(telephone) != 11 {
		response.Response(context,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}

	if len(password)<6{
		response.Response(context,http.StatusUnprocessableEntity,422,nil,"密码不能小于6位")
		return
	}

	//姓名为空传10位字符串
	if len(name)==0{
		name = utils.RandomString(10)
	}

	//判断手机号是否存在
	if isTelephoneExist(db,telephone){
		response.Response(context,http.StatusUnprocessableEntity,422,nil,"手机号已存在")
		return
	}

	//创建用户
	hasedPassword,err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil{
		response.Response(context,http.StatusInternalServerError,500,nil,"加密错误")
		return
	}
	newUser := model.User{
		Name: name,
		Telephone: telephone,
		Password: string(hasedPassword),
	}
	db.Create(&newUser)
	//发放Token
	token,err := common.ReleaseToken(newUser)
	if err != nil{
		response.Response(context,http.StatusInternalServerError,500,nil,"系统错误")
		log.Printf("token generate error: %v ",err)
		return
	}
	//返回结果
	response.Success(context,gin.H{"token": token},"注册成功")

}

func Login(context *gin.Context)  {
	db := common.GetDB()
	var requestUser = model.User{}
	context.Bind(&requestUser)
	//获取参数
	telephone := requestUser.Telephone
	password := requestUser.Password
	//数据验证
	if len(telephone) != 11 {
		response.Response(context,http.StatusUnprocessableEntity,422,nil,"手机号必须为11位")
		return
	}

	if len(password)<6{
		response.Response(context,http.StatusUnprocessableEntity,422,nil,"密码不能小于6位")
		return
	}

	//判断手机号是否存在
	var user model.User
	db.Where("telephone=?",telephone).First(&user)
	if user.ID == 0{
		response.Response(context,http.StatusUnprocessableEntity,422,nil,"用户不存在")
		return
	}
	//判断密码是否正确
	 err := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))
	 if err !=nil{
		 response.Response(context,http.StatusBadRequest,400,nil,"密码错误")
		return
	}
	//发放Token
	token,err := common.ReleaseToken(user)
	if err != nil{
		response.Response(context,http.StatusInternalServerError,500,nil,"系统错误")
		log.Printf("token generate error: %v ",err)
		return
	}
	//返回结果
	response.Success(context,gin.H{"token": token},"登录成功")


}

func Info(context *gin.Context)  {
	user,_ := context.Get("user")

	context.JSON(http.StatusOK,gin.H{"code":200,"data":gin.H{"user":dto.ToUserDto(user.(model.User))}})
}

func isTelephoneExist(db *gorm.DB,telephone string) bool {
	var user model.User
	db.Where("telephone=?",telephone).First(&user)
	if user.ID != 0{
		return true
	}
	return false
}
