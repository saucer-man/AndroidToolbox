package app

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"

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
	log.Debugf("adb devices命令返回结果: %s", out)
	sc := bufio.NewScanner(strings.NewReader(out))
	for sc.Scan() {
		line := sc.Text()
		if strings.Contains(line, "model") {
			id := strings.Split(line, ` `)
			devices = append(devices, id[0])
		}
	}
	log.Debugf("adb devices命令返回结果: %s", devices)
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

func (a *App) GetDeviceProp(deviceid string) DevicePropList {
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

	out, err, _ := RunCommand("adb", "-s", deviceid, "shell", "getprop")
	log.Debugf("adb getprop命令返回错误结果: %s", err)

	for _, propName := range propAttitudeOrderedMap.Keys() {

		propMean, _ := propAttitudeOrderedMap.Get(propName)
		log.Debugf("propName: %s", propName)
		if err != "" {
			res.DevicePropList = append(res.DevicePropList, DeviceProp{propName, "", propMean})
			continue
		}
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
	log.Infof("adb getprop命令返回结果: %s", res)

	return res
}

// 同步执行命令
func (a *App) Excute(commands []string) ExcuteResult {
	log.Infof("将要执行%+v", commands)
	excuteResult := ExcuteResult{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}

	command := commands[0]
	args := []string{}
	if len(command) > 1 {
		args = commands[1:]
	}
	// args := strings.Split(command, " ")
	// name := args[0]
	// args = args[1:]
	// 先检查命令是否存在
	if !CheckExists(command) {
		log.Errorf("didn't find '%s' executable\n", command)
		excuteResult.Stderr = fmt.Sprintf("%s 未安装", command)
		excuteResult.ExitCode = 1
		return excuteResult
	}
	stdout, stderr, exitCode := RunCommand(command, args...)
	log.Infof("执行结果为,stdout:%+v, stderr:%+v, exitCode:%+v", stdout, stderr, exitCode)
	excuteResult.Stdout = stdout
	excuteResult.Stderr = stderr
	excuteResult.ExitCode = exitCode
	return excuteResult
}

// 同步执行命令, 但是是使用的cmd /c的形式
func (a *App) ExcuteWithEnv(commands string) ExcuteResult {
	log.Infof("将要执行%+v", commands)
	excuteResult := ExcuteResult{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}

	stdout, stderr, exitCode := RunCommandWithEnv(commands)
	log.Debugf("执行结果为,stdout:%+v, stderr:%+v, exitCode:%+v", stdout, stderr, exitCode)
	excuteResult.Stdout = stdout
	excuteResult.Stderr = stderr
	excuteResult.ExitCode = exitCode
	return excuteResult
}

// 异步执行命令
func (a *App) ExcuteSync(commands []string) ExcuteResult {
	excuteResult := ExcuteResult{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}
	command := commands[0]
	args := []string{}
	if len(command) > 1 {
		args = commands[1:]
	}

	// 先检查命令是否存在
	if !CheckExists(command) {
		log.Errorf("didn't find '%s' executable\n", command)
		excuteResult.Stderr = fmt.Sprintf("%s 未安装", command)
		excuteResult.ExitCode = 1
		return excuteResult
	}
	log.Info("run command:", command, args)
	cmd := exec.Command(command, args...)
	PrepareBackgroundCommand(cmd)
	cmd.Start()
	return excuteResult
}

// GetImage 获取adb截图数据，速度更快
func (a *App) GetImage() (exitCode int) {
	out, _, exitCode := RunCommand("adb", "shell", "screencap", "-p")
	x := bytes.Replace([]byte(out), []byte("\r\n"), []byte("\n"), -1)
	path := "./screenshot/screenshot.png"
	if err := os.WriteFile(path, x, 0666); err != nil {
		log.Error(err)
	}

	return exitCode
}

type ListFileRes struct {
	FilesList FilesList
	StdErr    string
}

type FilesList struct {
	FilesList []FileInfo
	Dir       string
}

type FileInfo struct {
	// 结构可以参考https://blog.csdn.net/deniece1/article/details/102770824
	Permition  string
	SubCount   string
	Uid        string
	Gid        string
	Size       string
	ModifyDate string
	Filename   string
}

// 获取某个路径的文件信息
func (a *App) ListPath(deviceid, path string) ListFileRes {
	listFileRes := ListFileRes{
		FilesList: FilesList{
			FilesList: []FileInfo{},
			Dir:       path,
		},
		StdErr: "",
	}

	stdOut, stdErr, exitCode := RunCommand("adb", "-s", deviceid, "shell", "ls", "-la", path)
	if exitCode != 0 {
		log.Errorf("执行命令 adb shell ls -la %s 出错%s", path, stdErr)
		listFileRes.StdErr = stdErr
		return listFileRes
	}

	sc := bufio.NewScanner(strings.NewReader(stdOut))
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "total") {
			continue
		}
		log.Debugf("每一行的输出为：%s", line)

		r := regexp.MustCompile(`(.+?) {1,10}(\d+?) {1,10}(.+?) {1,10}(.+?) {1,10}(.+?) {1,10}(.+? .+?) (.+)`)

		matchs := r.FindStringSubmatch(line)
		var fileInfo FileInfo = FileInfo{
			Permition:  matchs[1],
			SubCount:   matchs[2],
			Uid:        matchs[3],
			Gid:        matchs[4],
			Size:       matchs[5],
			ModifyDate: matchs[6],
			Filename:   matchs[7],
		}
		listFileRes.FilesList.FilesList = append(listFileRes.FilesList.FilesList, fileInfo)
	}
	return listFileRes

}

