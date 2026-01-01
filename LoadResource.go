package GoMiniblink

import (
	"io/ioutil"
	url2 "net/url"
	"os"
	"strings"
)

type LoadResource interface {
	Domain() string
	ByUri(uri *url2.URL) []byte
}

// FileLoader 文件资源加载器实现
type FileLoader struct {
	domain string
	dir    string
}

// 初始化文件资源加载器
func (_this *FileLoader) Init(dir, domain string) *FileLoader {
	_this.dir = strings.TrimRight(dir, string(os.PathSeparator))
	_this.domain = strings.ToLower(strings.TrimRight(domain, "/"))
	return _this
}

// 获取文件加载器的域名
func (_this *FileLoader) Domain() string {
	return _this.domain
}

// 根据 URI 加载文件内容
func (_this *FileLoader) ByUri(uri *url2.URL) []byte {
	path := strings.Join([]string{_this.dir, uri.Path}, "")
	path = strings.ReplaceAll(path, "/", string(os.PathSeparator))
	if data, err := ioutil.ReadFile(path); err == nil {
		return data
	}
	return nil
}
