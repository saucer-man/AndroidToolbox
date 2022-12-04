import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import App from './App.vue'

// 引入router
import router from './route/route'


const app = createApp(App)
app.config.globalProperties.$deviceid = ''
app.use(ElementPlus)
app.use(router)
app.mount('#app')