type Actitity struct {
	PackageName  string
	ActivityName string
}

// 获取当前的activity
func (a *App) GetCurrentActivity(deviceid string) Actitity {

	excuteResult := Actitity{
		PackageName:  "",
		ActivityName: "",
	}
	stdOut, stdErr, exitCode := RunCommand("adb", "-s", deviceid, "shell", "dumpsys window | grep mCurrentFocus")
	if exitCode != 0 {
		log.Errorf("执行命令 adb shell dumpsys window | grep mCurrentFocus 出错: %s", stdErr)
		return excuteResult
	}

	r := regexp.MustCompile(`mCurrentFocus=Window{.*? .*? (.*?)(/.*?)?}`)

	matchs := r.FindStringSubmatch(stdOut)

	if len(matchs) == 2 {
		excuteResult.PackageName = matchs[1]
	}
	if len(matchs) == 3 {
		excuteResult.PackageName = matchs[1]
		excuteResult.ActivityName = strings.Replace(matchs[2], "/", "", 1)
	}
	log.Infof("获取当前包名的返回结果为:%+v", excuteResult)
	return excuteResult
}

func (a *App) UploadFile(deviceid string, dir string) ExcuteResult {
	excuteResult := ExcuteResult{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}
	toUploadFilePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		// DefaultDirectory: "C://",
		Title: "选择将要上传的文件",
	})
	log.Infof("选择了一个文件:%+v", toUploadFilePath)
	if err != nil {
		excuteResult.ExitCode = -1
		excuteResult.Stderr = fmt.Sprintf("failed opening dialog - %s", err.Error())
		return excuteResult
	}
	if toUploadFilePath == "" {
		excuteResult.Stderr = "您取消了该操作"
		excuteResult.ExitCode = -1
		return excuteResult
	}
	toSaveFileName := filepath.Base(toUploadFilePath)
	toSavePath := strings.ReplaceAll(filepath.Join(dir, toSaveFileName), "\\", "/") // 在手机上的路径无法使用\\，必须使用/

	return a.Excute([]string{"adb", "-s", deviceid, "push", toUploadFilePath, toSavePath})
}

func (a *App) DownloadFile(deviceid string, filePath string) ExcuteResult {
	excuteResult := ExcuteResult{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}
	toSaveDir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		// DefaultDirectory: "C://",
		Title: "选择将要保存的文件位置",
	})
	log.Infof("选择了一个目录:%+v", toSaveDir)
	if err != nil {
		excuteResult.Stderr = fmt.Sprintf("failed opening dialog - %s", err.Error())
		excuteResult.ExitCode = -1
		return excuteResult
	}
	if toSaveDir == "" {
		excuteResult.Stderr = "您取消了该操作"
		excuteResult.ExitCode = -1
		return excuteResult
	}
	toSaveFileName := filepath.Base(filePath)
	toSavePath := filepath.Join(toSaveDir, toSaveFileName)

	return a.Excute([]string{"adb", "-s", deviceid, "pull", filePath, toSavePath})

}

func (a *App) InstallPackage(deviceid string) ExcuteResult {
	excuteResult := ExcuteResult{
		Stdout:   "",
		Stderr:   "",
		ExitCode: 0,
	}
	toInstallFilePath, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		// DefaultDirectory: "C://",
		Title: "选择将要安装的apk",
	})
	log.Infof("选择了一个apk:%+v", toInstallFilePath)
	if err != nil {
		excuteResult.ExitCode = -1
		excuteResult.Stderr = fmt.Sprintf("failed opening dialog - %s", err.Error())
		return excuteResult
	}
	if toInstallFilePath == "" {
		excuteResult.Stderr = "您取消了该操作"
		excuteResult.ExitCode = -1
		return excuteResult
	}

	return a.Excute([]string{"adb", "-s", deviceid, "install", toInstallFilePath})
}
