package windows

import (
	br "github.com/reghtml/mblink/forms/bridge"
)

var (
	msgbox = new(winMsgBox).init()
)

// 创建新窗体
func (_this *Provider) NewForm(param br.FormParam) br.Form {
	return new(winForm).init(_this, param)
}

// 创建新控件
func (_this *Provider) NewControl() br.Control {
	return new(winControl).init(_this)
}

// 创建新消息框
func (_this *Provider) NewMsgBox() br.MsgBox {
	return msgbox
}
