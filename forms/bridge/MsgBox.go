package bridge

import fm "github.com/reghtml/mblink/forms"

// MsgBox 消息框接口
type MsgBox interface {
	Show(param fm.MsgBoxParam) fm.MsgBoxResult
}
