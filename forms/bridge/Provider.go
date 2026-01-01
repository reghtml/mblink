package bridge

import (
	f "github.com/reghtml/mblink/forms"
)

// FormParam 窗体创建参数
type FormParam struct {
	HideInTaskbar bool
	HideIcon      bool
}

// Provider 平台提供者接口，提供跨平台的窗体和控制功能
type Provider interface {
	RunMain(form Form)
	Exit(code int)
	SetIcon(file string)
	SetIconData(iconData []byte)
	SetMinSize(width, height int) // ✅ 设置最小尺寸
	GetScreen() f.Screen
	ModifierKeys() map[f.Keys]bool
	MouseIsDown() map[f.MouseButtons]bool
	MouseLocation() f.Point
	AppDir() string

	NewForm(param FormParam) Form
	NewControl() Control
	NewMsgBox() MsgBox
}
