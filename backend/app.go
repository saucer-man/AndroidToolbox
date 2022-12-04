package app

import (
	"bufio"
	"bytes"
	"context"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/elliotchance/orderedmap/v2"

	log "github.com/sirupsen/logrus"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// 获取手机的一些信息
func (a *App) GetDeviceList() []string {
	out, err := Exex("adb", "devices", "-l")
	log.Info("adb devices命令返回错误结果: %s", err)
	var devices []string
	log.Info("adb devices命令返回结果: %s", out)
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, "model") {
			id := strings.Split(line, ` `)
			devices = append(devices, id[0])
		}
	}
	log.Info("adb devices命令返回结果: %s", devices)
	return devices
}

type DevicePropList struct {
	DevicePropList []DeviceProp
}
type DeviceProp struct {
	Name  string
	Value string
	Mean  string
}

func (a *App) GetDeviceProp() DevicePropList {
	var res DevicePropList /*创建集合 */
	propAttitudeOrderedMap := orderedmap.NewOrderedMap[string, string]()

	propAttitudeOrderedMap.Set("设备品牌", "ro.product.brand")
	propAttitudeOrderedMap.Set("设备型号", "ro.product.model")
	propAttitudeOrderedMap.Set("上市名称", "ro.product.marketname")
	propAttitudeOrderedMap.Set("设备名", "ro.product.device")

	propAttitudeOrderedMap.Set("安卓版本", "ro.build.version.release")

	propAttitudeOrderedMap.Set("SDK版本", "ro.build.version.sdk")
	propAttitudeOrderedMap.Set("CPU架构", "ro.product.cpu.abi")
	propAttitudeOrderedMap.Set("版本id", "ro.build.id")
	propAttitudeOrderedMap.Set("fingerprint", "ro.build.fingerprint")

	// propAttitudeMap := map[string]string{"手机品牌": "ro.product.brand",
	// 	"设备名称":        "ro.product.device",
	// 	"设备内部代号":      "ro.product.model",
	// 	"版本id":        "ro.build.id",
	// 	"版本号":         "ro.build.display.id",
	// 	"SDK版本":       "ro.build.version.sdk",
	// 	"版本代号":        "ro.build.version.codename",
	// 	"编译时间":        "ro.build.date",
	// 	"编译类型":        "ro.build.type",
	// 	"编译者":         "ro.build.user",
	// 	"制造商":         "ro.product.manufacturer",
	// 	"fingerprint": "ro.build.fingerprint",
	// }

	out, err, _ := RunCommand("adb", "shell", "getprop")
	log.Info("adb getprop命令返回错误结果: %s", err)
	for _, propName := range propAttitudeOrderedMap.Keys() {
		propMean, _ := propAttitudeOrderedMap.Get(propName)
		log.Info("propName: %s", propName)
		sc := bufio.NewScanner(strings.NewReader(out))
		for sc.Scan() {
			line := sc.Text()
			tmp := strings.Split(line, `: `)
			value := tmp[len(tmp)-1]
			value = strings.Trim(value, " []")

			// fmt.Println(propName, "值是", propAttitudeMap[propName])
			if strings.Contains(line, propMean) {
				res.DevicePropList = append(res.DevicePropList, DeviceProp{propName, value, propMean})
				break
			}
		}
	}
	log.Info("adb getprop命令返回结果: %s", res)

	return res
}

// 同步执行命令
func (a *App) Excute(command string) (exitCode int) {
	args := strings.Split(command, " ")
	name := args[0]
	args = args[1:]
	log.Info("run command:", name, args)
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout := outbuf.String()
	stderr := errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Infof("Could not get exit code for failed program: %v, %v", name, args)
			exitCode = 0
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	log.Printf("command result, stdout: %v, stderr: %v, exitCode: %v", stdout, stderr, exitCode)
	return
}

// 异步执行命令
func (a *App) ExcuteSync(command string) {
	args := strings.Split(command, " ")
	name := args[0]
	args = args[1:]
	log.Info("run command:", name, args)
	cmd := exec.Command(name, args...)
	cmd.Start()
}

// GetImage 直接读取adb截图数据，速度更快
func (a *App) GetImage() (exitCode int) {
	out, _, exitCode := RunCommand("adb", "shell", "screencap", "-p")
	x := bytes.Replace([]byte(out), []byte("\r\n"), []byte("\n"), -1)
	path := "./screenshot/screenshot.png"
	if err := os.WriteFile(path, x, 0666); err != nil {
		log.Error(err)
	}

	return exitCode
}
