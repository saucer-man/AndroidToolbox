<template>
  <div>
    <el-container>
      <el-main>
        <el-form :inline="true">
          <h3 align="center">软件管理</h3>
          <el-form-item label="当前界面包名:">{{ currentPackage }} </el-form-item>
          <el-form-item label="当前界面活动:">{{ currentActivity }} </el-form-item>
          <el-form-item label="选择app">
            <el-select
              @change="getPackageInfo"
              v-model="selectPackage"
              filterable
              placeholder="{{currentActivity}}"
            >
              <el-option
                v-for="item in packages"
                :key="item.value"
                :label="item.value"
                :value="item.value"
              />
            </el-select>
          </el-form-item>
          <el-form-item label="app信息" prop="output">
            <el-input style="width: 600px" v-model="selectPackageInfo" :rows="12" type="textarea" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="downloadPackage">提取安装包apk</el-button>
            <el-button type="primary" @click="startPackage">启动app</el-button>
            <el-button type="primary" @click="stopPackage">停止app</el-button>
            <el-button type="primary" @click="installPackage">软件安装</el-button>

            <el-button type="primary" @click="clearPackage">清除数据</el-button>
            <el-button type="primary" @click="removePackage">卸载</el-button>
            <el-button type="primary" @click="removePackageWithData">保留数据卸载</el-button>
          </el-form-item>
        </el-form>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ElMessage, ElMessageBox } from 'element-plus';

import { watch, ref, onActivated } from 'vue';

import { GetCurrentActivity,Excute,InstallPackage } from '../../wailsjs/go/app/App.js';

const props = defineProps({
  selectdevice: {
    type: String,
    required: true,
  },
});
watch(
  () => props.selectdevice,
  (newProps) => {
    console.log('通过prop获取到的selectdevice更新成了:', props.selectdevice);
    updateValue();
  },
);

const currentPackage = ref('');
const currentActivity = ref('');
const selectPackage = ref('');
const selectPackageInfo = ref('');
const packages = ref([]);

onActivated(() => {
  updateValue();
});
const updateValue = async () => {
  updatePackages();
  GetCurrentActivity(props.selectdevice).then((result) => {
    currentPackage.value = result.PackageName;
    currentActivity.value = result.ActivityName;
    selectPackage.value = result.PackageName;
  });
  getPackageInfo();
};
const updatePackages = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'pm', 'list', 'packages']).then((result) => {
    console.log('adb shell pm list packages 返回值:', result);
    if (result.ExitCode == 0) {
      var packages_raw = result.Stdout.split('\n');
      for (var i = 0; i < packages_raw.length; i++) {
        if (packages_raw[i] != null && packages_raw[i].length > 0) {
          
          packages.value.push({
            value: packages_raw[i].replace('package:', '').replace('\r', ''),
          });
        }
      }
      
    } else {
      ElMessage.success('命令执行出错: ' + execResult.Stderr);
    }
  });
};
// 提取安装包apk
const downloadPackage = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'pm', 'path', selectPackage.value]).then(
    (result) => {
      console.log('adb shell pm path package 返回值:', result);
      if (result.ExitCode == 0) {
        var packagePath = result.Stdout.replace('package:', '').replace('\r\n', '');
        Excute([
          'adb',
          '-s',
          props.selectdevice,
          'pull',
          packagePath,
          `${selectPackage.value}.apk`,
        ]).then((result1) => {
          if (result1.ExitCode == 0) {
            ElMessage.success(`成功将${selectPackage.value}下载到当前目录`);
          } else {
            ElMessage.error('命令执行出错: ' + result1.Stderr);
          }
        });
      } else {
        ElMessage.error('命令执行出错: ' + result.Stderr);
      }
    },
  );
};
const clearPackage = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'pm', 'clear', selectPackage.value]).then(
    (result) => {
      console.log('adb shell pm clear 返回值:', result);
      if (result.ExitCode == 0) {
        ElMessage.success(`成功清除`);
      } else {
        ElMessage.success('命令执行出错: ' + result.Stderr);
      }
    },
  );
};
const removePackage = async () => {
  Excute(['adb', '-s', props.selectdevice, 'uninstall', selectPackage.value]).then((result) => {
    console.log('adb uninstall 返回值:', result);
    if (result.ExitCode == 0) {
      ElMessage.success(`成功卸载`);
    } else {
      ElMessage.success('命令执行出错: ' + result.Stderr);
    }
  });
};
const removePackageWithData = async () => {
  Excute(['adb', '-s', props.selectdevice, 'uninstall', '-k', selectPackage.value]).then(
    (result) => {
      console.log('adb uninstall 返回值:', result);
      if (result.ExitCode == 0) {
        ElMessage.success(`成功卸载`);
      } else {
        ElMessage.success('命令执行出错: ' + result.Stderr);
      }
    },
  );
};
const startPackage = async () => {
  Excute([
    'adb',
    'shell',
    'monkey',
    '-p',
    selectPackage.value,
    '-c',
    'android.intent.category.LAUNCHER',
    '1',
  ]).then((result) => {
    console.log('adb 启动app 返回值:', result);
    if (result.ExitCode == 0) {
      ElMessage.success(`成功启动`);
    } else {
      ElMessage.success('命令执行出错: ' + result.Stderr);
    }
  });
};
const stopPackage = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'am', 'force-stop', selectPackage.value]).then(
    (result) => {
      console.log('adb 启动app 返回值:', result);
      if (result.ExitCode == 0) {
        ElMessage.success(`成功启动`);
      } else {
        ElMessage.success('命令执行出错: ' + result.Stderr);
      }
    },
  );
};
const handleCommandResult = async (execResult) => {
  console.log('handleCommandResult返回值:', execResult);
  if (execResult.ExitCode == 0) {
    ElMessage.success('命令执行成功');
  } else {
    ElMessage.success('命令执行出错: ' + execResult.Stderr);
  }
};
const installPackage = async () => {
    InstallPackage(props.selectdevice).then((result) => {
    handleCommandResult(result);
  });
};
const getPackageInfo = async () => {
  Excute([
    'adb',
    '-s',
    props.selectdevice,
    'shell',
    'dumpsys',
    'package',
    selectPackage.value,
  ]).then((result) => {
    console.log('adb 启动app 返回值:', result);
    if (result.ExitCode == 0) {
      selectPackageInfo.value = result.Stdout;
    } else {
      ElMessage.success('命令执行出错: ' + result.Stderr);
    }
  });
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
