<template>
  <el-container>
    <el-main>
      <div align="center">
        current dir:
        <el-input
          v-model="currentDir"
          class="w-50 m-2"
          size="small"
          placeholder="Please Input path"
          @keyup.enter.native="updatePath(currentDir)" />
      </div>

      <el-table
        ref="singleTableRef"
        :data="tableData"
        highlight-current-row
        style="width: 100%"
        height="500"
        @current-change="handleCurrentChange"
        :row-style="{ height: '10px' }"
        :cell-style="{ padding: '0px' }"
        @cell-dblclick="bccelldblclick"
        @row-contextmenu="rowContextmenu">
        <!--  单击，这里就去掉好了 -->
        <el-table-column v-mouse-menu="options" type="index" width="50" />
        <el-table-column
          v-mouse-menu="options"
          property="Permition"
          label="权限"
          width="100" />
        <el-table-column property="Uid" label="Uid" width="100" />
        <el-table-column property="Gid" label="Gid" width="100" />
        <el-table-column property="Size" label="Size" width="100" />
        <el-table-column property="ModifyDate" label="ModifyDate" width="150" />
        <el-table-column property="Filename" label="Filename" />
      </el-table>
    </el-main>

    <!-- <el-footer>
      <div style="margin-top: 20px">
        <el-row>
          <el-col :span="4"
            ><el-button type="primary" @click="upload">上传</el-button></el-col
          >
          <el-col :span="4">
            <el-button type="primary" @click="download">下载</el-button>
          </el-col>
          <el-col :span="4"
            ><el-button type="primary" @click="copy">复制</el-button></el-col
          >
          <el-col :span="4"
            ><el-button type="primary" @click="paste">粘贴</el-button></el-col
          >
          <el-col :span="4"
            ><el-button type="primary" @click="move">移动</el-button></el-col
          >
          <el-col :span="4">
            <el-button type="primary" @click="delPath">删除</el-button>
          </el-col>
        </el-row>
      </div>
    </el-footer> -->
  </el-container>
</template>
<script setup>
import { ElMessage, ElMessageBox } from "element-plus";
import { CustomMouseMenu } from "@howdyjs/mouse-menu";

import {
  ListPath,
  Excute,
  UploadFile,
  DownloadFile,
} from "../../wailsjs/go/app/App.js";

import useGetGlobalProperties from "../hooks/global";
import buildPath from "../hooks/path.js";
import { onActivated, ref, watch } from "vue";
const props = defineProps({
  selectdevice: {
    type: String,
    required: true,
  },
});

const currentDir = ref("");
const selectPath = ref("");

const toCopyFilePath = ref("");
const toMoveFilePath = ref("");
const options = {
  menuWrapperCss: {
    background: "#3b393d",
  },
  menuItemCss: {
    labelColor: "#fff",
  },
  menuList: [
    {
      label: "从文件中导入",
      // tips: "Open",
      fn: (params, currentEl, bindingEl, e) => {
        console.log("从文件中导入", params, currentEl, bindingEl, e);
      },
    },
  ],
};

const menuList = [
  {
    label: "下载",

    fn: (params, currentEl, bindingEl, e) => {
      // params 就是传递进来的数，其实也就是row，但是我们这里用不到，因为在右键之后就一直处理了path
      // currentEl 点击右键时 所在的 HtmlElement 元素，这里就是event.srcElement
      // bindingEl 参数 3 为指令绑定的当前元素，这里为null
      // e就是 MouseEvent
      console.log("下载", params, currentEl, bindingEl, e);

      // 下面开始下载

      DownloadFile(props.selectdevice, selectPath.value).then((result) => {
        handleCommandResult(result);
        updatePath(currentDir.value);
      });
    },
  },
  {
    label: "上传",
    fn: (params, currentEl, bindingEl, e) => {
      UploadFile(props.selectdevice, currentDir.value).then((result) => {
        handleCommandResult(result);
        updatePath(currentDir.value);
      });
    },
  },
  {
    label: "复制",
    fn: (params, currentEl, bindingEl, e) => {
      toCopyFilePath.value = selectPath.value;
      ElMessage({
        message: "选择文件成功，请进入需要拷贝到的目录，然后点击拷贝",
        type: "success",
      });
    },
  },
  {
    label: "粘贴",
    fn: (params, currentEl, bindingEl, e) => {
      if (toCopyFilePath.value == "") {
        ElMessage({
          message: "请先选择需要拷贝的文件",
          type: "error",
        });
      } else {
        Excute([
          "adb",
          "-s",
          props.selectdevice,
          "shell",
          "cp",
          "-r",
          "\""+toCopyFilePath.value+"\"",
          "\""+currentDir.value+"\"",
        ]).then((result) => {
          handleCommandResult(result);
          updatePath(currentDir.value);
          toCopyFilePath.value = "";
        });
      }
    },
  },

  {
    label: "移动",
    fn: (params, currentEl, bindingEl, e) => {
      if (toMoveFilePath.value == "") {
        toMoveFilePath.value = selectPath.value;
        ElMessage({
          message: "选择文件成功，请进入需要移动到的目录，然后再点击移动",
          type: "success",
        });
      } else {
        Excute([
          "adb",
          "-s",
          props.selectdevice,
          "shell",
          "mv",
          "\""+toMoveFilePath.value+"\"",
          "\""+currentDir.value+"\"",
        ]).then((result) => {
          handleCommandResult(result);
          updatePath(currentDir.value);
          toMoveFilePath.value = "";
        });
      }
    },
  },

  {
    label: "删除",
    fn: (params, currentEl, bindingEl, e) => {
      ElMessageBox.confirm("是否确定删除" + selectPath.value, "Warning", {
    confirmButtonText: "OK",
    cancelButtonText: "Cancel",
    type: "warning",
  })
    .then(() => {
      Excute([
        "adb",
        "-s",
        props.selectdevice,
        "shell",
        "rm",
        "-rf",
        "\""+selectPath.value+"\"",
      ]).then((result) => {
        handleCommandResult(result);
        updatePath(currentDir.value);
      });
    })
    .catch(() => {});
    },
  },
];
const tableData = ref([
  {
    Permition: "",
    SubCount: "",
    Uid: "",
    Gid: "",
    Size: "",
    ModifyDate: "",
    Filename: "",
  },
]);

