<template>
  <el-container>
    <el-aside width="130px">
      <el-menu default-active="/index" class="el-menu-vertical-demo" router="true" active-text-color="#409eff">
        <el-menu-item index="/index">主页</el-menu-item>
        <el-menu-item index="/softmanager">软件管理</el-menu-item>
        <el-menu-item index="/filemanager">文件管理</el-menu-item>
        <el-menu-item index="/adb">ADB终端</el-menu-item>

        <el-dropdown @command="handleSelectdevice">

          <el-button type="primary">
            选择设备<br />{{ selectdevice.substring(0,9) }}
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-for="item in devices" :command="item">
                {{ item }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </el-menu>
    </el-aside>

    <router-view v-slot="{ Component }">
      <transition name="router-fade" mode="out-in">
        <keep-alive>
          <component :is="Component" v-model:selectdevice="selectdevice" />
        </keep-alive>
      </transition>
    </router-view>
  </el-container>
</template>

<script>
import { ArrowDown } from '@element-plus/icons-vue';
import { ElMessage } from 'element-plus';

import {
  GetDeviceList,
  GetDeviceProp,
  Excute,
  GetImage,
  ExcuteSync,
} from '../wailsjs/go/app/App.js';
import useGetGlobalProperties from './hooks/global';
import { reactive, ref, getCurrentInstance } from 'vue';
import { resolveUnref } from '@vueuse/shared';

export default {
  data() {
    return {
      devices: [],
      selectdevice: '',
    };
  },

  created() {

    setInterval(() => {
      // console.log("开始---");
      GetDeviceList().then((result) => {
        this.devices = result;
        // console.log('获取到的devices:', this.devices);
        if (result != null && result.length > 0 && this.selectdevice == "") {
          this.selectdevice = result[0];
        } else if (result == null) {
          this.devices = []
          this.selectdevice = ""
        }

      });
    }, 5000);

  },

  methods: {
    handleSelectdevice: function (device) {
      this.selectdevice = device;
    },
  },
  beforeUnmount() { },
};
</script>

<style scoped>
/* .el-menu-item {
      background-color: rgba(255, 0, 0, );
    }*/

.el-menu-item.is-active {
  background-color: rgba(64, 158, 255, 0.1);
}

.example-showcase .el-dropdown-link {
  cursor: pointer;
  color: var(--el-color-primary);
  display: flex;
  align-items: center;
}

.el-dropdown {
  margin-top: 270px;
  width: 100%;
}

.el-button-group {
  width: 100%;
}

.el-container {
  /*设置内部填充为0，几个布局元素之间没有间距*/
  padding: 0px;
  /*外部间距也是如此设置*/
  margin: 0px;
  /*统一设置高度为100%*/
  height: 100%;
}
</style>
