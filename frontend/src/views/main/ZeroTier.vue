<template>
  <div class="zerotier-view">
    <n-space vertical :size="20">
      <!-- ZeroTier 状态卡片 -->
      <n-card title="ZeroTier 状态">
        <n-descriptions :column="2" bordered>
          <n-descriptions-item label="节点ID">
            {{ ztInfo.address || '--' }}
          </n-descriptions-item>
          <n-descriptions-item label="在线状态">
            <n-tag :type="ztInfo.online ? 'success' : 'error'">
              {{ ztInfo.online ? '在线' : '离线' }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="版本">
            {{ ztInfo.version || '--' }}
          </n-descriptions-item>
          <n-descriptions-item label="操作">
            <n-button
              size="small"
              @click="refreshStatus"
              :loading="loadingStatus"
            >
              刷新状态
            </n-button>
          </n-descriptions-item>
        </n-descriptions>
      </n-card>

      <!-- 加入网络 -->
      <n-card title="加入网络">
        <n-space vertical>
          <n-alert type="info" :show-icon="false">
            输入 ZeroTier 网络 ID（16位十六进制字符）来加入网络
          </n-alert>
          <n-input-group>
            <n-input
              v-model:value="networkIdInput"
              placeholder="例如: 8056c2e21c000001"
              :maxlength="16"
              style="width: 70%"
            />
            <n-button
              type="primary"
              @click="joinNetwork"
              :loading="joiningNetwork"
              :disabled="!networkIdInput || networkIdInput.length !== 16"
              style="width: 30%"
            >
              加入网络
            </n-button>
          </n-input-group>
        </n-space>
      </n-card>

      <!-- 已加入的网络列表 -->
      <n-card title="已加入的网络">
        <template #header-extra>
          <n-button
            size="small"
            @click="loadNetworks"
            :loading="loadingNetworks"
          >
            刷新列表
          </n-button>
        </template>

        <n-spin :show="loadingNetworks">
          <n-empty
            v-if="networks.length === 0 && !loadingNetworks"
            description="暂无已加入的网络"
          />

          <n-list v-else bordered>
            <n-list-item v-for="network in networks" :key="network.id">
              <n-thing>
                <template #header>
                  <n-space align="center">
                    <n-text strong>{{ network.name || '未命名网络' }}</n-text>
                    <n-tag
                      :type="network.status === 'OK' ? 'success' : 'warning'"
                      size="small"
                    >
                      {{ network.status }}
                    </n-tag>
                  </n-space>
                </template>
                <template #description>
                  <n-space vertical :size="4">
                    <n-text depth="3">网络ID: {{ network.id }}</n-text>
                    <n-text depth="3" v-if="network.assignedAddresses && network.assignedAddresses.length > 0">
                      分配IP: {{ network.assignedAddresses.join(', ') }}
                    </n-text>
                    <n-text depth="3">类型: {{ network.type }}</n-text>
                  </n-space>
                </template>
                <template #action>
                  <n-popconfirm
                    @positive-click="leaveNetwork(network.id)"
                  >
                    <template #trigger>
                      <n-button
                        type="error"
                        size="small"
                        :loading="leavingNetworkId === network.id"
                      >
                        离开网络
                      </n-button>
                    </template>
                    确定要离开网络 "{{ network.name || network.id }}" 吗？
                  </n-popconfirm>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-spin>
      </n-card>
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import {
  GetZeroTierInfo,
  ListZeroTierNetworks,
  JoinZeroTierNetwork,
  LeaveZeroTierNetwork,
  EnsureZeroTierInstalled
} from '../../../wailsjs/go/main/App'

const message = useMessage()

// ZeroTier 状态
const ztInfo = ref({
  address: '',
  online: false,
  version: ''
})
const loadingStatus = ref(false)

// 网络列表
const networks = ref([])
const loadingNetworks = ref(false)

// 加入网络
const networkIdInput = ref('')
const joiningNetwork = ref(false)

// 离开网络
const leavingNetworkId = ref(null)

onMounted(async () => {
  // 确保 ZeroTier 已安装
  try {
    await EnsureZeroTierInstalled()
  } catch (error) {
    message.error('ZeroTier 安装失败: ' + error)
    return
  }

  // 加载状态和网络列表
  await refreshStatus()
  await loadNetworks()
})

const refreshStatus = async () => {
  loadingStatus.value = true
  try {
    const result = await GetZeroTierInfo()
    const data = JSON.parse(result)
    ztInfo.value = {
      address: data.address || '',
      online: data.online || false,
      version: data.version || ''
    }
  } catch (error) {
    console.error('获取 ZeroTier 状态失败:', error)
    message.error('获取状态失败: ' + error)
  } finally {
    loadingStatus.value = false
  }
}

const loadNetworks = async () => {
  loadingNetworks.value = true
  try {
    const result = await ListZeroTierNetworks()
    const data = JSON.parse(result)
    networks.value = Array.isArray(data) ? data : []
  } catch (error) {
    console.error('获取网络列表失败:', error)
    message.error('获取网络列表失败: ' + error)
    networks.value = []
  } finally {
    loadingNetworks.value = false
  }
}

const joinNetwork = async () => {
  if (!networkIdInput.value || networkIdInput.value.length !== 16) {
    message.warning('请输入有效的16位网络ID')
    return
  }

  joiningNetwork.value = true
  try {
    await JoinZeroTierNetwork(networkIdInput.value)
    message.success('成功加入网络')
    networkIdInput.value = ''

    // 等待一下再刷新列表
    setTimeout(async () => {
      await loadNetworks()
    }, 1000)
  } catch (error) {
    console.error('加入网络失败:', error)
    message.error('加入网络失败: ' + error)
  } finally {
    joiningNetwork.value = false
  }
}

const leaveNetwork = async (networkId) => {
  leavingNetworkId.value = networkId
  try {
    await LeaveZeroTierNetwork(networkId)
    message.success('已离开网络')
    await loadNetworks()
  } catch (error) {
    console.error('离开网络失败:', error)
    message.error('离开网络失败: ' + error)
  } finally {
    leavingNetworkId.value = null
  }
}
</script>

<style scoped>
.zerotier-view {
  max-width: 1000px;
}
</style>
