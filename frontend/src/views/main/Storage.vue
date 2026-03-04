<template>
  <div class="storage-view">
    <n-space vertical :size="20">
      <!-- 可挂载的共享 -->
      <n-card title="可挂载的共享">
        <template #header-extra>
          <n-button size="small" :loading="loadingShares" @click="loadShares">
            刷新
          </n-button>
        </template>
        <n-spin :show="loadingShares">
          <n-empty v-if="shares.length === 0" description="暂无可用共享" />
          <n-list v-else bordered>
            <n-list-item v-for="(share, index) in shares" :key="index">
              <n-thing>
                <template #header>
                  <n-text strong>{{ share.name }}</n-text>
                </template>
                <template #description>
                  <n-text depth="3">{{ share.comment || '无描述' }}</n-text>
                </template>
                <template #action>
                  <n-button
                    type="primary"
                    size="small"
                    @click="selectShare(share)"
                  >
                    挂载
                  </n-button>
                </template>
              </n-thing>
            </n-list-item>
          </n-list>
        </n-spin>
      </n-card>

      <!-- 挂载配置对话框 -->
      <n-modal v-model:show="showMountDialog">
        <n-card
          style="width: 500px"
          title="挂载配置"
          :bordered="false"
          size="huge"
        >
          <n-form :model="mountForm" label-placement="left" label-width="100">
            <n-form-item label="远程路径">
              <n-input
                v-model:value="mountForm.remotePath"
                disabled
              />
            </n-form-item>
            <n-form-item label="本地挂载点">
              <n-input-group>
                <n-input
                  v-model:value="mountForm.localPath"
                  placeholder="/home/user/zima"
                  style="width: 70%"
                />
                <n-button
                  type="primary"
                  @click="selectMountDirectory"
                  style="width: 30%"
                >
                  选择文件夹
                </n-button>
              </n-input-group>
            </n-form-item>
            <n-form-item label="用户名">
              <n-input
                v-model:value="mountForm.username"
                placeholder="username"
              />
            </n-form-item>
            <n-form-item label="密码">
              <n-input
                v-model:value="mountForm.password"
                type="password"
                placeholder="password"
                show-password-on="click"
              />
            </n-form-item>
            <n-form-item label="自动挂载">
              <n-switch v-model:value="mountForm.autoMount" />
            </n-form-item>
          </n-form>

          <n-space justify="end">
            <n-button @click="showMountDialog = false">取消</n-button>
            <n-button
              type="primary"
              :loading="mounting"
              @click="handleMount"
            >
              立即挂载
            </n-button>
          </n-space>
        </n-card>
      </n-modal>

      <!-- 已挂载的共享 -->
      <n-card title="已挂载的共享">
        <template #header-extra>
          <n-button size="small" :loading="loadingMounts" @click="loadMounts">
            刷新
          </n-button>
        </template>
        <n-spin :show="loadingMounts">
          <n-empty v-if="mounts.length === 0" description="暂无挂载" />
          <n-list v-else bordered>
            <n-list-item v-for="(mount, index) in mounts" :key="index">
              <n-thing>
                <template #header>
                  <n-text strong>{{ mount.remotePath }}</n-text>
                </template>
                <template #description>
                  <n-space vertical size="small">
                    <n-text depth="3">挂载点: {{ mount.localPath }}</n-text>
                    <n-tag size="small" type="success">已挂载</n-tag>
                  </n-space>
                </template>
                <template #action>
                  <n-button
                    type="error"
                    size="small"
                    :loading="unmounting"
                    @click="handleUnmount(mount)"
                  >
                    卸载
                  </n-button>
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
import { ref, onMounted, inject } from 'vue'
import {
  NCard,
  NSpace,
  NForm,
  NFormItem,
  NInput,
  NInputGroup,
  NButton,
  NSwitch,
  NList,
  NListItem,
  NThing,
  NText,
  NTag,
  NEmpty,
  NModal,
  NSpin,
  useMessage
} from 'naive-ui'
import { MountSMB, SaveCredentials, ListSMBShares, ListMounts, UnmountSMB } from '../../../wailsjs/go/main/App'

