package GoMiniblink

import (
	"image"
)

// ResponseCallback 响应回调函数类型
type ResponseCallback func(args ResponseEvArgs)

// RequestBeforeCallback 请求前回调函数类型
type RequestBeforeCallback func(args RequestBeforeEvArgs)

// JsReadyCallback JavaScript 就绪回调函数类型
type JsReadyCallback func(args JsReadyEvArgs)

// ConsoleCallback 控制台消息回调函数类型
type ConsoleCallback func(args ConsoleEvArgs)

// DocumentReadyCallback 文档就绪回调函数类型
type DocumentReadyCallback func(args DocumentReadyEvArgs)

// PaintUpdatedCallback 绘制更新回调函数类型
type PaintUpdatedCallback func(args PaintUpdatedEvArgs)

// Miniblink Miniblink 浏览器接口
type Miniblink interface {
	SetBmpPaintMode(b bool)
	SetProxy(info ProxyInfo)
	MouseIsEnable() bool
	MouseEnable(b bool)
	ToBitmap() *image.RGBA
	CallJsFunc(name string, param []interface{}) interface{}
	JsFunc(name string, fn GoFn, state interface{})
	RunJs(script string) interface{}
	SetOnConsole(callback ConsoleCallback)
	SetOnJsReady(callback JsReadyCallback)
	SetOnRequestBefore(callback RequestBeforeCallback)
	SetOnDocumentReady(callback DocumentReadyCallback)
	SetOnPaintUpdated(callback PaintUpdatedCallback)
	LoadUri(uri string)
	LoadHtmlWithBaseUrl(html, baseUrl string)
	GetHandle() wkeHandle
	ShowDevTools(path string)
	SetLocalStorageFullPath(path string)
	SetUserAgent(userAgent string)
}
