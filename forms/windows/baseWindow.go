package windows

import (
	"github.com/reghtml/mblink/forms/windows/win32"
)

// baseWindow 窗口基接口
type baseWindow interface {
	hWnd() win32.HWND
	onWndMsg(hWnd win32.HWND, msg uint32, wParam uintptr, lParam uintptr) uintptr
}
