<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const isLogin = ref(true)
const router = useRouter()

const loginForm = ref({
  username: '',
  password: ''
})

const registerForm = ref({
  username: '',
  password: '',
  confirmPassword: '',
  name: '',
  gender: 'M',
  birthday: '',
  phone: '',
  email: '',
  address: '',
  user_type: 1, // Default as normal user
})

// Add userTypeOptions
const userTypeOptions = [
  { value: 1, label: '普通用户' },
  { value: 2, label: '管理员用户' },
  { value: 3, label: '机构用户' },
]

const handleLogin = async () => {
  try {
    const response = await fetch('http://localhost:3000/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(loginForm.value)
    })
    
    const data = await response.json();
    if (!response.ok) throw new Error(data.error || '登录失败');
    
    localStorage.setItem('jwt', data.token)
    localStorage.setItem('uid', data.uid)
    
    // Get user type and store it
    try {
      const userResponse = await fetch(`http://localhost:3000/api/users/${data.uid}/profile`, {
        method: 'GET',
        headers: { Authorization: data.token },
      })
      if (userResponse.ok) {
        const userData = await userResponse.json()
        if (userData && userData.data && userData.data.user_type) {
          localStorage.setItem('userType', userData.data.user_type.toString())
        }
      }
    } catch (error) {
      console.error('Failed to get user type', error)
    }

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

const handleRegister = async () => {
  try {
    if (!validateRegisterForm()) {
      return;
    }

    // 确保所有必需字段都存在且格式正确
    const registrationData = {
      ...registerForm.value,
      gender: registerForm.value.gender || 'M'  // 确保默认值
    };

    console.log('Registration data:', registrationData);  // 添加日志以便调试

    const response = await fetch('http://localhost:3000/api/auth/register', {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(registrationData)
    });
    
    const data = await response.json();
    if (!response.ok) {
      throw new Error(data.error || '注册失败');
    }
    
    localStorage.setItem('jwt', data.token);
    alert('注册成功！请登录');
    isLogin.value = true;
    loginForm.value = { username: '', password: '' };
  } catch (error) {
    console.error('Registration error:', error);
    let errorMessage = '注册失败: 未知错误';
    if (error instanceof Error) {
      errorMessage = '注册失败: ' + error.message;
    }
    alert(errorMessage);
  }
}

const validateRegisterForm = () => {
  if (registerForm.value.password !== registerForm.value.confirmPassword) {
    alert('两次输入的密码不一致');
    return false;
  }
  if (!registerForm.value.phone.match(/^\d{11}$/)) {
    alert('请输入11位手机号码');
    return false;
  }
  if (registerForm.value.email && !registerForm.value.email.includes('@')) {
    alert('请输入有效的邮箱地址');
    return false;
  }
  if (!['M', 'F'].includes(registerForm.value.gender)) {
    alert('性别必须为 M(男) 或 F(女)');
    return false;
  }
  return true;
}

const toggleMode = () => {
  isLogin.value = !isLogin.value
  loginForm.value = { username: '', password: '' }
  registerForm.value = {
    username: '',
    password: '',
    confirmPassword: '',
    name: '',
    gender: 'M',
    birthday: '',
    phone: '',
    email: '',
    address: '',
    user_type: 1, // Reset to default user type
  }
}
</script>

<template>
  <div class="page-container">
    <h1 class="login-title">欢迎来到健康管理系统</h1>
    <div class="login-container">
      <h2>{{ isLogin ? '用户登录' : '用户注册' }}</h2>
      
      <form v-if="isLogin" @submit.prevent="handleLogin">
        <div class="form-group">
          <label>用户名</label>
          <input v-model="loginForm.username" type="text" required>
        </div>
        <div class="form-group">
          <label>密码</label>
          <input v-model="loginForm.password" type="password" required>
        </div>
        <button type="submit">登录</button>
      </form>

      <form v-else @submit.prevent="handleRegister" class="register-form">
        <div class="form-row">
          <div class="form-group">
            <label>用户名</label>
            <input v-model="registerForm.username" type="text" required>
          </div>
          <div class="form-group">
            <label>密码</label>
            <input v-model="registerForm.password" type="password" required>
          </div>
          <div class="form-group">
            <label>确认密码</label>
            <input v-model="registerForm.confirmPassword" type="password" required>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>姓名</label>
            <input v-model="registerForm.name" type="text" required>
          </div>
          <div class="form-group">
            <label>性别</label>
            <select v-model="registerForm.gender" required>
              <option value="">请选择性别</option>
              <option value="M">男</option>
              <option value="F">女</option>
            </select>
          </div>
          <div class="form-group">
            <label>生日</label>
            <input v-model="registerForm.birthday" type="date" required>
          </div>
        </div>

        <div class="form-row">
          <div class="form-group">
            <label>手机号码</label>
            <input v-model="registerForm.phone" type="tel" required maxlength="11" 
                   pattern="[0-9]{11}" placeholder="请输入11位手机号">
          </div>
          <div class="form-group">
            <label>邮箱 (选填)</label>
            <input v-model="registerForm.email" type="email" placeholder="请输入邮箱地址">
          </div>
          <div class="form-group">
            <label>用户类型</label>
            <select v-model="registerForm.user_type" required>
              <option v-for="option in userTypeOptions" :key="option.value" :value="option.value">
                {{ option.label }}
              </option>
            </select>
          </div>
        </div>

        <div class="form-row single-column">
          <div class="form-group">
            <label>地址 (选填)</label>
            <input v-model="registerForm.address" type="text" placeholder="请输入地址">
          </div>
        </div>

        <button type="submit">注册</button>
      </form>

      <button type="button" class="toggle-button" @click="toggleMode">
        {{ isLogin ? '没有账号？去注册' : '已有账号？去登录' }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.page-container {
  min-height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  box-sizing: border-box;
  position: absolute;
  top: 0;
  left: 0;
}

.login-title {
  text-align: center;
  font-size: 2.5rem;
  color: #2c3e50;
  margin-bottom: 2rem;
}

.login-container {
  width: 100%;
  max-width: 500px;
  margin: 0 auto;
  padding: 2.5rem;
  border: 1px solid #ccc;
  border-radius: 12px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  background-color: #ffffff;
}

h2 {
  text-align: center;
  color: #2c3e50;
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

input {
  width: 100%;
  padding: 0.8rem;
  margin-top: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
}

button {
  width: 100%;
  background: #007bff;
  color: white;
  padding: 0.8rem;
  border: none;
  border-radius: 6px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.3s ease;
}

button:hover {
  background: #0056b3;
}

.toggle-button {
  margin-top: 1rem;
  background: #6c757d;
}

.toggle-button:hover {
  background: #5a6268;
}

.form-group input {
  width: 100%;
  padding: 0.8rem;
  margin-top: 0.5rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 1rem;
  box-sizing: border-box;
}

.form-group input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0,123,255,0.25);
}

.form-group label {
  display: block;
  margin-bottom: 0.3rem;
  color: #495057;
}

.register-form {
  max-width: 1000px;
  margin: 0 auto;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-row.single-column {
  grid-template-columns: 1fr;
}

.form-group {
  margin-bottom: 0.8rem;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 0.6rem;
  margin-top: 0.3rem;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 0.9rem;
  box-sizing: border-box;
}

.form-group select {
  background-color: white;
  height: 2.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.2rem;
  color: #495057;
  font-weight: 500;
  font-size: 0.9rem;
}

@media (max-width: 900px) {
  .form-row {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .login-container {
    padding: 1.5rem;
  }
}
</style>