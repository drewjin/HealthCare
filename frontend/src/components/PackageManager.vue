<template>
  <div class="package-manager">
    <h2>管理体检套餐</h2>
    
    <el-alert
      v-if="saveSuccess"
      title="套餐信息保存成功"
      type="success"
      :closable="false"
      show-icon
      class="alert-message"
    />
    
    <el-alert
      v-if="saveError"
      :title="errorMessage"
      type="error"
      :closable="false"
      show-icon
      class="alert-message"
    />
    
    <el-button type="primary" @click="addNewPackage" class="add-button">
      添加新套餐
    </el-button>
    
    <div v-if="packages.length > 0" class="packages-list">
      <el-card v-for="(pkg, index) in packages" :key="index" class="package-card">
        <template #header>
          <div class="package-header">
            <el-input v-model="pkg.name" placeholder="套餐名称" />
            <el-button type="danger" @click="removePackage(index)" size="small">
              删除
            </el-button>
          </div>
        </template>
        
        <div class="package-form">
          <el-form label-position="top">
            <el-form-item label="套餐描述">
              <el-input v-model="pkg.description" type="textarea" :rows="3" placeholder="请输入套餐描述" />
            </el-form-item>
            
            <el-form-item label="适用人群">
              <el-input v-model="pkg.suitableFor" placeholder="例如：中老年人、白领人士等" />
            </el-form-item>
            
            <el-form-item label="检查项目">
              <el-input v-model="pkg.items" type="textarea" :rows="3" placeholder="请列出检查项目，可用逗号分隔" />
            </el-form-item>
            
            <el-form-item label="价格 (元)">
              <el-input-number v-model="pkg.price" :min="0" :precision="2" :step="100" />
            </el-form-item>
          </el-form>
        </div>
      </el-card>
    </div>
    
    <div v-else class="empty-packages">
      <el-empty description="暂无体检套餐，请添加" />
    </div>
    
    <div class="action-buttons">
      <el-button type="primary" @click="savePackages" :loading="saving">
        保存套餐信息
      </el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  institutionId: number | string
}>()

interface Package {
  name: string
  description: string
  suitableFor: string
  items: string
  price: number
}

const packages = ref<Package[]>([])
const saving = ref(false)
const saveSuccess = ref(false)
const saveError = ref(false)
const errorMessage = ref('')

const fetchPackages = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/institutions/${props.institutionId}/packages`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    // 解析套餐数据
    if (response.data.packages) {
      try {
        packages.value = JSON.parse(response.data.packages)
      } catch (e) {
        console.error('Failed to parse packages:', e)
        packages.value = []
      }
    } else {
      packages.value = []
    }
  } catch (error) {
    console.error('Failed to fetch packages:', error)
    ElMessage.error('获取套餐信息失败')
  }
}

const addNewPackage = () => {
  packages.value.push({
    name: '新套餐',
    description: '',
    suitableFor: '',
    items: '',
    price: 0
  })
}

const removePackage = (index: number) => {
  packages.value.splice(index, 1)
}

const savePackages = async () => {
  resetAlerts()
  
  // 验证套餐信息
  for (const pkg of packages.value) {
    if (!pkg.name || !pkg.description || !pkg.items) {
      ElMessage.warning('请填写所有必要的套餐信息')
      return
    }
  }
  
  saving.value = true
  try {
    const token = localStorage.getItem('jwt')
    await axios.put(`/api/institutions/${props.institutionId}/packages`, {
      packages: JSON.stringify(packages.value)
    }, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    saveSuccess.value = true
    ElMessage.success('套餐信息保存成功')
  } catch (error) {
    console.error('Failed to save packages:', error)
    saveError.value = true
    errorMessage.value = '保存套餐信息失败'
  } finally {
    saving.value = false
    
    // 3秒后隐藏成功消息
    if (saveSuccess.value) {
      setTimeout(() => {
        saveSuccess.value = false
      }, 3000)
    }
  }
}

const resetAlerts = () => {
  saveSuccess.value = false
  saveError.value = false
  errorMessage.value = ''
}

onMounted(() => {
  fetchPackages()
})
</script>

<style scoped>
.package-manager {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.alert-message {
  margin-bottom: 20px;
}

.add-button {
  margin-bottom: 20px;
}

.packages-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
  margin-bottom: 20px;
}

.package-card {
  margin-bottom: 15px;
}

.package-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.package-form {
  padding: 10px 0;
}

.action-buttons {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.empty-packages {
  margin: 30px 0;
  text-align: center;
}
</style>