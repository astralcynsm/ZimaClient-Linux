<template>
  <div class="login-page">
    <div class="content">
      <div class="header">
        <div class="device-icon">
          <n-icon size="48" :component="Server" />
        </div>
        <h1>{{ deviceName }}</h1>
        <p class="subtitle">{{ deviceIp }}</p>
      </div>

      <div class="login-card">
        <n-form ref="formRef" :model="formData" :rules="rules" size="large">
          <n-form-item path="username" label="用户名">
            <n-input
              v-model:value="formData.username"
              placeholder="输入用户名"
              @keyup.enter="handleLogin"
            />
          </n-form-item>

          <n-form-item path="password" label="密码">
            <n-input
              v-model:value="formData.password"
              type="password"
              show-password-on="click"
              placeholder="输入密码"
              @keyup.enter="handleLogin"
            />
          </n-form-item>

          <n-form-item>
            <n-checkbox v-model:checked="rememberPassword">
              记住密码
            </n-checkbox>
          </n-form-item>

          <div class="button-group">
            <n-button
              type="primary"
              size="large"
              block
              :loading="logging"
              @click="handleLogin"
            >
              登录
            </n-button>

            <n-button
              size="large"
              block
              quaternary
              @click="goBack"
            >
              返回
            </n-button>
          </div>
        </n-form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { Server } from '@vicons/ionicons5'
import { SaveCredentials } from '../../../wailsjs/go/main/App'

const router = useRouter()
const message = useMessage()

const deviceIp = ref('')
const deviceName = ref('')
const formRef = ref(null)
const logging = ref(false)
const rememberPassword = ref(true)

const formData = ref({
  username: '',
  password: ''
})

const rules = {
  username: {
    required: true,
    message: '请输入用户名',
    trigger: 'blur'
  },
  password: {
    required: true,
    message: '请输入密码',
    trigger: 'blur'
  }
}

onMounted(async () => {
  deviceIp.value = localStorage.getItem('rnctl_device_ip') || ''
  deviceName.value = localStorage.getItem('rnctl_device_hostname') || deviceIp.value

  if (!deviceIp.value) {
    message.warning('未选择设备')
    router.push('/wizard/connect')
    return
  }

  // 尝试加载保存的凭证
  try {
    const { GetCredentials } = await import('../../../wailsjs/go/main/App')
    const credData = await GetCredentials(deviceIp.value)
    if (credData) {
      const parsed = JSON.parse(credData)
      formData.value.username = parsed.username
      formData.value.password = parsed.password
      rememberPassword.value = true
    }
  } catch (error) {
    // 没有保存的凭证，尝试只加载用户名
    const savedUsername = localStorage.getItem('rnctl_username')
    if (savedUsername) {
      formData.value.username = savedUsername
    }
  }
})

async function handleLogin() {
  try {
    await formRef.value?.validate()
  } catch {
    return
  }

  logging.value = true

  try {
    // 这里应该调用实际的登录 API
    // 暂时模拟登录成功
    await new Promise(resolve => setTimeout(resolve, 1000))

    // 保存凭证
    if (rememberPassword.value) {
      const credData = JSON.stringify({
        username: formData.value.username,
        password: formData.value.password
      })
      await SaveCredentials(deviceIp.value, credData)
    }

    // 保存登录状态
    localStorage.setItem('rnctl_logged_in', 'true')
    localStorage.setItem('rnctl_username', formData.value.username)

    message.success('登录成功')
    // 跳转到自动挂载配置页面
    router.push('/wizard/auto-mount')
  } catch (error) {
    message.error('登录失败: ' + error)
  } finally {
    logging.value = false
  }
}

function goBack() {
  router.push('/wizard/connect')
}
</script>

<style scoped>
.login-page {
  width: 100%;
  min-height: 100vh;
  padding: 40px 20px;
  background: var(--bg-secondary);
  overflow-y: auto;
  display: flex;
  align-items: center;
  justify-content: center;
}

.content {
  width: 100%;
  max-width: 420px;
}

.header {
  text-align: center;
  margin-bottom: 32px;
}

.device-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 20px;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-400) 100%);
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.header h1 {
  font-size: 24px;
  font-weight: 600;
  margin: 0 0 8px 0;
  color: var(--text-primary);
}

.subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0;
}

.login-card {
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 32px;
}

@media (prefers-color-scheme: dark) {
  .login-card {
    background: rgba(23, 23, 23, 0.8);
    border-color: rgba(64, 64, 64, 0.5);
  }
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 8px;
}
</style>
