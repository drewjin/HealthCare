<template>
  <div class="health-data-form">
    <h2>添加用户体检数据</h2>
    
    <el-form v-if="loaded && !loadFailed" :model="formData" label-width="120px">
      <div class="form-header">
        <el-alert
          v-if="submitSuccess"
          title="体检数据保存成功"
          type="success"
          :closable="false"
          show-icon
          class="alert-message"
        />
      </div>
      
      <div class="user-info-section">
        <el-descriptions title="用户信息" border>
          <el-descriptions-item label="用户ID">{{ customerId }}</el-descriptions-item>
          <el-descriptions-item label="套餐ID">{{ planId }}</el-descriptions-item>
          <el-descriptions-item label="套餐名称">{{ planName }}</el-descriptions-item>
        </el-descriptions>
      </div>
      
      <el-divider content-position="center">体检项目</el-divider>
      
      <div v-if="healthItems.length === 0" class="empty-state">
        <el-empty description="暂无体检项目" />
      </div>
      
      <template v-else>
        <div v-for="(item, index) in healthItems" :key="index" class="health-item">
          <el-form-item :label="item.item_name">
            <div v-if="item.item_description" class="item-description">
              {{ item.item_description }}
            </div>
            <el-input v-model="healthItemValues[index].value" placeholder="请输入检测结果"></el-input>
          </el-form-item>
        </div>
        
        <el-form-item>
          <el-button type="primary" @click="submitHealthData" :loading="submitting">保存体检数据</el-button>
          
          <el-dropdown split-button type="success" @click="checkIfAllComplete" @command="handleStatusCommand">
            检查完成状态
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="mark-complete">手动标记完成</el-dropdown-item>
                <el-dropdown-item command="mark-incomplete">标记为未完成</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </template>
    </el-form>
    
    <!-- 加载失败时显示错误信息 -->
    <el-result
      v-else-if="loadFailed"
      icon="error"
      :title="errorTitle"
      :sub-title="errorMessage"
    >
      <template #extra>
        <el-button type="primary" @click="goBack">返回上一页</el-button>
        <el-button v-if="isAdmin && errorCode === 'PLAN_NOT_FOUND'" @click="recoverPlan">恢复套餐</el-button>
      </template>
    </el-result>
    
    <div v-else class="loading-state">
      <el-skeleton :rows="10" animated />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()

// Form state
const loaded = ref(false)
const loadFailed = ref(false)
const submitting = ref(false)
const submitSuccess = ref(false)
const errorMessage = ref('')
const errorTitle = ref('加载失败')
const errorCode = ref('')
const formData = ref({})

// Check if user is admin
const isAdmin = ref(false)
const checkIfAdmin = () => {
  const userType = localStorage.getItem('userType')
  isAdmin.value = userType === '2' // 2 = admin user
}

// Parameters from route
const customerId = computed(() => route.params.customer_id as string)
const planId = computed(() => route.params.plan_id as string)
const planName = ref('')

// Health items data
interface HealthItem {
  id: number
  plan_item_id: number
  item_name: string
  item_description: string
}

const healthItems = ref<HealthItem[]>([])
const healthItemValues = ref<{ item_id: number, value: string }[]>([])

// Fetch health items for this plan
const fetchHealthItems = async () => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('未登录或令牌丢失')
      router.push('/login')
      return
    }
    
    // 获取plan items using the userview/plan API instead, which uses PlanHeathItem relation correctly
    const response = await axios.get(`/api/userview/plan`, {
      params: { plan_id: planId.value },
      headers: { Authorization: token }
    })
    
    if (response.data && Array.isArray(response.data.plan_items) && response.data.plan_items.length > 0) {
      interface ApiHealthItem {
        plan_item_id: number   // PlanHeathItem的ID
        plan_id: number
        plan_name: string
        item_id: number        // HealthItem的ID
        item_name: string
        item_description: string
        item_value: string
      }
      
      // Extract plan name from first item
      if (response.data.plan_name) {
        planName.value = response.data.plan_name;
      } else if (response.data.plan_items.length > 0) {
        planName.value = response.data.plan_items[0].plan_name;
      }
      
      // Map the items
      healthItems.value = response.data.plan_items.map((item: ApiHealthItem) => ({
        id: item.item_id,
        plan_item_id: item.plan_item_id,
        item_name: item.item_name,
        item_description: item.item_description || ''
      }));
      
      // Initialize health item values
      healthItemValues.value = response.data.plan_items.map((item: ApiHealthItem) => ({
        item_id: item.item_id,
        value: item.item_value || ''
      }));
      
      loaded.value = true
    } else {
      // No plan items found, but we still show the UI with an empty state
      healthItems.value = []
      // Try to get the plan name from the response if available
      if (response.data && response.data.plan_name) {
        planName.value = response.data.plan_name;
      } else {
        planName.value = '未找到套餐名称';
      }
      loaded.value = true
      
      // Show a warning message to the user
      if (response.data && response.data.message) {
        ElMessage.warning(response.data.message);
      } else {
        ElMessage.warning('该套餐下没有任何体检项目');
      }
    }
  } catch (error: any) {
    console.error('获取体检项目出错：', error)
    
    // Set load failed state
    loadFailed.value = true
    loaded.value = true
    
    // Handle different error cases
    if (error.response) {
      // Handle specific status codes
      if (error.response.status === 404) {
        errorTitle.value = '套餐不存在'
        errorCode.value = 'PLAN_NOT_FOUND'
        
        if (error.response.data && error.response.data.error) {
          if (error.response.data.error.includes('Plan not found')) {
            errorMessage.value = '套餐不存在或已被删除，请联系管理员';
          } else {
            errorMessage.value = error.response.data.error;
          }
        } else {
          errorMessage.value = '套餐不存在或已被删除，请联系管理员';
        }
      } else if (error.response.status === 403) {
        errorTitle.value = '无权访问'
        errorCode.value = 'PERMISSION_DENIED'
        errorMessage.value = '您没有权限访问该套餐的体检项目';
      } else {
        // Generic error handling
        errorTitle.value = '加载失败'
        errorCode.value = 'UNKNOWN_ERROR'
        errorMessage.value = error.response.data?.error || '获取体检项目失败，请重试';
      }
    } else {
      // Network error or other issues
      errorTitle.value = '连接错误'
      errorCode.value = 'CONNECTION_ERROR'
      errorMessage.value = '无法连接到服务器，请检查网络连接后重试';
    }
  }
}

