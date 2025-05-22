<template>
  <div class="add-user-data-container">
    <div class="page-header-actions">
      <el-page-header @back="goBack" title="返回机构管理">
        <template #content>
          <span>添加用户体检数据</span>
        </template>
      </el-page-header>
      <el-button type="primary" plain icon="House" @click="goToDashboard">返回主页</el-button>
    </div>

    <div class="user-selector-section">
      <el-card shadow="hover">
        <template #header>
          <div class="card-header">
            <span>选择用户和套餐</span>
          </div>
        </template>
        
        <div v-if="step === 1">
          <el-form :model="formData" label-width="100px">
            <el-form-item label="用户ID">
              <el-input v-model="formData.customerId" placeholder="请输入用户ID"></el-input>
            </el-form-item>
            
            <el-form-item label="套餐ID">
              <el-input v-model="formData.planId" placeholder="请输入套餐ID"></el-input>
            </el-form-item>
            
            <el-divider></el-divider>
            
            <el-form-item>
              <el-button type="primary" @click="validateAndProceed" :disabled="!isFormValid">
                下一步：填写体检数据
              </el-button>
            </el-form-item>
          </el-form>
        </div>
        
        <div v-else-if="step === 2">
          <add-user-health-data />
        </div>
      </el-card>
    </div>
    
    <!-- 用户套餐列表 -->
    <div class="user-plans-section" v-if="step === 1 && showUserPlans">
      <el-card shadow="hover">
        <template #header>
          <div class="card-header">
            <span>用户已购套餐列表</span>
          </div>
        </template>
        
        <div v-if="loadingUserPlans" class="loading-state">
          <el-skeleton :rows="3" animated />
        </div>
        
        <div v-else-if="userPlans.length === 0" class="empty-state">
          <el-empty description="该用户暂无已购套餐" />
        </div>
        
        <el-table v-else :data="userPlans" style="width: 100%">
          <el-table-column prop="id" label="ID" width="80" />
          <el-table-column prop="plan_name" label="套餐名称" />
          <el-table-column prop="institution_name" label="体检机构" />
          <el-table-column prop="status" label="状态">
            <template #default="scope">
              <el-tag :type="scope.row.status === 0 ? 'warning' : 'success'">
                {{ scope.row.status === 0 ? '待检测' : '已完成' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button
                size="small"
                type="primary"
                @click="selectPlan(scope.row)"
                :disabled="scope.row.status === 1"
              >
                填写体检数据
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import AddUserHealthData from '@/components/AddUserHealthData.vue'

const router = useRouter()
const step = ref(1)

// Form data
const formData = ref({
  customerId: '',
  planId: ''
})

// Validation
const isFormValid = computed(() => {
  return formData.value.customerId.trim() !== '' && 
         formData.value.planId.trim() !== '' && 
         !isNaN(Number(formData.value.customerId)) && 
         !isNaN(Number(formData.value.planId))
})

// User plans data
const showUserPlans = ref(false)
const loadingUserPlans = ref(false)
const userPlans = ref<Array<{
  id: number;
  plan_name: string;
  institution_name: string;
  status: number;
}>>([])

// Watch for customer ID changes to fetch plans
watch(() => formData.value.customerId, (newValue) => {
  if (newValue && !isNaN(Number(newValue))) {
    fetchUserPlans(newValue)
  } else {
    showUserPlans.value = false
    userPlans.value = []
  }
})

// Fetch user plans
const fetchUserPlans = async (userId: string) => {
  if (!userId || userId.trim() === '') return
  
  try {
    loadingUserPlans.value = true
    showUserPlans.value = true
    
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('未登录或令牌丢失')
      router.push('/login')
      return
    }
    
    const response = await axios.get(`/api/users/${userId}/plans`, {
      headers: { Authorization: token }
    })
    
    if (response.data && Array.isArray(response.data)) {
      userPlans.value = response.data
    } else {
      userPlans.value = []
    }
  } catch (error) {
    console.error('获取用户套餐列表出错：', error)
    userPlans.value = []
    ElMessage.error('获取用户套餐列表失败')
  } finally {
    loadingUserPlans.value = false
  }
}

// Select a plan from the table
const selectPlan = (plan: any) => {
  formData.value.planId = plan.id.toString()
  ElMessage.success(`已选择套餐：${plan.plan_name}`)
}

// Validate form and proceed to next step
const validateAndProceed = () => {
  if (!isFormValid.value) {
    ElMessage.warning('请填写正确的用户ID和套餐ID')
    return
  }
  
  // Navigate to the detail route with the parameters
  router.push(`/add-user-data/${formData.value.customerId}/${formData.value.planId}`)
}

// Navigation functions
const goBack = () => {
  if (step.value === 2) {
    step.value = 1
  } else {
    router.push('/institution-manage')
  }
}

const goToDashboard = () => {
  router.push('/dashboard')
}
</script>

<style scoped>
.add-user-data-container {
  padding: 20px;
}

.page-header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.user-selector-section {
  margin-bottom: 30px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.loading-state {
  padding: 20px 0;
}

.empty-state {
  padding: 30px 0;
  display: flex;
  justify-content: center;
}
</style>
