
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
          
          <el-col :span="4">
           
            <input type="file" id="myfile" @change.prevent.stop="upload">
  </el-col>
          <el-col :span="4"><el-button type="primary" @click="returnHome">下载</el-button></el-col>
          <el-col :span="4"><el-button type="primary" @click="returnHome">复制</el-button></el-col>
          <el-col :span="4"><el-button type="primary" @click="returnHome">粘贴</el-button></el-col>
          <el-col :span="4"><el-button type="primary" @click="returnHome">移动</el-button></el-col>
          <el-col :span="4">
            <el-button type="primary" @click="delPath">删除</el-button>
          </el-col>

        </el-row>
      </div>
    </el-footer>
  </el-container>
</template>
  
<script>
import { ElMessage, ElMessageBox} from "element-plus";
import { handleError } from "vue";
import { ListPath, Excute } from "../../wailsjs/go/app/App.js";
import useGlobalProperties from '../hooks/globalVar.js'
import buildPath from '../hooks/path.js'

export default {
  data() {
    return {
      currentDir: "",
      selectPath: "",
      fileName: '',
      toUploadFileList:[],  // UploadUserFile
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
          },
            that.updatePath(this.currentDir)


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

    handleChange(uploadFile, files){
      console.log("执行了了handleChange，参数为",uploadFile,files,"此时toUploadFileList列表为：",this.toUploadFileList);
      filepath 
    },
    
    upload(event) {
      console.log("event",event)
        let files = event.target.files[0];
        console.log("files",files)
        this.fileName = this.getObjectUrl(files);
        console.log("this.fileName",this.fileName)
        console.log(files.path)
    },
    getObjectUrl(file) {
      let url = null;
      if (window.createObjectURL != undefined) {
        // basic
        url = window.createObjectURL(file);
      } else if (window.webkitURL != undefined) {
        // webkit or chrome
        url = window.webkitURL.createObjectURL(file);
      } else if (window.URL != undefined) {
        // mozilla(firefox)
        url = window.URL.createObjectURL(file);
      }
      return url;
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