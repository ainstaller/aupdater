package main

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"

	"gopkg.in/toast.v1"

	"os"

	"github.com/aInstaller/icon"
	"github.com/aInstaller/utils/steam"
	"github.com/franela/goreq"
	"github.com/getlantern/systray"
)

const (
	mainMenuOverride = "https://raw.githubusercontent.com/n0kk/ahud/master/resource/ui/mainmenuoverride.res"
)

var (
	wd string
)

type Update struct {
	Version string
}

func (u *Update) Check() error {
	err := u.Notify("Checking for ahud updates...")
	if err != nil {
		return err
	}

	req := goreq.Request{
		Method: "GET",
		Uri:    mainMenuOverride,
	}

	res, err := req.Do()
	if err != nil {
		return err
	}

	str, err := res.Body.ToString()
	if err != nil {
		return err
	}

	var (
		month = 0
		year  = 0
		day   = 0

		re = regexp.MustCompile(`\"v(\d+)\.(\d+)\"`)
	)

	m := re.FindAllStringSubmatch(str, -1)
	if len(m) == 0 {
		return nil
	}

	for i, match := range m[0] {
		var err error
		if i == 1 {
			year, err = strconv.Atoi(match)
		} else if i == 2 {
			month, err = strconv.Atoi(match[0:2])
			if err != nil {
				return err
			}

			day, err = strconv.Atoi(match[2:4])
		}

		if err != nil {
			return err
		}

		fmt.Println(match, "found at index", i)
	}

	steam.FindGame()

	return nil
}

func (u *Update) Notify(msg string) error {
	n := toast.Notification{
		AppID:   "aInstaller",
		Title:   "aInstaller",
		Message: msg,
		//Message: "Checking for ahud updates...",
		Icon: filepath.Join(wd, "./icon.png"),
	}

	return n.Push()
}

func main() {
	var err error
	wd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	systray.Run(onReady)
}

func onReady() {
	systray.SetIcon(icon.Data)
	systray.SetTitle("aInstaller")
	systray.SetTooltip("ahud auto updater")

	mCheck := systray.AddMenuItem("Check for updates", "Check if there is a new ahud version available for download")
	mQuit := systray.AddMenuItem("Quit", "Close the update checker")

	go func() {
		select {
		case <-mCheck.ClickedCh:
			log.Println("checking...")
			u := new(Update)
			if err := u.Check(); err != nil {
				_ = u.Notify("Error")
				panic(err)
			}

		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}()
}
