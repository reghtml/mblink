package GoMiniblink

import (
	"fmt"

	fm "github.com/reghtml/mblink/forms"
	br "github.com/reghtml/mblink/forms/bridge"
	cs "github.com/reghtml/mblink/forms/controls"
	gw "github.com/reghtml/mblink/forms/windows/win32"
)

var (
	_fnMax   = "__formMax"
	_fnMin   = "__formMin"
	_fnClose = "__formClose"
	_fnDrop  = "__formDrop"
)

// MiniblinkForm 包含 Miniblink 浏览器的窗体
type MiniblinkForm struct {
	cs.Form
	View *MiniblinkBrowser

	wke           wkeHandle
	isTransparent bool
	resizeState   int
}

// 初始化窗体（使用默认参数）
func (_this *MiniblinkForm) Init() *MiniblinkForm {
	return _this.InitEx(br.FormParam{})
}

// 使用指定参数初始化窗体
func (_this *MiniblinkForm) InitEx(param br.FormParam) *MiniblinkForm {
	_this.Form.InitEx(param)
	_this.View = new(MiniblinkBrowser).Init()
	_this.View.SetAnchor(fm.AnchorStyle_Fill)
	_this.AddChild(_this.View)
	_this.wke = wkeHandle(_this.View.GetMiniblinkHandle())
	_this.setOn()
	_this.View.OnFocus()
	_this.View.JsFuncEx(_fnMax, func() {
		if _this.GetState() == fm.FormState_Max {
			_this.SetState(fm.FormState_Normal)
		} else {
			_this.SetState(fm.FormState_Max)
		}
	})
	_this.View.JsFuncEx(_fnMin, func() {
		_this.SetState(fm.FormState_Min)
	})
	_this.View.JsFuncEx(_fnClose, func() {
		_this.Close()
	})
	_this.View.EvDocumentReady["__goMiniblink_init_js"] = func(_ *MiniblinkBrowser, e DocumentReadyEvArgs) {
		e.RunJs("window.setFormButton();window.mbFormDrop();")
	}
	_this.setDrop()
	// 自动启用 setSize JavaScript 函数
	_this.enableSetSizeJsFunc()
	return _this
}

func (_this *MiniblinkForm) setDrop() {
	isDrop := false
	var anchor fm.Point
	var base fm.Point
	var isDragging bool // 实际拖拽状态标志

	_this.View.JsFuncEx(_fnDrop, func() {
		// JavaScript 调用此函数表示鼠标在拖拽区域按下
		// 设置 isDrop 标志，等待 EvMouseDown 事件来开始实际拖拽
		isDrop = true
	})

	_this.View.EvMouseDown["__goMiniblink_drop"] = func(s cs.GUI, e *fm.MouseEvArgs) {
		if isDrop && e.Button&fm.MouseButtons_Left != 0 {
			// 开始拖拽
			isDragging = true
			base = _this.GetBound().Point
			anchor = fm.Point{
				X: e.ScreenX,
				Y: e.ScreenY,
			}
			// 不改变鼠标样式，保持默认
			_this.View.MouseEnable(false)
		}
	}

	_this.View.EvMouseUp["__goMiniblink_drop"] = func(s cs.GUI, e *fm.MouseEvArgs) {
		if isDragging {
			// 停止拖拽
			isDragging = false
			isDrop = false
			// 不改变鼠标样式，保持默认
			_this.View.MouseEnable(true)
		}
	}

	_this.View.EvMouseMove["__goMiniblink_drop"] = func(s cs.GUI, e *fm.MouseEvArgs) {
		if isDragging {
			// 使用 cs.App 的 MouseIsDown 来实时检查鼠标左键状态
			mouseDown := cs.App.MouseIsDown()
			if !mouseDown[fm.MouseButtons_Left] {
				// 鼠标左键已经抬起，立即停止拖拽
				isDragging = false
				isDrop = false
				// 不改变鼠标样式，保持默认
				_this.View.MouseEnable(true)
				return
			}
			// 继续拖拽
			var nx = e.ScreenX - anchor.X
			var ny = e.ScreenY - anchor.Y
			nx = base.X + nx
			ny = base.Y + ny
			_this.SetLocation(nx, ny)
			e.IsHandle = true
		} else if isDrop {
			// 如果 isDrop 为 true 但 isDragging 为 false，说明还没有开始拖拽
			// 检查鼠标左键是否按下，如果按下则开始拖拽
			if e.Button&fm.MouseButtons_Left != 0 {
				isDragging = true
				base = _this.GetBound().Point
				anchor = fm.Point{
					X: e.ScreenX,
					Y: e.ScreenY,
				}
				// 不改变鼠标样式，保持默认
				_this.View.MouseEnable(false)
			}
		}
	}
}

