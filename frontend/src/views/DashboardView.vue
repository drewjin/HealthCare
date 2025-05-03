<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const uid = ref(localStorage.getItem('uid'))

const handleLogout = () => {
  localStorage.removeItem('jwt')
  localStorage.removeItem('uid')
  router.push('/login')
}
</script>

<template>
  <div class="dashboard-container">
    <h2>健康管理控制台</h2>
    <div class="user-info">
      <p>用户ID：{{ uid }}</p>
    </div>
    <div class="actions">
      <button class="nav-btn" @click="$router.push('/profile')">个人资料</button>
      <button class="nav-btn" @click="$router.push('/records')">健康记录</button>
      <button class="logout-btn" @click="handleLogout">退出登录</button>
    </div>
  </div>
</template>

<style scoped>
.dashboard-container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.user-info {
  margin: 2rem 0;
  padding: 1rem;
  background: rgba(248, 249, 250, 0.9);
  border-radius: 8px;
  color: #2c3e50;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.actions {
  display: grid;
  gap: 1rem;
  grid-template-columns: repeat(2, 1fr);
}

.nav-btn {
  background: #28a745;
  color: white;
  padding: 1rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.nav-btn:hover {
  background: #218838;
  transform: translateY(-2px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.logout-btn {
  grid-column: span 2;
  background: #dc3545;
  padding: 1rem;
  transition: all 0.3s ease;
}

.logout-btn:hover {
  background: #c82333;
  transform: translateY(-2px);
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}
</style>