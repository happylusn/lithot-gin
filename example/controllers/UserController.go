package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/happylusn/lithot-gin/lithot"
)

type Person struct {
	Name string
}

type UserController struct {
	Person1 *Person `inject:"-"`
	Person2 *Person `inject:"TestConfig.PersonWithArgs('jlu1')"`
	//Db *gorm.DB `inject:"-"`
}

func NewUserController() *UserController {
	return &UserController{}
}
func (this *UserController) Index(c *gin.Context) string {
	return fmt.Sprintf("hello %s, this is 首页", this.Person2.Name)
}
func (this *UserController) Index1(c *gin.Context) lithot.Json {
	return gin.H{"code": 100, "data": nil}
}
func (this *UserController) UserList(c *gin.Context) lithot.SimpleQuery {
	return "select * from users"
}
func (this *UserController) UserDetail(c *gin.Context) lithot.Query {
	return lithot.SimpleQuery("select * from users where id=?").WithArgs(c.Param("id")).WithFirst().WithKey("data")
}
func (this *UserController) Name() string {
	return "IndexController"
}
func (this *UserController) Build(r *lithot.Lithot) {
	r.Handle("GET", "/", this.Index)
	r.Handle("GET", "/users", this.UserList)
	r.Handle("GET", "/users/:id", this.UserDetail)
}
