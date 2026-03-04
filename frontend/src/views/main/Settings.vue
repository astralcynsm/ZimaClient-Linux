<template>
  <div class="settings-view">
    <n-space vertical :size="20">
      <!-- 应用信息 -->
      <n-card title="应用信息">
        <n-descriptions bordered :column="1" size="small">
          <n-descriptions-item label="应用名称">
            ZimaClient - Zima 设备管理客户端
          </n-descriptions-item>
          <n-descriptions-item label="版本">
            v0.1.0-alpha
          </n-descriptions-item>
          <n-descriptions-item label="构建日期">
            {{ buildDate }}
          </n-descriptions-item>
          <n-descriptions-item label="当前设备">
            {{ deviceName }}
          </n-descriptions-item>
          <n-descriptions-item label="设备IP">
            {{ deviceIP }}
          </n-descriptions-item>
        </n-descriptions>
      </n-card>

      <!-- 连接设置 -->
      <n-card title="连接设置">
        <n-space vertical :size="16">
          <n-form-item label="连接超时时间">
            <n-input-number
              v-model:value="connectionTimeout"
              :min="5"
              :max="60"
              :step="5"
              style="width: 200px"
            >
              <template #suffix>秒</template>
            </n-input-number>
          </n-form-item>

          <n-form-item label="健康检查间隔">
            <n-input-number
              v-model:value="healthCheckInterval"
              :min="10"
              :max="300"
              :step="10"
              style="width: 200px"
            >
              <template #suffix>秒</template>
            </n-input-number>
          </n-form-item>

          <n-form-item label="自动重连">
            <n-switch v-model:value="autoReconnect">
              <template #checked>开启</template>
              <template #unchecked>关闭</template>
            </n-switch>
          </n-form-item>

          <n-space>
            <n-button type="primary" @click="saveConnectionSettings">
              保存设置
            </n-button>
            <n-button @click="resetConnectionSettings">
              恢复默认
            </n-button>
          </n-space>
        </n-space>
      </n-card>

      <!-- 存储设置 -->
      <n-card title="存储设置">
        <n-space vertical :size="16">
          <n-form-item label="默认挂载路径">
            <n-input-group>
              <n-input
                v-model:value="defaultMountPath"
                placeholder="/run/user/$(id -u)/zimaclient"
                readonly
                style="flex: 1"
              />
              <n-button disabled>
                运行时目录
              </n-button>
            </n-input-group>
          </n-form-item>

          <n-alert type="info" :show-icon="false" size="small">
            挂载路径使用运行时目录，避免影响文件管理器性能。格式：{{ defaultMountPath }}/设备名-共享名
          </n-alert>

          <n-form-item label="挂载后自动打开">
            <n-switch v-model:value="autoOpenAfterMount">
              <template #checked>开启</template>
              <template #unchecked>关闭</template>
            </n-switch>
          </n-form-item>

          <n-form-item label="退出时自动卸载">
            <n-switch v-model:value="autoUnmountOnExit">
              <template #checked>开启</template>
              <template #unchecked>关闭</template>
            </n-switch>
          </n-form-item>

          <n-space>
            <n-button type="primary" @click="saveStorageSettings">
              保存设置
            </n-button>
          </n-space>
        </n-space>
      </n-card>

      <!-- ZeroTier设置 -->
      <n-card title="ZeroTier设置">
        <n-space vertical :size="16">
          <n-alert type="info" :show-icon="false">
            ZeroTier状态: <n-tag :type="ztOnline ? 'success' : 'error'" size="small">
              {{ ztOnline ? '在线' : '离线' }}
            </n-tag>
          </n-alert>

          <n-form-item label="启动时自动连接">
            <n-switch v-model:value="ztAutoConnect">
              <template #checked>开启</template>
              <template #unchecked>关闭</template>
            </n-switch>
          </n-form-item>

          <n-space>
            <n-button type="primary" @click="saveZTSettings">
              保存设置
            </n-button>
            <n-button @click="openZTPage">
              管理网络
            </n-button>
          </n-space>
        </n-space>
      </n-card>

      <!-- 高级设置 -->
      <n-card title="高级设置">
        <n-space vertical :size="16">
          <n-form-item label="调试模式">
            <n-switch v-model:value="debugMode">
              <template #checked>开启</template>
              <template #unchecked>关闭</template>
            </n-switch>
          </n-form-item>

          <n-divider />

          <n-space>
            <n-button type="error" @click="clearCache">
              清除缓存
            </n-button>
            <n-button type="warning" @click="resetAllSettings">
              重置所有设置
            </n-button>
          </n-space>
        </n-space>
      </n-card>

      <!-- 关于 -->
      <n-card title="关于">
        <n-space vertical :size="12">
          <p style="margin: 0; color: var(--text-secondary)">
            ZimaClient 是一个用于远程管理 ZimaOS 设备的客户端工具，支持局域网和 ZeroTier 智能切换。
          </p>
          <n-divider style="margin: 12px 0" />
          <n-space :size="12">
            <n-button text tag="a" href="https://github.com" target="_blank">
              GitHub
            </n-button>
            <n-button text @click="showLicense">
              开源许可
            </n-button>
            <n-button text @click="checkUpdate">
              检查更新
            </n-button>
          </n-space>
        </n-space>
      </n-card>
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage, useDialog } from 'naive-ui'
import { GetZeroTierInfo } from '../../../wailsjs/go/main/App'

