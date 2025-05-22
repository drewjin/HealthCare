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
                <el-form-item label="套餐名称" class="form-item" required>
                  <el-input v-model="pkg.name" placeholder="套餐名称（必填）" />
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
              
              <el-form-item label="检查项目" required>
                <div class="health-items-container">
                  <el-input v-model="pkg.items" type="textarea" :rows="2" placeholder="请列出检查项目，可用逗号分隔（必填）" />
                  <div class="health-items-controls">
                    <el-button type="primary" size="small" @click="viewHealthItems(pkg)" :disabled="!pkg.id">
                      查看健康检查项目
                    </el-button>
                  </div>
                </div>
              </el-form-item>
              
              <div class="package-actions">
                <el-button type="primary" @click.stop="savePackage(index)" size="small" style="margin-right: 10px;">
                  保存套餐
                </el-button>
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
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
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
      const plansWithItems = await Promise.all(response.data.plans.map(async (plan: any) => {
        // 获取此套餐的所有检查项目
        const planItemsResponse = await axios.get(`/api/userview/plan?plan_id=${plan.ID || plan.id}`, {
          headers: { Authorization: `${token}` }
        })
        
        // 提取项目名称并组合成逗号分隔的字符串
        let itemsString = ''
        if (planItemsResponse.data && planItemsResponse.data.plan_items) {
          itemsString = planItemsResponse.data.plan_items
            .map((item: any) => item.item_name)
            .join(', ')
        }
        
        return {
          id: plan.ID || plan.id,
          name: plan.name || plan.plan_name || '',
          description: plan.description || '',
          suitableFor: plan.suitable_for || '',
          items: itemsString || plan.items || '',
          price: plan.price || plan.plan_price || 0
        }
      }))
      
      packages.value = plansWithItems
      
      // 如果有现有套餐，初始化状态
      if (packages.value.length > 0) {
        isEditing.value = false // 初始状态不显示保存按钮
      }
    } else {
      packages.value = []
    }
  } catch (error: any) {
    console.error('Failed to fetch packages:', error)
    ElMessage.error('获取套餐信息失败')
  }
}

const addNewPackage = () => {
  const newIndex = packages.value.length
  packages.value.push({
    name: '新套餐',  // 默认名称，用户需要修改
    description: '', 
    suitableFor: '所有人群',
    items: '体检项目1',  // 默认添加一个项目，提示用户该字段必填
    price: 0
  })
  // 自动展开新添加的套餐
  activePackages.value = [...activePackages.value, newIndex]
  // 显示保存按钮
  isEditing.value = true
  
  // 提示用户填写必要信息
  ElMessage({
    message: '请完善套餐信息，带*号的是必填项',
    type: 'info',
    duration: 3000
  })
}

const removePackage = async (index: number) => {
  const pkg = packages.value[index]
  if (pkg.id) {
    try {
      // Display a confirmation dialog before deleting
      await ElMessageBox.confirm(
        `确定要删除套餐 "${pkg.name}" 吗？此操作不可恢复。`,
        '确认删除',
        {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }
      )

      const token = localStorage.getItem('jwt')
      await axios.delete(`/api/institutions/plan`, {
        headers: { Authorization: `${token}` },
        data: { plan_id: pkg.id },
      })
      ElMessage.success(`套餐 "${pkg.name}" 删除成功`)
      packages.value.splice(index, 1)
      // 移除被删除套餐的激活状态
      activePackages.value = activePackages.value.filter(i => i !== index)
      // 更新索引大于被删除套餐的激活状态
      activePackages.value = activePackages.value.map(i => i > index ? i - 1 : i)
      // If no packages are left, set editing to false
      if (packages.value.length === 0) {
        isEditing.value = false
      }
    } catch (error: any) {
      if (error === 'cancel') {
        ElMessage.info('已取消删除')
      } else {
        console.error('Failed to delete package:', error)
        ElMessage.error(error.response?.data?.error || `删除套餐 "${pkg.name}" 失败`)
      }
    }
  } else {
    // If the package doesn't have an ID, it's a new package that hasn't been saved yet.
    // Just remove it from the local list.
    packages.value.splice(index, 1)
    // 移除被删除套餐的激活状态
    activePackages.value = activePackages.value.filter(i => i !== index)
    // 更新索引大于被删除套餐的激活状态
    activePackages.value = activePackages.value.map(i => i > index ? i - 1 : i)
    if (packages.value.length === 0) {
      isEditing.value = false
    }
  }
}

