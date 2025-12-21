package bridge

import fm "github.com/reghtml/mblink/forms"

type MsgBox interface {
	Show(param fm.MsgBoxParam) fm.MsgBoxResult
}
