/*
-*- encoding: utf-8 -*-
@File    :   gee.go
@Time    :   2022/10/19 22:30:32
@Author  :   victor2022
@Version :   1.0
@Desc    :   main structs of gee
*/
package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strings"
)

type (
	// 定义HandlerFunc
	HandlerFunc func(*Context)

	// 定义Group结构，实现分组功能
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc // 中间件扩展
		parent      *RouterGroup  // 父Group
		engine      *Engine       // 所有的Group持有一个Engine实例
	}

	// 实现ServeHTTP接口，保存路由信息
	// 将Engine抽象为一个顶层分组
	Engine struct {
		*RouterGroup
		router        *router
		groups        []*RouterGroup     // store all groups
		htmlTemplates *template.Template // for html render
		funcMap       template.FuncMap   // for html render
	}
)

// 创建Engine的方法，工厂模式
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	// 初始元素只有当前主Group
	engine.groups = []*RouterGroup{engine.RouterGroup}
	// 返回生成的Engine
	return engine
}

/*
@Time    :   2022/10/24 10:47:49
@Author  :   victor2022
@Desc    :   创建Engine，默认使用logger和recovery插件
*/
func Default() *Engine {
	engine := New()
	// 设置中间件
	engine.Use(Logger(), Recovery())
	return engine
}

// 向group中添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 添加分组的方法
// 所有的分组都共享一个engine实例
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		// 拼接当前前缀
		prefix: group.prefix + prefix,
		// 父分组
		parent: group,
		// 持有的engine
		engine: engine,
	}
	// 向group组中添加当前group
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 添加路由的方法，属于Group
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	// 拼接pattern
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// 添加Get请求的方法
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	// 调用公用方法添加路由
	group.addRoute("GET", pattern, handler)
}

// 添加POST请求的方法
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 开启服务器
func (engine *Engine) Run(addr string) (err error) {
	// 开始监听端口，提供服务
	res := http.ListenAndServe(addr, engine)
	// 在发生错误时输出日志
	log.Fatal(res)
	// 返回可能的错误
	return res
}

/*
@Time    :   2022/10/23 19:59:26
@Author  :   victor2022
@Desc    :   生成静态资源处理器
*/
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absPath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absPath, http.FileServer(fs))
	log.Printf("File Server - %s -> %v", absPath, fs)
	return func(c *Context) {
		file := c.Param("filepath")
		// 检查是否存在当前文件
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
		c.Status(http.StatusOK)
	}
}

/*
@Time    :   2022/10/23 19:58:00
@Author  :   victor2022
@Desc    :   注册文件服务
root为文件在系统中的路径
*/
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// 注册Get请求
	group.GET(urlPattern, handler)
}

/*
@Time    :   2022/10/23 21:24:48
@Author  :   victor2022
@Desc    :   设置页面渲染函数
*/
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

/*
@Time    :   2022/10/23 21:28:18
@Author  :   victor2022
@Desc    :   加载模板
*/
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

/*
@Time    :   2022/10/23 21:33:54
@Author  :   victor2022
@Desc    :   服务器接口方法
*/
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 获取全局handler(定义在engine中的handler)
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// 创建对应的上下文
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}
