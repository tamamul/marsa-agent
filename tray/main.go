// Program kecil terpisah dari agent.exe (service). Ini yang kelihatan di
// system tray Windows 11, jalan di sesi login user. Tugasnya cuma:
// buka dashboard di browser, dan kirim perintah start/stop/restart ke
// Windows Service "MarsaPrintAgent" lewat sc.exe.
package main

import (
	_ "embed"
	"os/exec"
	"runtime"

	"github.com/getlantern/systray"
)

//go:embed icon.ico
var iconBytes []byte

const dashboardURL = "http://127.0.0.1:9191"

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconBytes)
	systray.SetTitle("MARSA Print Agent")
	systray.SetTooltip("MARSA Print Agent - klik kanan untuk menu")

	mOpen := systray.AddMenuItem("Buka Dashboard", "Lihat status dan setting")
	systray.AddSeparator()
	mRestart := systray.AddMenuItem("Restart Service", "Mulai ulang MarsaPrintAgent")
	mStop := systray.AddMenuItem("Hentikan Service", "Stop MarsaPrintAgent")
	mStart := systray.AddMenuItem("Jalankan Service", "Start MarsaPrintAgent")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Keluar", "Tutup ikon tray (service tetap jalan di background)")

	go func() {
		for {
			select {
			case <-mOpen.ClickedCh:
				openBrowser(dashboardURL)
			case <-mRestart.ClickedCh:
				runSC("stop")
				runSC("start")
			case <-mStop.ClickedCh:
				runSC("stop")
			case <-mStart.ClickedCh:
				runSC("start")
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// tray ditutup, tapi Windows Service MarsaPrintAgent tetap jalan terus
	// di background - ini cuma "remote control"-nya saja yang ditutup.
}

func runSC(action string) {
	cmd := exec.Command("sc.exe", action, "MarsaPrintAgent")
	_ = cmd.Run()
}

func openBrowser(url string) {
	if runtime.GOOS != "windows" {
		return
	}
	_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
}
