package bridge

import (
	fm "github.com/reghtml/mblink/forms"
)

// FormStateProc 窗体状态变化回调函数类型
type FormStateProc func(state fm.FormState)

// FormActiveProc 窗体激活回调函数类型
type FormActiveProc func()

// Form 窗体接口
type Form interface {
	Controls

	Close()
	ShowDialog()
	SetTitle(title string)
	SetBorderStyle(style fm.FormBorder)
	ShowToMax()
	ShowToMin()
	/*
		允许在无边框模式下调整窗体大小
	*/
	NoneBorderResize()
	Active()

	SetMaximizeBox(isShow bool)
	SetMinimizeBox(isShow bool)
	SetIcon(iconFile string)
	SetIconData(iconData []byte)
	SetTopMost(isTop bool)

	SetOnState(proc FormStateProc) FormStateProc
	SetOnActive(proc FormActiveProc) FormActiveProc
}
