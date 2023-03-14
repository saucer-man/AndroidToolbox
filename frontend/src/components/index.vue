<template>
  <div class="common-layout">
    <el-container>
      <!-- <el-header>Header</el-header> -->
      <el-container>
        <el-aside width="80%">
          <h3 align="center">设备信息</h3>
          <div id="wrap">
            <el-table :data="tableData" max-height="430" width="600" empty-text="没有信息">
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
              <el-col :span="10"
                ><el-button type="primary" @click="return1">返回</el-button></el-col
              >
              <el-col :span="4"></el-col>
              <el-col :span="10"
                ><el-button type="primary" @click="returnHome">主页</el-button></el-col
              >
            </el-row>
            <el-row>
              <el-col :span="10"
                ><el-button type="primary" @click="soundAdd">音量+</el-button></el-col
              >
              <el-col :span="4"></el-col>
              <el-col :span="10"
                ><el-button type="primary" @click="soundDel">音量-</el-button></el-col
              >
            </el-row>
            <el-row>
              <el-col :span="10"><el-button type="primary" @click="lock">锁屏</el-button></el-col>
              <el-col :span="4"></el-col>
              <el-col :span="10"><el-button type="primary" @click="unlock">解锁</el-button></el-col>
            </el-row>
            <el-row>
              <el-col :span="10"
                ><el-button type="primary" @click="snapshot">截屏</el-button></el-col
              >
              <el-col :span="4"></el-col>
              <el-col :span="10"><el-button type="primary" @click="scrcpy">投屏</el-button></el-col>
            </el-row>
            <h3 align="center">重启</h3>
            <el-row>
              <el-col :span="10"><el-button type="primary" @click="reboot">重启</el-button></el-col>
              <el-col :span="4"></el-col>
              <el-col :span="10"
                ><el-button type="primary" @click="shutdown">关机</el-button></el-col
              >
            </el-row>
            <el-row>
              <el-col :span="10"
                ><el-button type="primary" @click="recovery">recovery</el-button></el-col
              >
              <el-col :span="4"></el-col>
              <el-col :span="10"
                ><el-button type="primary" @click="fastboot">fastboot</el-button></el-col
              >
            </el-row>
          </el-main>
        </el-container>
      </el-container>
    </el-container>
  </div>
</template>
<script setup>
import { ElMessage } from 'element-plus';
import { watch, onActivated, ref } from 'vue';
import { GetDeviceProp, Excute, GetImage, ExcuteSync } from '../../wailsjs/go/app/App.js';
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

const tableData = ref([
  {
    Name: '品牌',
    Value: 'vivo',
    Mean: 'ro.product.brand',
  },
]);

onActivated(() => {
  updateValue();
});
const updateValue = async () => {
  GetDeviceProp(props.selectdevice).then((result) => {
    tableData.value = result.DevicePropList;
    console.log('更新了tableData：', tableData.value);
  });
};
const return1 = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'input', 'keyevent', 'BACK']).then((result) =>
    handleCommandResult(result),
  );
};
const returnHome = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'input', 'keyevent', 'HOME']).then((result) =>
    handleCommandResult(result),
  );
};
const soundAdd = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'input', 'keyevent', '24']).then((result) =>
    handleCommandResult(result),
  );
};
const soundDel = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'input', 'keyevent', '25']).then((result) =>
    handleCommandResult(result),
  );
};
const lock = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shell', 'input', 'keyevent', '26']).then((result) =>
    handleCommandResult(result),
  );
};
const unlock = async () => {
  Excute([
    'adb',
    '-s',
    props.selectdevice,
    'shell',
    'input',
    'swipe',
    '200',
    '500',
    '200',
    '0',
  ]).then((result) => {
    handleCommandResult(result);
  });
};
const snapshot = async () => {
  GetImage().then((result) => {
    if (result == 0) {
      ElMessage.success('命令执行成功, 截图保存在screenshot目录');
    } else {
      ElMessage.success('命令执行失败');
    }
  });
};
const scrcpy = async () => {
  ExcuteSync(['scrcpy', '-s', props.selectdevice]).then((result) => handleCommandResult(result));
};
const reboot = async () => {
  Excute(['adb', '-s', props.selectdevice, 'reboot']).then((result) => {
    handleCommandResult(result);
  });
};
const shutdown = async () => {
  Excute(['adb', '-s', props.selectdevice, 'shutdown']).then((result) => {
    handleCommandResult(result);
  });
};
const recovery = async () => {
  Excute(['adb', '-s', props.selectdevice, 'reboot', 'recovery']).then((result) => {
    handleCommandResult(result);
  });
};
const fastboot = async () => {
  Excute(['adb', '-s', props.selectdevice, 'reboot', 'fastboot']).then((result) => {
    handleCommandResult(result);
  });
};
const handleCommandResult = async (execResult) => {
  console.log('handleCommandResult返回值:', execResult);
  if (execResult.ExitCode == 0) {
    ElMessage.success('命令执行成功');
  } else {
    ElMessage.success('命令执行出错: ' + execResult.Stderr);
  }
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
  margin-left: 10px;
  width: 180px;
}

</style>
