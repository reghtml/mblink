package GoMiniblink

import (
	"image"
	"strconv"

	"github.com/reghtml/mblink/forms"
)

type PaintUpdatedEvArgs interface {
	Bitmap() *image.RGBA
	Bound() forms.Bound
	Cancel()
	IsCancel() bool
}

// freePaintUpdatedEvArgs 绘制更新事件参数的实现
type freePaintUpdatedEvArgs struct {
	bitmap *image.RGBA
	bound  forms.Bound
	cancel bool
}

// 初始化绘制更新事件参数
func (_this *freePaintUpdatedEvArgs) init(bitmap *image.RGBA, bound forms.Bound) *freePaintUpdatedEvArgs {
	_this.bitmap = bitmap
	_this.bound = bound
	return _this
}

// 获取位图数据
func (_this *freePaintUpdatedEvArgs) Bitmap() *image.RGBA {
	return _this.bitmap
}

// 获取绘制区域边界
func (_this *freePaintUpdatedEvArgs) Bound() forms.Bound {
	return _this.bound
}

// 检查是否已取消
func (_this *freePaintUpdatedEvArgs) IsCancel() bool {
	return _this.cancel
}

// 取消绘制
func (_this *freePaintUpdatedEvArgs) Cancel() {
	_this.cancel = true
}

type DocumentReadyEvArgs interface {
	FrameContext
}

// freeDocumentReadyEvArgs 文档就绪事件参数的实现
type freeDocumentReadyEvArgs struct {
	*freeFrameContext
}

// 初始化文档就绪事件参数
func (_this *freeDocumentReadyEvArgs) init(mb Miniblink, frame wkeFrame) *freeDocumentReadyEvArgs {
	_this.freeFrameContext = new(freeFrameContext).init(mb, frame)
	return _this
}

type FrameContext interface {
	FrameId() uintptr
	IsMain() bool
	Url() string
	IsRemote() bool
	RunJs(script string) interface{}
}

// freeFrameContext 框架上下文的实现
type freeFrameContext struct {
	id       uintptr
	isMain   bool
	url      string
	isRemote bool
	core     Miniblink
}

// 初始化框架上下文
func (_this *freeFrameContext) init(mb Miniblink, frame wkeFrame) *freeFrameContext {
	_this.core = mb
	_this.id = uintptr(frame)
	_this.isMain = mbApi.wkeIsMainFrame(_this.core.GetHandle(), frame)
	_this.isRemote = mbApi.wkeIsWebRemoteFrame(_this.core.GetHandle(), frame)
	_this.url = mbApi.wkeGetFrameUrl(_this.core.GetHandle(), frame)
	return _this
}

// 在框架中执行 JavaScript 代码
func (_this *freeFrameContext) RunJs(script string) interface{} {
	if len(script) > 0 {
		es := mbApi.wkeGetGlobalExecByFrame(_this.core.GetHandle(), wkeFrame(_this.id))
		rs := mbApi.jsEval(es, script)
		return toGoValue(_this.core, es, rs)
	}
	return nil
}

// 判断框架是否为远程框架
func (_this *freeFrameContext) IsRemote() bool {
	return _this.isRemote
}

// 获取框架的 URL
func (_this *freeFrameContext) Url() string {
	return _this.url
}

// 判断是否为主框架
func (_this *freeFrameContext) IsMain() bool {
	return _this.isMain
}

// 获取框架 ID
func (_this *freeFrameContext) FrameId() uintptr {
	return _this.id
}

type ConsoleEvArgs interface {
	Level() string
	Message() string
	SourceName() string
	SourceLine() int
	StackTrace() string
}

// freeConsoleMessageEvArgs 控制台消息事件参数的实现
type freeConsoleMessageEvArgs struct {
	level   string
	message string
	name    string
	line    int
	stack   string
}

// 初始化控制台消息事件参数
func (_this *freeConsoleMessageEvArgs) init() *freeConsoleMessageEvArgs {
	return _this
}

// 获取控制台消息级别
func (_this *freeConsoleMessageEvArgs) Level() string {
	return _this.level
}

// 获取控制台消息内容
func (_this *freeConsoleMessageEvArgs) Message() string {
	return _this.message
}

// 获取消息源文件名
func (_this *freeConsoleMessageEvArgs) SourceName() string {
	return _this.name
}

// 获取消息源行号
func (_this *freeConsoleMessageEvArgs) SourceLine() int {
	return _this.line
}

// 获取堆栈跟踪信息
func (_this *freeConsoleMessageEvArgs) StackTrace() string {
	return _this.stack
}