func (_this *MiniblinkForm) setOn() {
	bakOnResize := _this.OnResize
	_this.OnResize = func(e fm.Rect) {
		_this.View.SetSize(e.Width, e.Height)
		bakOnResize(e)
	}
	bakOnLoad := _this.OnLoad
	_this.OnLoad = func() {
		if _this.isTransparent {
			hWnd := gw.HWND(_this.GetHandle())
			style := gw.GetWindowLong(hWnd, gw.GWL_EXSTYLE)
			if style&gw.WS_EX_LAYERED != gw.WS_EX_LAYERED {
				gw.SetWindowLong(hWnd, gw.GWL_EXSTYLE, style|gw.WS_EX_LAYERED)
			}
			b := _this.GetBound()
			_this.transparentPaint(b.Width, b.Height)
		}
		bakOnLoad()
	}
	bakOnJsReady := _this.View.OnJsReady
	_this.View.OnJsReady = func(e JsReadyEvArgs) {
		bakOnJsReady(e)
		_this.setFormFn(e)
	}
}

// 启用透明模式
func (_this *MiniblinkForm) TransparentMode() {
	_this.isTransparent = true
	_this.SetBorderStyle(fm.FormBorder_None)
	_this.View.OnPaintUpdated = func(e PaintUpdatedEvArgs) {
		_this.transparentPaint(e.Bound().Width, e.Bound().Height)
		e.Cancel()
	}
	mbApi.wkeSetTransparent(_this.wke, true)
}

func (_this *MiniblinkForm) transparentPaint(width, height int) {
	bn := _this.GetBound()
	hWnd := gw.HWND(_this.GetHandle())
	mdc := gw.HDC(mbApi.wkeGetViewDC(_this.View.core.GetHandle()))
	src := gw.POINT{}
	dst := gw.POINT{
		X: int32(bn.X),
		Y: int32(bn.Y),
	}
	size := gw.SIZE{
		CX: int32(width),
		CY: int32(height),
	}
	blend := gw.BLENDFUNCTION{
		SourceConstantAlpha: 255,
		AlphaFormat:         gw.AC_SRC_ALPHA,
	}
	gw.UpdateLayeredWindow(hWnd, 0, &dst, &size, mdc, &src, 0, &blend, 2)
}

