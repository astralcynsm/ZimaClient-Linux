<template>
  <n-layout style="height: 100vh;">
    <n-layout-header :style="headerStyle">
      <div style="display: flex; justify-content: space-between; align-items: center;">
        <h2 :style="titleStyle">ZimaClient - {{ deviceName }}</h2>
        <n-space align="center" :size="16">
          <ConnectionMonitor />
          <n-switch v-model:value="isDark" @update:value="toggleTheme">
            <template #checked>🌙</template>
            <template #unchecked>☀️</template>
          </n-switch>
        </n-space>
      </div>
    </n-layout-header>

    <n-layout has-sider style="height: calc(100vh - 64px);">
      <n-layout-sider
        bordered
        :width="240"
        :style="siderStyle"
      >
        <n-menu
          :value="activeKey"
          :options="menuOptions"
          :style="menuStyle"
          @update:value="handleMenuSelect"
        />
      </n-layout-sider>

      <n-layout-content style="padding: 24px;" :style="contentStyle">
        <router-view />
      </n-layout-content>
    </n-layout>
  </n-layout>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useMessage, useDialog } from 'naive-ui'
import { h } from 'vue'
import { NIcon } from 'naive-ui'
import {
  Home,
  Server,
  FolderOpen,
  Globe,
  Settings,
  LogOut
} from '@vicons/ionicons5'
import ConnectionMonitor from '../../components/ConnectionMonitor.vue'

const router = useRouter()
const route = useRoute()
const message = useMessage()
const dialog = useDialog()

const isDark = ref(true)
const deviceName = ref('')

onMounted(() => {
  // 读取主题偏好
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme) {
    isDark.value = savedTheme === 'dark'
  }

  // 读取设备名称
  deviceName.value = localStorage.getItem('rnctl_device_hostname') || '远程网络控制'
})

const toggleTheme = (value) => {
  isDark.value = value
  localStorage.setItem('theme', value ? 'dark' : 'light')
  // 触发全局主题更新
  window.dispatchEvent(new Event('storage'))
}

// 动态样式
const headerStyle = computed(() => ({
  height: '64px',
  padding: '12px 24px',
  backgroundColor: isDark.value ? '#18181c' : '#ffffff',
  borderBottom: `1px solid ${isDark.value ? '#333' : '#e0e0e0'}`
}))

const titleStyle = computed(() => ({
  color: isDark.value ? 'white' : '#333',
  margin: 0
}))

const siderStyle = computed(() => ({
  backgroundColor: isDark.value ? '#18181c' : '#ffffff'
}))

const menuStyle = computed(() => ({
  backgroundColor: isDark.value ? '#18181c' : '#ffffff',
  color: isDark.value ? 'white' : '#333'
}))

const contentStyle = computed(() => ({
  backgroundColor: isDark.value ? '#18181c' : '#f5f5f5'
}))

// 根据当前路由设置活动菜单项
const activeKey = computed(() => {
  const path = route.path
  if (path.includes('dashboard')) return 'dashboard'
  if (path.includes('storage')) return 'storage'
  if (path.includes('zerotier')) return 'zerotier'
  if (path.includes('settings')) return 'settings'
  return 'dashboard'
})

function renderIcon(icon) {
  return () => h(NIcon, null, { default: () => h(icon) })
}

const menuOptions = computed(() => [
  {
    label: '仪表板',
    key: 'dashboard',
    icon: renderIcon(Home)
  },
  {
    label: '存储管理',
    key: 'storage',
    icon: renderIcon(FolderOpen)
  },
  {
    label: 'ZeroTier',
    key: 'zerotier',
    icon: renderIcon(Globe)
  },
  {
    label: '设置',
    key: 'settings',
    icon: renderIcon(Settings)
  },
  {
    type: 'divider'
  },
  {
    label: '连接其他设备',
    key: 'connect-other',
    icon: renderIcon(Server)
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: renderIcon(LogOut)
  }
])

function handleMenuSelect(key) {
  if (key === 'logout') {
    handleLogout()
  } else if (key === 'connect-other') {
    handleConnectOther()
  } else {
    router.push(`/main/${key}`)
  }
}

function handleLogout() {
  dialog.warning({
    title: '确认退出',
    content: '确定要退出登录吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      // 检查是否需要自动卸载
      const autoUnmount = localStorage.getItem('rnctl_auto_unmount_on_exit') !== 'false'

      if (autoUnmount) {
        try {
          const { UnmountAll } = await import('../../../wailsjs/go/main/App')
          await UnmountAll()
          console.log('已卸载所有挂载点')
        } catch (error) {
          console.error('卸载失败:', error)
          // 卸载失败不影响退出流程
        }
      }

      localStorage.removeItem('rnctl_logged_in')
      message.success('已退出登录')
      router.push('/wizard/connect')
    }
  })
}

function handleConnectOther() {
  dialog.info({
    title: '连接其他设备',
    content: '这将返回到设备选择页面，当前连接将断开。是否继续？',
    positiveText: '继续',
    negativeText: '取消',
    onPositiveClick: () => {
      localStorage.removeItem('rnctl_logged_in')
      router.push('/wizard/connect')
    }
  })
}
</script>
