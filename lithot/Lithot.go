package lithot

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/happylusn/lithot-gin/injector"
	"log"
	"reflect"
	"strings"
	"sync"
)

type Bean interface {
	Name() string
}

var Empty = &struct{}{}
var innerRouter *LithotTree // inner tree node . backup httpmethod and path
var innerRouter_once sync.Once

func getInnerRouter() *LithotTree {
	innerRouter_once.Do(func() {
		innerRouter = NewLithotTree()
	})
	return innerRouter
}

type Lithot struct {
	*gin.Engine
	g            *gin.RouterGroup // 保存 group对象
	exprData     map[string]interface{}
	currentGroup string // temp-var for group string
	errorHandle  ErrorHandle
}

// 404处理
func HandleNotFound(c *gin.Context) {
	Throw("Not Found", 404, c)
}

func NewLithot(ginMiddlewares ...gin.HandlerFunc) *Lithot {
	g := &Lithot{Engine: gin.New(),
		exprData: map[string]interface{}{},
	}
	g.Use(gin.Logger(), ErrorHandler(g)) //强迫加载的异常处理中间件
	for _, handler := range ginMiddlewares {
		g.Use(handler)
	}
	config := InitConfig()
	injector.BeanFactory.Set(g)      // inject self
	injector.BeanFactory.Set(config) // add global into (new)BeanFactory
	injector.BeanFactory.Set(NewGPAUtil())
	if config.Server.Html != "" {
		g.LoadHTMLGlob(config.Server.Html)
	}
	return g
}
func (this *Lithot) Launch() {
	var port int32 = 8080
	if config := injector.BeanFactory.Get((*SysConfig)(nil)); config != nil {
		port = config.(*SysConfig).Server.Port
	}
	this.applyAll()
	getCronTask().Start()
	this.Run(fmt.Sprintf(":%d", port))
}
func (this *Lithot) LaunchWithPort(port int) {

	this.applyAll()
	getCronTask().Start()
	this.Run(fmt.Sprintf(":%d", port))
}

func (this *Lithot) getPath(relativePath string) string {
	g := "/" + this.currentGroup
	if g == "/" {
		g = ""
	}
	g = g + relativePath
	g = strings.Replace(g, "//", "/", -1)
	return g
}
func (this *Lithot) Handle(httpMethod, relativePath string, handler interface{}, middlewares ...Middleware) *Lithot {
	return this.handle(httpMethod, relativePath, handler, middlewares...)
}

func (this *Lithot) GET(relativePath string, handler interface{}, middlewares ...Middleware) *Lithot {
	return this.handle("GET", relativePath, handler, middlewares...)
}

func (this *Lithot) POST(relativePath string, handler interface{}, middlewares ...Middleware) *Lithot {
	return this.handle("POST", relativePath, handler, middlewares...)
}

func (this *Lithot) PUT(relativePath string, handler interface{}, middlewares ...Middleware) *Lithot {
	return this.handle("PUT", relativePath, handler, middlewares...)
}

func (this *Lithot) PATCH(relativePath string, handler interface{}, middlewares ...Middleware) *Lithot {
	return this.handle("PATCH", relativePath, handler, middlewares...)
}

func (this *Lithot) DELETE(relativePath string, handler interface{}, middlewares ...Middleware) *Lithot {
	return this.handle("DELETE", relativePath, handler, middlewares...)
}

func (this *Lithot) handle(httpMethod, relativePath string, handler interface{}, middlewares ...Middleware) *Lithot {
	if h := Convert(handler); h != nil {
		methods := strings.Split(httpMethod, ",")
		for _, f := range middlewares {
			injector.BeanFactory.Apply(f)
		}
		for _, method := range methods {
			getInnerRouter().addRoute(method, this.getPath(relativePath), middlewares) //for future
			this.g.Handle(method, relativePath, h)
		}
	}
	return this
}

// 注册全局中间件
func (this *Lithot) RegisterMiddleware(middlewares ...Middleware) *Lithot {
	for _, m := range middlewares {
		injector.BeanFactory.Set(m)
	}
	getMiddlewareHandler().AddMiddleware(middlewares...)
	return this
}

func (this *Lithot) Beans(beans ...Bean) *Lithot {
	for _, bean := range beans {
		this.exprData[bean.Name()] = bean
		injector.BeanFactory.Set(bean)
	}
	return this
}

func (this *Lithot) Configure(configurations ...interface{}) *Lithot {
	injector.BeanFactory.Config(configurations...)
	return this
}
func (this *Lithot) applyAll() {
	for t, v := range injector.BeanFactory.GetBeanMapper() {
		if t.Elem().Kind() == reflect.Struct {
			injector.BeanFactory.Apply(v.Interface())
		}
	}
}
func (this *Lithot) GetSysConfig() *SysConfig {
	if config := injector.BeanFactory.Get((*SysConfig)(nil)); config != nil {
		return config.(*SysConfig)
	}
	return nil
}

func (this *Lithot) Mount(group string, controllers ...Controller) *Lithot {
	this.g = this.Group(group)
	for _, controller := range controllers {
		this.currentGroup = group
		controller.Build(this)
		//this.beanFactory.inject(class)
		this.Beans(controller)
	}
	return this
}

//0/3 * * * * *  //增加定时任务
func (this *Lithot) Task(cron string, expr interface{}) *Lithot {
	var err error
	if f, ok := expr.(func()); ok {
		_, err = getCronTask().AddFunc(cron, f)
	} else if exp, ok := expr.(Expression); ok {
		_, err = getCronTask().AddFunc(cron, func() {
			_, expErr := ExecExpr(exp, this.exprData)
			if expErr != nil {
				log.Println(expErr)
			}
		})
	}
	if err != nil {
		log.Println(err)
	}
	return this
}

type ErrorHandle func(c *gin.Context, err interface{})

func (this *Lithot) SetErrorHandle(f ErrorHandle) *Lithot {
	this.errorHandle = f
	return this
}

func (this *Lithot) Static(relativePath, root string) *Lithot {
	this.Engine.Static(relativePath, root)
	return this
}

func (this *Lithot) LoadHTMLGlob(pattern string) *Lithot {
	this.Engine.LoadHTMLGlob(pattern)
	return this
}
