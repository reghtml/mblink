package windows

import (
	"unsafe"

	win "github.com/reghtml/mblink/forms/windows/win32"
)

// winControl Windows 控件实现
type winControl struct {
	winBase
}

// 初始化控件
func (_this *winControl) init(provider *Provider) *winControl {
	_this.winBase.init(provider)
	win.CreateWindowEx(
		0,
		sto16(provider.className),
		sto16(""),
		win.WS_CHILD, 0, 0, 100, 100, _this.app.defOwner, 0,
		provider.hInstance, unsafe.Pointer(_this))
	return _this
}
