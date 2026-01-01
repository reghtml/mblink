package controls

import (
	fm "github.com/reghtml/mblink/forms"
	br "github.com/reghtml/mblink/forms/bridge"
)

// Control 控件基类
type Control struct {
	BaseUI

	impl   br.Control
	anchor fm.AnchorStyle
}

// 初始化控件
func (_this *Control) Init() *Control {
	_this.impl = App.NewControl()
	_this.BaseUI.Init(_this, _this.impl)
	_this.anchor = fm.AnchorStyle_Left | fm.AnchorStyle_Top
	_this.SetSize(200, 200)
	return _this
}

// 获取控件的锚定样式
func (_this *Control) GetAnchor() fm.AnchorStyle {
	return _this.anchor
}

// 设置控件的锚定样式
func (_this *Control) SetAnchor(style fm.AnchorStyle) {
	_this.anchor = style
	if _this.parent != nil {
		if ct, ok := _this.parent.(Container); ok {
			bn := ct.GetBound()
			ct.SetSize(bn.Width, bn.Height-1)
			ct.SetSize(bn.Width, bn.Height)
		}
	}
}

func (_this *Control) toControl() br.Control {
	return _this.impl
}

func (_this *Control) setParent(parent GUI) {
	_this.parent = parent
}
func (_this *Control) setOwner(owner GUI) {
	_this.owner = owner
}