type JsReadyEvArgs interface {
	FrameContext
}

// wkeJsReadyEvArgs JavaScript 就绪事件参数的实现
type wkeJsReadyEvArgs struct {
	*freeFrameContext
}

// 初始化 JavaScript 就绪事件参数
func (_this *wkeJsReadyEvArgs) init(mb Miniblink, frame wkeFrame) *wkeJsReadyEvArgs {
	_this.freeFrameContext = new(freeFrameContext).init(mb, frame)
	return _this
}

type ResponseEvArgs interface {
	RequestBefore() RequestBeforeEvArgs
	ContentType() string
	SetContentType(contentType string)
	Data() []byte
	SetData(data []byte)
	Headers() map[string]string
}

// freeResponseEvArgs 响应事件参数的实现
type freeResponseEvArgs struct {
	_req  *freeRequestBeforeEvArgs
	_data []byte
}

// 初始化响应事件参数
func (_this *freeResponseEvArgs) init(request *freeRequestBeforeEvArgs, data []byte) *freeResponseEvArgs {
	_this._req = request
	_this._data = data
	return _this
}

// 获取响应头信息
func (_this *freeResponseEvArgs) Headers() map[string]string {
	return mbApi.wkeNetGetRawResponseHead(_this._req._job)
}

// 设置响应数据
func (_this *freeResponseEvArgs) SetData(data []byte) {
	_this._data = data
	mbApi.wkeNetSetData(_this._req._job, _this._data)
}

// 获取响应数据
func (_this *freeResponseEvArgs) Data() []byte {
	return _this._data
}

// 获取请求前事件参数
func (_this *freeResponseEvArgs) RequestBefore() RequestBeforeEvArgs {
	return _this._req
}

// 获取响应内容类型
func (_this *freeResponseEvArgs) ContentType() string {
	return mbApi.wkeNetGetMIMEType(_this._req._job)
}

// 设置响应内容类型
func (_this *freeResponseEvArgs) SetContentType(contentType string) {
	mbApi.wkeNetSetMIMEType(_this._req._job, contentType)
}

type LoadFailEvArgs interface {
	RequestBefore() RequestBeforeEvArgs
}

// freeLoadFailEvArgs 加载失败事件参数的实现
type freeLoadFailEvArgs struct {
	_req *freeRequestBeforeEvArgs
}

// 初始化加载失败事件参数
func (_this *freeLoadFailEvArgs) init(request *freeRequestBeforeEvArgs) *freeLoadFailEvArgs {
	_this._req = request
	return _this
}

// 获取请求前事件参数
func (_this *freeLoadFailEvArgs) RequestBefore() RequestBeforeEvArgs {
	return _this._req
}

type RequestBeforeEvArgs interface {
	Url() string
	Method() string
	SetData([]byte)
	Data() []byte
	SetCancel(b bool)
	ResetUrl(url string)
	SetHeader(name, value string)
	/**
	内容最终呈现时触发
	args:intf, ResponseEvArgs
	*/
	EvResponse() *EventDispatcher
	/**
	加载失败时触发
	args:intf, LoadFailEvArgs
	*/
	EvLoadFail() *EventDispatcher
	/**
	请求流程全部完成时触发
	args:intf, RequestBeforeEvArgs
	*/
	EvFinish() *EventDispatcher
}

// freeRequestBeforeEvArgs 请求前事件参数的实现
type freeRequestBeforeEvArgs struct {
	_wke    Miniblink
	_job    wkeNetJob
	_url    string
	_cancel bool
	_data   []byte
	//1=发送之前,2=异步处理,3=已发送,4=收到真实数据,5=完成
	_state         int
	_evResponseKey string
	_evResponse    *EventDispatcher
	_evLoadFailKey string
	_evLoadFail    *EventDispatcher
	_evFinishKey   string
	_evFinish      *EventDispatcher
}

// 初始化请求前事件参数
func (_this *freeRequestBeforeEvArgs) init(wke Miniblink, job wkeNetJob) *freeRequestBeforeEvArgs {
	_this._wke = wke
	_this._url = mbApi.wkeNetGetUrlByJob(job)
	_this._job = job
	_this._state = 1
	_this._evResponseKey = "evResp" + strconv.FormatUint(uint64(job), 10)
	_this._evResponse = new(EventDispatcher).Init(_this._evResponseKey)
	_this._evLoadFailKey = "evFail" + strconv.FormatUint(uint64(job), 10)
	_this._evLoadFail = new(EventDispatcher).Init(_this._evLoadFailKey)
	_this._evFinishKey = "evFsh" + strconv.FormatUint(uint64(job), 10)
	_this._evFinish = new(EventDispatcher).Init(_this._evFinishKey)
	return _this
}

