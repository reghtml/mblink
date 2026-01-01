package GoMiniblink

import (
	"image"
	url2 "net/url"
	"reflect"
	"strings"

	cs "github.com/reghtml/mblink/forms/controls"
)

// MiniblinkBrowser 浏览器控件，封装了 Miniblink 核心功能
type MiniblinkBrowser struct {
	cs.Control
	core   Miniblink
	fnlist map[string]reflect.Value

	EvRequestBefore map[string]func(sender *MiniblinkBrowser, e RequestBeforeEvArgs)
	OnRequestBefore func(e RequestBeforeEvArgs)

	EvJsReady map[string]func(sender *MiniblinkBrowser, e JsReadyEvArgs)
	OnJsReady func(e JsReadyEvArgs)

	EvConsole map[string]func(sender *MiniblinkBrowser, e ConsoleEvArgs)
	OnConsole func(e ConsoleEvArgs)

	EvDocumentReady map[string]func(sender *MiniblinkBrowser, e DocumentReadyEvArgs)
	OnDocumentReady func(e DocumentReadyEvArgs)

	EvPaintUpdated map[string]func(sender *MiniblinkBrowser, e PaintUpdatedEvArgs)
	OnPaintUpdated func(e PaintUpdatedEvArgs)

	ResourceLoader []LoadResource
}

// 初始化浏览器控件
func (_this *MiniblinkBrowser) Init() *MiniblinkBrowser {
	_this.Control.Init()
	_this.EvRequestBefore = make(map[string]func(*MiniblinkBrowser, RequestBeforeEvArgs))
	_this.EvJsReady = make(map[string]func(*MiniblinkBrowser, JsReadyEvArgs))
	_this.EvConsole = make(map[string]func(*MiniblinkBrowser, ConsoleEvArgs))
	_this.EvDocumentReady = make(map[string]func(*MiniblinkBrowser, DocumentReadyEvArgs))

	_this.OnRequestBefore = _this.defOnRequestBefore
	_this.OnJsReady = _this.defOnJsReady
	_this.OnConsole = _this.defOnConsole
	_this.OnDocumentReady = _this.defOnDocumentReady
	_this.OnPaintUpdated = _this.defOnPaintUpdated

	_this.EvRequestBefore["__goMiniblink"] = _this.loadRes
	_this.EvDestroy["__goMiniblink"] = _this.onClosed
	_this.core = new(freeMiniblink).init(&_this.Control)
	_this.mbInit()
	return _this
}

func (_this *MiniblinkBrowser) onClosed(_ cs.GUI) {
	destroyWebView(_this.core.GetHandle())
}

func (_this *MiniblinkBrowser) loadRes(_ *MiniblinkBrowser, e RequestBeforeEvArgs) {
	if len(_this.ResourceLoader) == 0 {
		return
	}
	url, err := url2.Parse(e.Url())
	if err != nil {
		return
	}
	host := strings.ToLower(url.Host)
	for i := range _this.ResourceLoader {
		loader := _this.ResourceLoader[i]
		if strings.HasPrefix(strings.ToLower(loader.Domain()), host) == false {
			continue
		}
		data := loader.ByUri(url)
		if data != nil {
			e.SetData(data)
			break
		}
	}
}

func (_this *MiniblinkBrowser) mbInit() {
	_this.core.SetOnRequestBefore(func(args RequestBeforeEvArgs) {
		if _this.OnRequestBefore != nil {
			_this.OnRequestBefore(args)
		}
	})
	_this.core.SetOnJsReady(func(args JsReadyEvArgs) {
		if _this.OnJsReady != nil {
			_this.OnJsReady(args)
		}
	})
	_this.core.SetOnConsole(func(args ConsoleEvArgs) {
		if _this.OnConsole != nil {
			_this.OnConsole(args)
		}
	})
	_this.core.SetOnDocumentReady(func(args DocumentReadyEvArgs) {
		if _this.OnDocumentReady != nil {
			_this.OnDocumentReady(args)
		}
	})
	_this.core.SetOnPaintUpdated(func(args PaintUpdatedEvArgs) {
		if _this.OnPaintUpdated != nil {
			_this.OnPaintUpdated(args)
		}
	})
}

// 设置代理服务器
func (_this *MiniblinkBrowser) SetProxy(info ProxyInfo) {
	_this.core.SetProxy(info)
}

// 加载指定的 URL
func (_this *MiniblinkBrowser) LoadUri(uri string) {
	_this.core.LoadUri(uri)
}

// 加载 HTML 内容并指定基础 URL
func (_this *MiniblinkBrowser) LoadHtmlWithBaseUrl(html, baseUrl string) {
	_this.core.LoadHtmlWithBaseUrl(html, baseUrl)
}

// 绑定 Go 函数为 JavaScript 全局函数
func (_this *MiniblinkBrowser) JsFunc(name string, fn GoFn, state interface{}) {
	_this.core.JsFunc(name, fn, state)
}

// 使用反射绑定任意 Go 函数为 JavaScript 全局函数
func (_this *MiniblinkBrowser) JsFuncEx(name string, fn interface{}) {
	p := reflect.TypeOf(fn)
	if p.Kind() != reflect.Func {
		return
	}
	_this.JsFunc(name, func(ctx GoFnContext) interface{} {
		rt := reflect.TypeOf(ctx.State)
		rv := reflect.ValueOf(ctx.State)
		var args []reflect.Value
		for i := 0; i < rt.NumIn() && i < len(ctx.Param); i++ {
			args = append(args, reflect.ValueOf(ctx.Param[i]))
		}
		rs := rv.Call(args)
		if rt.NumOut() > 0 {
			return rs[0].Interface()
		}
		return nil
	}, fn)
}

// 调用 JavaScript 函数并返回结果
func (_this *MiniblinkBrowser) CallJsFunc(name string, param ...interface{}) interface{} {
	return _this.core.CallJsFunc(name, param)
}

// 将 webview 内容转换为位图
func (_this *MiniblinkBrowser) ToBitmap() *image.RGBA {
	return _this.core.ToBitmap()
}

// 获取底层的 miniblink 句柄
func (_this *MiniblinkBrowser) GetMiniblinkHandle() uintptr {
	return uintptr(_this.core.GetHandle())
}

// 检查鼠标事件是否启用
func (_this *MiniblinkBrowser) MouseIsEnable() bool {
	return _this.core.MouseIsEnable()
}

// 启用或禁用鼠标事件
func (_this *MiniblinkBrowser) MouseEnable(b bool) {
	_this.core.MouseEnable(b)
}

// 设置是否使用位图绘制模式
func (_this *MiniblinkBrowser) SetBmpPaintMode(b bool) {
	_this.core.SetBmpPaintMode(b)
}

// 显示开发者工具
func (_this *MiniblinkBrowser) ShowDevTools(path string) {
	_this.core.ShowDevTools(path)
}

// 设置本地存储的完整路径
func (_this *MiniblinkBrowser) SetLocalStorageFullPath(path string) {
	_this.core.SetLocalStorageFullPath(path)
}

// 设置 HTTP 请求的 User-Agent 字符串
func (_this *MiniblinkBrowser) SetUserAgent(userAgent string) {
	_this.core.SetUserAgent(userAgent)
}
