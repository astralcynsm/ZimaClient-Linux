<template>
  <div class="dashboard">
    <!-- Device Info Card -->
    <div class="info-card">
      <div class="info-header">
        <h2 class="info-title">{{ deviceName }}</h2>
        <n-tag :type="isConnected ? 'success' : 'error'" size="small" round>
          {{ isConnected ? '已连接' : '未连接' }}
        </n-tag>
      </div>
      <div class="info-details">
        <div class="info-item">
          <span class="info-label">IP 地址</span>
          <span class="info-value">{{ deviceIp }}</span>
        </div>
        <div class="info-item">
          <span class="info-label">用户名</span>
          <span class="info-value">{{ username }}</span>
        </div>
      </div>
    </div>

    <!-- Quick Access -->
    <div v-if="hasQuickAccess" class="quick-access-card" @click="openQuickAccess">
      <div class="quick-access-icon">
        <n-icon size="32" :component="FolderOpen" />
      </div>
      <div class="quick-access-content">
        <div class="quick-access-label">快速访问</div>
        <div class="quick-access-name">{{ quickAccessName }}</div>
      </div>
      <div class="quick-access-arrow">
        <n-icon size="20" :component="ChevronForward" />
      </div>
    </div>

    <!-- Stats Grid -->
    <div class="stats-grid">
      <div class="stat-card" @click="navigateTo('storage')">
        <div class="stat-icon storage">
          <n-icon size="28" :component="FolderOpen" />
        </div>
        <div class="stat-content">
          <div class="stat-label">存储</div>
          <div class="stat-value">{{ mountedCount }} 个挂载</div>
        </div>
      </div>

      <div class="stat-card" @click="navigateTo('zerotier')">
        <div class="stat-icon network">
          <n-icon size="28" :component="Globe" />
        </div>
        <div class="stat-content">
          <div class="stat-label">ZeroTier</div>
          <div class="stat-value">{{ ztStatus }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon system">
          <n-icon size="28" :component="Server" />
        </div>
        <div class="stat-content">
          <div class="stat-label">系统</div>
          <div class="stat-value">{{ uptime }}</div>
        </div>
      </div>
    </div>

    <!-- Actions -->
    <div class="actions-section">
      <n-button
        type="primary"
        size="large"
        @click="navigateTo('storage')"
        style="flex: 1"
      >
        <template #icon>
          <n-icon><FolderOpen /></n-icon>
        </template>
        管理存储
      </n-button>
      <n-button
        size="large"
        @click="refreshStatus"
        quaternary
      >
        <template #icon>
          <n-icon><Refresh /></n-icon>
        </template>
        刷新
      </n-button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { FolderOpen, Globe, Server, Refresh, ChevronForward } from '@vicons/ionicons5'
import { GetZeroTierInfo, ListMounts, GetSystemUptime, OpenQuickAccess, GetAutoMountConfig } from '../../../wailsjs/go/main/App'

const router = useRouter()
const message = useMessage()

const deviceName = ref('')
const deviceIp = ref('')
const username = ref('')
const isConnected = ref(true)
const mountedCount = ref(0)
const ztStatus = ref('加载中...')
const uptime = ref('--')
const hasQuickAccess = ref(false)
const quickAccessName = ref('')

onMounted(() => {
  deviceName.value = localStorage.getItem('rnctl_device_hostname') || '未知'
  deviceIp.value = localStorage.getItem('rnctl_device_ip') || '未知'
  username.value = localStorage.getItem('rnctl_username') || '未知'

  loadStatus()
})

async function loadStatus() {
  // 加载ZeroTier状态
  try {
    const status = await GetZeroTierInfo()
    const parsed = JSON.parse(status)
    ztStatus.value = parsed.online ? '在线' : '离线'
  } catch (error) {
    console.error('Failed to load ZeroTier status:', error)
    ztStatus.value = '未知'
  }

  // 加载挂载数量
  try {
    const mounts = await ListMounts()
    const parsed = JSON.parse(mounts)
    mountedCount.value = Array.isArray(parsed) ? parsed.length : 0
  } catch (error) {
    console.error('Failed to load mounts:', error)
    mountedCount.value = 0
  }

  // 加载系统运行时间
  try {
    const uptimeData = await GetSystemUptime()
    const parsed = JSON.parse(uptimeData)
    uptime.value = parsed.uptime_text || '--'
  } catch (error) {
    console.error('Failed to load uptime:', error)
    uptime.value = '--'
  }

  // 检查是否有快速访问配置
  try {
    const config = await GetAutoMountConfig()
    const parsed = JSON.parse(config)
    if (parsed.quickAccess && parsed.quickAccess.name) {
      hasQuickAccess.value = true
      quickAccessName.value = parsed.quickAccess.name
    }
  } catch (error) {
    console.error('Failed to load auto mount config:', error)
  }
}

function navigateTo(path) {
  router.push(`/main/${path}`)
}

function refreshStatus() {
  loadStatus()
  message.success('状态已刷新')
}

async function openQuickAccess() {
  try {
    await OpenQuickAccess()
    message.success('已打开快速访问')
  } catch (error) {
    message.error('打开失败: ' + error)
  }
}
</script>

<style scoped>
.dashboard {
  max-width: 800px;
  margin: 0 auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

/* Info Card */
.info-card {
  background: rgba(var(--card-bg-rgb), 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 20px;
  transition: all 0.2s ease-in-out;
}

.info-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
}

.info-title {
  font-size: 20px;
  font-weight: 600;
  color: var(--text-primary);
}

.info-details {
  display: flex;
  gap: 32px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: var(--text-tertiary);
}

.info-value {
  font-size: 14px;
  color: var(--text-primary);
  font-weight: 500;
}

/* Quick Access Card */
.quick-access-card {
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-400) 100%);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out;
  color: white;
  will-change: transform;
}

.quick-access-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 87, 255, 0.3);
}

.quick-access-card:active {
  transform: translateY(0);
}

.quick-access-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border-radius: 12px;
}

.quick-access-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.quick-access-label {
  font-size: 12px;
  opacity: 0.9;
}

.quick-access-name {
  font-size: 18px;
  font-weight: 600;
}

.quick-access-arrow {
  opacity: 0.8;
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.stat-card {
  background: rgba(var(--card-bg-rgb), 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 12px;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  will-change: transform;
}

.stat-card:hover {
  background: rgba(var(--card-hover-bg-rgb), 0.9);
  border-color: var(--border-hover);
  transform: translateY(-2px);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.storage {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.network {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.system {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-content {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-label {
  font-size: 12px;
  color: var(--text-tertiary);
}

.stat-value {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

/* Actions Section */
.actions-section {
  display: flex;
  gap: 12px;
  margin-top: 8px;
}
</style>
