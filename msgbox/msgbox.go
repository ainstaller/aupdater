package msgbox

import (
	"syscall"

	"github.com/lxn/win"
)

type MsgBoxStyle uint

const (
	OK                MsgBoxStyle = win.MB_OK
	OKCancel          MsgBoxStyle = win.MB_OKCANCEL
	AbortRetryIgnore  MsgBoxStyle = win.MB_ABORTRETRYIGNORE
	YesNoCancel       MsgBoxStyle = win.MB_YESNOCANCEL
	YesNo             MsgBoxStyle = win.MB_YESNO
	RetryCancel       MsgBoxStyle = win.MB_RETRYCANCEL
	CancelTryContinue MsgBoxStyle = win.MB_CANCELTRYCONTINUE
	IconHand          MsgBoxStyle = win.MB_ICONHAND
	IconQuestion      MsgBoxStyle = win.MB_ICONQUESTION
	IconExclamation   MsgBoxStyle = win.MB_ICONEXCLAMATION
	IconAsterisk      MsgBoxStyle = win.MB_ICONASTERISK
	UserIcon          MsgBoxStyle = win.MB_USERICON
	IconWarning       MsgBoxStyle = win.MB_ICONWARNING
	IconError         MsgBoxStyle = win.MB_ICONERROR
	IconInformation   MsgBoxStyle = win.MB_ICONINFORMATION
	IconStop          MsgBoxStyle = win.MB_ICONSTOP
	DefButton1        MsgBoxStyle = win.MB_DEFBUTTON1
	DefButton2        MsgBoxStyle = win.MB_DEFBUTTON2
	DefButton3        MsgBoxStyle = win.MB_DEFBUTTON3
	DefButton4        MsgBoxStyle = win.MB_DEFBUTTON4
)

func New(title, message string, style MsgBoxStyle) int {
	var ownerHWnd win.HWND

	return int(win.MessageBox(
		ownerHWnd,
		syscall.StringToUTF16Ptr(message),
		syscall.StringToUTF16Ptr(title),
		uint32(style)))
}
