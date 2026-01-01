package windows

import (
	br "github.com/reghtml/mblink/forms/bridge"
	win "github.com/reghtml/mblink/forms/windows/win32"
)

// winContainer Windows 容器实现
type winContainer struct {
	winBase
	_self  br.Control
	_ctrls []br.Control
}

// 初始化容器
func (_this *winContainer) init(provider *Provider, self br.Control) *winContainer {
	_this.winBase.init(provider)
	_this._self = self
	return _this
}

// 获取所有子控件
func (_this *winContainer) GetChilds() []br.Control {
	return _this._ctrls
}

// 添加子控件
func (_this *winContainer) AddControl(control br.Control) {
	if win.SetParent(win.HWND(control.GetHandle()), win.HWND(_this.GetHandle())) != 0 {
		if wc, ok := control.(*winControl); ok {
			wc.parent = _this._self
			if ow, ok := _this._self.(br.Form); ok {
				wc.owner = ow
			} else {
				wc.owner = _this.GetOwner()
			}
		}
	}
	_this._ctrls = append(_this._ctrls, control)
	control.Show()
}

// 移除子控件
func (_this *winContainer) RemoveControl(control br.Control) {
	for i, n := range _this._ctrls {
		if n.GetHandle() == control.GetHandle() {
			hWnd := win.HWND(control.GetHandle())
			if win.SetParent(hWnd, _this.app.defOwner) != 0 {
				win.SendMessage(hWnd, win.WM_DESTROY, 0, 0)
				if wc, ok := control.(*winControl); ok {
					wc.parent = nil
					wc.owner = nil
				}
			}
			_this._ctrls = append(_this._ctrls[0:i], _this._ctrls[i+1:]...)
			break
		}
	}
}
