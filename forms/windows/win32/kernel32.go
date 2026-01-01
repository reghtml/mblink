// Copyright 2010 The win Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build windows
// +build windows

package win32

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// MAX_PATH 最大路径长度
const MAX_PATH = 260

// 错误代码
const (
	ERROR_SUCCESS             = 0
	ERROR_INVALID_FUNCTION    = 1
	ERROR_FILE_NOT_FOUND      = 2
	ERROR_INVALID_PARAMETER   = 87
	ERROR_INSUFFICIENT_BUFFER = 122
	ERROR_MORE_DATA           = 234
)

// GlobalAlloc 标志常量
const (
	GHND          = 0x0042
	GMEM_FIXED    = 0x0000
	GMEM_MOVEABLE = 0x0002
	GMEM_ZEROINIT = 0x0040
	GPTR          = GMEM_FIXED | GMEM_ZEROINIT
)

// 预定义区域设置 ID
const (
	LOCALE_CUSTOM_DEFAULT     LCID = 0x0c00
	LOCALE_CUSTOM_UI_DEFAULT  LCID = 0x1400
	LOCALE_CUSTOM_UNSPECIFIED LCID = 0x1000
	LOCALE_INVARIANT          LCID = 0x007f
	LOCALE_USER_DEFAULT       LCID = 0x0400
	LOCALE_SYSTEM_DEFAULT     LCID = 0x0800
)

// LCTYPE 常量
const (
	LOCALE_SDECIMAL          LCTYPE = 14
	LOCALE_STHOUSAND         LCTYPE = 15
	LOCALE_SISO3166CTRYNAME  LCTYPE = 0x5a
	LOCALE_SISO3166CTRYNAME2 LCTYPE = 0x68
	LOCALE_SISO639LANGNAME   LCTYPE = 0x59
	LOCALE_SISO639LANGNAME2  LCTYPE = 0x67
)

var (
	// Library
	libkernel32 *windows.LazyDLL

	// Functions
	activateActCtx                     *windows.LazyProc
	closeHandle                        *windows.LazyProc
	createActCtx                       *windows.LazyProc
	fileTimeToSystemTime               *windows.LazyProc
	findResource                       *windows.LazyProc
	getConsoleTitle                    *windows.LazyProc
	getConsoleWindow                   *windows.LazyProc
	getCurrentThreadId                 *windows.LazyProc
	getLastError                       *windows.LazyProc
	getLocaleInfo                      *windows.LazyProc
	getLogicalDriveStrings             *windows.LazyProc
	getModuleHandle                    *windows.LazyProc
	getNumberFormat                    *windows.LazyProc
	getPhysicallyInstalledSystemMemory *windows.LazyProc
	getProfileString                   *windows.LazyProc
	getThreadLocale                    *windows.LazyProc
	getThreadUILanguage                *windows.LazyProc
	getVersion                         *windows.LazyProc
	globalAlloc                        *windows.LazyProc
	globalFree                         *windows.LazyProc
	globalLock                         *windows.LazyProc
	globalUnlock                       *windows.LazyProc
	moveMemory                         *windows.LazyProc
	mulDiv                             *windows.LazyProc
	loadResource                       *windows.LazyProc
	lockResource                       *windows.LazyProc
	setLastError                       *windows.LazyProc
	sizeofResource                     *windows.LazyProc
	systemTimeToFileTime               *windows.LazyProc
)

// Windows 句柄和标识符类型
type (
	ATOM          uint16  // 原子类型
	HANDLE        uintptr // 句柄类型
	HGLOBAL       HANDLE  // 全局内存句柄
	HINSTANCE     HANDLE  // 实例句柄
	LCID          uint32  // 区域设置 ID
	LCTYPE        uint32  // 区域设置类型
	LANGID        uint16  // 语言 ID
	HMODULE       uintptr // 模块句柄
	HWINEVENTHOOK HANDLE  // Windows 事件钩子句柄
	HRSRC         uintptr // 资源句柄
)

// FILETIME 文件时间结构
type FILETIME struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}