const router = useRouter()
const message = useMessage()
const dialog = useDialog()

// 应用信息
const buildDate = ref(new Date().toLocaleDateString('zh-CN'))
const deviceName = ref('')
const deviceIP = ref('')

// 连接设置
const connectionTimeout = ref(15)
const healthCheckInterval = ref(30)
const autoReconnect = ref(true)

// 存储设置
const defaultMountPath = ref('')
const autoOpenAfterMount = ref(false)
const autoUnmountOnExit = ref(true) // 默认开启

// ZeroTier设置
const ztOnline = ref(false)
const ztAutoConnect = ref(false)

// 高级设置
const debugMode = ref(false)

onMounted(async () => {
  // 加载设备信息
  deviceName.value = localStorage.getItem('rnctl_device_hostname') || '未连接'
  deviceIP.value = localStorage.getItem('rnctl_device_ip') || '未知'

  // 加载运行时挂载目录
  try {
    const { GetRuntimeMountDir } = await import('../../../wailsjs/go/main/App')
    const runtimeDir = await GetRuntimeMountDir()
    defaultMountPath.value = runtimeDir
  } catch (error) {
    console.error('获取运行时目录失败:', error)
    defaultMountPath.value = '/run/user/$(id -u)/zimaclient'
  }

  // 加载设置
  loadSettings()

  // 检查ZeroTier状态
  try {
    const result = await GetZeroTierInfo()
    const data = JSON.parse(result)
    ztOnline.value = data.online || false
  } catch (error) {
    console.error('获取ZeroTier状态失败:', error)
  }
})

function loadSettings() {
  // 连接设置
  connectionTimeout.value = parseInt(localStorage.getItem('rnctl_connection_timeout') || '15')
  healthCheckInterval.value = parseInt(localStorage.getItem('rnctl_health_check_interval') || '30')
  autoReconnect.value = localStorage.getItem('rnctl_auto_reconnect') !== 'false'

  // 存储设置（defaultMountPath在onMounted中动态获取，不从localStorage读取）
  autoOpenAfterMount.value = localStorage.getItem('rnctl_auto_open_after_mount') === 'true'
  autoUnmountOnExit.value = localStorage.getItem('rnctl_auto_unmount_on_exit') !== 'false' // 默认true

  // ZeroTier设置
  ztAutoConnect.value = localStorage.getItem('rnctl_zt_auto_connect') === 'true'

  // 高级设置
  debugMode.value = localStorage.getItem('rnctl_debug_mode') === 'true'
}

function saveConnectionSettings() {
  localStorage.setItem('rnctl_connection_timeout', connectionTimeout.value.toString())
  localStorage.setItem('rnctl_health_check_interval', healthCheckInterval.value.toString())
  localStorage.setItem('rnctl_auto_reconnect', autoReconnect.value.toString())
  message.success('连接设置已保存')
}

function resetConnectionSettings() {
  connectionTimeout.value = 15
  healthCheckInterval.value = 30
  autoReconnect.value = true
  message.info('已恢复默认设置')
}

function saveStorageSettings() {
  localStorage.setItem('rnctl_default_mount_path', defaultMountPath.value)
  localStorage.setItem('rnctl_auto_open_after_mount', autoOpenAfterMount.value.toString())
  localStorage.setItem('rnctl_auto_unmount_on_exit', autoUnmountOnExit.value.toString())
  message.success('存储设置已保存')
}

function saveZTSettings() {
  localStorage.setItem('rnctl_zt_auto_connect', ztAutoConnect.value.toString())
  message.success('ZeroTier设置已保存')
}

function openZTPage() {
  router.push('/main/zerotier')
}

function clearCache() {
  dialog.warning({
    title: '清除缓存',
    content: '确定要清除所有缓存数据吗？这不会删除已保存的密码和设置。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      // 清除缓存但保留重要数据
      const keysToKeep = [
        'rnctl_device_info',
        'rnctl_device_ip',
        'rnctl_device_hostname',
        'rnctl_logged_in',
        'rnctl_username'
      ]

      const tempData = {}
      keysToKeep.forEach(key => {
        const value = localStorage.getItem(key)
        if (value) tempData[key] = value
      })

      localStorage.clear()

      Object.entries(tempData).forEach(([key, value]) => {
        localStorage.setItem(key, value)
      })

      message.success('缓存已清除')
    }
  })
}

function resetAllSettings() {
  dialog.error({
    title: '重置所有设置',
    content: '确定要重置所有设置吗？这将清除所有配置并退出登录。',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      localStorage.clear()
      message.success('设置已重置')
      setTimeout(() => {
        router.push('/wizard/connect')
      }, 1000)
    }
  })
}

function showLicense() {
  dialog.info({
    title: '开源许可',
    content: 'MIT License\n\nCopyright (c) 2025 ZimaClient\n\n本软件采用MIT许可证开源。',
    positiveText: '确定'
  })
}

function checkUpdate() {
  message.info('当前已是最新版本 v0.1.0-alpha')
}
</script>

<style scoped>
.settings-view {
  max-width: 800px;
  padding: 20px;
}
</style>
