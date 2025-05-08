<template>
  <div class="family-requests">
    <div v-if="requests.length > 0" class="requests-container">
      <h3>待处理的家庭关系请求</h3>
      <div v-for="request in requests" :key="request.id" class="request-item">
        <p>用户 {{ request.requester }} 请求建立 {{ request.relationship }} 关系</p>
        <div class="actions">
          <button @click="handleRequest(request.id, true)" class="accept">接受</button>
          <button @click="handleRequest(request.id, false)" class="reject">拒绝</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import axios from 'axios'

const requests = ref<Array<{
  id: number
  requester: string
  relationship: string
  created_at: string
}>>([])

const userId = localStorage.getItem('userId')
const CHECK_INTERVAL = 30000 // 30秒检查一次
let checkTimer: number | undefined

const checkPendingRequests = async () => {
  try {
    const response = await axios.get(`/api/family/pending/${userId}`)
    requests.value = response.data
  } catch (error) {
    console.error('Failed to fetch pending requests:', error)
  }
}

const handleRequest = async (requestId: number, accept: boolean) => {
  try {
    await axios.post(`/api/family/handle/${userId}/${requestId}`, { accept })
    requests.value = requests.value.filter(req => req.id !== requestId)
  } catch (error) {
    console.error('Failed to handle request:', error)
  }
}

onMounted(() => {
  checkPendingRequests()
  checkTimer = window.setInterval(checkPendingRequests, CHECK_INTERVAL)
})

onUnmounted(() => {
  if (checkTimer) {
    clearInterval(checkTimer)
  }
})
</script>

<style scoped>
.family-requests {
  margin: 1rem;
}

.requests-container {
  max-width: 600px;
  margin: 0 auto;
}

.request-item {
  border: 1px solid #ddd;
  padding: 1rem;
  margin-bottom: 1rem;
  border-radius: 4px;
}

.actions {
  display: flex;
  gap: 1rem;
  margin-top: 0.5rem;
}

button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.accept {
  background-color: #4CAF50;
  color: white;
}

.reject {
  background-color: #f44336;
  color: white;
}

button:hover {
  opacity: 0.9;
}
</style>