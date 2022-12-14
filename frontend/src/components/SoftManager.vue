<template>
  <div>
    <el-container>
      <el-main>
        <el-form :inline="true">
          <h3 align="center">软件管理</h3>
          <el-form-item label="当前界面包名:">{{ currentPackage }}

          </el-form-item>
          <el-form-item label="当前界面活动:">{{ currentActivity }}

          </el-form-item>
          <el-form-item label="选择app">
            <el-select @change="getPackageInfo" v-model="selectPackage" filterable placeholder={{currentActivity}}>
              <el-option v-for="item in packages" :key="item.value" :label="item.value" :value="item.value" />
            </el-select>

          </el-form-item>
          <el-form-item label="app信息" prop="output">
            <el-input style="width:600px" v-model="selectPackageInfo" :rows="12" type="textarea" />
          </el-form-item>
          <el-form-item>

            <el-button type="primary" @click="downloadPackage">提取安装包apk</el-button>
            <el-button type="primary" @click="startPackage">启动app</el-button>
            <el-button type="primary" @click="stopPackage">停止app</el-button>
            <el-button type="primary" @click="installPackage">软件安装</el-button>

            <el-button type="primary" @click="clearPackage">清除数据</el-button>
            <el-button type="primary" @click="removePackage">卸载</el-button>
            <el-button type="primary" @click="removePackageithData">保留数据卸载</el-button>

          </el-form-item>

        </el-form>

      </el-main>

    </el-container>

  </div>
</template>

<script>
import { ElMessage, ElMessageBox } from "element-plus";
import { timePanelSharedProps } from "element-plus/es/components/time-picker/src/props/shared";
import { handleError } from "vue";

import { GetDeviceList, GetDeviceProp, Excute, GetImage, ExcuteSync, GetCurrentActivity } from "../../wailsjs/go/app/App.js";
import useGlobalProperties from '../hooks/globalVar.js'


export default {
  data() {
    return {
      currentPackage: "",
      currentActivity: "",
      selectPackage: "",
      selectPackageInfo: "",
      packages: [],
    };
  },


  mounted() {
    const globalProperties = useGlobalProperties()
    console.log(globalProperties.$deviceid)
    this.updatePackages()
    GetCurrentActivity().then((result) => {
      this.currentPackage = result.PackageName
      this.currentActivity = result.ActivityName
      this.selectPackage = result.PackageName
    })
    this.getPackageInfo()


  },
  methods: {
    // 更新所有的包
    updatePackages: function () {
      const that = this
      Excute(["adb","shell","pm","list","packages"]).then((result) => {
        console.log("adb shell pm list packages 返回值:", result)
        if (result.ExitCode == 0) {
          var packages = result.Stdout.split('\n');
          for (var i = 0; i < packages.length; i++) {
            if (packages[i] != null && packages[i].length > 0)	//去掉空的
            {
              // console.log(packages[i].replace("package:",""));
              that.packages.push({ "value": packages[i].replace("package:", "").replace("\r", "") })
            }
          }
        } else {
          ElMessage.success("命令执行出错: " + execResult.Stderr)
        }
      }
      )
    },
    // 提取安装包apk
    downloadPackaged: function () {
      const that = this
      Excute(["adb","shell","pm","path",that.selectPackage]).then((result) => {
        console.log("adb shell pm path package 返回值:", result)
        if (result.ExitCode == 0) {
          var packagePath = result.Stdout.replace("package:", "").replace("\r\n", "");
          Excute(["adb","pull", packagePath,`${that.selectPackage}.apk`]).then((result1) => {
            if (result1.ExitCode == 0) {
              ElMessage.success(`成功将${that.selectPackage}下载到当前目录`)
            } else {
              ElMessage.error("命令执行出错: " + result1.Stderr)
            }
          })

        } else {
          ElMessage.error("命令执行出错: " + result.Stderr)
        }
      }
      )
    },
    clearPackage: function () {
      const that = this
      Excute(["adb","shell", "pm", "clear" ,that.selectPackage]).then((result) => {
        console.log("adb shell pm clear 返回值:", result)
        if (result.ExitCode == 0) {
          ElMessage.success(`成功清除`)
        } else {
          ElMessage.success("命令执行出错: " + result.Stderr)
        }
      }
      )
    },
    removePackage: function () {
      const that = this
      Excute(["adb" ,"uninstall" ,that.selectPackage]).then((result) => {
        console.log("adb uninstall 返回值:", result)
        if (result.ExitCode == 0) {
          ElMessage.success(`成功卸载`)
        } else {
          ElMessage.success("命令执行出错: " + result.Stderr)
        }
      }
      )
    },
    removePackageWithData: function () {
      const that = this
      Excute(["adb" ,"uninstall" ,"-k", that.selectPackage]).then((result) => {
        console.log("adb uninstall 返回值:", result)
        if (result.ExitCode == 0) {
          ElMessage.success(`成功卸载`)
        } else {
          ElMessage.success("命令执行出错: " + result.Stderr)
        }
      }
      )
    },
    startPackage: function () {
      const that = this
      Excute(["adb" ,"shell" ,"monkey", "-p", that.selectPackage,"-c", "android.intent.category.LAUNCHER", "1"]).then((result) => {
        console.log("adb 启动app 返回值:", result)
        if (result.ExitCode == 0) {
          ElMessage.success(`成功启动`)
        } else {
          ElMessage.success("命令执行出错: " + result.Stderr)
        }
      }
      )
    },
    stopPackage: function () {
      const that = this
      Excute(["adb", "shell", "am" ,"force-stop", that.selectPackage]).then((result) => {
        console.log("adb 启动app 返回值:", result)
        if (result.ExitCode == 0) {
          ElMessage.success(`成功启动`)
        } else {
          ElMessage.success("命令执行出错: " + result.Stderr)
        }
      }
      )
    },
    handleCommandResult: function (execResult) {
      console.log("handleCommandResult返回值:", execResult)
      if (execResult.ExitCode == 0) {
        ElMessage.success("命令执行成功")
      } else {
        ElMessage.success("命令执行出错: " + execResult.Stderr)
      }
    },
    installPackage() {
      var that = this
      ElMessageBox.prompt('请输入要安装的apk路径', '上传文件', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      })
        .then(({ value }) => {
          Excute(["adb","install",value]).then((result) => {
            that.handleCommandResult(result)

          })
        })
        .catch(() => {
        })
    },
    getPackageInfo: function () {
      const that = this
      Excute(["adb", "shell", "dumpsys", "package" ,that.selectPackage]).then((result) => {
        console.log("adb 启动app 返回值:", result)
        if (result.ExitCode == 0) {
          that.selectPackageInfo = result.Stdout
        } else {
          ElMessage.success("命令执行出错: " + result.Stderr)
        }
      }
      )
    },
  },
  beforeUnmount() {

  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style type="text/css" scoped>
.el-main {
  margin-left: 10px;
  width: 790px;

}


.el-main {
  /* padding-top:0 !important; */
  padding: 0 !important;
}
</style>