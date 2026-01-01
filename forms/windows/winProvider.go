package windows

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"syscall"
	"unsafe"

	fm "github.com/reghtml/mblink/forms"
	br "github.com/reghtml/mblink/forms/bridge"
	win "github.com/reghtml/mblink/forms/windows/win32"
)

// windowsMsgProc Windows 消息处理函数类型
type windowsMsgProc func(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr

// Provider Windows 平台提供者实现
type Provider struct {
	hInstance  win.HINSTANCE
	className  string
	wndClass   win.WNDCLASSEX
	mainId     win.HWND
	tmpWnd     map[uintptr]baseWindow
	handleWnds map[win.HWND]baseWindow
	defOwner   win.HWND
	defIcon    win.HICON
	watchAll   map[win.HWND][]windowsMsgProc
	forms      map[win.HWND]br.Form
	minWidth   int32 // 最小宽度
	minHeight  int32 // 最小高度
}

// 初始化提供者
func (_this *Provider) Init() *Provider {
	_this.forms = make(map[win.HWND]br.Form)
	_this.watchAll = make(map[win.HWND][]windowsMsgProc)
	_this.tmpWnd = make(map[uintptr]baseWindow)
	_this.handleWnds = make(map[win.HWND]baseWindow)
	_this.className = "YAN4/TOOLS"
	_this.hInstance = win.GetModuleHandle(nil)
	_this.registerWndClass()
	return _this
}

// 监视窗口消息
func (_this *Provider) watch(wnd baseWindow, proc windowsMsgProc) {
	_this.watchAll[wnd.hWnd()] = append(_this.watchAll[wnd.hWnd()], proc)
}

// 获取鼠标位置
func (_this *Provider) MouseLocation() fm.Point {
	pos := win.POINT{}
	win.GetCursorPos(&pos)
	return fm.Point{
		X: int(pos.X),
		Y: int(pos.Y),
	}
}

// 获取应用程序目录
func (_this *Provider) AppDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

// 获取修饰键状态
func (_this *Provider) ModifierKeys() map[fm.Keys]bool {
	keys := make(map[fm.Keys]bool)
	cs := win.GetKeyState(int32(win.VK_CONTROL))
	ss := win.GetKeyState(int32(win.VK_SHIFT))
	as := win.GetKeyState(int32(win.VK_MENU))
	keys[fm.Keys_Ctrl] = cs < 0
	keys[fm.Keys_Shift] = ss < 0
	keys[fm.Keys_Alt] = as < 0
	return keys
}

// 获取鼠标按键按下状态
func (_this *Provider) MouseIsDown() map[fm.MouseButtons]bool {
	keys := make(map[fm.MouseButtons]bool)
	ls := win.GetKeyState(int32(win.VK_LBUTTON))
	rs := win.GetKeyState(int32(win.VK_RBUTTON))
	ms := win.GetKeyState(int32(win.VK_MBUTTON))
	keys[fm.MouseButtons_Left] = ls < 0
	keys[fm.MouseButtons_Right] = rs < 0
	keys[fm.MouseButtons_Middle] = ms < 0
	return keys
}

// 获取屏幕信息
func (_this *Provider) GetScreen() fm.Screen {
	var s = fm.Screen{
		Full: fm.Rect{
			Width:  int(win.GetSystemMetrics(win.SM_CXSCREEN)),
			Height: int(win.GetSystemMetrics(win.SM_CYSCREEN)),
		},
		WorkArea: fm.Rect{
			Width:  int(win.GetSystemMetrics(win.SM_CXFULLSCREEN)),
			Height: int(win.GetSystemMetrics(win.SM_CYFULLSCREEN)),
		},
	}
	return s
}

// 设置应用程序图标
func (_this *Provider) SetIcon(file string) {
	h := win.LoadImage(_this.hInstance, sto16(file), win.IMAGE_ICON, 0, 0, win.LR_LOADFROMFILE)
	_this.defIcon = win.HICON(h)
}

// 设置图标数据
func (_this *Provider) SetIconData(iconData []byte) {
	if len(iconData) == 0 {
		return
	}
	hIcon := _this.createIconFromData(iconData)
	if hIcon != 0 {
		_this.defIcon = win.HICON(hIcon)
	}
}

// SetMinSize 设置全局窗口最小尺寸
func (_this *Provider) SetMinSize(width, height int) {
	_this.minWidth = int32(width)
	_this.minHeight = int32(height)
}

// 从图标数据创建图标句柄
func (_this *Provider) createIconFromData(iconData []byte) uintptr {
	if len(iconData) < 22 {
		return 0
	}

	// ICO文件格式：ICONDIR (6字节) + ICONDIRENTRY (16字节) + 图标数据
	iconCount := int(iconData[4]) | (int(iconData[5]) << 8)
	if iconCount == 0 || iconCount > 100 {
		return 0
	}

	firstEntryOffset := 6
	if len(iconData) < firstEntryOffset+16 {
		return 0
	}

	// 读取第一个图标的偏移量和大小
	iconOffset := int(iconData[firstEntryOffset+12]) |
		(int(iconData[firstEntryOffset+13]) << 8) |
		(int(iconData[firstEntryOffset+14]) << 16) |
		(int(iconData[firstEntryOffset+15]) << 24)

	iconSize := int(iconData[firstEntryOffset+8]) |
		(int(iconData[firstEntryOffset+9]) << 8) |
		(int(iconData[firstEntryOffset+10]) << 16) |
		(int(iconData[firstEntryOffset+11]) << 24)

	if iconOffset >= len(iconData) || iconOffset+iconSize > len(iconData) {
		return 0
	}

	iconResData := iconData[iconOffset : iconOffset+iconSize]

	// 使用CreateIconFromResourceEx从内存创建图标
	user32 := syscall.NewLazyDLL("user32.dll")
	createIconFromResourceEx := user32.NewProc("CreateIconFromResourceEx")

	ret, _, _ := createIconFromResourceEx.Call(
		uintptr(unsafe.Pointer(&iconResData[0])),
		uintptr(len(iconResData)),
		1,          // TRUE - 图标
		0x00030000, // 版本号
		0,          // 默认宽度
		0,          // 默认高度
		0,          // 默认标志
	)

	return ret
}

// 注册窗口类
func (_this *Provider) registerWndClass() {
	_this.wndClass = win.WNDCLASSEX{
		Style:         win.CS_HREDRAW | win.CS_VREDRAW,
		LpfnWndProc:   syscall.NewCallbackCDecl(_this.classMsgProc),
		HInstance:     _this.hInstance,
		LpszClassName: sto16(_this.className),
		HCursor:       win.LoadCursor(0, win.MAKEINTRESOURCE(win.IDC_ARROW)),
		HbrBackground: win.GetSysColorBrush(win.COLOR_WINDOW),
	}
	_this.wndClass.CbSize = uint32(unsafe.Sizeof(_this.wndClass))
	win.RegisterClassEx(&_this.wndClass)
	_this.defOwner = win.CreateWindowEx(0,
		sto16(_this.className), sto16(""),
		win.WS_OVERLAPPED, 0, 0, 0, 0,
		0, 0, _this.hInstance, unsafe.Pointer(nil))
}

// 添加窗口到临时映射表
func (_this *Provider) add(wnd baseWindow) {
	ref := reflect.ValueOf(wnd).Pointer()
	_this.tmpWnd[ref] = wnd
}

// 窗口类消息处理过程
func (_this *Provider) classMsgProc(hWnd win.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr {
	runtime.LockOSThread()
	switch msg {
	case win.WM_CREATE:
		cs := *((*win.CREATESTRUCT)(unsafe.Pointer(lParam)))
		if wnd, es := _this.tmpWnd[cs.CreateParams]; es {
			delete(_this.tmpWnd, cs.CreateParams)
			_this.handleWnds[hWnd] = wnd
			if _this.mainId == 0 {
				_this.mainId = hWnd
			}
		}
	}
	for _, list := range _this.watchAll {
		for _, proc := range list {
			if rs := proc(hWnd, msg, wParam, lParam); rs != 0 {
				runtime.UnlockOSThread()
				return rs
			}
		}
	}
	if wnd, ok := _this.handleWnds[hWnd]; ok {
		if rs := wnd.onWndMsg(hWnd, msg, wParam, lParam); rs != 0 {
			runtime.UnlockOSThread()
			return rs
		}
	}
	switch msg {
	case win.WM_GETMINMAXINFO:
		mmi := (*win.MINMAXINFO)(unsafe.Pointer(lParam))

		// 使用 SetMinSize 设置的全局最小尺寸
		if _this.minWidth > 0 {
			mmi.PtMinTrackSize.X = _this.minWidth
		}
		if _this.minHeight > 0 {
			mmi.PtMinTrackSize.Y = _this.minHeight
		}

		return 0
	case win.WM_DESTROY:
		delete(_this.handleWnds, hWnd)
		delete(_this.forms, hWnd)
		delete(_this.watchAll, hWnd)
		if hWnd == _this.mainId {
			_this.Exit(0)
		}
	}
	rs := win.DefWindowProc(hWnd, msg, wParam, lParam)
	runtime.UnlockOSThread()
	return rs
}

// 退出应用程序
func (_this *Provider) Exit(code int) {
	win.PostQuitMessage(int32(code))
}

// 运行主窗体消息循环
func (_this *Provider) RunMain(form br.Form) {
	runtime.LockOSThread()
	form.Show()
	var message win.MSG
	for {
		if win.GetMessage(&message, 0, 0, 0) {
			win.TranslateMessage(&message)
			win.DispatchMessage(&message)
		} else {
			break
		}
	}
	runtime.UnlockOSThread()
	os.Exit(0)
}
