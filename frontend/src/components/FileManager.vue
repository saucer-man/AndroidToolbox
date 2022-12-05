
<template>
  <el-container>
    
    <el-main>
      <div align="center">current dir: <el-input
      v-model="currentDir"
      class="w-50 m-2"
      size="small"
      placeholder="Please Input path"
      @keyup.enter.native="updatePath(currentDir)"
    /></div>
      <el-table ref="singleTableRef" :data="tableData" highlight-current-row style="width: 100%"
        @current-change="handleCurrentChange" height="450" :row-style="{height:'10px'}" :cell-style="{padding:'0px'}"
        @cell-dblclick="bccelldblclick">
        <el-table-column  type="index" width="50" />
        <el-table-column property="Permition" label="权限"  width="100"/>
        <el-table-column property="Uid" label="Uid" width="50" />
        <el-table-column property="Gid" label="Gid" width="100"/>
        <el-table-column property="Size" label="Size" width="100" />
        <el-table-column property="ModifyDate" label="ModifyDate" width="150"/>
        <el-table-column property="Filename" label="Filename" />
      </el-table>
    </el-main>

    <el-footer>
      <div style="margin-top: 20px">
        <el-row>
              <el-col :span="4"><el-button type="primary" @click="return1">上传</el-button></el-col>
              <el-col :span="4"><el-button type="primary" @click="returnHome">下载</el-button></el-col>
              <el-col :span="4"><el-button type="primary" @click="returnHome">复制</el-button></el-col>
              <el-col :span="4"><el-button type="primary" @click="returnHome">粘贴</el-button></el-col>
              <el-col :span="4"><el-button type="primary" @click="returnHome">移动</el-button></el-col>
              <el-col :span="4"><el-button type="primary" @click="returnHome">删除</el-button></el-col>

            </el-row>
      </div>
    </el-footer>
  </el-container>
</template>
  
<script>
import { ElMessage } from "element-plus";
import { handleError } from "vue";
import { ListPath } from "../../wailsjs/go/app/App.js";
import useGlobalProperties from '../hooks/globalVar.js'
import buildPath from '../hooks/path.js'

export default {
  data() {
    return {
      currentDir: "",
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
    updatePath(path){
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

      if (!row.Permition.startsWith('d')){
        ElMessage.info("你双击了一个非目录文件")
      }else{
        var path = buildPath(this.currentDir,row.Filename)
        this.updatePath(path)
      }
      console.log(this.currentDir)


      
 
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