package bridge

import (
	fm "github.com/reghtml/mblink/forms"
)

// WindowLostFocusProc 窗口失去焦点回调函数类型
type WindowLostFocusProc func() bool

// WindowFocusProc 窗口获得焦点回调函数类型
type WindowFocusProc func() bool

// WindowImeStartCompositionProc 输入法开始组合回调函数类型
type WindowImeStartCompositionProc func() bool

// WindowSetCursorProc 设置光标回调函数类型
type WindowSetCursorProc func() bool

// WindowShowProc 窗口显示回调函数类型
type WindowShowProc func()

// WindowCreateProc 窗口创建回调函数类型
type WindowCreateProc func(handle uintptr)

// WindowDestroyProc 窗口销毁回调函数类型
type WindowDestroyProc func()

// WindowResizeProc 窗口大小改变回调函数类型
type WindowResizeProc func(e fm.Rect)

// WindowMoveProc 窗口移动回调函数类型
type WindowMoveProc func(e fm.Point) bool

// WindowMouseMoveProc 鼠标移动回调函数类型
type WindowMouseMoveProc func(e *fm.MouseEvArgs)

// WindowMouseDownProc 鼠标按下回调函数类型
type WindowMouseDownProc func(e *fm.MouseEvArgs)

// WindowMouseUpProc 鼠标释放回调函数类型
type WindowMouseUpProc func(e *fm.MouseEvArgs)

// WindowMouseWheelProc 鼠标滚轮回调函数类型
type WindowMouseWheelProc func(e *fm.MouseEvArgs)

// WindowMouseClickProc 鼠标点击回调函数类型
type WindowMouseClickProc func(e *fm.MouseEvArgs)

// WindowPaintProc 窗口绘制回调函数类型
type WindowPaintProc func(e fm.PaintEvArgs) bool

// WindowKeyDownProc 按键按下回调函数类型
type WindowKeyDownProc func(e *fm.KeyEvArgs)

// WindowKeyUpProc 按键释放回调函数类型
type WindowKeyUpProc func(e *fm.KeyEvArgs)

// WindowKeyPressProc 按键按下回调函数类型
type WindowKeyPressProc func(e *fm.KeyPressEvArgs)

// Window 窗口接口，提供窗口的基本操作和事件处理
type Window interface {
	GetHandle() uintptr
	SetOnCreate(proc WindowCreateProc) WindowCreateProc
	SetOnDestroy(proc WindowDestroyProc) WindowDestroyProc
	SetOnResize(proc WindowResizeProc) WindowResizeProc
	SetOnMove(proc WindowMoveProc) WindowMoveProc
	SetOnMouseMove(proc WindowMouseMoveProc) WindowMouseMoveProc
	SetOnMouseDown(proc WindowMouseDownProc) WindowMouseDownProc
	SetOnMouseUp(proc WindowMouseUpProc) WindowMouseUpProc
	SetOnMouseWheel(proc WindowMouseWheelProc) WindowMouseWheelProc
	SetOnMouseClick(proc WindowMouseClickProc) WindowMouseClickProc
	SetOnPaint(proc WindowPaintProc) WindowPaintProc
	SetOnKeyDown(proc WindowKeyDownProc) WindowKeyDownProc
	SetOnKeyUp(proc WindowKeyUpProc) WindowKeyUpProc
	SetOnKeyPress(proc WindowKeyPressProc) WindowKeyPressProc
	SetOnCursor(proc WindowSetCursorProc) WindowSetCursorProc
	SetOnImeStartComposition(proc WindowImeStartCompositionProc) WindowImeStartCompositionProc
	SetOnFocus(proc WindowFocusProc) WindowFocusProc
	SetOnLostFocus(proc WindowLostFocusProc) WindowLostFocusProc
	SetOnShow(proc WindowShowProc) WindowShowProc

	GetProvider() Provider
	Invoke(fn func())
	InvokeEx(fn func(state interface{}), state interface{})
	SetSize(w int, h int)
	SetLocation(x int, y int)
	GetBound() fm.Bound
	Show()
	Hide()
	SetBgColor(color int32)
	CreateGraphics() fm.Graphics
	SetCursor(cursor fm.CursorType)
	GetCursor() fm.CursorType
	GetParent() Control
	GetOwner() Form
	ToClientPoint(p fm.Point) fm.Point
	IsEnable() bool
	Enable(b bool)
}
