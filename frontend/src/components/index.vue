
<template>






  <div class="common-layout">
    <el-container>
      <!-- <el-header>Header</el-header> -->
      <el-container>
        <el-aside width="80%">
          <h3 align="center">设备信息</h3>
          <div id="wrap">
            <el-table :data="tableData" max-height="430">
              <el-table-column prop="Name" label="属性" width="200" />
              <el-table-column prop="Value" label="值" width="200" />
              <el-table-column prop="Mean" label="获取方式" width="200" />
            </el-table>
          </div>
        </el-aside>
        <el-container>
          <el-main>
            <h3 align="center">按键模拟</h3>

            <el-row>
              <el-col :span="10"><el-button type="primary" @click="return1">返回</el-button></el-col>
              <el-col :span="4"></el-col>
              <el-col :span="10"><el-button type="primary" @click="returnHome">主页</el-button></el-col>
            </el-row>
            <el-row>
              <el-col :span="10"><el-button type="primary" @click="soundAdd">音量+</el-button></el-col>
              <el-col :span="4"></el-col>
              <el-col :span="10"><el-button type="primary" @click="soundDel">音量-</el-button></el-col>
            </el-row>
            <el-row>
              <el-col :span="10"><el-button type="primary"  @click="lock">锁屏</el-button></el-col>
              <el-col :span="4"></el-col>
              <el-col :span="10"><el-button type="primary"  @click="unlock">解锁</el-button></el-col>
              
            </el-row>
            <el-row>
              <el-col :span="10"><el-button type="primary"  @click="snapshot">截屏</el-button></el-col>
              <el-col :span="4"></el-col>
              <el-col :span="10"><el-button type="primary"  @click="scrcpy">投屏</el-button></el-col>
            </el-row>
            <h3 align="center">重启</h3>
            <el-row>
              <el-col :span="10"><el-button type="primary" @click="reboot">重启</el-button></el-col>
              <el-col :span="4"></el-col>
              <el-col :span="10"><el-button type="primary" @click="shutdown">关机</el-button></el-col>

            </el-row>
            <el-row>
              <el-col :span="10"><el-button type="primary" @click="bootloader">bootloader</el-button></el-col>
              <el-col :span="4"></el-col>
              <el-col :span="10"><el-button type="primary" @click="fastboot" >fastboot</el-button></el-col>
            </el-row>
          </el-main>



        </el-container>
      </el-container>
    </el-container>
  </div>
</template>
<script>

import { ElMessage } from "element-plus";
import { handleError } from "vue";
import { GetDeviceList, GetDeviceProp, Excute, GetImage,ExcuteSync } from "../../wailsjs/go/app/App.js";
import useGlobalProperties from '../hooks/globalVar.js'

export default {
  data() {
    return {
      devices: [],

      tableData: [
        {
          Name: '品牌',
          Value: 'vivo',
          Mean: "ro.product.brand"
        }
      ]

    };
  },


  mounted() {
    const globalProperties = useGlobalProperties()
    console.log(globalProperties.$deviceid)
    GetDeviceList().then((result) => (
      this.devices = result,
      this.deviceid,
      globalProperties.$deviceid = this.devices[0],
      console.log(globalProperties.$deviceid)
    ));
    GetDeviceProp().then((result) => (console.log("GetDeviceProp返回值", result),
      this.tableData = result.DevicePropList,
      console.log(typeof result.DevicePropList),
      console.log("下面是tables"),
      console.log(this.tableData)
    ))
    // ExcuteTest("adb shell ls -F /storage/emulated/0/").then({

    // })
  },
  methods: {
    return1: function () {
      Excute("adb shell input keyevent BACK").then((result) => (
        this.handleCommandResult(result)
      )
      )


    },
    returnHome: function () {
      Excute("adb shell input keyevent HOME").then((result) => (
        this.handleCommandResult(result)
      )
      )
    },
    soundAdd: function () {
      Excute("adb shell input keyevent 24").then((result) => (
        this.handleCommandResult(result)
      )
      )
    },
    soundDel: function () {
      Excute("adb shell input keyevent 25").then((result) => (
        this.handleCommandResult(result)
      )
      )
    },
    lock: function () {
      Excute("adb shell input keyevent 26").then((result) => (
        this.handleCommandResult(result)
      )
      )
    },
    unlock: function () {
      Excute("adb shell input swipe 200 500 200 0").then((result) => {
        this.handleCommandResult(result)
      }
      )
    },
    snapshot: function () {
      GetImage().then((result) => {
        if (result==0){
        ElMessage.success("命令执行成功, 截图保存在screenshot目录")
      }else{
        ElMessage.success("命令执行失败")
      }
      }
      )
    },
    scrcpy: function () {
      ExcuteSync("scrcpy").then((result) => (
        this.handleCommandResult(result)
      )
      )
    },
    reboot: function () {
      Excute("adb reboot").then((result) => {
        this.handleCommandResult(result)
      }
      )
    },
    shutdown: function () {
      Excute("adb shutdown").then((result) => {
        this.handleCommandResult(result)
      }
      )
    },
    bootloader: function () {
      Excute("adb reboot bootloader").then((result) => {
        this.handleCommandResult(result)
      }
      )
    },
    fastboot: function () {
      Excute("adb reboot fastboot").then((result) => {
        this.handleCommandResult(result)
      }
      )
    },
    handleCommandResult: function (execResult) {
      console.log("handleCommandResult返回值:",execResult)
      if (execResult.ExitCode == 0) {
        ElMessage.success("命令执行成功")
      } else {
        ElMessage.success("命令执行出错: "+execResult.Stderr)
      }
    }
  },
  beforeUnmount() {

  },
};
</script>

<style type="text/css" scoped>
.el-button {
  width: 60px;
  margin-bottom: 20px;
}

.el-main {
  /* padding-top:0 !important; */
  padding: 0 !important;
}
</style>