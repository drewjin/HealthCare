<template>
  <div class="package-manager">
    <div class="header-section">
      <el-alert
        v-if="saveSuccess"
        :title="successMessage"
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
      
      <el-alert
        v-if="isEditing && !saveSuccess && !saveError"
        title="您正在编辑套餐信息，完成后请点击保存按钮"
        type="info"
        :closable="false"
        show-icon
        class="alert-message"
      />
    </div>
    
    <div class="actions-row">
      <el-button 
        v-if="!isEditing" 
        type="primary" 
        @click="addNewPackage" 
        class="add-button"
      >
        创建新套餐
      </el-button>
      
      <el-button 
        v-if="isEditing"
        type="primary" 
        @click="savePackages" 
        :loading="saving" 
        class="save-button"
      >
        保存套餐信息
      </el-button>
    </div>
    
    <div v-if="packages.length > 0" class="packages-list">
      <el-collapse v-model="activePackages" accordion @change="handleCollapseChange">
        <el-collapse-item 
          v-for="(pkg, index) in packages" 
          :key="index"
          :name="index"
        >
          <template #title>
            <div class="package-collapse-header">
              <span>{{ pkg.name || '新套餐' }}</span>
              <span class="package-price">{{ pkg.price }}元</span>
            </div>
          </template>
          
          <div class="package-content">
            <el-form :model="pkg" label-width="80px" class="package-form">
              <div class="form-grid">
                <el-form-item label="套餐名称" class="form-item">
                  <el-input v-model="pkg.name" placeholder="套餐名称" />
                </el-form-item>
                
                <el-form-item label="适用人群" class="form-item">
                  <el-input v-model="pkg.suitableFor" placeholder="例如：中老年人、白领人士等" />
                </el-form-item>
                
                <el-form-item label="价格" class="form-item">
                  <el-input-number v-model="pkg.price" :min="0" :precision="2" :step="100" style="width: 100%;" controls-position="right" />
                </el-form-item>
              </div>
              
              <el-form-item label="套餐描述">
                <el-input v-model="pkg.description" type="textarea" :rows="2" placeholder="请输入套餐描述" />
              </el-form-item>
              
              <el-form-item label="检查项目">
                <el-input v-model="pkg.items" type="textarea" :rows="2" placeholder="请列出检查项目，可用逗号分隔" />
              </el-form-item>
              
              <div class="package-actions">
                <el-button type="danger" @click.stop="removePackage(index)" size="small">
                  删除套餐
                </el-button>
              </div>
            </el-form>
          </div>
        </el-collapse-item>
      </el-collapse>
    </div>
    
    <div v-else class="empty-packages">
      <el-empty description="暂无体检套餐" />
      <div class="empty-action">
        <el-button type="primary" @click="addNewPackage" class="add-button">
          添加新套餐
        </el-button>
      </div>
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
  id?: number
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
const successMessage = ref('套餐信息保存成功')
const activePackages = ref<number[]>([]) // 用于控制折叠面板的展开状态
const isEditing = ref(false) // 控制是否显示保存按钮

