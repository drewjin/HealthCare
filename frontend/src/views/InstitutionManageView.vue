<template>
  <div class="institution-manage-container">
    <h1 class="page-title">机构管理</h1>
    
    <el-tabs v-model="activeTab" class="institution-tabs">
      <el-tab-pane label="我的机构" name="myInstitution">
        <el-card v-if="loading" class="loading-card">
          <el-skeleton :rows="5" animated />
        </el-card>
        
        <div v-else class="institution-info-section">
          <!-- 显示已创建的机构信息 -->
          <div v-if="hasInstitution">
            <el-card class="institution-info-card">
              <template #header>
                <div class="card-header">
                  <span>机构信息</span>
                  <el-tag :type="getStatusType(institution.status)" effect="plain">
                    {{ getStatusText(institution.status) }}
                  </el-tag>
                </div>
              </template>
              <div class="institution-details">
                <p><strong>机构名称：</strong> {{ institution.institution_name }}</p>
                <p><strong>地址：</strong> {{ institution.institution_address }}</p>
                <p><strong>联系电话：</strong> {{ institution.institution_phone || '未提供' }}</p>
                <p><strong>资质证明：</strong> {{ institution.institution_qualification }}</p>
              </div>
              <div class="institution-actions" v-if="institution.status === 1">
                <el-button type="primary" @click="goToDetail">查看详情</el-button>
              </div>
            </el-card>
            
            <!-- 机构已被批准，显示套餐管理 -->
            <div v-if="institution.status === 1" class="package-section">
              <h2>套餐管理</h2>
              <package-manager :institution-id="institution.ID" />
            </div>
            
            <!-- 机构未被批准，显示状态信息 -->
            <div v-else-if="institution.status === 0" class="pending-message">
              <el-alert
                title="您的机构信息正在审核中，请耐心等待"
                type="warning"
                :closable="false"
                show-icon
              />
            </div>
            
            <div v-else-if="institution.status === 2" class="rejected-message">
              <el-alert
                title="您的机构申请已被拒绝，请联系管理员了解详情"
                type="error"
                :closable="false"
                show-icon
              />
            </div>
          </div>
          
          <!-- 未创建机构，显示创建表单 -->
          <div v-else class="create-institution-form">
            <h2>创建机构信息</h2>
            <el-form :model="institutionForm" label-width="120px">
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
                  rows="3"
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
import { ElMessage } from 'element-plus'
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
    const userResponse = await axios.get(`/api/user/${uid}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    isAdmin.value = userResponse.data.user_type === 2
    
    // 获取用户的机构信息
    try {
      const institutionResponse = await axios.get(`/api/user/${uid}/institution`, {
        headers: { Authorization: `Bearer ${token}` }
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
      headers: { Authorization: `Bearer ${token}` }
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

// 获取状态文本
const getStatusText = (status: number): string => {
  switch (status) {
    case 0: return '待审核'
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
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-title {
  margin-bottom: 20px;
  color: #303133;
}

.institution-tabs {
  margin-top: 20px;
}

.loading-card {
  padding: 20px;
}

.institution-info-section {
  margin-top: 20px;
}

.institution-info-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.institution-details {
  padding: 10px 0;
}

.institution-details p {
  margin-bottom: 10px;
}

.institution-actions {
  margin-top: 15px;
  display: flex;
  justify-content: flex-end;
}

.package-section {
  margin-top: 30px;
}

.package-section h2 {
  margin-bottom: 15px;
  font-size: 18px;
  color: #303133;
}

.pending-message,
.rejected-message {
  margin: 30px 0;
}

.create-institution-form {
  margin-top: 20px;
}

.create-institution-form h2 {
  margin-bottom: 20px;
  font-size: 18px;
  color: #303133;
}
</style>
