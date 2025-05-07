<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElContainer, ElAside, ElMain, ElHeader, ElMenu, ElMenuItem, ElCard, ElAvatar, ElRow, ElCol, ElButton, ElDivider, ElForm, ElFormItem, ElInput, ElMessage, ElDialog } from 'element-plus'
import { User, List, Setting } from '@element-plus/icons-vue'

const router = useRouter()
const uid = ref(localStorage.getItem('uid'))
const activeMenu = ref('1')
const userInfo = ref({
  username: '',
  name: '',
  email: '',
  phone: '',
  gender: '',
  birthday: '',
  address: '',
})

const showResetPwdDialog = ref(false)
const resetPwdForm = ref({
  prevPassword: '',
  newPassword: '',
  newPasswordConfirm: ''
})

onMounted(async () => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('未登录')

    const response = await fetch(`http://localhost:3000/api/users/${uid.value}/profile`, {
      method: 'GET',
      headers: { Authorization: token },
    })
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || '获取用户信息失败')
    }
    userInfo.value = (await response.json()).data
  } catch (error) {
    console.error(error)
    router.push('/login')
  }
})

const handleLogout = () => {
  localStorage.removeItem('jwt')
  localStorage.removeItem('uid')
  router.push('/login')
}

const handleGoToProfile = () => {
  activeMenu.value = '1'
}

const handleResetPwd = async () => {
  try {
    if (resetPwdForm.value.newPassword !== resetPwdForm.value.newPasswordConfirm) {
      ElMessage.error('新密码不一致')
      return
    }

    const response = await fetch(`http://localhost:3000/api/users/${uid.value}/reset_pwd`, {
      method: 'PUT',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': localStorage.getItem('jwt') || ''
      },
      body: JSON.stringify({
        prev_password: resetPwdForm.value.prevPassword,
        new_password: resetPwdForm.value.newPassword,
        new_password_confirm: resetPwdForm.value.newPasswordConfirm
      })
    })

    const data = await response.json()
    if (!response.ok) throw new Error(data.error)

    ElMessage.success('密码修改成功，请重新登录')
    handleLogout()
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '密码修改失败')
  }
}
</script>

<template>
  <ElContainer class="dashboard-container">
    <ElHeader class="dashboard-header">
      <div class="header-content">
        <div class="header-left">
          <h2>健康管理系统</h2>
          <ElButton type="primary" @click="handleGoToProfile" plain>
            <el-icon><User /></el-icon>
            个人信息
          </ElButton>
        </div>
        <ElButton type="danger" @click="handleLogout">退出登录</ElButton>
      </div>
    </ElHeader>
    
    <ElContainer>
      <ElAside width="200px">
        <ElMenu
          default-active="1"
          class="dashboard-menu"
          @select="activeMenu = $event"
        >
          <ElMenuItem index="1">
            <el-icon><User /></el-icon>
            <span>个人资料</span>
          </ElMenuItem>
          <ElMenuItem index="2">
            <el-icon><List /></el-icon>
            <span>健康记录</span>
          </ElMenuItem>
          <ElMenuItem index="3">
            <el-icon><Setting /></el-icon>
            <span>修改密码</span>
          </ElMenuItem>
        </ElMenu>
      </ElAside>
      
      <ElMain>
        <!-- 个人资料面板 -->
        <div v-if="activeMenu === '1'">
          <ElRow :gutter="20">
            <ElCol :span="8">
              <ElCard shadow="hover">
                <div class="user-profile-card">
                  <ElAvatar :size="100" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
                  <h3>{{ userInfo.name }}</h3>
                  <p>ID: {{ uid }}</p>
                  <ElButton type="primary" @click="showResetPwdDialog = true">修改密码</ElButton>
                </div>
              </ElCard>
            </ElCol>
            <ElCol :span="16">
              <ElCard shadow="hover">
                <template #header>
                  <div class="card-header">
                    <span>基本信息</span>
                  </div>
                </template>
                <div class="user-info-list">
                  <div class="info-item">
                    <span class="label">用户名：</span>
                    <span>{{ userInfo.username }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">邮箱：</span>
                    <span>{{ userInfo.email }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">电话：</span>
                    <span>{{ userInfo.phone }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">性别：</span>
                    <span>{{ userInfo.gender === 'M' ? '男' : '女' }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">生日：</span>
                    <span>{{ userInfo.birthday }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">地址：</span>
                    <span>{{ userInfo.address }}</span>
                  </div>
                </div>
              </ElCard>
            </ElCol>
          </ElRow>

          <!-- 修改密码对话框 -->
          <ElDialog
            v-model="showResetPwdDialog"
            title="修改密码"
            width="30%"
            :close-on-click-modal="false"
          >
            <ElForm :model="resetPwdForm" label-width="100px">
              <ElFormItem label="当前密码">
                <ElInput v-model="resetPwdForm.prevPassword" type="password" show-password />
              </ElFormItem>
              <ElFormItem label="新密码">
                <ElInput v-model="resetPwdForm.newPassword" type="password" show-password />
              </ElFormItem>
              <ElFormItem label="确认新密码">
                <ElInput v-model="resetPwdForm.newPasswordConfirm" type="password" show-password />
              </ElFormItem>
            </ElForm>
            <template #footer>
              <span class="dialog-footer">
                <ElButton @click="showResetPwdDialog = false">取消</ElButton>
                <ElButton type="primary" @click="handleResetPwd">
                  确认修改
                </ElButton>
              </span>
            </template>
          </ElDialog>
        </div>

        <!-- 健康记录面板 -->
        <div v-if="activeMenu === '2'">
          <ElCard shadow="hover">
            <template #header>
              <div class="card-header">
                <span>健康记录</span>
                <ElButton type="primary">添加记录</ElButton>
              </div>
            </template>
            <p>健康记录内容将在此显示</p>
          </ElCard>
        </div>

        <!-- 系统设置面板 -->
        <div v-if="activeMenu === '3'">
          <ElCard shadow="hover">
            <template #header>
              <div class="card-header">
                <span>系统设置</span>
              </div>
            </template>
            <p>系统设置内容将在此显示</p>
          </ElCard>
        </div>
      </ElMain>
    </ElContainer>
  </ElContainer>
</template>

<style scoped>
.dashboard-container {
  height: 100vh;
}

.dashboard-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.dashboard-menu {
  height: 100%;
  border-right: none;
}

.user-profile-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 20px 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info-list {
  padding: 0 20px;
}

.info-item {
  display: flex;
  padding: 10px 0;
}

.label {
  font-weight: bold;
  width: 80px;
  color: #606266;
}

.el-divider {
  margin: 8px 0;
}

.settings-container {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>