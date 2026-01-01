package controls

import (
	"github.com/reghtml/mblink/forms/bridge"
)

// MainForm 主窗体接口
type MainForm interface {
	getFormImpl() bridge.Form
}

var App bridge.Provider

func Run(form MainForm) {
	App.RunMain(form.getFormImpl())
}
