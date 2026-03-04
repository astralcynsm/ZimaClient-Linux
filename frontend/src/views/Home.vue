<template>
  <div class="home-view">
    <n-space vertical :size="32">
      <n-card title="欢迎使用 rnctl">
        <n-space vertical>
          <p>rnctl 是一个强大的远程网络控制工具，帮助你轻松管理 ZimaOS 设备和网络。</p>
          <n-divider />
          <n-steps :current="currentStep" :status="stepStatus">
            <n-step title="安装 ZeroTier" description="确保 ZeroTier 已安装" />
            <n-step title="加入网络" description="连接到 ZeroTier 网络" />
            <n-step title="发现设备" description="扫描局域网设备" />
          </n-steps>
        </n-space>
      </n-card>

      <n-card title="快速开始">
        <n-space vertical :size="16">
          <n-alert v-if="!ztInstalled" type="warning" title="ZeroTier 未安装">
            请先安装 ZeroTier 才能使用完整功能
          </n-alert>

          <n-button
            v-if="!ztInstalled"
            type="primary"
            size="large"
            :loading="installing"
            @click="installZeroTier"
          >
            {{ installing ? '正在安装...' : '安装 ZeroTier' }}
          </n-button>

          <n-button
            v-else
            type="success"
            size="large"
            @click="goToZeroTier"
          >
            前往 ZeroTier 管理
          </n-button>
        </n-space>
      </n-card>
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  NSpace,
  NCard,
  NDivider,
  NSteps,
  NStep,
  NAlert,
  NButton,
  useMessage
} from 'naive-ui'
import { EnsureZeroTierInstalled, GetZeroTierInfo } from '../../wailsjs/go/main/App'

const message = useMessage()
const ztInstalled = ref(false)
const installing = ref(false)
const currentStep = ref(1)
const stepStatus = ref('process')

const emit = defineEmits(['navigate'])

onMounted(async () => {
  await checkZeroTier()
})

async function checkZeroTier() {
  try {
    await GetZeroTierInfo()
    ztInstalled.value = true
    currentStep.value = 2
  } catch (e) {
    ztInstalled.value = false
    currentStep.value = 1
  }
}

async function installZeroTier() {
  installing.value = true
  try {
    await EnsureZeroTierInstalled()
    message.success('ZeroTier 安装成功！')
    ztInstalled.value = true
    currentStep.value = 2
  } catch (e) {
    message.error('安装失败: ' + e)
  } finally {
    installing.value = false
  }
}

function goToZeroTier() {
  emit('navigate', 'zerotier')
}
</script>
