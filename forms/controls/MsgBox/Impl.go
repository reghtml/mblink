package MsgBox

import (
	fm "github.com/reghtml/mblink/forms"
	cs "github.com/reghtml/mblink/forms/controls"
)

// 显示消息框
func Show(param fm.MsgBoxParam) fm.MsgBoxResult {
	return cs.App.NewMsgBox().Show(param)
}

// 显示信息消息框
func ShowInfo(title, text string) {
	Show(fm.MsgBoxParam{
		Title:  title,
		Text:   text,
		Icon:   fm.MsgBoxIcon_Info,
		Button: fm.MsgBoxButton_Ok,
	})
}

// 显示问题消息框
func ShowQuestion(title, text string) fm.MsgBoxResult {
	return Show(fm.MsgBoxParam{
		Title:  title,
		Text:   text,
		Icon:   fm.MsgBoxIcon_Question,
		Button: fm.MsgBoxButton_YesNo,
	})
}

// 显示警告消息框
func ShowWarn(title, text string) {
	Show(fm.MsgBoxParam{
		Title:  title,
		Text:   text,
		Icon:   fm.MsgBoxIcon_Warn,
		Button: fm.MsgBoxButton_Ok,
	})
}

// 显示错误消息框
func ShowError(title, text string) {
	Show(fm.MsgBoxParam{
		Title:  title,
		Text:   text,
		Icon:   fm.MsgBoxIcon_Error,
		Button: fm.MsgBoxButton_Ok,
	})
}