// NUMBERFMT 数字格式化结构
type NUMBERFMT struct {
	NumDigits     uint32
	LeadingZero   uint32
	Grouping      uint32
	LpDecimalSep  *uint16
	LpThousandSep *uint16
	NegativeOrder uint32
}

// SYSTEMTIME 系统时间结构
type SYSTEMTIME struct {
	WYear         uint16
	WMonth        uint16
	WDayOfWeek    uint16
	WDay          uint16
	WHour         uint16
	WMinute       uint16
	WSecond       uint16
	WMilliseconds uint16
}

// ACTCTX 激活上下文结构
type ACTCTX struct {
	size                  uint32
	Flags                 uint32
	Source                *uint16 // UTF-16 string
	ProcessorArchitecture uint16
	LangID                uint16
	AssemblyDirectory     *uint16 // UTF-16 string
	ResourceName          *uint16 // UTF-16 string
	ApplicationName       *uint16 // UTF-16 string
	Module                HMODULE
}

func init() {
	// Library
	libkernel32 = windows.NewLazySystemDLL("kernel32.dll")

	// Functions
	activateActCtx = libkernel32.NewProc("ActivateActCtx")
	closeHandle = libkernel32.NewProc("CloseHandle")
	createActCtx = libkernel32.NewProc("CreateActCtxW")
	fileTimeToSystemTime = libkernel32.NewProc("FileTimeToSystemTime")
	findResource = libkernel32.NewProc("FindResourceW")
	getConsoleTitle = libkernel32.NewProc("GetConsoleTitleW")
	getConsoleWindow = libkernel32.NewProc("GetConsoleWindow")
	getCurrentThreadId = libkernel32.NewProc("GetCurrentThreadId")
	getLastError = libkernel32.NewProc("GetLastError")
	getLocaleInfo = libkernel32.NewProc("GetLocaleInfoW")
	getLogicalDriveStrings = libkernel32.NewProc("GetLogicalDriveStringsW")
	getModuleHandle = libkernel32.NewProc("GetModuleHandleW")
	getNumberFormat = libkernel32.NewProc("GetNumberFormatW")
	getPhysicallyInstalledSystemMemory = libkernel32.NewProc("GetPhysicallyInstalledSystemMemory")
	getProfileString = libkernel32.NewProc("GetProfileStringW")
	getThreadLocale = libkernel32.NewProc("GetThreadLocale")
	getThreadUILanguage = libkernel32.NewProc("GetThreadUILanguage")
	getVersion = libkernel32.NewProc("GetVersion")
	globalAlloc = libkernel32.NewProc("GlobalAlloc")
	globalFree = libkernel32.NewProc("GlobalFree")
	globalLock = libkernel32.NewProc("GlobalLock")
	globalUnlock = libkernel32.NewProc("GlobalUnlock")
	moveMemory = libkernel32.NewProc("RtlMoveMemory")
	mulDiv = libkernel32.NewProc("MulDiv")
	loadResource = libkernel32.NewProc("LoadResource")
	lockResource = libkernel32.NewProc("LockResource")
	setLastError = libkernel32.NewProc("SetLastError")
	sizeofResource = libkernel32.NewProc("SizeofResource")
	systemTimeToFileTime = libkernel32.NewProc("SystemTimeToFileTime")
}

// 激活激活上下文
func ActivateActCtx(ctx HANDLE) (uintptr, bool) {
	var cookie uintptr
	ret, _, _ := syscall.Syscall(activateActCtx.Addr(), 2,
		uintptr(ctx),
		uintptr(unsafe.Pointer(&cookie)),
		0)
	return cookie, ret != 0
}

// 关闭对象句柄
func CloseHandle(hObject HANDLE) bool {
	ret, _, _ := syscall.Syscall(closeHandle.Addr(), 1,
		uintptr(hObject),
		0,
		0)

	return ret != 0
}

// 创建激活上下文
func CreateActCtx(ctx *ACTCTX) HANDLE {
	if ctx != nil {
		ctx.size = uint32(unsafe.Sizeof(*ctx))
	}
	ret, _, _ := syscall.Syscall(
		createActCtx.Addr(),
		1,
		uintptr(unsafe.Pointer(ctx)),
		0,
		0)
	return HANDLE(ret)
}