const message = useMessage()
const checkSudo = inject('checkSudo')
const mounting = ref(false)
const unmounting = ref(false)
const loadingMounts = ref(false)
const saving = ref(false)
const mounts = ref([])
const selectedHost = ref('') // 从localStorage自动获取
const shares = ref([])
const loadingShares = ref(false)
const showMountDialog = ref(false)

const mountForm = ref({
  remotePath: '',
  localPath: '/home/user/zima',
  username: '',
  password: '',
  autoMount: false,
  deviceName: '',
  shareName: ''
})

const selectMountDirectory = async () => {
  try {
    const { SelectMountDirectory } = await import('../../../wailsjs/go/main/App')
    const selectedPath = await SelectMountDirectory()
    if (selectedPath) {
      mountForm.value.localPath = selectedPath
      message.success('已选择挂载目录')
    }
  } catch (error) {
    console.error('选择文件夹失败:', error)
    message.error('选择文件夹失败: ' + error)
  }
}

const loadShares = async () => {
  if (!selectedHost.value) {
    message.warning('未找到已连接的设备')
    return
  }

  loadingShares.value = true
  shares.value = []
  console.log('[Storage] 开始加载共享列表，主机:', selectedHost.value)

  // 获取用户凭证
  const username = localStorage.getItem('rnctl_username') || ''
  const deviceIp = localStorage.getItem('rnctl_device_ip') || selectedHost.value

  let password = ''
  try {
    const { GetCredentials } = await import('../../../wailsjs/go/main/App')
    const credData = await GetCredentials(deviceIp)
    if (credData) {
      const parsed = JSON.parse(credData)
      password = parsed.password || ''
    }
  } catch (error) {
    console.log('未找到保存的密码，尝试使用空密码')
  }

  const timeoutPromise = new Promise((_, reject) => {
    setTimeout(() => reject(new Error('加载超时')), 15000)
  })

  try {
    const result = await Promise.race([
      ListSMBShares(selectedHost.value, username, password),
      timeoutPromise
    ])
    console.log('[Storage] ListSMBShares返回:', result)
    shares.value = JSON.parse(result)
    console.log('[Storage] 解析后的共享列表:', shares.value)
    if (shares.value.length === 0) {
      message.warning('未找到共享，请检查设备连接')
    } else {
      message.success(`找到 ${shares.value.length} 个共享`)
    }
  } catch (error) {
    console.error('[Storage] 加载共享列表失败:', error)
    const errorMsg = String(error)
    if (errorMsg.includes('超时') || error.message === '加载超时') {
      message.error('连接超时，请检查设备连接')
    } else {
      message.error('加载共享列表失败: ' + error)
    }
    shares.value = []
  } finally {
    loadingShares.value = false
    console.log('[Storage] loadingShares设置为false')
  }
}

const selectShare = async (share) => {
  // 获取运行时挂载目录
  let runtimeDir = '/tmp/zimaclient'
  try {
    const { GetRuntimeMountDir } = await import('../../../wailsjs/go/main/App')
    runtimeDir = await GetRuntimeMountDir()
  } catch (error) {
    console.error('获取运行时目录失败:', error)
  }

  // 从localStorage获取用户名和密码
  const username = localStorage.getItem('rnctl_username') || ''
  const deviceIp = localStorage.getItem('rnctl_device_ip') || selectedHost.value
  const deviceName = localStorage.getItem('rnctl_device_hostname') || selectedHost.value

  // 尝试获取保存的密码
  let password = ''
  try {
    const { GetCredentials } = await import('../../../wailsjs/go/main/App')
    const credData = await GetCredentials(deviceIp)
    if (credData) {
      const parsed = JSON.parse(credData)
      password = parsed.password || ''
    }
  } catch (error) {
    console.log('未找到保存的密码')
  }

  // 构建挂载路径：/run/user/$(id -u)/zimaclient/<设备名-共享名>
  const mountLabel = `${deviceName}-${share.name}`
  const mountPath = `${runtimeDir}/${mountLabel}`

  mountForm.value.remotePath = `//${selectedHost.value}/${share.name}`
  mountForm.value.localPath = mountPath
  mountForm.value.shareName = share.name
  mountForm.value.username = username
  mountForm.value.password = password
  mountForm.value.deviceName = deviceName

  showMountDialog.value = true
}

