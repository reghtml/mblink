package windows

import (
	"unsafe"

	fm "github.com/reghtml/mblink/forms"
	br "github.com/reghtml/mblink/forms/bridge"
	win "github.com/reghtml/mblink/forms/windows/win32"
)

// winForm Windows 窗体实现
type winForm struct {
	winContainer
	_onClose           func() (cancel bool)
	_onState           br.FormStateProc
	_onActive          br.FormActiveProc
	border             fm.FormBorder
	_isModal           bool
	restoreRect        win.RECT // 保存恢复时使用的窗口位置和大小（用于无边框窗口）
	hasRestoreRect     bool     // 是否有保存的恢复位置和大小
	_isCustomMaximized bool     // 是否处于自定义最大化状态（用于无边框窗口）
}

// 初始化窗体
func (_this *winForm) init(provider *Provider, param br.FormParam) *winForm {
	_this.winContainer.init(provider, _this)
	_this.onWndProc = _this.msgProc
	parent := win.HWND(0)
	exStyle := win.WS_EX_APPWINDOW | win.WS_EX_CONTROLPARENT
	if param.HideInTaskbar {
		exStyle &= ^win.WS_EX_APPWINDOW
		parent = _this.app.defOwner
	}
	if param.HideIcon {
		exStyle |= win.WS_EX_DLGMODALFRAME
	}
	x := 100 + len(_this.app.forms)*25
	y := 100 + len(_this.app.forms)*25
	win.CreateWindowEx(
		uint64(exStyle),
		sto16(_this.app.className),
		sto16(""),
		win.WS_OVERLAPPEDWINDOW,
		int32(x), int32(y), 200, 300, parent, 0, _this.app.hInstance, unsafe.Pointer(_this))
	return _this
}

// 处理窗口消息
func (_this *winForm) msgProc(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	rs := _this.winBase.msgProc(hWnd, msg, wParam, lParam)
	if rs != 0 {
		return rs
	}
	switch msg {
	case win.WM_CREATE:
		_this.app.forms[hWnd] = _this
	case win.WM_DESTROY:
		if _this._isModal {
			win.PostQuitMessage(0)
			return 1
		}
	case win.WM_GETMINMAXINFO:
		// 对于无边框窗口，设置最大化尺寸为工作区域（不覆盖任务栏）
		if _this.border == fm.FormBorder_None {
			mmi := (*win.MINMAXINFO)(unsafe.Pointer(lParam))
			hMonitor := win.MonitorFromWindow(hWnd, win.MONITOR_DEFAULTTONEAREST)
			if hMonitor != 0 {
				var mi win.MONITORINFO
				mi.CbSize = uint32(unsafe.Sizeof(mi))
				if win.GetMonitorInfo(hMonitor, &mi) {
					rcWork := mi.RcWork
					rcMonitor := mi.RcMonitor
					// 计算工作区域相对于监视器的位置
					mmi.PtMaxPosition.X = rcWork.Left - rcMonitor.Left
					mmi.PtMaxPosition.Y = rcWork.Top - rcMonitor.Top
					// 设置最大尺寸为工作区域大小
					mmi.PtMaxSize.X = rcWork.Right - rcWork.Left
					mmi.PtMaxSize.Y = rcWork.Bottom - rcWork.Top
					return 0
				}
			}
		}
	case win.WM_SIZE:
		if _this._onState != nil {
			switch int(wParam) {
			case win.SIZE_RESTORED:
				_this._onState(fm.FormState_Normal)
			case win.SIZE_MAXIMIZED:
				_this._onState(fm.FormState_Max)
			case win.SIZE_MINIMIZED:
				_this._onState(fm.FormState_Min)
			}
		}
	case win.WM_ACTIVATE:
		if _this._onActive != nil {
			_this._onActive()
		}
	case win.WM_CLOSE:
		if _this._onClose != nil && _this._onClose() {
			rs = 1
		}
		if rs != 0 && _this._isModal {
			_this._isModal = false
		}
	}
	return rs
}