watch(
  () => props.selectdevice,
  (newProps) => {
    console.log("通过prop获取到的selectdevice更新成了:", props.selectdevice);
    updatePath("/storage/emulated/0/");
  }
);

onActivated(() => {
  console.log("通过prop获取到的selectdevice", props.selectdevice);
  updatePath("/storage/emulated/0/");
});

// 双击一个目录文件
const bccelldblclick = async (row, column, cell, event) => {
  if (!row.Permition.startsWith("d")) {
    ElMessage.info("你双击了一个非目录文件");
  } else {
    var path = buildPath(currentDir.value, row.Filename);
    updatePath(path);
  }
  console.log(currentDir.value);
};

// path是一个目录，双击该目录后，需要更新页面
const updatePath = async (path) => {
  ListPath(props.selectdevice, path).then((result) => {
    if (result.StdErr != "") {
      ElMessage.error("执行出错:" + result.StdErr);
    } else {
      currentDir.value = result.FilesList.Dir;
      tableData.value = result.FilesList.FilesList;
    }
  });
};

const handleCurrentChange = async (currentRow, oldCurrentRow) => {
  console.log("handleCurrentChange");
  console.log(currentRow);
  if (currentRow != null) {
    if (currentRow.hasOwnProperty("Filename")) {
      selectPath.value = buildPath(currentDir.value, currentRow.Filename);
    }
  }
};
const delPath = async () => {
  ElMessageBox.confirm("是否确定删除" + selectPath.value, "Warning", {
    confirmButtonText: "OK",
    cancelButtonText: "Cancel",
    type: "warning",
  })
    .then(() => {
      Excute([
        "adb",
        "-s",
        props.selectdevice,
        "shell",
        "rm",
        "-rf",
        selectPath.value,
      ]).then((result) => {
        handleCommandResult(result);
        updatePath(currentDir.value);
      });
    })
    .catch(() => {});
};
const handleCommandResult = async (execResult) => {
  console.log("handleCommandResult返回值:", execResult);
  if (execResult.ExitCode == 0) {
    ElMessage.success("命令执行成功");
  } else {
    ElMessage.success("命令执行出错: " + execResult.Stderr);
  }
};
const upload = async () => {
  UploadFile(props.selectdevice, currentDir.value).then((result) => {
    handleCommandResult(result);
    updatePath(currentDir.value);
  });
};
const download = async () => {
  DownloadFile(props.selectdevice, selectPath.value).then((result) => {
    handleCommandResult(result);
    updatePath(currentDir.value);
  });
};
const copy = async () => {
  toCopyFilePath.value = selectPath;
  ElMessage({
    message: "选择文件成功，请进入需要拷贝到的目录，然后点击拷贝",
    type: "success",
  });
};
const paste = async () => {
  if (toCopyFilePath.value == "") {
    ElMessage({
      message: "请先选择需要拷贝的文件",
      type: "error",
    });
  } else {
    Excute([
      "adb",
      "-s",
      props.selectdevice,
      "shell",
      "cp",
      "-r",
      toCopyFilePath.value,
      currentDir.value,
    ]).then((result) => {
      handleCommandResult(result);
      updatePath(currentDir.value);
      toCopyFilePath.value = "";
    });
  }
};
const updamovetePath = async () => {
  if (toMoveFilePath.value == "") {
    toMoveFilePath.value = selectPath.value;
    ElMessage({
      message: "选择文件成功，请进入需要移动到的目录，然后再点击移动",
      type: "success",
    });
  } else {
    Excute([
      "adb",
      "-s",
      props.selectdevice,
      "shell",
      "mv",
      toMoveFilePath.value,
      currentDir.value,
    ]).then((result) => {
      handleCommandResult(result);
      updatePath(currentDir.value);
      toMoveFilePath.value = "";
    });
  }
};

const rowContextmenu = async (row, column, event) => {
  console.log("点击了右键");
  console.log(row);
  console.log(column);
  console.log(event);
  selectPath.value = buildPath(currentDir.value, row.Filename);

  console.log("设置selectPath.value为:", selectPath.value);
  // todo: 将path拼接
  // const dom = ref()
  event.preventDefault();
  const MouseMenuCtx = CustomMouseMenu({
    // el: dom.value,
    el: event.srcElement, // 点击右键时所在的 HtmlElement 元素
    menuList,
    params: row, // 把row当参数传进去
  });

  MouseMenuCtx.show(event.clientX, event.clientY);
};
</script>

<style type="text/css" scoped>
.el-main {
  /* padding-top:0 !important; */
  padding: 0 !important;
}

.el-input {
  width: 300px;
}
</style>
