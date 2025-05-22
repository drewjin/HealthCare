import { createApp } from 'vue'
import './style.css'
import App from './App.vue'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import axios from 'axios'

const app = createApp(App)

// Register all icons globally
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// Intercept requests to attach JWT
axios.interceptors.request.use(config => {
  const token = localStorage.getItem('jwt')
  if (token && config.headers) {
    config.headers['Authorization'] = `${token}`
  }
  return config
}, error => Promise.reject(error))

app.use(ElementPlus)
  .use(router)
  .mount('#app')