const fetchPackages = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/institutions/${props.institutionId}/plans`, {
      headers: { Authorization: `${token}` }
    })
    
    // 解析套餐数据
    if (response.data.plans && response.data.plans.length > 0) {
      packages.value = response.data.plans.map((plan: any) => {
        return {
          id: plan.ID || plan.id,
          name: plan.name || plan.plan_name || '',
          description: plan.description || '',
          suitableFor: plan.suitable_for || '',
          items: plan.items || '',
          price: plan.price || plan.plan_price || 0
        }
      })
      
      // 如果有现有套餐，初始化状态
      if (packages.value.length > 0) {
        isEditing.value = false // 初始状态不显示保存按钮
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
  const newIndex = packages.value.length
  packages.value.push({
    name: '新套餐',
    description: '',
    suitableFor: '',
    items: '',
    price: 0
  })
  // 自动展开新添加的套餐
  activePackages.value = [...activePackages.value, newIndex]
  // 显示保存按钮
  isEditing.value = true
}

const removePackage = (index: number) => {
  packages.value.splice(index, 1)
  // 移除被删除套餐的激活状态
  activePackages.value = activePackages.value.filter(i => i !== index)
  // 更新索引大于被删除套餐的激活状态
  activePackages.value = activePackages.value.map(i => i > index ? i - 1 : i)
}

const savePackages = async () => {
  resetAlerts()
  
  // 验证套餐信息
  for (const pkg of packages.value) {
    if (!pkg.name || !pkg.items) {
      ElMessage.warning('请填写所有必要的套餐信息')
      return
    }
  }
  
  saving.value = true
  try {
    const token = localStorage.getItem('jwt')
    const currentPackage = packages.value[0]
    
    // 检查是否有ID，如果有则是更新操作
    if (currentPackage.id) {
      // 使用PATCH更新现有套餐
      await axios.patch(`/api/institutions/${props.institutionId}/item`, {
        plan_id: currentPackage.id,
        plan_name: currentPackage.name,
        plan_price: currentPackage.price,
        description: currentPackage.description,
        suitable_for: currentPackage.suitableFor
      }, {
        headers: { Authorization: `${token}` }
      })
      
      successMessage.value = '套餐信息更新成功'
    } else {
      // 创建新套餐
      await axios.post(`/api/institutions/${props.institutionId}/plans`, {
        plan_name: currentPackage.name,
        health_item: currentPackage.items.split(',')[0],
        item_description: currentPackage.description,
        plan_price: currentPackage.price,
        description: currentPackage.description,
        suitable_for: currentPackage.suitableFor
      }, {
        headers: { Authorization: `${token}` }
      })
      
      // 如果有多个体检项目，为套餐添加剩余项目
      if (currentPackage.items.split(',').length > 1) {
        const planResponse = await axios.get(`/api/institutions/${props.institutionId}/plans`, {
          headers: { Authorization: `${token}` }
        })
        
        const planId = planResponse.data.plans[0].ID
        
        for (let i = 1; i < currentPackage.items.split(',').length; i++) {
          await axios.post(`/api/institutions/${props.institutionId}/${planId}/item`, {
            health_item: currentPackage.items.split(',')[i].trim(),
            item_description: currentPackage.description
          }, {
            headers: { Authorization: `${token}` }
          })
        }
      }
      
      successMessage.value = '套餐信息创建成功'
    }
    
    saveSuccess.value = true
    
    // 刷新套餐数据
    await fetchPackages()
    
    // 保存成功后隐藏保存按钮
    isEditing.value = false
    
  } catch (error: any) {
    console.error('Failed to save packages:', error)
    saveError.value = true
    errorMessage.value = error.response?.data?.error || '保存套餐信息失败'
    ElMessage.error(errorMessage.value)
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
  successMessage.value = '套餐信息保存成功'
}

const handleCollapseChange = (activeNames: number[]) => {
  // 当打开一个套餐时，启用编辑模式
  if (activeNames.length > 0) {
    isEditing.value = true
  }
}

onMounted(() => {
  fetchPackages()
})
</script>

<style scoped>
.package-manager {
  width: 100%;
  padding: 10px 0;
  box-sizing: border-box;
}

.header-section {
  margin-bottom: 15px;
}

.alert-message {
  margin-bottom: 10px;
}

.actions-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 15px;
}

.add-button, .save-button {
  min-width: 110px;
}

.packages-list {
  margin-bottom: 15px;
}

.package-collapse-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.package-price {
  font-weight: bold;
  color: #409EFF;
}

.package-content {
  padding: 5px 0;
}

.package-form {
  width: 100%;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 15px;
  margin-bottom: 10px;
}

.package-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.empty-packages {
  margin: 20px 0;
  text-align: center;
}

.empty-action {
  margin-top: 20px;
}

:deep(.el-collapse) {
  border: none;
}

:deep(.el-collapse-item__header) {
  font-size: 16px;
  padding: 0 15px;
  font-weight: bold;
  border-radius: 4px;
  background-color: #f5f7fa;
}

:deep(.el-collapse-item__wrap) {
  padding: 10px;
  border-bottom: none;
}

:deep(.el-collapse-item__content) {
  padding: 10px;
}

:deep(.el-form-item__label) {
  font-weight: bold;
}

:deep(.el-input-number .el-input__inner) {
  text-align: left;
}
</style>