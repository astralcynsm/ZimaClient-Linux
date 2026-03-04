<template>
  <n-modal v-model:show="showModal" :mask-closable="false" :close-on-esc="false">
    <n-card
      style="width: 400px"
      title="需要管理员权限"
      :bordered="false"
      size="huge"
      role="dialog"
      aria-modal="true"
    >
      <n-space vertical>
        <n-alert type="info">
          应用需要管理员权限来挂载网络存储。
        </n-alert>

        <n-radio-group v-model:value="authMethod">
          <n-space vertical>
            <n-radio value="polkit">
              使用系统授权对话框（推荐）
            </n-radio>
            <n-radio value="password">
              输入 sudo 密码
            </n-radio>
          </n-space>
        </n-radio-group>

        <n-input
          v-if="authMethod === 'password'"
          v-model:value="password"
          type="password"
          placeholder="请输入 sudo 密码"
          show-password-on="click"
          @keyup.enter="handleSubmit"
        />

        <n-space justify="end">
          <n-button @click="handleCancel">取消</n-button>
          <n-button type="primary" :loading="loading" @click="handleSubmit">
            确认
          </n-button>
        </n-space>
      </n-space>
    </n-card>
  </n-modal>
</template>

<script setup>
import { ref } from 'vue'
import { useMessage } from 'naive-ui'
import { InitSudoPassword } from '../../wailsjs/go/main/App'

const message = useMessage()
const showModal = ref(false)
const password = ref('')
const loading = ref(false)
const authMethod = ref('polkit')

const emit = defineEmits(['success', 'cancel'])

function show() {
  showModal.value = true
  password.value = ''
  authMethod.value = 'polkit'
}

async function handleSubmit() {
  if (authMethod.value === 'password' && !password.value) {
    message.warning('请输入密码')
    return
  }

  loading.value = true
  try {
    // 如果选择 polkit，传空密码（后端会自动使用 pkexec）
    const pwd = authMethod.value === 'polkit' ? '' : password.value
    await InitSudoPassword(pwd)
    message.success('授权成功')
    showModal.value = false
    emit('success')
  } catch (error) {
    message.error('授权失败: ' + error)
  } finally {
    loading.value = false
  }
}

function handleCancel() {
  showModal.value = false
  emit('cancel')
}

defineExpose({
  show
})
</script>