// Submit user health data
const submitHealthData = async () => {
  try {
    submitting.value = true
    submitSuccess.value = false
    
    // Validate that at least some values are filled in
    const hasValues = healthItemValues.value.some(item => item.value.trim() !== '')
    if (!hasValues) {
      ElMessage.warning('请至少填写一项检测结果')
      submitting.value = false
      return
    }
    
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('未登录或令牌丢失')
      router.push('/login')
      return
    }
    
    // Format data for API
    const payload = healthItemValues.value
      .filter(item => item.value.trim() !== '')
      .map(item => ({
        item_id: item.item_id,
        item_value: item.value.trim()
      }))
    
    // Send data to API
    const response = await axios.post(
      `/api/adduserdata/${customerId.value}/${planId.value}`, 
      payload,
      { headers: { Authorization: token } }
    )
    
    if (response.status === 200) {
      submitSuccess.value = true
      ElMessage.success('体检数据保存成功')
      
      // Reset form after successful submission
      setTimeout(() => {
        submitSuccess.value = false
        // Update package status - Only if all items are complete
        checkIfAllComplete()
      }, 1500)
    } else {
      throw new Error('保存体检数据失败')
    }
  } catch (error) {
    console.error('提交体检数据出错：', error)
    errorMessage.value = '保存体检数据失败，请重试'
  } finally {
    submitting.value = false
  }
}

// Check if all items are complete and prompt to update package status if needed
const checkIfAllComplete = async () => {
  try {
    // Reload the items to get the latest values
    await fetchHealthItems()
    
    // Check if all items have values
    const allComplete = healthItemValues.value.every(item => item.value.trim() !== '')
    
    if (allComplete) {
      // Prompt the user to mark the package as complete
      const { isConfirmed } = await ElMessageBox.confirm(
        '所有体检项目已完成。是否将套餐标记为已完成？',
        '完成确认',
        {
          confirmButtonText: '确认完成',
          cancelButtonText: '保持进行中',
          type: 'success'
        }
      ).catch(() => ({ isConfirmed: false }))
      
      if (isConfirmed) {
        updatePackageStatus(1) // 1 = completed
      }
    }
  } catch (error) {
    console.error('检查完成状态出错：', error)
  }
}

// Update package status to completed if applicable
const updatePackageStatus = async (status: number) => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) return
    
    const response = await axios.patch(
      `/api/user-packages/${customerId.value}/${planId.value}/status`,
      { status }, // 1 = completed, 0 = pending
      { headers: { Authorization: token } }
    )
    
    if (response.status === 200) {
      ElMessage.success(status === 1 ? '套餐已标记为完成' : '套餐已标记为进行中')
    }
  } catch (error) {
    console.error('更新套餐状态失败：', error)
    ElMessage.error('更新套餐状态失败，请重试')
  }
}

// Recover a deleted plan (admin only function)
const recoverPlan = async () => {
  try {
    // Confirm with the user first
    await ElMessageBox.confirm(
      '确定要恢复此套餐吗？这将使套餐重新可用。',
      '恢复套餐确认',
      {
        confirmButtonText: '确认恢复',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('未登录或令牌丢失')
      router.push('/login')
      return
    }
    
    // Call the recovery API
    const response = await axios.post(
      `/api/admin/plans/recover/${planId.value}`,
      {},
      { headers: { Authorization: token } }
    )
    
    if (response.status === 200) {
      ElMessage.success('套餐恢复成功')
      // Reload the data
      loadFailed.value = false
      loaded.value = false
      fetchHealthItems()
    } else {
      throw new Error('恢复套餐失败')
    }
  } catch (error: any) {
    if (error.toString().includes('cancel')) {
      // User cancelled the operation
      return
    }
    
    console.error('恢复套餐失败：', error)
    ElMessage.error(error.response?.data?.error || '恢复套餐失败，请重试')
  }
}

// Navigate back
const goBack = () => {
  router.back()
}

onMounted(() => {
  // Check if the user is an admin
  checkIfAdmin()
  
  // Fetch health items for the specified plan
  fetchHealthItems()
})
</script>

<style scoped>
.health-data-form {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.form-header {
  margin-bottom: 20px;
}

.alert-message {
  margin-bottom: 15px;
}

.user-info-section {
  margin-bottom: 20px;
}

.health-item {
  margin-bottom: 15px;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.item-description {
  color: #666;
  font-size: 0.9em;
  margin-bottom: 5px;
  font-style: italic;
}

.loading-state {
  padding: 30px;
}

.empty-state {
  margin: 40px 0;
  display: flex;
  justify-content: center;
}
</style>