// 将文件时间转换为系统时间
func FileTimeToSystemTime(lpFileTime *FILETIME, lpSystemTime *SYSTEMTIME) bool {
	ret, _, _ := syscall.Syscall(fileTimeToSystemTime.Addr(), 2,
		uintptr(unsafe.Pointer(lpFileTime)),
		uintptr(unsafe.Pointer(lpSystemTime)),
		0)

	return ret != 0
}

// 查找资源
func FindResource(hModule HMODULE, lpName, lpType *uint16) HRSRC {
	ret, _, _ := syscall.Syscall(findResource.Addr(), 3,
		uintptr(hModule),
		uintptr(unsafe.Pointer(lpName)),
		uintptr(unsafe.Pointer(lpType)))

	return HRSRC(ret)
}

// 获取控制台窗口标题
func GetConsoleTitle(lpConsoleTitle *uint16, nSize uint32) uint32 {
	ret, _, _ := syscall.Syscall(getConsoleTitle.Addr(), 2,
		uintptr(unsafe.Pointer(lpConsoleTitle)),
		uintptr(nSize),
		0)

	return uint32(ret)
}

// 获取控制台窗口句柄
func GetConsoleWindow() HWND {
	ret, _, _ := syscall.Syscall(getConsoleWindow.Addr(), 0,
		0,
		0,
		0)

	return HWND(ret)
}

// 获取当前线程 ID
func GetCurrentThreadId() uint32 {
	ret, _, _ := syscall.Syscall(getCurrentThreadId.Addr(), 0,
		0,
		0,
		0)

	return uint32(ret)
}

// 获取最后一个错误代码
func GetLastError() uint32 {
	ret, _, _ := syscall.Syscall(getLastError.Addr(), 0,
		0,
		0,
		0)

	return uint32(ret)
}

// 获取区域设置信息
func GetLocaleInfo(Locale LCID, LCType LCTYPE, lpLCData *uint16, cchData int32) int32 {
	ret, _, _ := syscall.Syscall6(getLocaleInfo.Addr(), 4,
		uintptr(Locale),
		uintptr(LCType),
		uintptr(unsafe.Pointer(lpLCData)),
		uintptr(cchData),
		0,
		0)

	return int32(ret)
}

// 获取逻辑驱动器字符串
func GetLogicalDriveStrings(nBufferLength uint32, lpBuffer *uint16) uint32 {
	ret, _, _ := syscall.Syscall(getLogicalDriveStrings.Addr(), 2,
		uintptr(nBufferLength),
		uintptr(unsafe.Pointer(lpBuffer)),
		0)

	return uint32(ret)
}

// 获取模块句柄
func GetModuleHandle(lpModuleName *uint16) HINSTANCE {
	ret, _, _ := syscall.Syscall(getModuleHandle.Addr(), 1,
		uintptr(unsafe.Pointer(lpModuleName)),
		0,
		0)

	return HINSTANCE(ret)
}

// 格式化数字字符串
func GetNumberFormat(Locale LCID, dwFlags uint32, lpValue *uint16, lpFormat *NUMBERFMT, lpNumberStr *uint16, cchNumber int32) int32 {
	ret, _, _ := syscall.Syscall6(getNumberFormat.Addr(), 6,
		uintptr(Locale),
		uintptr(dwFlags),
		uintptr(unsafe.Pointer(lpValue)),
		uintptr(unsafe.Pointer(lpFormat)),
		uintptr(unsafe.Pointer(lpNumberStr)),
		uintptr(cchNumber))

	return int32(ret)
}

// 获取物理安装的系统内存大小
func GetPhysicallyInstalledSystemMemory(totalMemoryInKilobytes *uint64) bool {
	if getPhysicallyInstalledSystemMemory.Find() != nil {
		return false
	}
	ret, _, _ := syscall.Syscall(getPhysicallyInstalledSystemMemory.Addr(), 1,
		uintptr(unsafe.Pointer(totalMemoryInKilobytes)),
		0,
		0)

	return ret != 0
}

