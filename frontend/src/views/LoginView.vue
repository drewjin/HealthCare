<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const form = ref({ username: '', password: '' })
const router = useRouter()

const handleLogin = async () => {
  try {
    const response = await fetch('http://localhost:3000/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form.value)
    })
    
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || '登录失败');
    
    localStorage.setItem('jwt', data.token)
    localStorage.setItem('uid', data.uid)  

    // console.log('登录成功:', data.token, data.uid)
    router.push('/dashboard')  
  } catch (error) {
    console.error(error)
    let errorMessage = '登录失败: 未知错误';
    if (error instanceof Error) {
      errorMessage = '登录失败: ' + error.message;
    }
    alert(errorMessage);
  }
}
</script>

<template>
  <div class="login-container">
    <h2>用户登录</h2>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label>用户名</label>
        <input v-model="form.username" type="username" required>
      </div>
      <div class="form-group">
        <label>密码</label>
        <input v-model="form.password" type="password" required>
      </div>
      <button type="submit">登录</button>
    </form>
  </div>
</template>

<style scoped>
.login-container {
  max-width: 400px;
  margin: 2rem auto;
  padding: 2rem;
  border: 1px solid #ddd;
  border-radius: 8px;
}

.form-group {
  margin-bottom: 1rem;
}

input {
  width: 100%;
  padding: 0.5rem;
  margin-top: 0.3rem;
}

button {
  background: #007bff;
  color: white;
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
</style>