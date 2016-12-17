package main

import (
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"

	"gopkg.in/toast.v1"

	"os"

	"github.com/aInstaller/aupdater/msgbox"
	"github.com/aInstaller/aupdater/version"
	"github.com/aInstaller/icon"
	"github.com/aInstaller/utils/steam"
	"github.com/franela/goreq"
	"github.com/getlantern/systray"
	"github.com/matishsiao/goInfo"
	"github.com/mitchellh/go-ps"
	"github.com/ungerik/go-dry"
)

const (
	mainMenuOverride = "https://raw.githubusercontent.com/n0kk/ahud/master/resource/ui/mainmenuoverride.res"
	installerPath    = `/../Local/Programs/aInstaller/aInstaller.exe`
)

var (
	wd       string
	iconPath string

	psapi                 = syscall.NewLazyDLL("psapi.dll")
	getModuleFileNameProc = psapi.NewProc("GetProcessImageFileNameW")
)

type Update struct{}

func (u *Update) Should() bool {
	// doesn't check for updates when tf2 is running
	if tf2, _ := IsTF2Running(); tf2 {
		return false
	}

	next := time.Time{}

	if dry.FileExists("./aupdater.json") {
		res, _ := dry.FileGetJSON("./aupdater.json", 10*time.Second)
		next, _ = time.Parse("", res.(map[string]interface{})["next"].(string))
	}

	if next.Before(time.Now()) {
		dry.FileSetJSON("./aupdater.json", map[string]interface{}{
			"next": time.Now().Add(8 * time.Hour),
		})

		return true
	}

	return false
}

func (u *Update) Check() error {
	if !u.Should() {
		return nil
	}

	err := u.Notify("Checking for ahud updates...")
	if err != nil {
		return err
	}

	req := goreq.Request{
		Method: "GET",
		Uri:    mainMenuOverride,
	}

	// send request
	res, err := req.Do()
	if err != nil {
		return err
	}

	// get body as string
	str, err := res.Body.ToString()
	if err != nil {
		return err
	}

	// parse latest version
	latest, err := version.New(str)
	if err != nil {
		return err
	}

	// get game path
	gamePath, err := steam.FindGame()
	if err != nil {
		return err
	}

	// path for local main menu override file
	mmoPath := filepath.Join(gamePath, "./custom/ahud/resource/ui/mainmenuoverride.res")
	log.Println(latest.Year, latest.Month, latest.Day)

	// read main menu override's data
	mmo, err := dry.FileGetString(mmoPath, 10*time.Second)
	if err != nil {
		return err
	}

	// parse current version
	current, err := version.New(mmo)
	if err != nil {
		return err
	}

	log.Println(latest.After(current))
	// compare version
	if latest.After(current) {
		if u.Ask("Update available, do you want to update ahud now?") {
			go func() {
				// build full path for ainstaller.exe
				ai := filepath.Join(os.Getenv("appdata"), installerPath)
				dir := filepath.Dir(ai)

				// change working dir to ainstaller's dir
				if err := os.Chdir(dir); err != nil {
					log.Fatal(err)
				}

				// run ainstaller.exe
				cmd := exec.Command(ai)
				err := cmd.Start()
				if err != nil {
					log.Fatal(err)
				}
			}()
		}
	} else {
		u.Notify("No update available yet.")
	}

	return nil
}

func (u *Update) Ask(question string) bool {
	m := msgbox.New("aInstaller", question, msgbox.YesNo)
	log.Println(question, m == 6)

	return m == 6 // 6 == yes
}

func (u *Update) Notify(msg string) error {
	return nil
	// use native msg box for windows 7
	if IsWindows7() {
		m := msgbox.New("aInstaller", msg, msgbox.IconInformation)
		log.Println(m)
		return nil
	}

	// windows 8 - 10 notifications
	n := toast.Notification{
		AppID:   "aInstaller",
		Title:   "aInstaller",
		Message: msg,
		Icon:    iconPath,
	}

	return n.Push()
}

func main() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	// icon for the notification
	iconPath = filepath.Join(wd, "./icon.png")

	systray.Run(onReady)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("aInstaller")
	systray.SetTooltip("ahud auto updater")

	mCheck := systray.AddMenuItem("Check for updates", "Check if there is a new ahud version available for download")
	mQuit := systray.AddMenuItem("Quit", "Close the update checker")

	timer := time.NewTicker(30 * time.Minute)
	u := new(Update)
	check := func() {
		log.Println("checking...")
		if err := u.Check(); err != nil {
			_ = u.Notify("Error")
			panic(err)
		}
	}

	go check()
	go func() {
		for {
			select {
			case <-timer.C:
				go check()
			case <-mCheck.ClickedCh:
				go check()
			case <-mQuit.ClickedCh:
				systray.Quit()
				os.Exit(0)
			}
		}
	}()
}

func getModuleFileName(pid int) (string, error) {
	var n uint32
	b := make([]uint16, syscall.MAX_PATH)
	size := uint32(len(b))

	hProcess, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return "", os.NewSyscallError("OpenProcess", err)
	}

	defer syscall.CloseHandle(hProcess)
	path, _, err := getModuleFileNameProc.Call(uintptr(hProcess), uintptr(unsafe.Pointer(&b[0])), uintptr(size))

	n = uint32(path)
	if n == 0 {
		return "", err
	}

	return string(utf16.Decode(b[0:n])), nil
}

// IsTF2Running checks if tf2's process is found in the native running process list
func IsTF2Running() (bool, error) {
	list, err := ps.Processes()
	if err != nil {
		return false, err
	}

	// loop through process list
	for _, proc := range list {
		if strings.Compare(strings.ToLower(proc.Executable()), "hl2.exe") == 0 {
			// get full file path from process id
			path, err := getModuleFileName(proc.Pid())
			if err != nil {
				return false, err
			}

			log.Println(path)
			return true, nil
		}
	}

	return false, nil
}

// IsWindows7 checks if the current OS is windows 7
func IsWindows7() bool {
	gi := goInfo.GetInfo()
	return strings.HasPrefix(gi.Core, "7")
}
