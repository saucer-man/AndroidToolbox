import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'

// 引入router
import router from './route/route'
import { MouseMenuDirective } from '@howdyjs/mouse-menu';

const app = createApp(App)
// app.config.globalProperties.$devices = []
// app.config.globalProperties.$selectdevice = "12345"

app.use(ElementPlus)
app.use(router)

app.directive('MouseMenu', MouseMenuDirective);
app.mount('#app')