// 设置窗体是否置顶
func (_this *winForm) SetTopMost(isTop bool) {
	if isTop {
		win.SetWindowPos(_this.handle, win.HWND_TOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
	} else {
		win.SetWindowPos(_this.handle, win.HWND_NOTOPMOST, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE)
	}
}

// 设置窗体激活回调函数
func (_this *winForm) SetOnActive(proc br.FormActiveProc) br.FormActiveProc {
	pre := _this._onActive
	_this._onActive = proc
	return pre
}

// 启用无边框窗口的调整大小功能
func (_this *winForm) NoneBorderResize() {
	padd := 5
	rsState := new(int)
	_this.app.watch(_this, func(hWnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
		if _this.border != fm.FormBorder_None {
			return 0
		}
		switch msg {
		case win.WM_MOUSEMOVE:
			if hWnd != _this.handle {
				if wnd, ok := _this.app.handleWnds[hWnd].(br.Control); ok {
					if wnd.GetOwner().GetHandle() != uintptr(_this.handle) {
						return 0
					}
				} else {
					return 0
				}
			}
			sz := _this.GetBound().Rect
			p := _this.ToClientPoint(_this.app.MouseLocation())
			if p.X <= padd {
				if p.Y <= padd {
					*rsState = 7
				} else if p.Y+padd >= sz.Height {
					*rsState = 1
				} else {
					*rsState = 4
				}
			} else if p.Y <= padd {
				if p.X <= padd {
					*rsState = 7
				} else if p.X+padd >= sz.Width {
					*rsState = 9
				} else {
					*rsState = 8
				}
			} else if p.X+padd >= sz.Width {
				if p.Y <= padd {
					*rsState = 9
				} else if p.Y+padd >= sz.Height {
					*rsState = 3
				} else {
					*rsState = 6
				}
			} else if p.Y+padd >= sz.Height {
				if p.X <= padd {
					*rsState = 1
				} else if p.X+padd >= sz.Width {
					*rsState = 3
				} else {
					*rsState = 2
				}
			} else {
				*rsState = 0
			}
		case win.WM_SETCURSOR:
			cur := fm.CursorType_Default
			switch *rsState {
			case 8, 2:
				cur = fm.CursorType_SIZENS
			case 4, 6:
				cur = fm.CursorType_SIZEWE
			case 7, 3:
				cur = fm.CursorType_SIZENWSE
			case 9, 1:
				cur = fm.CursorType_SIZENESW
			}
			if cur != fm.CursorType_Default {
				res := win.MAKEINTRESOURCE(uintptr(toWinCursor(cur)))
				win.SetCursor(win.LoadCursor(0, res))
				return 1
			} else {
				return 0
			}
		case win.WM_LBUTTONDOWN:
			switch *rsState {
			case 4:
				win.SendMessage(_this.handle, win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|0xF001), lParam)
			case 6:
				win.SendMessage(_this.handle, win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|0xF002), lParam)
			case 8:
				win.SendMessage(_this.handle, win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|0xF003), lParam)
			case 7:
				win.SendMessage(_this.handle, win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|0xF004), lParam)
			case 9:
				win.SendMessage(_this.handle, win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|0xF005), lParam)
			case 2:
				win.SendMessage(_this.handle, win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|0xF006), lParam)
			case 1:
				win.SendMessage(_this.handle, win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|0xF007), lParam)
			case 3:
				win.SendMessage(_this.handle, win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|0xF008), lParam)
			default:
				return 0
			}
			return 1
		}
		return 0
	})
}

// 显示窗体
func (_this *winForm) Show() {
	// 如果是无边框窗口且已经最大化，使用保存的位置恢复
	if _this.border == fm.FormBorder_None && _this._isCustomMaximized && _this.hasRestoreRect {
		restoreWidth := _this.restoreRect.Right - _this.restoreRect.Left
		restoreHeight := _this.restoreRect.Bottom - _this.restoreRect.Top
		win.SetWindowPos(
			_this.handle,
			0,
			_this.restoreRect.Left,
			_this.restoreRect.Top,
			restoreWidth,
			restoreHeight,
			win.SWP_NOZORDER|win.SWP_NOACTIVATE|win.SWP_FRAMECHANGED,
		)
		win.UpdateWindow(_this.handle)
		_this._isCustomMaximized = false
		return
	}

	isMax := win.IsZoomed(_this.handle)
	isMin := win.IsIconic(_this.handle)
	if isMax || isMin {
		win.ShowWindow(_this.handle, win.SW_RESTORE)
	} else {
		win.ShowWindow(_this.handle, win.SW_SHOW)
	}
	win.UpdateWindow(_this.handle)
}

