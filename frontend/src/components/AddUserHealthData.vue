<template>
  <div class="health-data-form">
    <h2>添加用户体检数据</h2>
    
    <el-form v-if="loaded" :model="formData" label-width="120px">
      <div class="form-header">
        <el-alert
          v-if="submitSuccess"
          title="体检数据保存成功"
          type="success"
          :closable="false"
          show-icon
          class="alert-message"
        />
        
        <el-alert
          v-if="submitError"
          :title="errorMessage"
          type="error"
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
            <el-input v-model="healthItemValues[index].value" placeholder="请输入检测结果"></el-input>
          </el-form-item>
        </div>
        
        <el-form-item>
          <el-button type="primary" @click="submitHealthData" :loading="submitting">保存体检数据</el-button>
          <el-button @click="goBack">返回</el-button>
        </el-form-item>
      </template>
    </el-form>
    
    <div v-else class="loading-state">
      <el-skeleton :rows="10" animated />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()

// Form state
const loaded = ref(false)
const submitting = ref(false)
const submitSuccess = ref(false)
const submitError = ref(false)
const errorMessage = ref('')
const formData = ref({})

// Parameters from route
const customerId = computed(() => route.params.customer_id as string)
const planId = computed(() => route.params.plan_id as string)
const planName = ref('')

// Health items data
interface HealthItem {
  id: number
  item_name: string
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
    
    // First get plan items
    const response = await axios.get(`/api/plans/${planId.value}/items`, {
      headers: { Authorization: token }
    })
    
    if (response.data && Array.isArray(response.data.items)) {
      interface ApiHealthItem {
        id?: number
        item_id?: number
        item_name: string
      }
      
      healthItems.value = response.data.items.map((item: ApiHealthItem) => ({
        id: item.id || item.item_id || 0,
        item_name: item.item_name
      }))
      
      // Initialize health item values
      healthItemValues.value = healthItems.value.map((item: HealthItem) => ({
        item_id: item.id,
        value: ''
      }))
      
      // Get plan name
      if (response.data.plan_name) {
        planName.value = response.data.plan_name
      } else {
        // If plan name not included in items response, fetch it separately
        const planResponse = await axios.get(`/api/plans/${planId.value}`, {
          headers: { Authorization: token }
        })
        
        if (planResponse.data && planResponse.data.plan_name) {
          planName.value = planResponse.data.plan_name
        }
      }
      
      loaded.value = true
    } else {
      throw new Error('获取体检项目失败')
    }
  } catch (error) {
    console.error('获取体检项目出错：', error)
    errorMessage.value = '获取体检项目失败，请重试'
    submitError.value = true
    loaded.value = true
  }
}

// Submit user health data
const submitHealthData = async () => {
  try {
    submitting.value = true
    submitSuccess.value = false
    submitError.value = false
    
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
      }, 3000)
    } else {
      throw new Error('保存体检数据失败')
    }
  } catch (error) {
    console.error('提交体检数据出错：', error)
    submitError.value = true
    errorMessage.value = '保存体检数据失败，请重试'
  } finally {
    submitting.value = false
  }
}

// Navigate back
const goBack = () => {
  router.back()
}

onMounted(() => {
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

.loading-state {
  padding: 30px;
}

.empty-state {
  margin: 40px 0;
  display: flex;
  justify-content: center;
}
</style>
