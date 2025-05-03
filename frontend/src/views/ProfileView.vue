<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const userInfo = ref({ 
  username: '', 
  name: '', 
  email: '', 
  phone: '', 
  gender: '', 
  birthday: '', 
  address: '' 
})
const router = useRouter()

onMounted(async () => {
  try {
    const token = localStorage.getItem('jwt')
    const uid = String(localStorage.getItem('uid'))
    if (!token) throw new Error('未登录')
    // console.log(uid, '\n', token)
    const response = await fetch(`http://localhost:3000/api/users/${uid}/profile`, {
      method: 'GET',
      headers: { 'Authorization': token },
    })
    if (!response.ok) {
      const errorData = await response.json();
      throw new Error(errorData.error || '获取用户信息失败');
    }
    userInfo.value = (await response.json()).data
  } catch (error) {
    console.error(error)
    router.push('/dashboard')
  }
})

const handleLogout = () => {
  localStorage.removeItem('jwt')
  router.push('/login')
}
</script>

<template>
  <div class="profile-container">
    <h2>个人信息</h2>
    <div class="user-info-card">
      <p class="info-item">用户名：{{ userInfo.username }}</p>
      <p class="info-item">姓名：{{ userInfo.name }}</p>
      <p class="info-item">邮箱：{{ userInfo.email }}</p>
      <p class="info-item">电话：{{ userInfo.phone }}</p>
      <p class="info-item">性别：{{ userInfo.gender }}</p>
      <p class="info-item">生日：{{ userInfo.birthday }}</p>
      <p class="info-item">地址：{{ userInfo.address }}</p>
    </div>
    <div class="button-group">
  <button class="nav-btn" @click="$router.push('/dashboard')">返回主页</button>
  <button class="logout-btn" @click="handleLogout">退出登录</button>
</div>
  </div>
</template>

<style scoped>
.info-item {
  color: #333;
  background: #f8f9fa;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.user-info-card {
  background: white;
  box-shadow: 0 4px 6px rgba(0,0,0,0.1);
}

.logout-btn {
  background: #dc3545;
  color: white;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.nav-btn {
  background: #28a745;
  color: white;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  margin-right: 1rem;
}

.button-group {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}
</style>