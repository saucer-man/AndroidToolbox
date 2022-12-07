
<template>
  <el-container>

    <el-main>
      <div align="center">current dir: <el-input v-model="currentDir" class="w-50 m-2" size="small"
          placeholder="Please Input path" @keyup.enter.native="updatePath(currentDir)" /></div>
      <el-table ref="singleTableRef" :data="tableData" highlight-current-row style="width: 100%"
        @current-change="handleCurrentChange" height="450" :row-style="{ height: '10px' }"
        :cell-style="{ padding: '0px' }" @cell-dblclick="bccelldblclick">
        <el-table-column type="index" width="50" />
        <el-table-column property="Permition" label="权限" width="100" />
        <el-table-column property="Uid" label="Uid" width="50" />
        <el-table-column property="Gid" label="Gid" width="100" />
        <el-table-column property="Size" label="Size" width="100" />
        <el-table-column property="ModifyDate" label="ModifyDate" width="150" />
        <el-table-column property="Filename" label="Filename" />
      </el-table>
    </el-main>

    <el-footer>
      <div style="margin-top: 20px">
        <el-row>
          <el-col :span="4"><el-button type="primary" @click="upload">上传</el-button></el-col>
          <el-col :span="4">

            <el-button type="primary" @click="download">下载</el-button>
          </el-col>
          <el-col :span="4"><el-button type="primary" @click="copy">复制</el-button></el-col>
          <el-col :span="4"><el-button type="primary" @click="paste">粘贴</el-button></el-col>
          <el-col :span="4"><el-button type="primary" @click="move">移动</el-button></el-col>
          <el-col :span="4">
            <el-button type="primary" @click="delPath">删除</el-button>
          </el-col>

        </el-row>
      </div>
    </el-footer>
  </el-container>
</template>
  
<script>
import { ElMessage, ElMessageBox } from "element-plus";
import { handleError } from "vue";
import { ListPath, Excute } from "../../wailsjs/go/app/App.js";
import useGlobalProperties from '../hooks/globalVar.js'
import buildPath from '../hooks/path.js'

export default {
  data() {
    return {
      currentDir: "",
      selectPath: "",
      toCopyFilePath: "", // 将要拷贝的文件或者目录
      toMoveFilePath: "", // 将要拷贝的文件或者目录
      tableData: [
        {
          Permition: "",
          SubCount: "",
          Uid: "",
          Gid: "",
          Size: "",
          ModifyDate: "",
          Filename: "",
        },
      ]

    };
  },


  mounted() {
    const globalProperties = useGlobalProperties()
    console.log(globalProperties.$deviceid)
    this.updatePath("/storage/emulated/0/")

  },
  methods: {
    updatePath(path) {
      ListPath(path).then((result) => {
        if (result.StdErr != "") {
          ElMessage.error("执行出错:" + result.StdErr)
        } else {
          this.currentDir = result.FilesList.Dir
          this.tableData = result.FilesList.FilesList
        }
      })
    },
    //班次单元格双击
    bccelldblclick(row, column, cell, event) {

      if (!row.Permition.startsWith('d')) {
        ElMessage.info("你双击了一个非目录文件")
      } else {
        var path = buildPath(this.currentDir, row.Filename)
        this.updatePath(path)
      }
      console.log(this.currentDir)




    },
    handleCurrentChange(currentRow, oldCurrentRow) {
      console.log("handleCurrentChange")
      console.log(currentRow)
      if (currentRow != null) {
        if (currentRow.hasOwnProperty('Filename')) {
          this.selectPath = buildPath(this.currentDir, currentRow.Filename)
        }
      }


    },
    delPath() {
      var that = this
      ElMessageBox.confirm(
        '是否确定删除' + this.selectPath,
        'Warning',
        {
          confirmButtonText: 'OK',
          cancelButtonText: 'Cancel',
          type: 'warning',
        }
      )
        .then(() => {
          Excute("adb shell rm -rf " + this.selectPath).then((result) => {
            this.handleCommandResult(result)
            that.updatePath(this.currentDir)
          }



          )
        })
        .catch(() => {

        })


    },

    handleCommandResult: function (execResult) {
      console.log("handleCommandResult返回值:", execResult)
      if (execResult.ExitCode == 0) {
        ElMessage.success("命令执行成功")
      } else {
        ElMessage.success("命令执行出错: " + execResult.Stderr)
      }
    },


    upload() {
      var that = this
      ElMessageBox.prompt('请输入要上传的文件路径', '上传文件', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      })
        .then(({ value }) => {
          Excute(`adb push ${value} ${that.currentDir}`).then((result) => {
            that.handleCommandResult(result)
            that.updatePath(that.currentDir)
          })
        })
        .catch(() => {
        })
    },
    download() {
      var that = this
      ElMessageBox.prompt('请输入要下载的文件目录路径， 留空则为当前路径', '下载文件', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      })
        .then(({ value }) => {
          if (value == null) {
            value = "."
          }
          Excute(`adb pull ${that.selectPath} ${value}`).then((result) => {
            that.handleCommandResult(result)
            that.updatePath(that.currentDir)
          })
        })
        .catch(() => {
        })
    },

    copy() {
      this.toCopyFilePath = this.selectPath;
      ElMessage({
        message: '选择文件成功，请进入需要拷贝到的目录，然后点击拷贝',
        type: 'success',
      })
    },
    paste() {
      var that = this
      if (that.toCopyFilePath == "") {
        ElMessage({
          message: '请先选择需要拷贝的文件',
          type: 'error',
        })
      } else {
        Excute(`adb shell cp -r ${that.toCopyFilePath} ${that.currentDir}`).then((result) => {
          that.handleCommandResult(result)
          that.updatePath(that.currentDir)
          that.toCopyFilePath = ""
        })
      }
    },
    move() {
      var that = this
      if (that.toMoveFilePath == "") {
        that.toMoveFilePath = that.selectPath;
        ElMessage({
          message: '选择文件成功，请进入需要移动到的目录，然后再点击移动',
          type: 'success',
        })
      } else {
        Excute(`adb shell mv ${that.toMoveFilePath} ${that.currentDir}`).then((result) => {
          that.handleCommandResult(result)
          that.updatePath(that.currentDir)
          that.toMoveFilePath = ""
        })
      }
    }
  },
  beforeUnmount() {

  },
};
</script>

<style type="text/css">
.el-main {
  /* padding-top:0 !important; */
  padding: 0 !important;
}

.el-input {
  width: 300px;
}
</style>