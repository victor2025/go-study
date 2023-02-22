/*
-*- encoding: utf-8 -*-
@File    :   discovery.go
@Time    :   2022/12/03 14:39:21
@Author  :   victor2022
@Version :   1.0
@Desc    :   service discovery module
*/
package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect     SelectMode = iota // random select
	RoundRobinSelect                   // round robin
)

/*
@Time    :   2022/12/01 22:57:09
@Author  :   victor2022
@Desc    :   发现者
*/
type Discovery interface {
	Refresh() error                      // 从注册中心重新拉取服务列表
	Update(server []string) error        // 手动更新服务列表
	Get(mode SelectMode) (string, error) // 通过指定的方式从列表中获取调用地址
	GetAll() ([]string, error)           // 获取服务列表
}

/*
@Time    :   2022/12/01 23:03:49
@Author  :   victor2022
@Desc    :   基于本地服务列表的服务选择器
*/
type MultiServersDiscovery struct {
	r       *rand.Rand   // 生成随机序列号
	mu      sync.RWMutex // 保证线程安全
	servers []string
	index   int // 用来记录已经轮询到的位置
}

/*
@Time    :   2022/12/01 23:07:11
@Author  :   victor2022
@Desc    :   创建新的服务发现者
*/
func NewMultiServerDiscovery(servers []string) *MultiServersDiscovery {
	d := &MultiServersDiscovery{
		servers: servers,
		r:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.index = d.r.Intn(math.MaxInt32 - 1)
	return d
}

/*
@Time    :   2022/12/03 14:29:02
@Author  :   victor2022
@Desc    :   本地多服务发现器无需刷新
*/
func (d *MultiServersDiscovery) Refresh() error {
	return nil
}

/*
@Time    :   2022/12/03 14:30:11
@Author  :   victor2022
@Desc    :   更新服务列表
*/
func (d *MultiServersDiscovery) Update(servers []string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

/*
@Time    :   2022/12/03 14:30:49
@Author  :   victor2022
@Desc    :   通过传入的选择方式选择一个服务
*/
func (d *MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// 若服务列表为空，则直接返回错误
	n := len(d.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}
	// 根据不同模式选择服务
	switch mode {
	case RandomSelect:
		// 随机选择服务
		return d.servers[d.r.Intn(n)], nil
	case RoundRobinSelect:
		// 轮询选择服务
		s := d.servers[d.index%n] // 通过取余的方式来处理服务列表更新之后，服务数目变化的问题
		d.index = (d.index + 1) % n
		return s, nil
	default:
		return "", errors.New("rpc discovery: not supported select mode")
	}
}

/*
@Time    :   2022/12/03 14:37:26
@Author  :   victor2022
@Desc    :   返回发现器中所有的服务
*/
func (d *MultiServersDiscovery) GetAll() ([]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	// 返回服务列表的副本
	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

var _ Discovery = (*MultiServersDiscovery)(nil)