// 关闭窗体
func (_this *winForm) Close() {
	win.SendMessage(_this.handle, win.WM_CLOSE, 0, 0)
}

// 以模态方式显示窗体
func (_this *winForm) ShowDialog() {
	acHwnd := win.GetActiveWindow()
	if acfm, ok := _this.app.forms[acHwnd].(*winForm); ok {
		acfm.Enable(false)
		_this._isModal = true
		_this.Show()
		var msg win.MSG
		for {
			if win.GetMessage(&msg, 0, 0, 0) && _this._isModal {
				win.TranslateMessage(&msg)
				win.DispatchMessage(&msg)
			} else {
				break
			}
		}
		acfm.Enable(true)
		acfm.Active()
	}
}

// getWorkArea 获取当前窗口所在监视器的工作区域（确保不覆盖任务栏）
func (_this *winForm) getWorkArea() win.RECT {
	var workRect win.RECT

	// 获取屏幕总尺寸，用于验证
	screenWidth := int32(win.GetSystemMetrics(win.SM_CXSCREEN))
	screenHeight := int32(win.GetSystemMetrics(win.SM_CYSCREEN))

	// 优先使用 GetMonitorInfo 获取当前窗口所在监视器的工作区域（支持多显示器，最准确）
	hMonitor := win.MonitorFromWindow(_this.handle, win.MONITOR_DEFAULTTONEAREST)
	if hMonitor != 0 {
		var mi win.MONITORINFO
		mi.CbSize = uint32(unsafe.Sizeof(mi))
		if win.GetMonitorInfo(hMonitor, &mi) {
			workRect = mi.RcWork
			workWidth := workRect.Right - workRect.Left
			workHeight := workRect.Bottom - workRect.Top
			monitorWidth := mi.RcMonitor.Right - mi.RcMonitor.Left
			monitorHeight := mi.RcMonitor.Bottom - mi.RcMonitor.Top
			// 验证：工作区域必须小于监视器区域（确保排除了任务栏）
			// 并且工作区域的坐标应该在监视器范围内
			if workWidth > 0 && workHeight > 0 &&
				(workWidth < monitorWidth || workHeight < monitorHeight) &&
				workRect.Left >= mi.RcMonitor.Left && workRect.Top >= mi.RcMonitor.Top &&
				workRect.Right <= mi.RcMonitor.Right && workRect.Bottom <= mi.RcMonitor.Bottom {
				return workRect
			}
		}
	}

	// 如果 GetMonitorInfo 失败或结果不合理，使用 SystemParametersInfo（只适用于主显示器）
	var testRect win.RECT
	if win.SystemParametersInfo(win.SPI_GETWORKAREA, 0, unsafe.Pointer(&testRect), 0) {
		testWidth := testRect.Right - testRect.Left
		testHeight := testRect.Bottom - testRect.Top
		// 验证：工作区域必须小于屏幕尺寸（确保排除了任务栏）
		// 并且坐标应该在屏幕范围内（通常从0,0或很小的正数开始）
		if testWidth > 0 && testHeight > 0 &&
			(testWidth < screenWidth || testHeight < screenHeight) &&
			testRect.Left >= 0 && testRect.Top >= 0 &&
			testRect.Right <= screenWidth && testRect.Bottom <= screenHeight {
			return testRect
		}
	}

	// 最后的后备方案：使用系统指标，如果等于屏幕尺寸则减去任务栏空间
	fullScreenWidth := win.GetSystemMetrics(win.SM_CXFULLSCREEN)
	fullScreenHeight := win.GetSystemMetrics(win.SM_CYFULLSCREEN)

	// 如果系统指标等于屏幕尺寸，说明没有排除任务栏，需要手动减去
	// 通常任务栏在底部，高度约为40-48像素；如果在侧边，宽度约为40-48像素
	taskbarSize := int32(48) // 保守的任务栏尺寸
	if fullScreenWidth == screenWidth {
		// 任务栏可能在右侧或左侧，减小宽度
		fullScreenWidth = screenWidth - taskbarSize
	}
	if fullScreenHeight == screenHeight {
		// 任务栏可能在底部或顶部，减小高度
		fullScreenHeight = screenHeight - taskbarSize
	}
	// 确保值有效
	if fullScreenWidth <= 0 {
		fullScreenWidth = screenWidth - taskbarSize
	}
	if fullScreenHeight <= 0 {
		fullScreenHeight = screenHeight - taskbarSize
	}

	return win.RECT{
		Left:   0,
		Top:    0,
		Right:  fullScreenWidth,
		Bottom: fullScreenHeight,
	}
}

