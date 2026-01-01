package GoMiniblink

import (
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// wkeProxy 代理服务器配置结构体
type wkeProxy struct {
	Type     int32
	HostName [100]byte
	Port     uint16
	UserName [50]byte
	Password [50]byte
}

// 将两个32位整数合并为一个64位整数
func _toInt64(low, high int32) int64 {
	var l = int64(high)<<32 + int64(low)
	return *(&l)
}

// 将两个指针值合并为一个 jsValue
func _toJsValue(low, high uintptr) jsValue {
	return jsValue(uintptr(_toInt64(int32(low), int32(high))))
}

// 将 jsValue 拆分为两个32位整数（用于32位系统）
func _toLH(value jsValue) (low, high int32) {
	if is64 {
		return 0, 0
	}
	return int32(value), int32(int64(value) >> 32 & 0xffffffff)
}

var is64 bool

// winFreeApi Windows 平台下的 Miniblink API 实现
type winFreeApi struct {
	_dll *windows.LazyDLL

	_wkeInitialize               *windows.LazyProc
	_wkeCreateWebView            *windows.LazyProc
	_wkeSetHandle                *windows.LazyProc
	_wkeOnPaintBitUpdated        *windows.LazyProc
	_wkeLoadURL                  *windows.LazyProc
	_wkeLoadHtmlWithBaseUrl      *windows.LazyProc
	_wkeResize                   *windows.LazyProc
	_wkeNetOnResponse            *windows.LazyProc
	_wkeOnLoadUrlBegin           *windows.LazyProc
	_wkePaint                    *windows.LazyProc
	_wkeGetWidth                 *windows.LazyProc
	_wkeGetHeight                *windows.LazyProc
	_wkeFireMouseEvent           *windows.LazyProc
	_wkeFireMouseWheelEvent      *windows.LazyProc
	_wkeGetCursorInfoType        *windows.LazyProc
	_wkeFireKeyUpEvent           *windows.LazyProc
	_wkeFireKeyDownEvent         *windows.LazyProc
	_wkeFireKeyPressEvent        *windows.LazyProc
	_wkeGetCaretRect             *windows.LazyProc
	_wkeSetFocus                 *windows.LazyProc
	_wkeNetGetRequestMethod      *windows.LazyProc
	_wkeNetSetData               *windows.LazyProc
	_wkeNetCancelRequest         *windows.LazyProc
	_wkeJsBindFunction           *windows.LazyProc
	_jsArgCount                  *windows.LazyProc
	_jsArg                       *windows.LazyProc
	_jsTypeOf                    *windows.LazyProc
	_jsToTempString              *windows.LazyProc
	_jsToDoubleString            *windows.LazyProc
	_jsToInt                     *windows.LazyProc
	_jsToBoolean                 *windows.LazyProc
	_jsGetLength                 *windows.LazyProc
	_jsGetAt                     *windows.LazyProc
	_jsGetKeys                   *windows.LazyProc
	_jsGet                       *windows.LazyProc
	_jsSetGlobal                 *windows.LazyProc
	_jsGetGlobal                 *windows.LazyProc
	_wkeGlobalExec               *windows.LazyProc
	_jsCall                      *windows.LazyProc
	_jsUndefined                 *windows.LazyProc
	_jsInt                       *windows.LazyProc
	_jsBoolean                   *windows.LazyProc
	_jsDoubleString              *windows.LazyProc
	_jsString                    *windows.LazyProc
	_jsEmptyArray                *windows.LazyProc
	_jsSetLength                 *windows.LazyProc
	_jsSetAt                     *windows.LazyProc
	_jsFunction                  *windows.LazyProc
	_jsEmptyObject               *windows.LazyProc
	_jsSet                       *windows.LazyProc
	_wkeDestroyWebView           *windows.LazyProc
	_jsGetWebView                *windows.LazyProc
	_wkeKillFocus                *windows.LazyProc
	_jsEval                      *windows.LazyProc
	_wkeOnDidCreateScriptContext *windows.LazyProc
	_wkeIsMainFrame              *windows.LazyProc
	_wkeGetFrameUrl              *windows.LazyProc
	_wkeIsWebRemoteFrame         *windows.LazyProc
	_wkeGetGlobalExecByFrame     *windows.LazyProc
	_wkeOnConsole                *windows.LazyProc
	_wkeGetString                *windows.LazyProc
	_wkeNetSetHTTPHeaderField    *windows.LazyProc
	_wkeNetChangeRequestUrl      *windows.LazyProc
	_wkeNetHookRequest           *windows.LazyProc
	_wkeNetHoldJobToAsynCommit   *windows.LazyProc
	_wkeNetContinueJob           *windows.LazyProc
	_wkeOnLoadUrlEnd             *windows.LazyProc
	_wkeOnLoadUrlFail            *windows.LazyProc
	_wkeNetGetUrlByJob           *windows.LazyProc
	_wkeNetGetMIMEType           *windows.LazyProc
	_wkeNetSetMIMEType           *windows.LazyProc
	_wkeNetGetRawResponseHead    *windows.LazyProc
	_wkeOnDocumentReady2         *windows.LazyProc
	_wkeSetTransparent           *windows.LazyProc
	_wkeSetViewProxy             *windows.LazyProc
	_wkeGetViewDC                *windows.LazyProc
	_wkeSetDebugConfig           *windows.LazyProc
	_wkeSetLocalStorageFullPath  *windows.LazyProc
	_wkeSetUserAgent             *windows.LazyProc
}

// 初始化 Windows DLL 绑定并加载所有函数指针
func (_this *winFreeApi) init() *winFreeApi {
	is64 = unsafe.Sizeof(uintptr(0)) == 8
	var lib *windows.LazyDLL

	// 优先检查 C:\mb\mb.db 是否存在
	dllName := "mb.db"
	if !is64 {
		dllName = "mb86.db"
	}
	customDllPath := filepath.Join("C:", "mb", dllName)
	if _, err := os.Stat(customDllPath); err == nil {
		// 如果 C:\mb\mb.db 存在，使用完整路径
		lib = windows.NewLazyDLL(customDllPath)
	} else {
		// 否则使用相对路径（会在exe目录搜索）
		lib = windows.NewLazyDLL(dllName)
	}
	_this._wkeSetViewProxy = lib.NewProc("wkeSetViewProxy")
	_this._wkeSetTransparent = lib.NewProc("wkeSetTransparent")
	_this._wkeOnDocumentReady2 = lib.NewProc("wkeOnDocumentReady2")
	_this._wkeNetGetRawResponseHead = lib.NewProc("wkeNetGetRawResponseHead")
	_this._wkeNetSetMIMEType = lib.NewProc("wkeNetSetMIMEType")
	_this._wkeNetGetMIMEType = lib.NewProc("wkeNetGetMIMEType")
	_this._wkeNetGetUrlByJob = lib.NewProc("wkeNetGetUrlByJob")
	_this._wkeOnLoadUrlFail = lib.NewProc("wkeOnLoadUrlFail")
	_this._wkeOnLoadUrlEnd = lib.NewProc("wkeOnLoadUrlEnd")
	_this._wkeNetContinueJob = lib.NewProc("wkeNetContinueJob")
	_this._wkeNetHoldJobToAsynCommit = lib.NewProc("wkeNetHoldJobToAsynCommit")
	_this._wkeNetHookRequest = lib.NewProc("wkeNetHookRequest")
	_this._wkeNetChangeRequestUrl = lib.NewProc("wkeNetChangeRequestUrl")
	_this._wkeNetSetHTTPHeaderField = lib.NewProc("wkeNetSetHTTPHeaderField")
	_this._wkeGetString = lib.NewProc("wkeGetString")
	_this._wkeOnConsole = lib.NewProc("wkeOnConsole")
	_this._wkeGetGlobalExecByFrame = lib.NewProc("wkeGetGlobalExecByFrame")
	_this._wkeIsWebRemoteFrame = lib.NewProc("wkeIsWebRemoteFrame")
	_this._wkeGetFrameUrl = lib.NewProc("wkeGetFrameUrl")
	_this._wkeIsMainFrame = lib.NewProc("wkeIsMainFrame")
	_this._wkeOnDidCreateScriptContext = lib.NewProc("wkeOnDidCreateScriptContext")
	_this._jsEval = lib.NewProc("jsEval")
	_this._wkeKillFocus = lib.NewProc("wkeKillFocus")
	_this._jsToInt = lib.NewProc("jsToInt")
	_this._jsSet = lib.NewProc("jsSet")
	_this._jsEmptyObject = lib.NewProc("jsEmptyObject")
	_this._jsFunction = lib.NewProc("jsFunction")
	_this._jsSetAt = lib.NewProc("jsSetAt")
	_this._jsSetLength = lib.NewProc("jsSetLength")
	_this._jsEmptyArray = lib.NewProc("jsEmptyArray")
	_this._jsString = lib.NewProc("jsString")
	_this._jsDoubleString = lib.NewProc("jsDoubleString")
	_this._jsBoolean = lib.NewProc("jsBoolean")
	_this._jsInt = lib.NewProc("jsInt")
	_this._jsUndefined = lib.NewProc("jsUndefined")
	_this._jsCall = lib.NewProc("jsCall")
	_this._wkeGlobalExec = lib.NewProc("wkeGlobalExec")
	_this._jsGetGlobal = lib.NewProc("jsGetGlobal")
	_this._jsSetGlobal = lib.NewProc("jsSetGlobal")
	_this._jsGet = lib.NewProc("jsGet")
	_this._jsGetKeys = lib.NewProc("jsGetKeys")
	_this._jsGetAt = lib.NewProc("jsGetAt")
	_this._jsGetLength = lib.NewProc("jsGetLength")
	_this._jsToBoolean = lib.NewProc("jsToBoolean")
	_this._jsToDoubleString = lib.NewProc("jsToDoubleString")
	_this._jsToTempString = lib.NewProc("jsToTempString")
	_this._jsTypeOf = lib.NewProc("jsTypeOf")
	_this._jsArg = lib.NewProc("jsArg")
	_this._jsArgCount = lib.NewProc("jsArgCount")
	_this._wkeJsBindFunction = lib.NewProc("wkeJsBindFunction")
	_this._wkeNetCancelRequest = lib.NewProc("wkeNetCancelRequest")
	_this._wkeNetSetData = lib.NewProc("wkeNetSetData")
	_this._wkeNetGetRequestMethod = lib.NewProc("wkeNetGetRequestMethod")
	_this._wkeFireKeyPressEvent = lib.NewProc("wkeFireKeyPressEvent")
	_this._wkeFireKeyUpEvent = lib.NewProc("wkeFireKeyUpEvent")
	_this._wkeFireKeyDownEvent = lib.NewProc("wkeFireKeyDownEvent")
	_this._wkeGetCursorInfoType = lib.NewProc("wkeGetCursorInfoType")
	_this._wkeFireMouseWheelEvent = lib.NewProc("wkeFireMouseWheelEvent")
	_this._wkeFireMouseEvent = lib.NewProc("wkeFireMouseEvent")
	_this._wkeGetHeight = lib.NewProc("wkeGetHeight")
	_this._wkeGetWidth = lib.NewProc("wkeGetWidth")
	_this._wkePaint = lib.NewProc("wkePaint")
	_this._wkeOnLoadUrlBegin = lib.NewProc("wkeOnLoadUrlBegin")
	_this._wkeNetOnResponse = lib.NewProc("wkeNetOnResponse")
	_this._wkeLoadURL = lib.NewProc("wkeLoadURL")
	_this._wkeLoadHtmlWithBaseUrl = lib.NewProc("wkeLoadHtmlWithBaseUrl")
	_this._wkeResize = lib.NewProc("wkeResize")
	_this._wkeOnPaintBitUpdated = lib.NewProc("wkeOnPaintBitUpdated")
	_this._wkeSetHandle = lib.NewProc("wkeSetHandle")
	_this._wkeCreateWebView = lib.NewProc("wkeCreateWebView")
	_this._wkeInitialize = lib.NewProc("wkeInitialize")
	_this._wkeGetCaretRect = lib.NewProc("wkeGetCaretRect2")
	_this._wkeSetFocus = lib.NewProc("wkeSetFocus")
	_this._wkeDestroyWebView = lib.NewProc("wkeDestroyWebView")
	_this._jsGetWebView = lib.NewProc("jsGetWebView")
	_this._wkeGetViewDC = lib.NewProc("wkeGetViewDC")
	_this._wkeSetDebugConfig = lib.NewProc("wkeSetDebugConfig")
	_this._wkeSetLocalStorageFullPath = lib.NewProc("wkeSetLocalStorageFullPath")
	_this._wkeSetUserAgent = lib.NewProc("wkeSetUserAgent")

	_this._wkeInitialize.Call()
	return _this
}

// 获取 webview 的设备上下文句柄
func (_this *winFreeApi) wkeGetViewDC(wke wkeHandle) uintptr {
	r, _, _ := _this._wkeGetViewDC.Call(uintptr(wke))
	return r
}

// 为 webview 设置代理服务器
func (_this *winFreeApi) wkeSetViewProxy(wke wkeHandle, proxy ProxyInfo) {
	px := wkeProxy{
		Type: int32(proxy.Type),
		Port: uint16(proxy.Port),
	}
	for i, c := range proxy.HostName {
		px.HostName[i] = byte(c)
	}
	if proxy.UserName != "" {
		for i, c := range proxy.UserName {
			px.UserName[i] = byte(c)
		}
	}
	if proxy.Password != "" {
		for i, c := range proxy.Password {
			px.Password[i] = byte(c)
		}
	}
	_this._wkeSetViewProxy.Call(uintptr(wke), uintptr(unsafe.Pointer(&px)))
}

// 设置 webview 背景是否透明
func (_this *winFreeApi) wkeSetTransparent(wke wkeHandle, enable bool) {
	_this._wkeSetTransparent.Call(uintptr(wke), uintptr(toBool(enable)))
}

// 设置文档加载完成时的回调函数
func (_this *winFreeApi) wkeOnDocumentReady(wke wkeHandle, callback wkeDocumentReady2Callback, param uintptr) {
	_this._wkeOnDocumentReady2.Call(uintptr(wke), syscall.NewCallbackCDecl(callback), param)
}

// 获取网络请求的原始响应头
func (_this *winFreeApi) wkeNetGetRawResponseHead(job wkeNetJob) map[string]string {
	r, _, _ := _this._wkeNetGetRawResponseHead.Call(uintptr(job))
	var list []string
	slist := *((*wkeSlist)(unsafe.Pointer(r)))
	for slist.str != 0 {
		list = append(list, ptrToUtf8(slist.str))
		if slist.next == 0 {
			break
		} else {
			slist = *((*wkeSlist)(unsafe.Pointer(slist.next)))
		}
	}
	hMap := make(map[string]string)
	for i := 0; i < len(list); i += 2 {
		hMap[list[i]] = list[i+1]
	}
	return hMap
}

// 设置网络响应的 MIME 类型
func (_this *winFreeApi) wkeNetSetMIMEType(job wkeNetJob, mime string) {
	p := toCallStr(mime)
	_this._wkeNetSetMIMEType.Call(uintptr(job), uintptr(unsafe.Pointer(&p[0])))
}

// 获取网络响应的 MIME 类型
func (_this *winFreeApi) wkeNetGetMIMEType(job wkeNetJob) string {
	r, _, _ := _this._wkeNetGetMIMEType.Call(uintptr(job))
	return ptrToUtf8(r)
}

// 通过网络任务获取请求的 URL
func (_this *winFreeApi) wkeNetGetUrlByJob(job wkeNetJob) string {
	r, _, _ := _this._wkeNetGetUrlByJob.Call(uintptr(job))
	return ptrToUtf8(r)
}

// 设置 URL 加载失败时的回调函数
func (_this *winFreeApi) wkeOnLoadUrlFail(wke wkeHandle, callback wkeLoadUrlFailCallback, param uintptr) {
	_this._wkeOnLoadUrlFail.Call(uintptr(wke), syscall.NewCallbackCDecl(callback), param)
}

// 设置 URL 加载完成时的回调函数
func (_this *winFreeApi) wkeOnLoadUrlEnd(wke wkeHandle, callback wkeLoadUrlEndCallback, param uintptr) {
	_this._wkeOnLoadUrlEnd.Call(uintptr(wke), syscall.NewCallbackCDecl(callback), param)
}

// 继续执行被挂起的网络任务
func (_this *winFreeApi) wkeNetContinueJob(job wkeNetJob) {
	_this._wkeNetContinueJob.Call(uintptr(job))
}

// 挂起网络任务以便异步提交
func (_this *winFreeApi) wkeNetHoldJobToAsynCommit(job wkeNetJob) {
	_this._wkeNetHoldJobToAsynCommit.Call(uintptr(job))
}

// 拦截网络请求以便修改
func (_this *winFreeApi) wkeNetHookRequest(job wkeNetJob) {
	_this._wkeNetHookRequest.Call(uintptr(job))
}

// 修改网络请求的 URL
func (_this *winFreeApi) wkeNetChangeRequestUrl(job wkeNetJob, url string) {
	p := toCallStr(url)
	_this._wkeNetChangeRequestUrl.Call(uintptr(job), uintptr(unsafe.Pointer(&p[0])))
}

// 设置网络请求的 HTTP 头字段
func (_this *winFreeApi) wkeNetSetHTTPHeaderField(job wkeNetJob, name, value string) {
	np := toCallStr(name)
	vp := toCallStr(value)
	_this._wkeNetSetHTTPHeaderField.Call(uintptr(job), uintptr(unsafe.Pointer(&np[0])), uintptr(unsafe.Pointer(&vp[0])))
}

// 将 wkeString 转换为 Go 字符串
func (_this *winFreeApi) wkeGetString(str wkeString) string {
	r, _, _ := _this._wkeGetString.Call(uintptr(str))
	return ptrToUtf8(r)
}

// 设置控制台消息输出的回调函数
func (_this *winFreeApi) wkeOnConsole(wke wkeHandle, callback wkeConsoleCallback, param uintptr) {
	_this._wkeOnConsole.Call(uintptr(wke), syscall.NewCallbackCDecl(callback), param)
}

// 获取指定框架的全局 JavaScript 执行上下文
func (_this *winFreeApi) wkeGetGlobalExecByFrame(wke wkeHandle, frame wkeFrame) jsExecState {
	r, _, _ := _this._wkeGetGlobalExecByFrame.Call(uintptr(wke), uintptr(frame))
	return jsExecState(r)
}

// 判断指定框架是否为远程框架
func (_this *winFreeApi) wkeIsWebRemoteFrame(wke wkeHandle, frame wkeFrame) bool {
	r, _, _ := _this._wkeIsWebRemoteFrame.Call(uintptr(wke), uintptr(frame))
	return r != 0
}

// 获取指定框架的 URL
func (_this *winFreeApi) wkeGetFrameUrl(wke wkeHandle, frame wkeFrame) string {
	r, _, _ := _this._wkeGetFrameUrl.Call(uintptr(wke), uintptr(frame))
	return ptrToUtf8(r)
}

// 判断指定框架是否为主框架
func (_this *winFreeApi) wkeIsMainFrame(wke wkeHandle, frame wkeFrame) bool {
	r, _, _ := _this._wkeIsMainFrame.Call(uintptr(wke), uintptr(frame))
	return r != 0
}

// 设置脚本上下文创建时的回调函数
func (_this *winFreeApi) wkeOnDidCreateScriptContext(wke wkeHandle, callback wkeDidCreateScriptContextCallback, param uintptr) {
	_this._wkeOnDidCreateScriptContext.Call(uintptr(wke), syscall.NewCallbackCDecl(callback), param)
}

// 取消 webview 的焦点
func (_this *winFreeApi) wkeKillFocus(wke wkeHandle) {
	_this._wkeKillFocus.Call(uintptr(wke))
}

// 从 JavaScript 执行上下文获取关联的 webview 句柄
func (_this *winFreeApi) jsGetWebView(es jsExecState) wkeHandle {
	r, _, _ := _this._jsGetWebView.Call(uintptr(es))
	return wkeHandle(r)
}

// 销毁 webview 实例
func (_this *winFreeApi) wkeDestroyWebView(wke wkeHandle) {
	_this._wkeDestroyWebView.Call(uintptr(wke))
}

// 设置 JavaScript 对象的属性值
func (_this *winFreeApi) jsSet(es jsExecState, obj jsValue, name string, value jsValue) {
	ptr := []byte(name)
	if is64 {
		_this._jsSet.Call(uintptr(es), uintptr(obj), uintptr(unsafe.Pointer(&ptr[0])), uintptr(value))
	} else {
		ol, oh := _toLH(obj)
		vl, vh := _toLH(value)
		_this._jsSet.Call(uintptr(es), uintptr(ol), uintptr(oh), uintptr(unsafe.Pointer(&ptr[0])), uintptr(vl), uintptr(vh))
	}
}

// 创建一个空的 JavaScript 对象
func (_this *winFreeApi) jsEmptyObject(es jsExecState) jsValue {
	if is64 {
		r, _, _ := _this._jsEmptyObject.Call(uintptr(es))
		return jsValue(r)
	}
	l, h, _ := _this._jsEmptyObject.Call(uintptr(es))
	return _toJsValue(l, h)
}

// 创建一个 JavaScript 函数对象
func (_this *winFreeApi) jsFunction(es jsExecState, data *jsData) jsValue {
	if is64 {
		r, _, _ := _this._jsFunction.Call(uintptr(es), uintptr(unsafe.Pointer(data)))
		return jsValue(r)
	}
	l, h, _ := _this._jsFunction.Call(uintptr(es), uintptr(unsafe.Pointer(data)))
	return _toJsValue(l, h)
}

// 设置 JavaScript 数组指定索引的值
func (_this *winFreeApi) jsSetAt(es jsExecState, array jsValue, index uint32, value jsValue) {
	if is64 {
		_this._jsSetAt.Call(uintptr(es), uintptr(array), uintptr(index), uintptr(value))
	} else {
		al, ah := _toLH(array)
		vl, vh := _toLH(value)
		_this._jsSetAt.Call(uintptr(es), uintptr(al), uintptr(ah), uintptr(index), uintptr(vl), uintptr(vh))
	}
}

// 设置 JavaScript 数组的长度
func (_this *winFreeApi) jsSetLength(es jsExecState, array jsValue, length uint32) {
	if is64 {
		_this._jsSetLength.Call(uintptr(es), uintptr(array), uintptr(length))
	} else {
		l, h := _toLH(array)
		_this._jsSetLength.Call(uintptr(es), uintptr(l), uintptr(h), uintptr(length))
	}
}

// 创建一个空的 JavaScript 数组
func (_this *winFreeApi) jsEmptyArray(es jsExecState) jsValue {
	if is64 {
		r, _, _ := _this._jsEmptyArray.Call(uintptr(es))
		return jsValue(r)
	}
	l, h, _ := _this._jsEmptyArray.Call(uintptr(es))
	return _toJsValue(l, h)
}

// 将 Go 字符串转换为 JavaScript 字符串值
func (_this *winFreeApi) jsString(es jsExecState, value string) jsValue {
	ptr := toCallStr(value)
	if is64 {
		r, _, _ := _this._jsString.Call(uintptr(es), uintptr(unsafe.Pointer(&ptr[0])))
		return jsValue(r)
	}
	l, h, _ := _this._jsString.Call(uintptr(es), uintptr(unsafe.Pointer(&ptr[0])))
	return _toJsValue(l, h)
}

// 将浮点数转换为 JavaScript 数值
func (_this *winFreeApi) jsDouble(value float64) jsValue {
	ptr := toCallStr(strconv.FormatFloat(value, 'f', 9, 64))
	if is64 {
		r, _, _ := _this._jsDoubleString.Call(uintptr(unsafe.Pointer(&ptr[0])))
		return jsValue(r)
	}
	l, h, _ := _this._jsDoubleString.Call(uintptr(unsafe.Pointer(&ptr[0])))
	return _toJsValue(l, h)
}

// 将布尔值转换为 JavaScript 布尔值
func (_this *winFreeApi) jsBoolean(value bool) jsValue {
	if is64 {
		r, _, _ := _this._jsBoolean.Call(uintptr(toBool(value)))
		return jsValue(r)
	}
	l, h, _ := _this._jsBoolean.Call(uintptr(toBool(value)))
	return _toJsValue(l, h)
}

// 将整数转换为 JavaScript 数值
func (_this *winFreeApi) jsInt(value int32) jsValue {
	if is64 {
		r, _, _ := _this._jsInt.Call(uintptr(value))
		return jsValue(r)
	}
	l, h, _ := _this._jsInt.Call(uintptr(value))
	return _toJsValue(l, h)
}

// 调用 JavaScript 函数
func (_this *winFreeApi) jsCall(es jsExecState, fn, thisObject jsValue, args []jsValue) jsValue {
	var ptr = uintptr(0)
	if len(args) > 0 {
		ptr = uintptr(unsafe.Pointer(&args[0]))
	}
	if is64 {
		r, _, _ := _this._jsCall.Call(uintptr(es), uintptr(fn), uintptr(thisObject), ptr, uintptr(len(args)))
		return jsValue(r)
	}
	fl, fh := _toLH(fn)
	ol, oh := _toLH(thisObject)
	l, h, _ := _this._jsCall.Call(uintptr(es), uintptr(fl), uintptr(fh), uintptr(ol), uintptr(oh), ptr, uintptr(len(args)))
	return _toJsValue(l, h)
}

// 获取 webview 的全局 JavaScript 执行上下文
func (_this *winFreeApi) wkeGlobalExec(wke wkeHandle) jsExecState {
	r, _, _ := _this._wkeGlobalExec.Call(uintptr(wke))
	return jsExecState(r)
}

// 获取全局 JavaScript 变量
func (_this *winFreeApi) jsGetGlobal(es jsExecState, name string) jsValue {
	ptr := toCallStr(name)
	if is64 {
		r, _, _ := _this._jsGetGlobal.Call(uintptr(es), uintptr(unsafe.Pointer(&ptr[0])))
		return jsValue(r)
	}
	l, h, _ := _this._jsGetGlobal.Call(uintptr(es), uintptr(unsafe.Pointer(&ptr[0])))
	return _toJsValue(l, h)
}

// 设置全局 JavaScript 变量
func (_this *winFreeApi) jsSetGlobal(es jsExecState, name string, value jsValue) {
	ptr := toCallStr(name)
	if is64 {
		_this._jsSetGlobal.Call(uintptr(es), uintptr(unsafe.Pointer(&ptr[0])), uintptr(value))
	} else {
		l, h := _toLH(value)
		_this._jsSetGlobal.Call(uintptr(es), uintptr(unsafe.Pointer(&ptr[0])), uintptr(l), uintptr(h))
	}
}

// 执行 JavaScript 代码并返回结果
func (_this *winFreeApi) jsEval(es jsExecState, js string) jsValue {
	ptr := toCallStr(js)
	if is64 {
		rs, _, _ := _this._jsEval.Call(uintptr(es), uintptr(unsafe.Pointer(&ptr[0])))
		return jsValue(rs)
	}
	l, h, _ := _this._jsEval.Call(uintptr(es), uintptr(unsafe.Pointer(&ptr[0])))
	return _toJsValue(l, h)
}

// 获取 JavaScript 对象的所有属性名
func (_this *winFreeApi) jsGetKeys(es jsExecState, value jsValue) []string {
	var rs uintptr
	if is64 {
		r, _, _ := _this._jsGetKeys.Call(uintptr(es), uintptr(value))
		rs = r
	} else {
		l, h := _toLH(value)
		r, _, _ := _this._jsGetKeys.Call(uintptr(es), uintptr(l), uintptr(h))
		rs = r
	}
	keys := *((*jsKeys)(unsafe.Pointer(rs)))
	items := make([]string, keys.length)
	for i := 0; i < len(items); i++ {
		items[i] = ptrToUtf8(*((*uintptr)(unsafe.Pointer(keys.first))))
		keys.first += unsafe.Sizeof(uintptr(0))
	}
	return items

	//_this._jsGetKeys.Call(uintptr(es), uintptr(value))
	//return []string{"n1", "n2"}

	//json := _this.jsGetGlobal(es, "Object")
	//stringify := _this.jsGet(es, json, "keys")
	//rs := _this.jsCall(es, stringify, _this.jsUndefined(), []jsValue{value})
	//alen := _this.jsGetLength(es, rs)
	//items := make([]string, alen)
	//for i := 0; i < len(items); i++ {
	//	v := _this.jsGetAt(es, rs, uint32(i))
	//	str := _this.jsToTempString(es, v)
	//	items[i] = str
	//}
	//return items
}

// 获取 JavaScript 对象的属性值
func (_this *winFreeApi) jsGet(es jsExecState, value jsValue, name string) jsValue {
	ptr := toCallStr(name)
	if is64 {
		r, _, _ := _this._jsGet.Call(uintptr(es), uintptr(value), uintptr(unsafe.Pointer(&ptr[0])))
		return jsValue(r)
	}
	pl, ph := _toLH(value)
	l, h, _ := _this._jsGet.Call(uintptr(es), uintptr(pl), uintptr(ph), uintptr(unsafe.Pointer(&ptr[0])))
	return _toJsValue(l, h)
}

// 获取 JavaScript 数组指定索引的值
func (_this *winFreeApi) jsGetAt(es jsExecState, value jsValue, index uint32) jsValue {
	if is64 {
		r, _, _ := _this._jsGetAt.Call(uintptr(es), uintptr(value), uintptr(index))
		return jsValue(r)
	}
	pl, ph := _toLH(value)
	rl, rh, _ := _this._jsGetAt.Call(uintptr(es), uintptr(pl), uintptr(ph), uintptr(index))
	return _toJsValue(rl, rh)
}

// 获取 JavaScript 数组或对象的长度
func (_this *winFreeApi) jsGetLength(es jsExecState, value jsValue) int {
	if is64 {
		r, _, _ := _this._jsGetLength.Call(uintptr(es), uintptr(value))
		return int(r)
	}
	l, h := _toLH(value)
	r, _, _ := _this._jsGetLength.Call(uintptr(es), uintptr(l), uintptr(h))
	return int(r)
}

// 返回 JavaScript undefined 值
func (_this *winFreeApi) jsUndefined() jsValue {
	if is64 {
		r, _, _ := _this._jsUndefined.Call()
		return jsValue(r)
	}
	l, h, _ := _this._jsUndefined.Call()
	return _toJsValue(l, h)
}

// 将 JavaScript 值转换为 Go 布尔值
func (_this *winFreeApi) jsToBoolean(es jsExecState, value jsValue) bool {
	if is64 {
		r, _, _ := _this._jsToBoolean.Call(uintptr(es), uintptr(value))
		return byte(r) != 0
	}
	pl, ph := _toLH(value)
	r, _, _ := _this._jsToBoolean.Call(uintptr(es), uintptr(pl), uintptr(ph))
	return byte(r) != 0
}

// 将 JavaScript 值转换为 Go 浮点数
func (_this *winFreeApi) jsToDouble(es jsExecState, value jsValue) float64 {
	var rs uintptr
	if is64 {
		r, _, _ := _this._jsToDoubleString.Call(uintptr(es), uintptr(value))
		rs = r
	} else {
		l, h := _toLH(value)
		r, _, _ := _this._jsToDoubleString.Call(uintptr(es), uintptr(l), uintptr(h))
		rs = r
	}
	str := ptrToUtf8(rs)
	n, _ := strconv.ParseFloat(str, 10)
	return n
}

// 将 JavaScript 值转换为临时字符串
func (_this *winFreeApi) jsToTempString(es jsExecState, value jsValue) string {
	if is64 {
		r, _, _ := _this._jsToTempString.Call(uintptr(es), uintptr(value))
		return ptrToUtf8(r)
	}
	l, h := _toLH(value)
	r, _, _ := _this._jsToTempString.Call(uintptr(es), uintptr(l), uintptr(h))
	return ptrToUtf8(r)
}

// 获取 JavaScript 值的类型
func (_this *winFreeApi) jsTypeOf(value jsValue) jsType {
	if is64 {
		r, _, _ := _this._jsTypeOf.Call(uintptr(value))
		return jsType(r)
	}
	l, h := _toLH(value)
	r, _, _ := _this._jsTypeOf.Call(uintptr(l), uintptr(h))
	return jsType(r)
}

// 获取 JavaScript 函数调用参数列表中指定索引的参数
func (_this *winFreeApi) jsArg(es jsExecState, index uint32) jsValue {
	if is64 {
		r, _, _ := _this._jsArg.Call(uintptr(es), uintptr(index))
		return jsValue(r)
	}
	l, h, _ := _this._jsArg.Call(uintptr(es), uintptr(index))
	return jsValue(uintptr(_toInt64(int32(l), int32(h))))
}

// 获取 JavaScript 函数调用参数的数量
func (_this *winFreeApi) jsArgCount(es jsExecState) uint32 {
	r, _, _ := _this._jsArgCount.Call(uintptr(es))
	return uint32(r)
}

// 绑定 Go 函数为 JavaScript 全局函数
func (_this *winFreeApi) wkeJsBindFunction(name string, fn wkeJsNativeFunction, param uintptr, argCount uint32) {
	ptr := toCallStr(name)
	_this._wkeJsBindFunction.Call(uintptr(unsafe.Pointer(&ptr[0])), syscall.NewCallbackCDecl(fn), param, uintptr(argCount))
}

// 取消网络请求
func (_this *winFreeApi) wkeNetCancelRequest(job wkeNetJob) {
	_this._wkeNetCancelRequest.Call(uintptr(job))
}

// 设置网络响应接收时的回调函数
func (_this *winFreeApi) wkeNetOnResponse(wke wkeHandle, callback wkeNetResponseCallback, param uintptr) {
	_this._wkeNetOnResponse.Call(uintptr(wke), syscall.NewCallbackCDecl(callback), param)
}

// 设置 URL 开始加载时的回调函数
func (_this *winFreeApi) wkeOnLoadUrlBegin(wke wkeHandle, callback wkeLoadUrlBeginCallback, param uintptr) {
	_this._wkeOnLoadUrlBegin.Call(uintptr(wke), syscall.NewCallbackCDecl(callback), param)
}

// 获取网络请求的方法类型
func (_this *winFreeApi) wkeNetGetRequestMethod(job wkeNetJob) wkeRequestType {
	r, _, _ := _this._wkeNetGetRequestMethod.Call(uintptr(job))
	return wkeRequestType(r)
}

// 设置网络请求的响应数据
func (_this *winFreeApi) wkeNetSetData(job wkeNetJob, buf []byte) {
	if len(buf) == 0 {
		buf = []byte{0}
	}
	_this._wkeNetSetData.Call(uintptr(job), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
}

// 获取输入光标的位置和大小
func (_this *winFreeApi) wkeGetCaretRect(wke wkeHandle) wkeRect {
	r, _, _ := _this._wkeGetCaretRect.Call(uintptr(wke))
	return *((*wkeRect)(unsafe.Pointer(r)))
}

// 设置 webview 获得焦点
func (_this *winFreeApi) wkeSetFocus(wke wkeHandle) {
	_this._wkeSetFocus.Call(uintptr(wke))
}

// 触发按键按下事件
func (_this *winFreeApi) wkeFireKeyPressEvent(wke wkeHandle, code int, flags uint32, isSysKey bool) bool {
	ret, _, _ := _this._wkeFireKeyPressEvent.Call(
		uintptr(wke),
		uintptr(code),
		uintptr(flags),
		uintptr(toBool(isSysKey)))
	return byte(ret) != 0
}

// 触发按键按下事件
func (_this *winFreeApi) wkeFireKeyDownEvent(wke wkeHandle, code, flags uint32, isSysKey bool) bool {
	ret, _, _ := _this._wkeFireKeyDownEvent.Call(
		uintptr(wke),
		uintptr(code),
		uintptr(flags),
		uintptr(toBool(isSysKey)))
	return byte(ret) != 0
}

// 触发按键释放事件
func (_this *winFreeApi) wkeFireKeyUpEvent(wke wkeHandle, code, flags uint32, isSysKey bool) bool {
	ret, _, _ := _this._wkeFireKeyUpEvent.Call(
		uintptr(wke),
		uintptr(code),
		uintptr(flags),
		uintptr(toBool(isSysKey)))
	return byte(ret) != 0
}

// 获取当前鼠标光标类型
func (_this *winFreeApi) wkeGetCursorInfoType(wke wkeHandle) wkeCursorType {
	r, _, _ := _this._wkeGetCursorInfoType.Call(uintptr(wke))
	return wkeCursorType(r)
}

// 触发鼠标滚轮事件
func (_this *winFreeApi) wkeFireMouseWheelEvent(wke wkeHandle, x, y, delta, flags int32) bool {
	r, _, _ := _this._wkeFireMouseWheelEvent.Call(
		uintptr(wke),
		uintptr(x),
		uintptr(y),
		uintptr(delta),
		uintptr(flags))
	return byte(r) != 0
}

// 触发鼠标事件
func (_this *winFreeApi) wkeFireMouseEvent(wke wkeHandle, message, x, y, flags int32) bool {
	r, _, _ := _this._wkeFireMouseEvent.Call(
		uintptr(wke),
		uintptr(message),
		uintptr(x),
		uintptr(y),
		uintptr(flags))
	return byte(r) != 0
}

// 将位图数据绘制到 webview
func (_this *winFreeApi) wkePaint(wke wkeHandle, bits []byte, pitch uint32) {
	_this._wkePaint.Call(uintptr(wke), uintptr(unsafe.Pointer(&bits[0])), uintptr(pitch))
}

// 获取 webview 的高度
func (_this *winFreeApi) wkeGetHeight(wke wkeHandle) uint32 {
	r, _, _ := _this._wkeGetHeight.Call(uintptr(wke))
	return uint32(r)
}

// 获取 webview 的宽度
func (_this *winFreeApi) wkeGetWidth(wke wkeHandle) uint32 {
	r, _, _ := _this._wkeGetWidth.Call(uintptr(wke))
	return uint32(r)
}

// 调整 webview 的大小
func (_this *winFreeApi) wkeResize(wke wkeHandle, w, h uint32) {
	_this._wkeResize.Call(uintptr(wke), uintptr(w), uintptr(h))
}

// 加载指定的 URL
func (_this *winFreeApi) wkeLoadURL(wke wkeHandle, url string) {
	ptr := toCallStr(url)
	_this._wkeLoadURL.Call(uintptr(wke), uintptr(unsafe.Pointer(&ptr[0])))
}

// 加载 HTML 内容并指定基础 URL
func (_this *winFreeApi) wkeLoadHtmlWithBaseUrl(wke wkeHandle, html, baseUrl string) {
	htmlPtr := toCallStr(html)
	baseUrlPtr := toCallStr(baseUrl)
	_this._wkeLoadHtmlWithBaseUrl.Call(uintptr(wke), uintptr(unsafe.Pointer(&htmlPtr[0])), uintptr(unsafe.Pointer(&baseUrlPtr[0])))
}

// 设置位图更新时的回调函数
func (_this *winFreeApi) wkeOnPaintBitUpdated(wke wkeHandle, callback wkePaintBitUpdatedCallback, param uintptr) {
	_this._wkeOnPaintBitUpdated.Call(uintptr(wke), syscall.NewCallbackCDecl(callback), param)
}

// 将 webview 绑定到指定的窗口句柄
func (_this *winFreeApi) wkeSetHandle(wke wkeHandle, handle uintptr) {
	_this._wkeSetHandle.Call(uintptr(wke), handle)
}

// 创建一个新的 webview 实例
func (_this *winFreeApi) wkeCreateWebView() wkeHandle {
	r, _, _ := _this._wkeCreateWebView.Call()
	return wkeHandle(r)
}

// 设置调试配置参数
func (_this *winFreeApi) wkeSetDebugConfig(wke wkeHandle, debugString, param string) {
	debugPtr := toCallStr(debugString)
	paramPtr := toCallStr(param)
	_this._wkeSetDebugConfig.Call(uintptr(wke), uintptr(unsafe.Pointer(&debugPtr[0])), uintptr(unsafe.Pointer(&paramPtr[0])))
}

// 设置本地存储的完整路径
func (_this *winFreeApi) wkeSetLocalStorageFullPath(wke wkeHandle, path string) {
	// 将字符串转换为 UTF-16 (WCHAR*) 指针
	pathPtr, _ := syscall.UTF16PtrFromString(path)
	_this._wkeSetLocalStorageFullPath.Call(uintptr(wke), uintptr(unsafe.Pointer(pathPtr)))
}

// 设置 HTTP 请求的 User-Agent 字符串
func (_this *winFreeApi) wkeSetUserAgent(wke wkeHandle, userAgent string) {
	ptr := toCallStr(userAgent)
	_this._wkeSetUserAgent.Call(uintptr(wke), uintptr(unsafe.Pointer(&ptr[0])))
}
