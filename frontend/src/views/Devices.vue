<template>
  <div class="devices-view">
    <n-space vertical :size="20">
      <!-- ZeroTier 状态 -->
      <n-card title="ZeroTier 状态">
        <n-space vertical>
          <n-space align="center">
            <n-button
              type="primary"
              size="small"
              @click="loadZTStatus"
              :loading="loadingZT"
            >
              刷新状态
            </n-button>
            <n-tag v-if="ztStatus" :type="ztStatus.online ? 'success' : 'error'">
              {{ ztStatus.online ? '在线' : '离线' }}
            </n-tag>
          </n-space>

          <n-descriptions v-if="ztStatus" bordered :column="2" size="small">
            <n-descriptions-item label="节点ID">
              {{ ztStatus.address || 'N/A' }}
            </n-descriptions-item>
            <n-descriptions-item label="版本">
              {{ ztStatus.version || 'N/A' }}
            </n-descriptions-item>
            <n-descriptions-item label="系统时间">
              {{ formatTime(ztStatus.clock) }}
            </n-descriptions-item>
          </n-descriptions>

          <n-divider v-if="ztNetworks.length > 0">已加入的网络</n-divider>

          <n-list v-if="ztNetworks.length > 0" bordered size="small">
            <n-list-item v-for="(network, idx) in ztNetworks" :key="idx">
              <n-thing>
                <template #header>
                  <n-text strong>{{ network.name || network.nwid }}</n-text>
                </template>
                <template #description>
                  <n-space vertical size="small">
                    <n-text depth="3">网络ID: {{ network.nwid }}</n-text>
                    <n-text depth="3">状态: {{ network.status }}</n-text>
                    <n-text v-if="network.assignedAddresses" depth="3">
                      IP: {{ network.assignedAddresses.join(', ') }}
                    </n-text>
                  </n-space>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-space>
      </n-card>

      <!-- 设备发现 -->
      <n-card title="设备发现">
        <n-space vertical>
          <n-alert type="info">
            扫描局域网内的 Zima 设备
          </n-alert>
          <n-space>
            <n-button
              type="primary"
              :loading="scanning"
              @click="startScan"
            >
              {{ scanning ? '扫描中...' : '开始扫描' }}
            </n-button>
            <n-text v-if="devices.length > 0" type="success">
              找到 {{ devices.length }} 台设备
            </n-text>
          </n-space>

          <n-divider v-if="devices.length > 0" />

          <n-list v-if="devices.length > 0" bordered>
            <n-list-item v-for="(device, index) in devices" :key="index">
              <n-thing>
                <template #header>
                  <n-text strong>{{ device.device_name || '未知设备' }}</n-text>
                </template>
                <template #description>
                  <n-space vertical size="small">
                    <n-text depth="3">型号: {{ device.device_model || 'N/A' }}</n-text>
                    <n-text depth="3">系统: {{ device.os_version || 'N/A' }}</n-text>
                    <n-text depth="3">IP: {{ device.ActualIP || device.request_ip || 'N/A' }}</n-text>
                    <n-text depth="3">端口: {{ device.port || 'N/A' }}</n-text>
                    <n-tag :type="device.initialized ? 'success' : 'warning'\" size="small">
                      {{ device.initialized ? '已初始化' : '未初始化' }}
                    </n-tag>
                  </n-space>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>

          <n-empty
            v-if="!scanning && devices.length === 0 && scanned"
            description="未发现设备"
          />
        </n-space>
      </n-card>
    </n-space>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  NCard, NSpace, NAlert, NButton, NText, NDivider, NList, NListItem,
  NThing, NTag, NEmpty, NDescriptions, NDescriptionsItem, useMessage
} from 'naive-ui'
import { ScanLocalDevices, GetZeroTierInfo, ListZeroTierNetworks } from '../../wailsjs/go/main/App'

const message = useMessage()
const scanning = ref(false)
const devices = ref([])
const scanned = ref(false)
const loadingZT = ref(false)
const ztStatus = ref(null)
const ztNetworks = ref([])

const formatTime = (timestamp) => {
  if (!timestamp) return 'N/A'
  return new Date(timestamp).toLocaleString('zh-CN')
}

const loadZTStatus = async () => {
  loadingZT.value = true
  try {
    const info = await GetZeroTierInfo()
    ztStatus.value = JSON.parse(info)

    const networks = await ListZeroTierNetworks()
    ztNetworks.value = JSON.parse(networks)
  } catch (error) {
    message.error('加载 ZeroTier 状态失败: ' + error)
  } finally {
    loadingZT.value = false
  }
}

const startScan = async () => {
  scanning.value = true
  devices.value = []
  scanned.value = false

  try {
    const result = await ScanLocalDevices()
    const parsed = JSON.parse(result)
    devices.value = parsed || []
    scanned.value = true

    if (devices.value.length === 0) {
      message.warning('未发现设备，请确保设备在同一局域网内')
    } else {
      message.success(`发现 ${devices.value.length} 台设备`)
    }
  } catch (error) {
    message.error('扫描失败: ' + error)
    scanned.value = true
  } finally {
    scanning.value = false
  }
}

onMounted(() => {
  loadZTStatus()
})
</script>

<style scoped>
.devices-view {
  padding: 20px;
}
</style>
