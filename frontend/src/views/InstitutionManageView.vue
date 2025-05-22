<template>
  <div class="institution-manage-container">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">机构管理</h1>
        <el-button type="primary" plain icon="House" @click="goToDashboard">返回主页</el-button>
      </div>
    </div>
    
    <el-tabs v-model="activeTab" class="institution-tabs">
      <el-tab-pane label="我的机构" name="myInstitution">
        <el-card v-if="loading" class="loading-card">
          <el-skeleton :rows="5" animated />
        </el-card>
        
        <div v-else class="institution-info-section">
          <!-- 显示已创建的机构信息 -->
          <div v-if="hasInstitution" class="full-height-layout">
            <!-- 左侧边栏：机构信息 -->
            <div class="sidebar">
              <div class="sidebar-header">
                <h2>机构信息</h2>
                <el-tag :type="getStatusType(institution.status)" effect="plain">
                  {{ getStatusText(institution.status) }}
                </el-tag>
              </div>
              
              <div class="institution-details">
                <div class="info-row">
                  <span class="info-label">机构名称：</span>
                  <span class="info-value">{{ institution.institution_name }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">地址：</span>
                  <span class="info-value">{{ institution.institution_address }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">联系电话：</span>
                  <span class="info-value">{{ institution.institution_phone || '未提供' }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">资质证明：</span>
                  <span class="info-value">{{ institution.institution_qualification }}</span>
                </div>
              </div>
              
              <div class="institution-actions" v-if="institution.status === 1">
                <el-button type="primary" @click="goToDetail">查看详情页面</el-button>
                <el-button type="success" @click="goToAddUserData">添加用户体检数据</el-button>
                <el-button type="danger" @click="confirmDeleteInstitution">删除机构</el-button>
              </div>
              
              <div v-else class="institution-actions">
                <el-button type="primary" @click="goToDetail">查看详情页面</el-button>
                <el-button type="danger" @click="confirmDeleteInstitution">删除机构</el-button>
              </div>
              
              <!-- 机构未被批准，显示状态信息 -->
              <div v-if="institution.status === 0" class="status-alert">
                <el-alert
                  title="您的机构信息正在审核中，请耐心等待"
                  type="warning"
                  :closable="false"
                  show-icon
                />
              </div>
              
              <div v-if="institution.status === 2" class="status-alert">
                <el-alert
                  title="您的机构申请已被拒绝，请联系管理员了解详情"
                  type="error"
                  :closable="false"
                  show-icon
                />
              </div>
            </div>
            
            <!-- 右侧主内容：套餐管理 -->
            <div class="main-content">
              <!-- 机构已被批准，显示套餐管理 -->
              <div v-if="institution.status === 1" class="package-section">
                <div class="content-header">
                  <h2>套餐管理</h2>
                </div>
                <package-manager :institution-id="institution.ID" />
              </div>
              
              <div v-else class="pending-content">
                <el-empty :description="institution.status === 0 ? '机构审核期间无法管理套餐' : '机构未获批准，无法管理套餐'" />
              </div>
            </div>
          </div>
          
          <!-- 未创建机构，显示创建表单 -->
          <div v-else class="create-institution-form">
            <div class="form-header">
              <h2>创建机构信息</h2>
              <p class="form-hint">请填写以下信息申请成为医疗机构，审核通过后可发布体检套餐</p>
            </div>
            
            <el-form :model="institutionForm" label-width="120px" class="institution-form">
              <el-form-item label="机构名称" required>
                <el-input v-model="institutionForm.institution_name" placeholder="请输入机构名称" />
              </el-form-item>
              <el-form-item label="机构地址" required>
                <el-input v-model="institutionForm.institution_address" placeholder="请输入机构地址" />
              </el-form-item>
              <el-form-item label="联系电话">
                <el-input v-model="institutionForm.institution_phone" placeholder="请输入联系电话" />
              </el-form-item>
              <el-form-item label="资质证明" required>
                <el-input 
                  v-model="institutionForm.institution_qualification" 
                  type="textarea" 
                  rows="5"
                  placeholder="请输入机构资质证明信息"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" :loading="submitting" @click="submitInstitution">提交申请</el-button>
              </el-form-item>
            </el-form>
          </div>
        </div>
      </el-tab-pane>
      
      <el-tab-pane v-if="isAdmin" label="审核申请" name="reviewApplications">
        <admin-review />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import PackageManager from '@/components/PackageManager.vue'
import AdminReview from '@/components/AdminReview.vue'

interface Institution {
  ID: number
  institution_name: string
  institution_address: string
  institution_qualification: string
  institution_phone?: string
  status: number
}

const router = useRouter()
const activeTab = ref('myInstitution')
const institution = ref<Institution>({
  ID: 0,
  institution_name: '',
  institution_address: '',
  institution_qualification: '',
  status: 0
})
const hasInstitution = ref(false)
const loading = ref(true)
const isAdmin = ref(false)
const submitting = ref(false)

const institutionForm = ref({
  institution_name: '',
  institution_address: '',
  institution_phone: '',
  institution_qualification: ''
})

// 获取用户角色和机构信息
const getUserInfo = async () => {
  try {
    loading.value = true
    const token = localStorage.getItem('jwt')
    const uid = localStorage.getItem('uid')
    
    if (!token || !uid) {
      router.push('/login')
      return
    }
    
    // 获取用户角色
    const userResponse = await axios.get(`${uid}`, {
      headers: { Authorization: `${token}` }
    })
    
    isAdmin.value = userResponse.data.user_type === 2
    
    // 获取用户的机构信息
    try {
      const institutionResponse = await axios.get(`/api/users/${uid}/institution`, {
        headers: { Authorization: `${token}` }
      })
      
      if (institutionResponse.data && institutionResponse.data.ID) {
        institution.value = institutionResponse.data
        hasInstitution.value = true
      }
    } catch (error) {
      console.log('No institution found for this user')
      hasInstitution.value = false
    }
  } catch (error) {
    console.error('Error fetching user info:', error)
    ElMessage.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
}

// 提交机构申请
const submitInstitution = async () => {
  // 验证表单
  if (!institutionForm.value.institution_name || 
      !institutionForm.value.institution_address || 
      !institutionForm.value.institution_qualification) {
    ElMessage.warning('请填写所有必填字段')
    return
  }
  
  submitting.value = true
  try {
    const token = localStorage.getItem('jwt')
    const uid = localStorage.getItem('uid')
    
    if (!token || !uid) {
      router.push('/login')
      return
    }
    
    const response = await axios.post(`/api/institutions/${uid}`, {
      institution_name: institutionForm.value.institution_name,
      institution_address: institutionForm.value.institution_address,
      institution_phone: institutionForm.value.institution_phone,
      institution_qualification: institutionForm.value.institution_qualification
    }, {
      headers: { Authorization: `${token}` }
    })
    
    ElMessage.success('机构信息提交成功，等待管理员审核')
    
    // 更新页面信息
    institution.value = response.data.institution
    hasInstitution.value = true
  } catch (error) {
    console.error('Error submitting institution:', error)
    ElMessage.error('提交机构信息失败')
  } finally {
    submitting.value = false
  }
}

// 前往机构详情页
const goToDetail = () => {
  router.push(`/institutions/${institution.value.ID}`)
}

// 前往主页
const goToDashboard = () => {
  router.push('/dashboard')
}

// 前往添加用户体检数据页面
const goToAddUserData = () => {
  router.push('/add-user-data')
}

// 删除机构确认
const confirmDeleteInstitution = () => {
  ElMessageBox.confirm(
    '确定要删除此机构吗？此操作不可逆。',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(async () => {
      try {
        const token = localStorage.getItem('jwt')
        await axios.delete(`/api/institutions/${institution.value.ID}`, {
          headers: { Authorization: `${token}` }
        })
        ElMessage.success('机构已成功删除')
        hasInstitution.value = false
        institution.value = {
          ID: 0,
          institution_name: '',
          institution_address: '',
          institution_qualification: '',
          status: 0
        }
      } catch (error) {
        console.error('Error deleting institution:', error)
        ElMessage.error('删除机构失败')
      }
    })
    .catch(() => {
      ElMessage.info('已取消删除')
    })
}

// 获取状态文本
const getStatusText = (status: number): string => {
  switch (status) {
    case 0: return '审核中'
    case 1: return '已批准'
    case 2: return '已拒绝'
    default: return '未知状态'
  }
}

// 获取状态类型
const getStatusType = (status: number): string => {
  switch (status) {
    case 0: return 'warning'
    case 1: return 'success'
    case 2: return 'danger'
    default: return 'info'
  }
}

onMounted(() => {
  getUserInfo()
})
</script>

<style scoped>
.institution-manage-container {
  width: 100vw;
  height: 100vh;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.page-header {
  margin-bottom: 20px;
  border-bottom: 1px solid #eaeaea;
  padding-bottom: 15px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  color: #303133;
  margin: 0;
}

.institution-tabs {
  margin-top: 20px;
  width: 100%;
  flex: 1;
  display: flex;
  flex-direction: column;
}

.loading-card {
  padding: 20px;
  min-height: 300px;
}

.institution-info-section {
  flex: 1;
  display: flex;
  flex-direction: column;
}

/* 全高布局 */
.full-height-layout {
  display: flex;
  width: 100%;
  height: calc(100vh - 200px);
  margin-top: 10px;
}

/* 左侧边栏 */
.sidebar {
  flex: 0 0 300px;
  background-color: #f9f9f9;
  border-right: 1px solid #eaeaea;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 15px;
  border-bottom: 1px solid #eaeaea;
  margin-bottom: 15px;
}

.sidebar-header h2 {
  margin: 0;
  font-size: 18px;
}

.institution-details {
  padding: 10px 0;
  flex: 1;
}

.info-row {
  margin-bottom: 20px;
}

.info-label {
  display: block;
  font-weight: bold;
  color: #606266;
  margin-bottom: 5px;
}

.info-value {
  display: block;
  color: #303133;
}

.institution-actions {
  margin-top: auto;
  padding-top: 15px;
  display: flex;
  justify-content: center;
}

.status-alert {
  margin-top: auto;
  padding-top: 15px;
}

/* 右侧主内容 */
.main-content {
  flex: 1;
  padding: 0 20px 20px 20px;
  overflow-y: auto;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #eaeaea;
  margin-bottom: 20px;
}

.content-header h2 {
  margin: 0;
  font-size: 18px;
}

.package-section {
  width: 100%;
}

.pending-content {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background-color: #f9f9f9;
  border-radius: 4px;
}

/* 创建机构表单 */
.create-institution-form {
  max-width: 800px;
  margin: 0 auto;
  padding: 30px;
  background-color: #ffffff;
  border-radius: 4px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.form-header {
  margin-bottom: 25px;
  text-align: center;
}

.form-header h2 {
  margin-bottom: 10px;
  color: #303133;
}

.form-hint {
  color: #909399;
  font-size: 14px;
  margin: 0;
}

.institution-form {
  max-width: 600px;
  margin: 0 auto;
}

/* 响应式设计 */
@media screen and (max-width: 768px) {
  .full-height-layout {
    flex-direction: column;
    height: auto;
  }
  
  .sidebar {
    flex: none;
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #eaeaea;
    margin-bottom: 20px;
  }
  
  .main-content {
    padding: 0;
  }
}
</style>
