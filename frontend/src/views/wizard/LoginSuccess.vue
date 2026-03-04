<template>
  <div class="success-page">
    <div class="content">
      <div class="icon-wrapper">
        <n-icon size="80" color="#18a058">
          <CheckmarkCircle />
        </n-icon>
      </div>

      <h1>登录成功</h1>
      <p class="message">正在为您准备工作空间...</p>

      <n-progress
        type="line"
        :percentage="progress"
        :show-indicator="false"
        status="success"
      />

      <n-space vertical :size="12" style="margin-top: 40px;">
        <n-button
          type="primary"
          size="large"
          block
          @click="goToDashboard"
        >
          进入仪表板
        </n-button>
      </n-space>
    </div>

    <!-- Sudo 密码对话框 -->
    <SudoPasswordDialog
      ref="sudoDialog"
      @success="onSudoSuccess"
      @cancel="onSudoCancel"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { CheckmarkCircle } from '@vicons/ionicons5'
import SudoPasswordDialog from '../../components/SudoPasswordDialog.vue'
import { CheckSudoInitialized, ProcessAutoMount } from '../../../wailsjs/go/main/App'

const router = useRouter()
const message = useMessage()
const progress = ref(0)
const sudoDialog = ref(null)
const sudoInitialized = ref(false)

onMounted(async () => {
  // 检查 sudo 是否已初始化
  try {
    sudoInitialized.value = await CheckSudoInitialized()
  } catch (error) {
    console.error('检查 sudo 状态失败:', error)
  }

  // 尝试执行自动挂载
  try {
    await ProcessAutoMount()
    console.log('自动挂载处理完成')
  } catch (error) {
    console.error('自动挂载失败:', error)
    // 自动挂载失败不影响登录流程
  }

  // 模拟加载进度
  const interval = setInterval(() => {
    progress.value += 10
    if (progress.value >= 100) {
      clearInterval(interval)
      // 加载完成后，检查是否需要输入密码
      setTimeout(() => {
        if (!sudoInitialized.value) {
          sudoDialog.value?.show()
        } else {
          goToDashboard()
        }
      }, 500)
    }
  }, 200)
})

function onSudoSuccess() {
  sudoInitialized.value = true
  message.success('权限验证成功')
  goToDashboard()
}

function onSudoCancel() {
  message.warning('需要管理员权限才能使用存储功能')
  // 仍然可以进入仪表板，但存储功能会受限
  goToDashboard()
}

function goToDashboard() {
  router.push('/main/dashboard')
}
</script>

<style scoped>
.success-page {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 100vh;
  padding: 40px;
}

.content {
  width: 100%;
  max-width: 400px;
  text-align: center;
}

.icon-wrapper {
  margin-bottom: 24px;
}

h1 {
  font-size: 28px;
  font-weight: 600;
  margin: 0 0 12px 0;
}

.message {
  font-size: 16px;
  opacity: 0.7;
  margin: 0 0 32px 0;
}
</style>
