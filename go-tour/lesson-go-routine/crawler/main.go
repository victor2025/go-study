package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch 返回 URL 的 body 内容，并且将在这个页面上找到的 URL 放到一个 slice 中。
	Fetch(url string) (body string, urls []string, err error)
}

// 创建标记器
var marker = Marker{
	vis: make(map[string]int),
	mut: sync.Mutex{},
}

// Crawl 使用 fetcher 从某个 URL 开始递归的爬取页面，直到达到最大深度。
func Crawl(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup) {
	// 等待解锁当前routine
	defer wg.Done()
	// TODO: 并行的抓取 URL。
	// TODO: 不重复抓取页面。
	// 检查深度是否达到要求
	if depth <= 0 {
		return
	}
	// 检查标记，若已经被标记过，那么直接返回
	if !marker.Mark(url) {
		return
	}
	// 抓取当前页面内容以及下一层内容
	body, urls, err := fetcher.Fetch(url)
	// 判断是否抓取成功
	if err != nil {
		fmt.Println(err)
		return
	}
	// 打印数据
	fmt.Printf("found: %s %q\n", url, body)
	// 递归抓取接下来的页面，此处可以使用并行处理
	for _, u := range urls {
		// 跳过已经标记的url
		if marker.IsMarked(u) {
			continue
		}
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, wg)
	}
	return
}

func main() {
	// 记录运行中go-routine的数目
	wg := sync.WaitGroup{}
	// 开始爬取数据
	// 每次启动一个routine之前标记一次
	wg.Add(1)
	// 在线程执行完成后去除标记
	go Crawl("https://golang.org/", 4, fetcher, &wg)
	// 等待完成
	wg.Wait()
}

// fakeFetcher 是返回若干结果的 Fetcher。
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

// 抓取方法的实现
func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher 是填充后的 fakeFetcher。
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

// url标记器
type Marker struct {
	vis map[string]int
	mut sync.Mutex
}

// 标记url
func (m *Marker) Mark(url string) bool {
	m.mut.Lock()
	if m.vis[url] == 1 {
		return false
	}
	m.vis[url] = 1
	defer m.mut.Unlock()
	return true
}

// 获取访问状态
func (m *Marker) IsMarked(url string) bool {
	m.mut.Lock()
	defer m.mut.Unlock()
	return m.vis[url] == 1
}
