<template>
  <div style="background: white; color: black; padding: 20px;">
    <h1>测试页面</h1>
    <p>如果你能看到这个，说明 Vue 正常工作</p>
    <n-button type="primary" @click="handleClick">点击测试</n-button>
  </div>
</template>

<script setup>
import { useMessage } from 'naive-ui'

const message = useMessage()

function handleClick() {
  message.success('按钮点击成功！')
}
</script>