// 获取请求完成事件分发器
func (_this *freeRequestBeforeEvArgs) EvFinish() *EventDispatcher {
	return _this._evFinish
}

// 获取加载失败事件分发器
func (_this *freeRequestBeforeEvArgs) EvLoadFail() *EventDispatcher {
	return _this._evLoadFail
}

// 获取响应事件分发器
func (_this *freeRequestBeforeEvArgs) EvResponse() *EventDispatcher {
	return _this._evResponse
}

// 修改请求的 URL
func (_this *freeRequestBeforeEvArgs) ResetUrl(url string) {
	mbApi.wkeNetChangeRequestUrl(_this._job, url)
	_this._url = url
}

// 设置请求头字段
func (_this *freeRequestBeforeEvArgs) SetHeader(name, value string) {
	mbApi.wkeNetSetHTTPHeaderField(_this._job, name, value)
}

// 设置请求数据
func (_this *freeRequestBeforeEvArgs) SetData(data []byte) {
	_this._data = data
}

// 获取请求数据
func (_this *freeRequestBeforeEvArgs) Data() []byte {
	return _this._data
}

// 获取请求方法
func (_this *freeRequestBeforeEvArgs) Method() string {
	t := mbApi.wkeNetGetRequestMethod(_this._job)
	switch t {
	case wkeRequestType_Get:
		return "GET"
	case wkeRequestType_Post:
		return "POST"
	case wkeRequestType_Put:
		return "PUT"
	default:
		return "UNKNOW"
	}
}

// 获取请求的 URL
func (_this *freeRequestBeforeEvArgs) Url() string {
	return _this._url
}

// 设置是否取消请求
func (_this *freeRequestBeforeEvArgs) SetCancel(b bool) {
	_this._cancel = b
}

func (_this *freeRequestBeforeEvArgs) onBegin() {
	if _this._state == 2 {
		return
	}
	if _this._state == 1 && _this._data != nil {
		mbApi.wkeNetSetData(_this._job, _this._data)
		_this._cancel = true
		_this._state = 5
	} else if _this._evResponse.IsEmtpy() == false {
		mbApi.wkeNetHookRequest(_this._job)
		_this._cancel = false
	}
	if _this._cancel {
		mbApi.wkeNetCancelRequest(_this._job)
		if _this._data != nil {
			_this.onResponse(_this._data)
		} else {
			_this._evFinish.Fire(_this._evFinishKey, _this, _this)
		}
	} else {
		_this._state = 3
	}
}

func (_this *freeRequestBeforeEvArgs) onResponse(data []byte) {
	_this._state = 5
	args := new(freeResponseEvArgs).init(_this, data)
	_this._evResponse.Fire(_this._evResponseKey, _this, args)
	_this._evFinish.Fire(_this._evFinishKey, _this, _this)
}

func (_this *freeRequestBeforeEvArgs) onFail() {
	_this._state = 4
	args := new(freeLoadFailEvArgs).init(_this)
	_this._evLoadFail.Fire(_this._evLoadFailKey, _this, args)
	if _this._evResponse.IsEmtpy() {
		_this._evFinish.Fire(_this._evFinishKey, _this, _this)
	}
}

// 默认请求前事件处理
func (_this *MiniblinkBrowser) defOnRequestBefore(e RequestBeforeEvArgs) {
	for _, v := range _this.EvRequestBefore {
		v(_this, e)
	}
}

// 默认 JavaScript 就绪事件处理
func (_this *MiniblinkBrowser) defOnJsReady(e JsReadyEvArgs) {
	for _, v := range _this.EvJsReady {
		v(_this, e)
	}
}

// 默认控制台事件处理
func (_this *MiniblinkBrowser) defOnConsole(e ConsoleEvArgs) {
	for _, v := range _this.EvConsole {
		v(_this, e)
	}
}

// 默认文档就绪事件处理
func (_this *MiniblinkBrowser) defOnDocumentReady(e DocumentReadyEvArgs) {
	for _, v := range _this.EvDocumentReady {
		v(_this, e)
	}
}

// 默认绘制更新事件处理
func (_this *MiniblinkBrowser) defOnPaintUpdated(e PaintUpdatedEvArgs) {
	for _, v := range _this.EvPaintUpdated {
		v(_this, e)
	}
}
