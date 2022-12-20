/*
-*- encoding: utf-8 -*-
@File    :   debug.go
@Time    :   2022/11/11 16:36:14
@Author  :   victor2022
@Version :   1.0
@Desc    :   debug page for gee rpc
*/
package geerpc

import (
	"fmt"
	"html/template"
	"net/http"
)

// debug页面展示部分
const debugText = `<html>
	<body>
	<title>GeeRPC Services</title>
	{{range .}}
	<hr>
	Service {{.Name}}
	<hr>
		<table>
		<th align=center>Method</th><th align=center>Calls</th>
		{{range $name, $mtype := .Method}}
			<tr>
			<td align=left font=fixed>{{$name}}({{$mtype.ArgType}}, {{$mtype.ReplyType}}) error</td>
			<td align=center>{{$mtype.NumCalls}}</td>
			</tr>
		{{end}}
		</table>
	{{end}}
	</body>
	</html>`

// 由debugText创建HTML模板
var debug = template.Must(template.New("RPC debug").Parse(debugText))

/*
@Time    :   2022/11/11 16:45:24
@Author  :   victor2022
@Desc    :   debug页面结构体
*/
type debugHTTP struct {
	*Server
}

/*
@Time    :   2022/11/11 16:46:20
@Author  :   victor2022
@Desc    :   debug服务
*/
type debugService struct {
	Name   string
	Method map[string]*methodType
}

// Runs at /debug/geerpc
func (server debugHTTP) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Build a sorted version of the data.
	var services []debugService
	// 将服务列表中的所有服务的状态都映射到debug页面中
	server.serviceMap.Range(func(namei, svci interface{}) bool {
		svc := svci.(*service)
		services = append(services, debugService{
			Name:   namei.(string),
			Method: svc.method,
		})
		return true
	})
	// 使用模板解析页面
	err := debug.Execute(w, services)
	if err != nil {
		_, _ = fmt.Fprintln(w, "rpc: error executing template:", err.Error())
	}
}