func (_this *MiniblinkForm) setFormFn(frame FrameContext) {
	js := `
			var fnMax=window[%q];
			var fnMin=window[%q];
			var fnClose=window[%q];
			var fnDrop=window[%q];
			var mbFormClickHandler=null;
			var mbFormDblClickHandler=null;
			window.mbFormDrop=function(){
				document.getElementsByTagName("body")[0].addEventListener("mousedown",
					function (e) {
						var obj = e.target || e.srcElement;
						if ({ "INPUT": 1, "SELECT": 1 }[obj.tagName.toUpperCase()])
							return;
					
						while (obj) {
							for (var i = 0; i < obj.classList.length; i++) {
								if (obj.classList[i] === "mbform-nodrag")
									return;
								if (obj.classList[i] === "mbform-drag") {
									fnDrop();
									return;
								}
							}
							obj = obj.parentElement;
						}
					});
			};
			window.mbFormMax=function(obj){
				// 兼容旧代码，实际事件委托在 setFormButton 中统一处理
			};
			window.mbFormMin=function(obj){
				// 兼容旧代码，实际事件委托在 setFormButton 中统一处理
			};
			window.mbFormClose=function(obj){
				// 兼容旧代码，实际事件委托在 setFormButton 中统一处理
			};
			window.setFormButton=function(){
				// 使用事件委托处理所有按钮点击，避免 DOM 更新后事件丢失的问题
				// 每次调用时重新绑定，确保事件监听器始终有效
				if(document.addEventListener){
					// 移除旧的监听器（如果存在）
					if(mbFormClickHandler){
						document.removeEventListener("click", mbFormClickHandler, true);
					}
					mbFormClickHandler = function(e){
						var target = e.target || e.srcElement;
						if(!target.classList) return;
						// 检查目标元素或其父元素是否包含按钮类名
						while(target && target !== document){
							if(target.classList){
								if(target.classList.contains("mbform-max")){
									e.preventDefault();
									e.stopPropagation();
									if(fnMax) fnMax();
									return false;
								}else if(target.classList.contains("mbform-min")){
									e.preventDefault();
									e.stopPropagation();
									if(fnMin) fnMin();
									return false;
								}else if(target.classList.contains("mbform-close")){
									e.preventDefault();
									e.stopPropagation();
									if(fnClose) fnClose();
									return false;
								}
							}
							target = target.parentElement;
						}
					};
					document.addEventListener("click", mbFormClickHandler, true);
				}
				// 处理双击 .mbform-dbmax 元素，执行与单击 .mbform-max 相同的功能
				if(document.addEventListener){
					// 移除旧的监听器（如果存在）
					if(mbFormDblClickHandler){
						document.removeEventListener("dblclick", mbFormDblClickHandler, true);
					}
					mbFormDblClickHandler = function(e){
						var target = e.target || e.srcElement;
						if(!target.classList) return;
						// 检查目标元素或其父元素是否包含 mbform-dbmax 类名
						while(target && target !== document){
							if(target.classList && target.classList.contains("mbform-dbmax")){
								e.preventDefault();
								e.stopPropagation();
								if(fnMax) fnMax();
								return false;
							}
							target = target.parentElement;
						}
					};
					document.addEventListener("dblclick", mbFormDblClickHandler, true);
				}
				// 定期重新绑定，防止长时间不操作后失效（每3分钟重新绑定一次）
				if(!window.mbFormButtonTimer){
					window.mbFormButtonTimer = setInterval(function(){
						window.setFormButton();
					}, 3 * 60 * 1000);
				}
			};
	`
	js = fmt.Sprintf(js, _fnMax, _fnMin, _fnClose, _fnDrop)
	frame.RunJs(js)
}

// SetSizeAndPosition 根据数字小键盘布局设置窗口大小和位置
// position: 1-9 对应数字小键盘布局
// 7 8 9    左上  上中  右上
// 4 5 6    左中  中间  右中
// 1 2 3    左下  下中  右下
// 设置窗体大小和位置
func (_this *MiniblinkForm) SetSizeAndPosition(width, height, position int) {
	_this.SetSize(width, height)
	screen := cs.App.GetScreen()
	screenW := screen.WorkArea.Width
	screenH := screen.WorkArea.Height
	var x, y int

	switch position {
	case 1: // 左下角
		x = 0
		y = screenH - height
	case 2: // 下中
		x = (screenW - width) / 2
		y = screenH - height
	case 3: // 右下角
		x = screenW - width
		y = screenH - height
	case 4: // 左中
		x = 0
		y = (screenH - height) / 2
	case 5: // 中间
		x = (screenW - width) / 2
		y = (screenH - height) / 2
	case 6: // 右中
		x = screenW - width
		y = (screenH - height) / 2
	case 7: // 左上角
		x = 0
		y = 0
	case 8: // 上中
		x = (screenW - width) / 2
		y = 0
	case 9: // 右上角
		x = screenW - width
		y = 0
	default: // 默认居中
		x = (screenW - width) / 2
		y = (screenH - height) / 2
	}

	_this.SetLocation(x, y)
}

// enableSetSizeJsFunc 内部方法，启用 setSize JavaScript 函数
// 调用方式: setSize(width, height, position)
// position: 1-9 对应数字小键盘布局
func (_this *MiniblinkForm) enableSetSizeJsFunc() {
	_this.View.JsFuncEx("setSize", func(width, height, position float64) {
		w := int(width)
		h := int(height)
		pos := int(position)
		_this.SetSizeAndPosition(w, h, pos)
	})
}