const handleMount = async () => {
  if (!mountForm.value.remotePath) {
    message.warning('请输入远程路径')
    return
  }

  // 检查并请求sudo权限
  try {
    await checkSudo()
  } catch (error) {
    message.error('需要sudo权限才能挂载')
    return
  }

  mounting.value = true
  try {
    await MountSMB(
      mountForm.value.remotePath,
      mountForm.value.localPath,
      mountForm.value.username,
      mountForm.value.password,
      mountForm.value.deviceName || '',
      mountForm.value.shareName || ''
    )
    message.success('挂载成功，已添加到文件管理器侧边栏')

    // 重新加载挂载列表
    await loadMounts()

    // 保存自动挂载配置
    if (mountForm.value.autoMount) {
      await saveAutoMount()
    }

    showMountDialog.value = false
  } catch (error) {
    message.error('挂载失败: ' + error)
  } finally {
    mounting.value = false
  }
}

const saveAutoMount = async () => {
  saving.value = true
  try {
    const config = JSON.stringify({
      remotePath: mountForm.value.remotePath,
      localPath: mountForm.value.localPath,
      username: mountForm.value.username,
      password: mountForm.value.password
    })

    await SaveCredentials('auto_mount_config', config)
    message.success('自动挂载配置已保存')
  } catch (error) {
    message.error('保存失败: ' + error)
  } finally {
    saving.value = false
  }
}

const loadMounts = async () => {
  loadingMounts.value = true
  console.log('[Storage] 开始加载挂载列表...')

  const timeoutPromise = new Promise((_, reject) => {
    setTimeout(() => reject(new Error('加载超时')), 5000)
  })

  try {
    const result = await Promise.race([
      ListMounts(),
      timeoutPromise
    ])
    console.log('[Storage] ListMounts返回:', result)
    const parsed = JSON.parse(result)
    console.log('[Storage] 解析后的挂载列表:', parsed)
    mounts.value = Array.isArray(parsed) ? parsed : []
    console.log('[Storage] 挂载列表加载成功，数量:', mounts.value.length)
  } catch (error) {
    console.error('[Storage] 加载挂载列表失败:', error)
    if (error.message === '加载超时') {
      message.error('加载挂载列表超时，请检查系统状态')
    } else {
      message.error('加载挂载列表失败: ' + error)
    }
    mounts.value = []
  } finally {
    loadingMounts.value = false
    console.log('[Storage] loadingMounts设置为false')
  }
}

const handleUnmount = async (mount) => {
  try {
    await checkSudo()
  } catch (error) {
    message.error('需要sudo权限才能卸载')
    return
  }

  unmounting.value = true
  try {
    await UnmountSMB(mount.localPath, false)
    message.success('卸载成功，已从文件管理器侧边栏移除')
    await loadMounts()
  } catch (error) {
    message.error('卸载失败: ' + error)
  } finally {
    unmounting.value = false
  }
}

onMounted(async () => {
  // 从localStorage获取当前连接的设备IP
  selectedHost.value = localStorage.getItem('rnctl_device_ip') || ''

  if (!selectedHost.value) {
    message.warning('未找到已连接的设备，请先连接设备')
    return
  }

  // 获取运行时挂载目录并设置默认挂载路径
  try {
    const { GetRuntimeMountDir } = await import('../../../wailsjs/go/main/App')
    const runtimeDir = await GetRuntimeMountDir()
    const deviceName = localStorage.getItem('rnctl_device_hostname') || 'device'
    mountForm.value.localPath = `${runtimeDir}/${deviceName}-share`
  } catch (error) {
    console.error('获取运行时目录失败:', error)
    mountForm.value.localPath = '/tmp/zimaclient/share'
  }

  // 自动加载共享列表
  await loadShares()

  // 加载已挂载列表
  await loadMounts()
})
</script>

<style scoped>
.storage-view {
  padding: 20px;
}
</style>
