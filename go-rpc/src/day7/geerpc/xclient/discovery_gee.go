/*
-*- encoding: utf-8 -*-
@File    :   discovery_gee.go
@Time    :   2023/02/24 14:36:35
@Author  :   victor2022
@Version :   1.0
@Desc    :   discovery of geerpc registry
*/
package xclient

import (
	"log"
	"net/http"
	"strings"
	"time"
)

// 基于注册中心的服务发现器结构体
type GeeRegistryDiscovery struct {
	*MultiServersDiscovery               // 本地多任务发现器
	registry               string        // 注册中心地址
	timeout                time.Duration // 服务列表超时时间
	lastUpdate             time.Time     // 上一次更新时间，超时后，需要再次从服务中心获取服务列表
}

const defaultUpdateTimeout = time.Second * 10 // 默认超时时间，10s

// 创建新的服务发现器
func NewGeeRegistryDiscovery(registerAddr string, timeout time.Duration) *GeeRegistryDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTimeout
	}

	d := &GeeRegistryDiscovery{
		MultiServersDiscovery: NewMultiServerDiscovery(make([]string, 0)),
		registry:              registerAddr,
		timeout:               timeout,
	}
	return d
}

// 更新服务列表到本地
func (d *GeeRegistryDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	d.lastUpdate = time.Now()
	return nil
}

// 刷新服务列表
func (d *GeeRegistryDiscovery) Refresh() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.lastUpdate.Add(d.timeout).After(time.Now()) {
		// 若没有到刷新时间，则什么都不做
		return nil
	}
	log.Println("rpc registry: refresh servers from registry", d.registry)
	// 从注册中心获取服务列表
	resp, err := http.Get(d.registry)
	if err != nil {
		log.Println("rpc registry refresh err:", err)
		return err
	}
	// 将字符串分割为服务列表
	servers := strings.Split(resp.Header.Get("X-Geerpc-Servers"), ",")
	// 重新初始化服务列表
	d.servers = make([]string, 0, len(servers))
	// 将从注册中心获取的服务列表放入发现器的服务列表中
	for _, server := range servers {
		if strings.TrimSpace(server) != "" {
			d.servers = append(d.servers, strings.TrimSpace(server))
		}
	}
	d.lastUpdate = time.Now()
	return nil
}

// 按照传入方式获取一个服务地址
func (d *GeeRegistryDiscovery) Get(mode SelectMode) (string, error) {
	// 每次获取服务时，先刷新一下，保证服务列表没有过期
	if err := d.Refresh(); err != nil {
		return "", err
	}
	return d.MultiServersDiscovery.Get(mode)
}

// 获取所有服务地址
func (d *GeeRegistryDiscovery) GetAll() ([]string, error) {
	if err := d.Refresh(); err != nil {
		return nil, err
	}
	return d.MultiServersDiscovery.GetAll()
}