// 从配置文件获取字符串
func GetProfileString(lpAppName, lpKeyName, lpDefault *uint16, lpReturnedString uintptr, nSize uint32) bool {
	ret, _, _ := syscall.Syscall6(getProfileString.Addr(), 5,
		uintptr(unsafe.Pointer(lpAppName)),
		uintptr(unsafe.Pointer(lpKeyName)),
		uintptr(unsafe.Pointer(lpDefault)),
		lpReturnedString,
		uintptr(nSize),
		0)
	return ret != 0
}

// 获取线程区域设置
func GetThreadLocale() LCID {
	ret, _, _ := syscall.Syscall(getThreadLocale.Addr(), 0,
		0,
		0,
		0)

	return LCID(ret)
}

// 获取线程 UI 语言
func GetThreadUILanguage() LANGID {
	if getThreadUILanguage.Find() != nil {
		return 0
	}

	ret, _, _ := syscall.Syscall(getThreadUILanguage.Addr(), 0,
		0,
		0,
		0)

	return LANGID(ret)
}

// 获取操作系统版本信息
func GetVersion() uint32 {
	ret, _, _ := syscall.Syscall(getVersion.Addr(), 0,
		0,
		0,
		0)
	return uint32(ret)
}

// 分配全局内存
func GlobalAlloc(uFlags uint32, dwBytes uintptr) HGLOBAL {
	ret, _, _ := syscall.Syscall(globalAlloc.Addr(), 2,
		uintptr(uFlags),
		dwBytes,
		0)

	return HGLOBAL(ret)
}

// 释放全局内存
func GlobalFree(hMem HGLOBAL) HGLOBAL {
	ret, _, _ := syscall.Syscall(globalFree.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return HGLOBAL(ret)
}

// 锁定全局内存并返回指针
func GlobalLock(hMem HGLOBAL) unsafe.Pointer {
	ret, _, _ := syscall.Syscall(globalLock.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return unsafe.Pointer(ret)
}

// 解锁全局内存
func GlobalUnlock(hMem HGLOBAL) bool {
	ret, _, _ := syscall.Syscall(globalUnlock.Addr(), 1,
		uintptr(hMem),
		0,
		0)

	return ret != 0
}

// 移动内存块
func MoveMemory(destination, source unsafe.Pointer, length uintptr) {
	syscall.Syscall(moveMemory.Addr(), 3,
		uintptr(unsafe.Pointer(destination)),
		uintptr(source),
		uintptr(length))
}

// 执行整数乘法和除法运算
func MulDiv(nNumber, nNumerator, nDenominator int32) int32 {
	ret, _, _ := syscall.Syscall(mulDiv.Addr(), 3,
		uintptr(nNumber),
		uintptr(nNumerator),
		uintptr(nDenominator))

	return int32(ret)
}

// 加载资源
func LoadResource(hModule HMODULE, hResInfo HRSRC) HGLOBAL {
	ret, _, _ := syscall.Syscall(loadResource.Addr(), 2,
		uintptr(hModule),
		uintptr(hResInfo),
		0)

	return HGLOBAL(ret)
}

// 锁定资源内存
func LockResource(hResData HGLOBAL) uintptr {
	ret, _, _ := syscall.Syscall(lockResource.Addr(), 1,
		uintptr(hResData),
		0,
		0)

	return ret
}

// 设置最后一个错误代码
func SetLastError(dwErrorCode uint32) {
	syscall.Syscall(setLastError.Addr(), 1,
		uintptr(dwErrorCode),
		0,
		0)
}

// 获取资源大小
func SizeofResource(hModule HMODULE, hResInfo HRSRC) uint32 {
	ret, _, _ := syscall.Syscall(sizeofResource.Addr(), 2,
		uintptr(hModule),
		uintptr(hResInfo),
		0)

	return uint32(ret)
}

// 将系统时间转换为文件时间
func SystemTimeToFileTime(lpSystemTime *SYSTEMTIME, lpFileTime *FILETIME) bool {
	ret, _, _ := syscall.Syscall(systemTimeToFileTime.Addr(), 2,
		uintptr(unsafe.Pointer(lpSystemTime)),
		uintptr(unsafe.Pointer(lpFileTime)),
		0)

	return ret != 0
}
