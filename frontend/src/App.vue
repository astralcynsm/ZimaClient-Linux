<template>
  <n-config-provider :theme="currentTheme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-notification-provider>
        <n-dialog-provider>
          <n-global-style />
          <SudoPasswordDialog ref="sudoDialog" @success="handleSudoSuccess" />
          <div :style="containerStyle">
            <router-view />
          </div>
        </n-dialog-provider>
      </n-notification-provider>
    </n-message-provider>
  </n-config-provider>
</template>

<script setup>
import { ref, computed, onMounted, watch, provide, inject } from 'vue'
import {
  NConfigProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
  NGlobalStyle,
  darkTheme
} from 'naive-ui'
import SudoPasswordDialog from './components/SudoPasswordDialog.vue'
import { CheckSudoInitialized } from '../wailsjs/go/main/App'

// 主题配置
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
  },
  Dialog: {
    borderRadius: '12px'
  },
  Modal: {
    borderRadius: '12px'
  }
}

const isDark = ref(false) // 默认浅色模式
const sudoDialog = ref(null)
const sudoInitialized = ref(false)

// 提供全局sudo检查函数
provide('checkSudo', async () => {
  if (!sudoInitialized.value) {
    const initialized = await CheckSudoInitialized()
    if (!initialized) {
      return new Promise((resolve, reject) => {
        sudoDialog.value.show()
        const handleSuccess = () => {
          sudoInitialized.value = true
          resolve()
        }
        const handleCancel = () => {
          reject(new Error('用户取消了sudo授权'))
        }
        sudoDialog.value.$once('success', handleSuccess)
        sudoDialog.value.$once('cancel', handleCancel)
      })
    }
  }
  return Promise.resolve()
})

const handleSudoSuccess = () => {
  sudoInitialized.value = true
}

// 从 localStorage 读取主题偏好
onMounted(async () => {
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    isDark.value = savedTheme === 'dark'
  }

  // 检查sudo是否已初始化
  try {
    sudoInitialized.value = await CheckSudoInitialized()
  } catch (error) {
    console.error('检查sudo状态失败:', error)
  }

  // 监听storage事件以实时更新主题
  window.addEventListener('storage', () => {
    const theme = localStorage.getItem('theme')
    if (theme) {
      isDark.value = theme === 'dark'
    }
  })
})

// 监听主题变化（从子组件触发）
watch(() => localStorage.getItem('theme'), (newTheme) => {
  if (newTheme) {
    isDark.value = newTheme === 'dark'
  }
})

const currentTheme = computed(() => isDark.value ? darkTheme : null)

const containerStyle = computed(() => ({
  backgroundColor: 'var(--bg-primary)',
  minHeight: '100vh',
  color: 'var(--text-primary)'
}))
</script>
