<template>
  <n-space align="center" :size="8">
    <n-tag :type="connectionType" size="small" round>
      <template #icon>
        <n-icon :component="connectionIcon" />
      </template>
      {{ connectionLabel }}
    </n-tag>
    <n-popover v-if="lastSwitch" trigger="hover">
      <template #trigger>
        <n-icon :component="InformationCircle" size="16" style="cursor: pointer; opacity: 0.6" />
      </template>
      <div style="max-width: 300px">
        <p style="margin: 0 0 8px 0; font-weight: 600">连接切换记录</p>
        <p style="margin: 0; font-size: 13px">
          {{ lastSwitch.time }}<br>
          从 {{ lastSwitch.from }} 切换到 {{ lastSwitch.to }}<br>
          原因: {{ lastSwitch.reason }}
        </p>
      </div>
    </n-popover>
  </n-space>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Wifi, Globe, Close, InformationCircle } from '@vicons/ionicons5'
import { TestDeviceConnection } from '../../wailsjs/go/main/App'

const currentIP = ref('')
const isLAN = ref(true)
const isConnected = ref(true)
const lastSwitch = ref(null)
const checkInterval = ref(null)

const connectionType = computed(() => {
  if (!isConnected.value) return 'error'
  return isLAN.value ? 'success' : 'info'
})

const connectionLabel = computed(() => {
  if (!isConnected.value) return '未连接'
  return isLAN.value ? '局域网' : 'ZeroTier'
})

const connectionIcon = computed(() => {
  if (!isConnected.value) return Close
  return isLAN.value ? Wifi : Globe
})

// 检查连接状态
async function checkConnection() {
  try {
    const deviceInfoStr = localStorage.getItem('rnctl_device_info')
    if (!deviceInfoStr) return

    const deviceInfo = JSON.parse(deviceInfoStr)
    const result = await TestDeviceConnection(JSON.stringify(deviceInfo))
    const data = JSON.parse(result)

    if (data.success) {
      const newIP = data.connected_ip
      const oldIP = currentIP.value

      // 检测到IP切换
      if (oldIP && oldIP !== newIP) {
        const isNewLAN = deviceInfo.lan_ips.includes(newIP)
        const reason = isNewLAN ? 'LAN连接恢复' : '通过ZeroTier连接'

        lastSwitch.value = {
          from: oldIP,
          to: newIP,
          reason: reason,
          time: new Date().toLocaleTimeString('zh-CN')
        }

        // 显示切换通知
        window.$message?.info(`连接已切换: ${reason}`)
      }

      currentIP.value = newIP
      isLAN.value = deviceInfo.lan_ips.includes(newIP)
      isConnected.value = true

      // 更新localStorage中的当前IP
      localStorage.setItem('rnctl_device_ip', newIP)
    } else {
      isConnected.value = false
    }
  } catch (error) {
    console.error('连接检查失败:', error)
    isConnected.value = false
  }
}

onMounted(() => {
  // 初始检查
  checkConnection()

  // 每30秒检查一次
  checkInterval.value = setInterval(checkConnection, 30000)
})

onUnmounted(() => {
  if (checkInterval.value) {
    clearInterval(checkInterval.value)
  }
})
</script>
