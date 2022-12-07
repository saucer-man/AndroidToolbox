

<template>
  <el-container>
    <el-aside width="130px">
      <el-menu default-active="index" class="el-menu-vertical-demo" router=true>
        <el-menu-item index="/index">主页</el-menu-item>
        <el-menu-item index="/helloworld">软件管理</el-menu-item>
        <el-menu-item index="/filemanager">文件管理</el-menu-item>
        <el-menu-item index="/adb">ADB终端</el-menu-item>
      </el-menu>
    </el-aside>
    <router-view v-slot="{ Component }">
      <transition name="router-fade" mode="out-in">
        <keep-alive>
          <component :is="Component" />
        </keep-alive>
      </transition>
    </router-view>
  </el-container>
</template>



<style scoped>
.el-menu-item {
  background-color: rgba(255, 0, 0, 0.1);
}

.el-menu-item.is-active {
  background-color: #247bb4 !important;
}
</style>
