<template>
  <div class="zerotier-view">
    <n-space vertical :size="24">
      <n-card title="ZeroTier 状态">
        <n-space vertical>
          <n-descriptions :column="2" bordered>
            <n-descriptions-item label="状态">
              <n-tag :type="isOnline ? 'success' : 'error'">
                {{ isOnline ? '在线' : '离线' }}
              </n-tag>
            </n-descriptions-item>
            <n-descriptions-item label="节点 ID">
              {{ nodeId || '未知' }}
            </n-descriptions-item>
          </n-descriptions>
          <n-button @click="refreshStatus" :loading="loading">
            刷新状态
          </n-button>
        </n-space>
      </n-card>

      <n-card title="加入网络">
        <n-space vertical>
          <n-input
            v-model:value="networkId"
            placeholder="输入 ZeroTier 网络 ID (16位)"
            :maxlength="16"
          />
          <n-button
            type="primary"
            :loading="joining"
            :disabled="networkId.length !== 16"
            @click="joinNetwork"
          >
            加入网络
          </n-button>
        </n-space>
      </n-card>

      <n-card title="已加入的网络">
        <n-space vertical>
          <n-button @click="refreshNetworks" :loading="loadingNetworks">
            刷新列表
          </n-button>
          <n-list v-if="networks.length > 0" bordered>
            <n-list-item v-for="net in networks" :key="net.id">
              <n-thing :title="net.name || net.id">
                <template #description>
                  <n-space>
                    <n-tag size="small">{{ net.status }}</n-tag>
                    <span>{{ net.assignedAddresses?.join(', ') }}</span>
                  </n-space>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
          <n-empty v-else description="暂无网络" />
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
  NDescriptions,
  NDescriptionsItem,
  NTag,
  NButton,
  NInput,
  NList,
  NListItem,
  NThing,
  NEmpty,
  useMessage
} from 'naive-ui'
import {
  GetZeroTierInfo,
  JoinZeroTierNetwork,
  ListZeroTierNetworks
} from '../../wailsjs/go/main/App'

const message = useMessage()
const loading = ref(false)
const joining = ref(false)
const loadingNetworks = ref(false)
const isOnline = ref(false)
const nodeId = ref('')
const networkId = ref('')
const networks = ref([])

onMounted(async () => {
  await refreshStatus()
  await refreshNetworks()
})

async function refreshStatus() {
  loading.value = true
  try {
    const info = await GetZeroTierInfo()
    // 解析状态信息
    if (info.includes('ONLINE')) {
      isOnline.value = true
    }
    const match = info.match(/([0-9a-f]{10})/)
    if (match) {
      nodeId.value = match[1]
    }
  } catch (e) {
    message.error('获取状态失败: ' + e)
  } finally {
    loading.value = false
  }
}

async function joinNetwork() {
  joining.value = true
  try {
    await JoinZeroTierNetwork(networkId.value)
    message.success('成功加入网络！')
    networkId.value = ''
    await refreshNetworks()
  } catch (e) {
    message.error('加入失败: ' + e)
  } finally {
    joining.value = false
  }
}

async function refreshNetworks() {
  loadingNetworks.value = true
  try {
    const result = await ListZeroTierNetworks()
    networks.value = JSON.parse(result)
  } catch (e) {
    networks.value = []
  } finally {
    loadingNetworks.value = false
  }
}
</script>

