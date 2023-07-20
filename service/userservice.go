package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"imessage/models"
	"imessage/utils"
	"math/rand"
	"strconv"
	"time"
)

// GetUserList
// @Summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	// 从数据库获得数据,将所有的数据存储成数据,然后返回
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": data,
	})
}

// CreateUser
// @Summary 新增用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param phone query string false "手机号"
// @param email query string false "邮箱"
// @param passWord query string false "密码"
// @param rePassWord query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据,在这儿进行组创,然后进行存储
	// 先判断是否有冲突
	user.Name = c.Query("name")
	user.Phone = c.Query("phone")
	user.Email = c.Query("email")

	if !models.IsUniqueCreateUser(user) {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "请检查你的用户名或者电话,邮箱,它们已被注册!",
		})
		return
	}
	// 如果没有冲突,继续
	passWord := c.Query("passWord")
	rePassWord := c.Query("rePassWord")
	// 获得一个随机数
	salt := fmt.Sprintf("%06d", rand.Int31())
	if passWord != rePassWord {
		c.JSON(-4, gin.H{
			"code":    -1,
			"message": "两次密码不一致!",
		})

	} else {
		user.Salt = salt
		// 这里暂时存入一个不准确的时间

		user.PassWord = utils.MakePassword(passWord, user.Salt)
		user.HeartBeatTime = time.Now()
		user.LoginTime = time.Now()
		user.LoginOutTime = time.Now()

		models.CreateUser(user)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "新增用户成功!",
			"data":    user,
		})
	}

}

// DeleteUser
// @Summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据id,然后由于 id 是主要键值,再进行查找\删除操作
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "删除用户成功!",
	})
}

// UpdateUser
// @Summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "name"
// @param password formData string false "password"
// @param phone formData string false "phone"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	// 获取前端数据id,然后由于 id 是主要键值,再进行查找,删除操作
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)

	user.Name = c.PostForm("name")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	// 判断相异性

	if !models.IsUniqueUpdateUser(user) {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "请检查你的用户名或者电话,邮箱,它们已被注册!",
		})
		return
	}
	// 这个要和 salt 一起更新,这个很关键,不然修改密码后,无法登陆
	// 获得一个随机数
	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt
	passwordRaw := c.PostForm("password")
	PWD := utils.MakePassword(passwordRaw, user.Salt)
	user.PassWord = PWD

	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		fmt.Println(err)
		c.JSON(200, gin.H{
			"code":    -1, // 0 : 成功; -1 : 失败
			"message": "修改用户失败!",
		})

	} else {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0, // 0 : 成功; -1 : 失败
			"message": "修改用户成功!",
			"data":    user,
		})
	}

}

// UserLogin
// @Summary 用户登录
// @Tags 用户模块
// @param name formData string false "name"
// @param password formData string false "password"
// @Success 200 {string} json{"code","message"}
// @Router /user/userLogin [post]
func UserLogin(c *gin.Context) {
	// 拿到前端传来的用户名和密码
	Name := c.PostForm("name")
	PassWord := c.PostForm("password")
	// 根据用户名查询用户,如果查到进行密码判断
	user := models.FindUserByName(Name)

	if user.Name != "" {
		// 如果查询到用户名,则进行密码验证
		// 加密
		PWD := utils.MakePassword(PassWord, user.Salt)
		// 用密码和用户名来查询用户
		data := models.FindUserByNameAndPwd(Name, PWD)
		if data.Name != "" {
			c.JSON(200, gin.H{
				"code":    0, // 0 : 成功; -1 : 失败
				"message": "登陆成功!",
				"data":    data,
			})

		} else {
			c.JSON(200, gin.H{
				"code":    -1, // 0 : 成功; -1 : 失败
				"message": "登陆失败,请检查用户名或者登录密码!",
			})

		}

	} else {
		c.JSON(200, gin.H{
			"code":    -1, // 0 : 成功; -1 : 失败
			"message": "登陆失败,请检查用户名或者登录密码!",
		})
	}

}
