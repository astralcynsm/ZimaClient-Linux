<template>
  <n-space align="center" :size="8">
    <n-badge :dot="true" :type="statusType" :processing="connecting">
      <n-icon :size="20" :color="statusColor">
        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
        </svg>
      </n-icon>
    </n-badge>
    <n-text :depth="3" style="font-size: 12px;">
      {{ statusText }}
    </n-text>
  </n-space>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { NSpace, NBadge, NIcon, NText } from 'naive-ui'
import { GetZeroTierInfo } from '../../wailsjs/go/main/App'

const props = defineProps({
  autoRefresh: {
    type: Boolean,
    default: true
  }
})

const status = ref('disconnected') // disconnected, connecting, connected
const connecting = ref(false)
let refreshTimer = null

const statusType = computed(() => {
  switch (status.value) {
    case 'connected':
      return 'success'
    case 'connecting':
      return 'warning'
    default:
      return 'error'
  }
})

const statusColor = computed(() => {
  switch (status.value) {
    case 'connected':
      return '#18a058'
    case 'connecting':
      return '#f0a020'
    default:
      return '#d03050'
  }
})

const statusText = computed(() => {
  switch (status.value) {
    case 'connected':
      return '已连接'
    case 'connecting':
      return '连接中...'
    default:
      return '未连接'
  }
})

const checkStatus = async () => {
  try {
    connecting.value = true
    const info = await GetZeroTierInfo()
    if (info && info.includes('ONLINE')) {
      status.value = 'connected'
    } else {
      status.value = 'disconnected'
    }
  } catch (error) {
    status.value = 'disconnected'
  } finally {
    connecting.value = false
  }
}

onMounted(() => {
  checkStatus()

  if (props.autoRefresh) {
    refreshTimer = setInterval(checkStatus, 10000) // 每 10 秒刷新一次
  }
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
})
</script>
