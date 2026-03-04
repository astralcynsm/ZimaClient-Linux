import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './style.css'

// 导入 Naive UI
import naive from 'naive-ui'

// 创建并挂载 Vue 应用
const app = createApp(App)

// 配置 Naive UI 主题
const themeOverrides = {
  common: {
    primaryColor: '#0057FF',
    primaryColorHover: '#337CFF',
    primaryColorPressed: '#0046CC',
    primaryColorSuppl: '#669FFF',
    successColor: '#22c55e',
    warningColor: '#f59e0b',
    errorColor: '#ef4444',
    infoColor: '#0057FF',
    borderRadius: '8px',
    borderRadiusSmall: '6px',
    fontFamily: 'ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif'
  },
  Button: {
    borderRadiusMedium: '8px',
    borderRadiusLarge: '10px',
    heightMedium: '40px',
    heightLarge: '48px',
    fontSizeMedium: '14px',
    fontSizeLarge: '16px'
  },
  Input: {
    borderRadius: '8px',
    heightMedium: '40px',
    heightLarge: '48px'
  },
  Card: {
    borderRadius: '12px'
  }
}

app.use(router)
app.use(naive)

// 将主题配置注入到全局
app.provide('themeOverrides', themeOverrides)

app.mount('#app')

console.log('Vue app with Naive UI and Router mounted')
