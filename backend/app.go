package app

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	out, _, _ := RunCommand("adb", "devices", "-l")

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
type ExcuteResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
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
func (a *App) Excute(command string) ExcuteResult {
	excuteResult := ExcuteResult{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}
	args := strings.Split(command, " ")
	name := args[0]
	args = args[1:]
	// 先检查命令是否存在
	if !CheckExists(name) {
		log.Errorf("didn't find '%s' executable\n", name)
		excuteResult.Stderr = fmt.Sprintf("%s 未安装", name)
		excuteResult.ExitCode = 1
		return excuteResult
	}
	stdout, stderr, exitCode := RunCommand(name, args...)

	excuteResult.Stdout = stdout
	excuteResult.Stderr = stderr
	excuteResult.ExitCode = exitCode
	return excuteResult
}

// 异步执行命令
func (a *App) ExcuteSync(command string) ExcuteResult {
	excuteResult := ExcuteResult{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}
	args := strings.Split(command, " ")
	name := args[0]
	args = args[1:]
	log.Info("run command:", name, args)
	// 先检查命令是否存在
	if !CheckExists(name) {
		log.Errorf("didn't find '%s' executable\n", name)
		excuteResult.Stderr = fmt.Sprintf("%s 未安装", name)
		excuteResult.ExitCode = 1
		return excuteResult
	}
	cmd := exec.Command(name, args...)
	cmd.Start()
	return excuteResult
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