<template>
  <div class="connect-page">
    <div class="content">
      <!-- Header with refresh button -->
      <div class="page-header">
        <h1 v-if="devices.length === 0 && scanned">未找到设备。</h1>
        <h1 v-else-if="scanning">正在扫描...</h1>
        <h1 v-else-if="devices.length > 0">找到 {{ devices.length }} 台设备</h1>
        <h1 v-else>连接到设备</h1>
        <n-button
          quaternary
          size="large"
          :loading="scanning"
          @click="scanDevices"
        >
          刷新
        </n-button>
      </div>

      <!-- Help text when no devices found -->
      <div v-if="devices.length === 0 && scanned" class="help-text">
        <p>请确保以下事项：</p>
        <ol>
          <li>设备已开启。</li>
          <li>网络电缆已牢固连接，指示灯正在闪烁。</li>
          <li>设备和您的计算机接到同一路由器。</li>
        </ol>
        <p class="hint">
          如果您有网络ID，可以尝试
          <a href="#" @click.prevent="showNetworkIdInput = true">使用网络ID连接</a>。
        </p>
      </div>

      <!-- Device list -->
      <div v-if="devices.length > 0" class="device-list">
        <div
          v-for="device in devices"
          :key="device.device_id || device.actual_ip"
          class="device-card"
          @click="selectDevice(device)"
        >
          <div class="device-icon">
            <n-icon size="32" :component="Server" />
          </div>
          <div class="device-info">
            <div class="device-name">{{ device.device_name || device.actual_ip }}</div>
            <div class="device-ips">
              <n-space :size="4">
                <n-tag v-for="ip in device.lan_ips" :key="ip" size="small" type="success">
                  LAN: {{ ip }}
                </n-tag>
                <n-tag v-for="ip in device.zerotier_ips" :key="ip" size="small" type="info">
                  ZT: {{ ip }}
                </n-tag>
              </n-space>
            </div>
            <div v-if="device.device_model" class="device-model">{{ device.device_model }}</div>
          </div>
          <div class="device-arrow">
            <n-icon size="24" :component="ChevronForward" />
          </div>
        </div>
      </div>

      <!-- Network ID input (collapsed by default) -->
      <div v-if="showNetworkIdInput" class="network-id-section">
        <n-card title="使用网络ID连接" size="small">
          <n-space vertical :size="12">
            <n-alert type="info" size="small" :show-icon="false">
              输入 ZeroTier 网络 ID，系统将自动加入网络并发现设备
            </n-alert>
            <n-input
              v-model:value="ztNetworkId"
              placeholder="输入 16 位网络 ID"
              size="large"
              maxlength="16"
            />
            <n-space>
              <n-button
                type="primary"
                :loading="joiningZt"
                @click="connectViaZeroTier"
              >
                连接
              </n-button>
              <n-button @click="showNetworkIdInput = false">
                取消
              </n-button>
            </n-space>
          </n-space>
        </n-card>
      </div>

      <!-- Manual IP input link (bottom left) -->
      <div class="bottom-actions">
        <a
          v-if="!showManualInput"
          href="#"
          class="link-button"
          @click.prevent="showManualInput = true"
        >
          使用网络ID连接
        </a>
        <a
          v-if="!showManualIpInput"
          href="#"
          class="link-button"
          @click.prevent="showManualIpInput = true"
        >
          手动输入IP地址
        </a>
      </div>

      <!-- Manual IP input (collapsed) -->
      <div v-if="showManualIpInput" class="manual-ip-section">
        <n-card title="手动连接" size="small">
          <n-space vertical :size="12">
            <n-input
              v-model:value="manualIp"
              placeholder="例如: 192.168.1.100"
              size="large"
            />
            <n-space>
              <n-button
                type="primary"
                @click="connectManual"
              >
                连接
              </n-button>
              <n-button @click="showManualIpInput = false">
                取消
              </n-button>
            </n-space>
          </n-space>
        </n-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import { Server, ChevronForward } from '@vicons/ionicons5'
import { ScanLocalDevices, JoinZeroTierNetwork, ListZeroTierNetworks } from '../../../wailsjs/go/main/App'

const router = useRouter()
const message = useMessage()

const scanning = ref(false)
const scanned = ref(false)
const devices = ref([])
const manualIp = ref('')
const ztNetworkId = ref('')
const joiningZt = ref(false)
const showNetworkIdInput = ref(false)
const showManualInput = ref(false)
const showManualIpInput = ref(false)

// 首次进入自动扫描
onMounted(() => {
  scanDevices()
})

async function scanDevices() {
  scanning.value = true
  scanned.value = false
  devices.value = []

  try {
    const result = await ScanLocalDevices()
    console.log('扫描结果:', result)

    // 后端返回的是数组，不是对象
    const parsed = JSON.parse(result)
    devices.value = Array.isArray(parsed) ? parsed : []
    scanned.value = true

    if (devices.value.length === 0) {
      message.warning('未发现设备')
    } else {
      message.success(`发现 ${devices.value.length} 个设备`)
    }
  } catch (error) {
    console.error('扫描错误:', error)
    message.error('扫描失败: ' + error)
    scanned.value = true
  } finally {
    scanning.value = false
  }
}

