package main

import (
	gm "github.com/reghtml/mblink"
	fm "github.com/reghtml/mblink/forms"
	cs "github.com/reghtml/mblink/forms/controls"
	gw "github.com/reghtml/mblink/forms/windows"
)

func main() {
	cs.App = new(gw.Provider).Init()
	cs.App.SetIcon("app.ico")

	//绑定了miniblink的窗体，内部实现了一些用js控制窗体的功能
	frm := new(gm.MiniblinkForm).Init()
	frm.SetTitle("miniblink窗口")
	frm.SetLocation(100, 100)
	frm.SetSize(800, 500)
	frm.SetBorderStyle(fm.FormBorder_None)
	frm.NoneBorderResize()
	frm.View.ResourceLoader = append(frm.View.ResourceLoader, new(gm.FileLoader).Init("Res", "local"))
	frm.EvLoad["show"] = func(s cs.GUI) {
		frm.View.LoadUri("http://local/window.html")
	}
	cs.Run(&frm.Form)
}
