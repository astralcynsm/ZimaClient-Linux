<template>
  <div class="auto-mount-container">
    <div class="wizard-content">
      <div class="header">
        <h1 class="title">配置自动挂载</h1>
        <p class="subtitle">选择您希望自动挂载的共享文件夹，下次登录时将自动挂载</p>
      </div>

      <n-spin :show="loading">
        <div class="content">
          <!-- 共享列表 -->
          <div class="shares-section">
            <n-checkbox-group v-model:value="selectedMounts">
              <n-space vertical :size="12">
                <n-checkbox
                  v-for="share in shares"
                  :key="share.name"
                  :value="share.name"
                  :disabled="selectedMounts.length === 1 && selectedMounts.includes(share.name)"
                  class="share-checkbox"
                >
                  <div class="share-item">
                    <n-icon size="32" color="#0ea5e9">
                      <FolderOpen />
                    </n-icon>
                    <div class="share-info">
                      <n-text strong>{{ share.name }}</n-text>
                      <n-text depth="3" style="font-size: 12px">
                        {{ share.comment || '共享文件夹' }}
                      </n-text>
                    </div>
                  </div>
                </n-checkbox>
              </n-space>
            </n-checkbox-group>
          </div>

          <n-divider />

          <!-- 快速访问 -->
          <div class="quick-access-section">
            <div class="section-header">
              <n-text strong>快速访问</n-text>
              <n-tooltip trigger="hover">
                <template #trigger>
                  <n-icon size="16" style="cursor: help">
                    <HelpCircle />
                  </n-icon>
                </template>
                选择一个共享作为快速访问，可以从Dashboard一键打开
              </n-tooltip>
            </div>

            <n-select
              v-model:value="quickAccess"
              :options="mountOptions"
              placeholder="选择快速访问的文件夹"
              :disabled="selectedMounts.length === 0"
            />
          </div>
        </div>
      </n-spin>

      <!-- 操作按钮 -->
      <div class="actions">
        <n-button @click="skip" :disabled="saving">
          跳过
        </n-button>
        <n-button
          type="primary"
          @click="saveAndContinue"
          :loading="saving"
          :disabled="selectedMounts.length === 0"
        >
          完成配置
        </n-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { FolderOpen, HelpCircle } from '@vicons/ionicons5'
import {
  ListSMBShares,
  SaveAutoMountConfig,
  GetHomeDir
} from '../../../wailsjs/go/main/App'

const router = useRouter()
const message = useMessage()

const loading = ref(true)
const saving = ref(false)
const shares = ref([])
const selectedMounts = ref([])
const quickAccess = ref('')
const deviceIp = ref('')

const mountOptions = computed(() => {
  return selectedMounts.value.map(name => ({
    label: name,
    value: name
  }))
})

onMounted(async () => {
  deviceIp.value = localStorage.getItem('rnctl_device_ip') || ''

  if (!deviceIp.value) {
    message.error('未找到设备信息')
    router.push('/wizard/connect')
    return
  }

  // 加载共享列表
  try {
    const result = await ListSMBShares(deviceIp.value)
    shares.value = JSON.parse(result)

    // 默认选择第一个共享
    if (shares.value.length > 0) {
      selectedMounts.value = [shares.value[0].name]
      quickAccess.value = shares.value[0].name
    }
  } catch (error) {
    console.error('加载共享列表失败:', error)
    message.error('加载共享列表失败: ' + error)
  } finally {
    loading.value = false
  }
})

async function saveAndContinue() {
  if (selectedMounts.value.length === 0) {
    message.warning('请至少选择一个共享')
    return
  }

  saving.value = true

  try {
    // 获取用户主目录
    const homeDir = await GetHomeDir()

    // 构建配置
    const config = {
      enabled: true,
      configured: true,
      deviceIp: deviceIp.value,
      mounts: selectedMounts.value.map(name => ({
        name: name,
        remotePath: `//${deviceIp.value}/${name}`,
        localPath: `${homeDir}/${name}`
      })),
      quickAccess: quickAccess.value ? {
        name: quickAccess.value,
        localPath: `${homeDir}/${quickAccess.value}`
      } : null
    }

    // 保存配置
    await SaveAutoMountConfig(JSON.stringify(config))

    message.success('配置已保存，将在下次登录时自动挂载')

    // 跳转到登录成功页面
    router.push('/wizard/login-success')
  } catch (error) {
    console.error('保存配置失败:', error)
    message.error('保存配置失败: ' + error)
  } finally {
    saving.value = false
  }
}

function skip() {
  // 保存一个空配置，表示用户选择跳过
  const config = {
    enabled: false,
    configured: true,
    deviceIp: deviceIp.value,
    mounts: [],
    quickAccess: null
  }

  SaveAutoMountConfig(JSON.stringify(config))
    .then(() => {
      router.push('/wizard/login-success')
    })
    .catch(error => {
      console.error('保存配置失败:', error)
      router.push('/wizard/login-success')
    })
}
</script>

<style scoped>
.auto-mount-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.wizard-content {
  width: 100%;
  max-width: 600px;
  background: white;
  border-radius: 16px;
  padding: 40px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.header {
  text-align: center;
  margin-bottom: 32px;
}

.title {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 12px 0;
  color: #1a1a1a;
}

.subtitle {
  font-size: 14px;
  color: #666;
  margin: 0;
  line-height: 1.6;
}

.content {
  margin-bottom: 32px;
}

.shares-section {
  margin-bottom: 24px;
}

.share-checkbox {
  width: 100%;
}

.share-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f9fafb;
  transition: all 0.2s;
}

.share-item:hover {
  background: #f3f4f6;
  border-color: #d1d5db;
}

.share-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.quick-access-section {
  margin-top: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 12px;
}

.actions {
  display: flex;
  justify-content: space-between;
  gap: 16px;
}

.actions .n-button {
  flex: 1;
}
</style>