function selectDevice(device) {
  // 存储完整的设备信息（用于智能重连）
  const deviceInfo = {
    device_id: device.device_id,
    device_name: device.device_name,
    lan_ips: device.lan_ips || [],
    zerotier_ips: device.zerotier_ips || [],
    preferred_ip: device.preferred_ip || device.actual_ip,
    all_ips: device.all_ips || []
  }

  localStorage.setItem('rnctl_device_info', JSON.stringify(deviceInfo))

  // 兼容旧代码，仍然存储单一IP
  localStorage.setItem('rnctl_device_ip', device.preferred_ip || device.actual_ip)
  localStorage.setItem('rnctl_device_hostname', device.device_name || device.preferred_ip || device.actual_ip)

  router.push('/wizard/login')
}

function connectManual() {
  if (!manualIp.value) {
    message.warning('请输入设备 IP')
    return
  }

  localStorage.setItem('rnctl_device_ip', manualIp.value)
  localStorage.setItem('rnctl_device_hostname', manualIp.value)
  router.push('/wizard/login')
}

async function connectViaZeroTier() {
  if (!ztNetworkId.value || ztNetworkId.value.length !== 16) {
    message.warning('请输入有效的 16 位网络 ID')
    return
  }

  joiningZt.value = true

  try {
    // 加入 ZeroTier 网络
    await JoinZeroTierNetwork(ztNetworkId.value)
    message.success('已加入网络，正在查找设备...')

    // 等待网络连接建立
    await new Promise(resolve => setTimeout(resolve, 3000))

    // 获取网络信息
    const networksResult = await ListZeroTierNetworks()
    const networks = JSON.parse(networksResult)

    // 查找刚加入的网络
    const network = networks.find(n => n.nwid === ztNetworkId.value)

    if (network && network.assignedAddresses && network.assignedAddresses.length > 0) {
      // 扫描 ZeroTier 网络中的设备
      const scanResult = await ScanLocalDevices()
      const devices = JSON.parse(scanResult)

      if (devices.length > 0) {
        // 自动选择第一个设备
        selectDevice(devices[0])
      } else {
        message.warning('未在网络中发现设备，请手动输入 IP')
      }
    } else {
      message.warning('网络连接中，请稍后手动输入设备 IP')
    }
  } catch (error) {
    console.error('ZeroTier 连接失败:', error)
    message.error('连接失败: ' + error)
  } finally {
    joiningZt.value = false
  }
}
</script>

<style scoped>
.connect-page {
  width: 100%;
  min-height: 100vh;
  padding: 40px;
  background: var(--bg-primary);
  overflow-y: auto;
}

.content {
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  gap: 32px;
}

/* Page Header */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.page-header h1 {
  font-size: 36px;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
}

/* Help Text */
.help-text {
  color: var(--text-secondary);
  font-size: 15px;
  line-height: 1.8;
}

.help-text p {
  margin: 0 0 16px 0;
  color: var(--text-secondary);
}

.help-text ol {
  margin: 0 0 24px 0;
  padding-left: 24px;
  color: var(--text-secondary);
}

.help-text li {
  margin-bottom: 8px;
  color: var(--text-secondary);
}

.help-text .hint {
  margin-top: 24px;
  font-size: 14px;
  color: var(--text-secondary);
}

.help-text a {
  color: var(--primary-500);
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s;
}

.help-text a:hover {
  border-bottom-color: var(--primary-500);
}

/* Device List */
.device-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.device-card {
  background: rgba(245, 245, 245, 0.8);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: all 0.2s ease-in-out;
  will-change: transform;
}

@media (prefers-color-scheme: dark) {
  .device-card {
    background: rgba(23, 23, 23, 0.8);
    border-color: rgba(64, 64, 64, 0.5);
  }

  .device-card:hover {
    background: rgba(38, 38, 38, 0.9);
    border-color: var(--primary-500);
  }
}

.device-card:hover {
  background: rgba(229, 229, 229, 0.9);
  border-color: var(--primary-500);
  transform: translateX(4px);
}

.device-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--primary-500) 0%, var(--primary-400) 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.device-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.device-name {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary);
}

.device-ip {
  font-size: 14px;
  color: var(--text-secondary);
}

.device-model {
  font-size: 13px;
  color: var(--text-tertiary);
}

.device-arrow {
  color: var(--text-tertiary);
  transition: all 0.2s ease-in-out;
}

.device-card:hover .device-arrow {
  color: var(--primary-500);
  transform: translateX(4px);
}

/* Bottom Actions */
.bottom-actions {
  display: flex;
  gap: 24px;
  margin-top: 20px;
}

.link-button {
  color: var(--primary-500);
  text-decoration: none;
  font-size: 14px;
  border-bottom: 1px solid transparent;
  transition: border-color 0.2s;
}

.link-button:hover {
  border-bottom-color: var(--primary-500);
}

/* Collapsed sections */
.network-id-section,
.manual-ip-section {
  margin-top: 20px;
}
</style>
