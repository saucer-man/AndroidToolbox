import { createRouter, createWebHashHistory, createWebHistory } from "vue-router"

// 1. 定义路由组件， 注意，这里一定要使用 文件的全名（包含文件后缀名）
import index from "../components/index.vue";
import helloworld from "../components/HelloWorld.vue";
import filemanager from "../components/FileManager.vue";
import adb from "../components/Adb.vue"
// 2. 定义路由配置
const routes = [
  {
    path: "/",
    redirect: '/index'
  },
  { // 主页
    path: "/index", 
    name: 'index',
    component: index,
    meta: {
      keepAlive: false   // 不需要缓存
    }
  },
  { 
    path: "/helloworld", 
    name: 'helloworld',
    component: helloworld
  },
  { 
    path: "/filemanager", 
    name: 'filemanager',
    component: filemanager,
    meta: {
      keepAlive: true   // 不需要缓存
    }
  },
  { 
    path: "/adb", 
    name: 'adb',
    component: adb,
    meta: {
      keepAlive: true   // 不需要缓存
    }
  },
];

// 3. 创建路由实例
const router = createRouter({
  // 4. 采用hash 模式
  history: createWebHashHistory(),
  // 采用 history 模式
  // history: createWebHistory(),
  routes, // short for `routes: routes`
});

export default router