const savePackages = async () => {
  resetAlerts()
  
  // 验证套餐信息
  for (let i = 0; i < packages.value.length; i++) {
    const pkg = packages.value[i]
    // 确保处理空值情况
    const name = pkg.name || ''
    const items = pkg.items || ''
    
    const trimmedName = name.trim()
    const trimmedItems = items.trim()
    
    const missingFields = []
    
    if (!trimmedName) {
      missingFields.push('套餐名称')
    }
    
    if (!trimmedItems) {
      missingFields.push('检查项目')
    }
    
    if (missingFields.length > 0) {
      // 展开对应的套餐面板，方便用户查看和修改
      activePackages.value = [i]
      
      ElMessage({
        message: `第 ${i+1} 个套餐信息不完整，请填写：${missingFields.join('、')}`,
        type: 'warning',
        duration: 5000
      })
      return
    }
  }
  
  saving.value = true
  try {
    const token = localStorage.getItem('jwt')
    const currentPackage = packages.value[0]
    
    // 检查是否有ID，如果有则是更新操作
    if (currentPackage.id) {
      // 首先更新套餐基本信息
      await axios.patch(`/api/institutions/${props.institutionId}/item`, {
        plan_id: currentPackage.id,
        plan_name: currentPackage.name,
        plan_price: currentPackage.price,
        description: currentPackage.description,
        suitable_for: currentPackage.suitableFor
      }, {
        headers: { Authorization: `${token}` }
      })
      
      // 删除套餐原有检查项目并重新添加
      // 1. 先获取所有现有项目
      const planItemsResponse = await axios.get(`/api/userview/plan?plan_id=${currentPackage.id}`, {
        headers: { Authorization: `${token}` }
      })
      
      // 2. 删除所有现有项目
      if (planItemsResponse.data && planItemsResponse.data.plan_items) {
        for (const item of planItemsResponse.data.plan_items) {
          try {
            await axios.delete(`/api/institutions/plan/item`, {
              headers: { Authorization: `${token}` },
              data: {
                plan_id: currentPackage.id,
                item_id: item.item_id
              }
            })
          } catch (error: any) {
            console.error('无法删除原有项目:', error)
          }
        }
      }
      
      // 3. 添加新的项目
      const items = currentPackage.items.split(',').map(item => item.trim()).filter(item => item !== '')
      for (const item of items) {
        try {
          console.log(`正在添加项目 "${item}" 到套餐 #${currentPackage.id}`)
          console.log(`项目名称编码: ${encodeURIComponent(item)}`)
          await axios.post(`/api/institutions/${props.institutionId}/${currentPackage.id}/item`, {
            health_item: item.trim(),
            item_description: currentPackage.description || ''
          }, {
            headers: { 
              Authorization: `${token}`,
              'Content-Type': 'application/json; charset=utf-8'
            }
          })
        } catch (error: any) {
          console.error('添加项目失败:', error, '项目名称:', item)
          console.error('错误详情:', error.response?.data || '无详细信息')
          console.error('请求URL:', `/api/institutions/${props.institutionId}/${currentPackage.id}/item`)
          console.error('请求数据:', { health_item: item, item_description: currentPackage.description || '' })
          // 显示错误消息给用户
          ElMessage.error(`添加项目 "${item}" 失败: ${error.response?.data?.error || '未知错误'}`)
        }
      }
      
      successMessage.value = '套餐信息更新成功'
    } else {
      // 创建新套餐
      const firstItem = currentPackage.items.split(',')[0].trim()
      console.log(`创建套餐，第一个项目: "${firstItem}"，编码: ${encodeURIComponent(firstItem)}`)
      
      await axios.post(`/api/institutions/${props.institutionId}/plans`, {
        plan_name: currentPackage.name,
        health_item: firstItem,
        item_description: currentPackage.description,
        plan_price: currentPackage.price,
        description: currentPackage.description,
        suitable_for: currentPackage.suitableFor
      }, {
        headers: { 
          Authorization: `${token}`,
          'Content-Type': 'application/json; charset=utf-8'
        }
      })
      
      // 如果有多个体检项目，为套餐添加剩余项目
      if (currentPackage.items.split(',').length > 1) {
        const planResponse = await axios.get(`/api/institutions/${props.institutionId}/plans`, {
          headers: { Authorization: `${token}` }
        })
        
        const planId = planResponse.data.plans[0].ID
        
        for (let i = 1; i < currentPackage.items.split(',').length; i++) {
          const itemName = currentPackage.items.split(',')[i].trim()
          console.log(`添加附加项目: "${itemName}", 编码: ${encodeURIComponent(itemName)}`)
          
          await axios.post(`/api/institutions/${props.institutionId}/${planId}/item`, {
            health_item: itemName,
            item_description: currentPackage.description
          }, {
            headers: { 
              Authorization: `${token}`,
              'Content-Type': 'application/json; charset=utf-8'
            }
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

// 查看套餐中的健康检查项目
const viewHealthItems = (pkg: Package) => {
  if (!pkg.id) {
    ElMessage.warning('请先保存套餐信息，然后才能查看健康检查项目')
    return
  }
  
  // 跳转到健康检查项目管理页面
  const routeData = router.resolve({ 
    name: 'health-items',
    query: { plan_id: pkg.id.toString() }
  })
  window.open(routeData.href, '_blank')
}

// 保存单个套餐
const savePackage = async (index: number) => {
  resetAlerts()
  const pkg = packages.value[index]
  
  // 验证套餐信息
  const name = pkg.name || ''
  const items = pkg.items || ''
  
  const trimmedName = name.trim()
  const trimmedItems = items.trim()
  
  const missingFields = []
  
  if (!trimmedName) {
    missingFields.push('套餐名称')
  }
  
  if (!trimmedItems) {
    missingFields.push('检查项目')
  }
  
  if (missingFields.length > 0) {
    ElMessage({
      message: `套餐信息不完整，请填写：${missingFields.join('、')}`,
      type: 'warning',
      duration: 5000
    })
    return
  }
  
  saving.value = true
  try {
    const token = localStorage.getItem('jwt')
    
    // 检查是否有ID，如果有则是更新操作
    if (pkg.id) {
      // 首先更新套餐基本信息
      await axios.patch(`/api/institutions/${props.institutionId}/item`, {
        plan_id: pkg.id,
        plan_name: pkg.name,
        plan_price: pkg.price,
        description: pkg.description,
        suitable_for: pkg.suitableFor
      }, {
        headers: { Authorization: `${token}` }
      })
      
      // 删除套餐原有检查项目并重新添加
      // 1. 先获取所有现有项目
      const planItemsResponse = await axios.get(`/api/userview/plan?plan_id=${pkg.id}`, {
        headers: { Authorization: `${token}` }
      })
      
      // 2. 删除所有现有项目
      if (planItemsResponse.data && planItemsResponse.data.plan_items) {
        for (const item of planItemsResponse.data.plan_items) {
          try {
            await axios.delete(`/api/institutions/plan/item`, {
              headers: { Authorization: `${token}` },
              data: {
                plan_id: pkg.id,
                item_id: item.item_id
              }
            })
          } catch (error: any) {
            console.error('无法删除原有项目:', error)
          }
        }
      }
      
      // 3. 添加新的项目
      const items = pkg.items.split(',').map(item => item.trim()).filter(item => item !== '')
      for (const item of items) {
        try {
          console.log(`正在添加项目 "${item}" 到套餐 #${pkg.id}`)
          await axios.post(`/api/institutions/${props.institutionId}/${pkg.id}/item`, {
            health_item: item.trim(),
            item_description: pkg.description || ''
          }, {
            headers: { 
              Authorization: `${token}`,
              'Content-Type': 'application/json; charset=utf-8'
            }
          })
        } catch (error: any) {
          console.error('添加项目失败:', error, '项目名称:', item)
          console.error('错误详情:', error.response?.data || '无详细信息')
          ElMessage.error(`添加项目 "${item}" 失败: ${error.response?.data?.error || '未知错误'}`)
        }
      }
      
      successMessage.value = `套餐 "${pkg.name}" 更新成功`
    } else {
      // 创建新套餐
      const firstItem = pkg.items.split(',')[0].trim()
      
      await axios.post(`/api/institutions/${props.institutionId}/plans`, {
        plan_name: pkg.name,
        health_item: firstItem,
        item_description: pkg.description,
        plan_price: pkg.price,
        description: pkg.description,
        suitable_for: pkg.suitableFor
      }, {
        headers: { 
          Authorization: `${token}`,
          'Content-Type': 'application/json; charset=utf-8'
        }
      })
      
      // 获取新创建的套餐ID
      const planResponse = await axios.get(`/api/institutions/${props.institutionId}/plans`, {
        headers: { Authorization: `${token}` }
      })
      
      // 查找匹配的套餐
      let newPlanId = null
      if (planResponse.data.plans && planResponse.data.plans.length > 0) {
        const matchingPlan = planResponse.data.plans.find((plan: any) => 
          plan.name === pkg.name || plan.plan_name === pkg.name
        )
        if (matchingPlan) {
          newPlanId = matchingPlan.ID || matchingPlan.id
        } else {
          // 如果找不到匹配的，使用最后创建的套餐ID
          newPlanId = planResponse.data.plans[0].ID || planResponse.data.plans[0].id
        }
      }
      
      // 如果有多个体检项目，为套餐添加剩余项目
      if (newPlanId && pkg.items.split(',').length > 1) {
        for (let i = 1; i < pkg.items.split(',').length; i++) {
          const itemName = pkg.items.split(',')[i].trim()
          if (itemName) {
            await axios.post(`/api/institutions/${props.institutionId}/${newPlanId}/item`, {
              health_item: itemName,
              item_description: pkg.description
            }, {
              headers: { 
                Authorization: `${token}`,
                'Content-Type': 'application/json; charset=utf-8'
              }
            })
          }
        }
      }
      
      successMessage.value = `套餐 "${pkg.name}" 创建成功`
    }
    
    saveSuccess.value = true
    
    // 刷新套餐数据
    await fetchPackages()
    
  } catch (error: any) {
    console.error('Failed to save package:', error)
    saveError.value = true
    errorMessage.value = error.response?.data?.error || `保存套餐 "${pkg.name}" 失败`
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

.health-items-container {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.health-items-controls {
  display: flex;
  justify-content: flex-end;
}

/* 为必填字段标签添加样式 */
:deep(.el-form-item.is-required .el-form-item__label:before) {
  content: '*';
  color: #f56c6c;
  margin-right: 4px;
}

/* 为校验失败的输入框添加提示样式 */
:deep(.el-form-item.is-error .el-input__inner),
:deep(.el-form-item.is-error .el-textarea__inner) {
  border-color: #f56c6c;
}
</style>