// 最大化窗体
func (_this *winForm) ShowToMax() {
	// 现在通过 WM_GETMINMAXINFO 处理无边框窗口的最大化（自动使用工作区域）
	// 直接使用标准的 ShowWindow，系统会自动使用我们在 WM_GETMINMAXINFO 中设置的值
	win.ShowWindow(_this.handle, win.SW_MAXIMIZE)
	win.UpdateWindow(_this.handle)
}

// 最小化窗体
func (_this *winForm) ShowToMin() {
	win.ShowWindow(_this.handle, win.SW_MINIMIZE)
}

// 设置窗体标题
func (_this *winForm) SetTitle(title string) {
	win.SetWindowText(_this.handle, title)
}

// 设置窗体边框样式
func (_this *winForm) SetBorderStyle(border fm.FormBorder) {
	style := win.GetWindowLong(_this.handle, win.GWL_STYLE)
	switch border {
	case fm.FormBorder_Default:
		style |= win.WS_OVERLAPPEDWINDOW
	case fm.FormBorder_None:
		style &= ^win.WS_SIZEBOX & ^win.WS_CAPTION
	case fm.FormBorder_Disable_Resize:
		style &= ^win.WS_SIZEBOX
	}
	win.SetWindowLong(_this.handle, win.GWL_STYLE, style)
	bn := _this.GetBound()
	_this.SetSize(bn.Width, bn.Height-1)
	_this.SetSize(bn.Width, bn.Height)
	_this.border = border
}

// 设置窗体状态变化回调函数
func (_this *winForm) SetOnState(proc br.FormStateProc) br.FormStateProc {
	pre := _this._onState
	_this._onState = proc
	return pre
}

// 设置是否显示最大化按钮
func (_this *winForm) SetMaximizeBox(isShow bool) {
	style := win.GetWindowLong(_this.handle, win.GWL_STYLE)
	if isShow {
		style |= win.WS_MAXIMIZEBOX
	} else {
		style &= ^win.WS_MAXIMIZEBOX
	}
	win.SetWindowLong(_this.handle, win.GWL_STYLE, style)
}

// 设置是否显示最小化按钮
func (_this *winForm) SetMinimizeBox(isShow bool) {
	style := win.GetWindowLong(_this.handle, win.GWL_STYLE)
	if isShow {
		style |= win.WS_MINIMIZEBOX
	} else {
		style &= ^win.WS_MINIMIZEBOX
	}
	win.SetWindowLong(_this.handle, win.GWL_STYLE, style)
}

// 设置窗体图标（从文件）
func (_this *winForm) SetIcon(iconFile string) {
	style := win.GetWindowLong(_this.handle, win.GWL_EXSTYLE)
	if style&win.WS_EX_DLGMODALFRAME != 0 {
		return
	}
	h := win.LoadImage(_this.app.hInstance, sto16(iconFile), win.IMAGE_ICON, 0, 0, win.LR_LOADFROMFILE)
	if h != 0 {
		win.SendMessage(_this.handle, win.WM_SETICON, 1, uintptr(h))
		win.SendMessage(_this.handle, win.WM_SETICON, 0, uintptr(h))
	}
}

// 设置窗体图标（从数据）
func (_this *winForm) SetIconData(iconData []byte) {
	style := win.GetWindowLong(_this.handle, win.GWL_EXSTYLE)
	if style&win.WS_EX_DLGMODALFRAME != 0 {
		return
	}
	if len(iconData) == 0 {
		return
	}
	hIcon := _this.app.createIconFromData(iconData)
	if hIcon != 0 {
		win.SendMessage(_this.handle, win.WM_SETICON, 1, hIcon)
		win.SendMessage(_this.handle, win.WM_SETICON, 0, hIcon)
	}
}

// 启用或禁用窗体
func (_this *winForm) Enable(b bool) {
	_this.isEnable = b
	win.EnableWindow(_this.handle, b)
}

// 激活窗体
func (_this *winForm) Active() {
	win.SetActiveWindow(_this.handle)
}
