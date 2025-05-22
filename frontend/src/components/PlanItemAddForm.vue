<template>
  <div class="plan-item-add-form">
    <el-form :model="form" label-width="140px">
      <el-form-item label="套餐">
        <el-select 
          v-model="form.planId" 
          placeholder="请选择套餐" 
          filterable
          @change="onPlanChange"
        >
          <el-option 
            v-for="plan in plans" 
            :key="plan.id" 
            :label="plan.name" 
            :value="plan.id" 
          />
        </el-select>
      </el-form-item>
      
      <div v-if="form.planId" class="items-section">
        <el-divider content-position="left">关联健康项目</el-divider>
        
        <div class="search-bar">
          <el-input 
            v-model="searchKeyword" 
            placeholder="搜索健康项目" 
            clearable
            @input="searchHealthItems"
          >
            <template #append>
              <el-button :icon="Search" @click="searchHealthItems"></el-button>
            </template>
          </el-input>
        </div>
        
        <div v-if="loading" class="loading-state">
          <el-skeleton :rows="5" animated />
        </div>
        
        <div v-else-if="availableItems.length === 0" class="empty-state">
          <el-empty description="未找到可添加的健康项目" />
        </div>
        
        <el-table 
          v-else
          :data="availableItems" 
          style="width: 100%" 
          border
          height="300px"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="55" />
          <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="ItemName" label="项目名称" show-overflow-tooltip />
          <el-table-column label="操作" width="120">
            <template #default="scope">
              <el-button 
                type="primary" 
                link
                @click="addSingleItem(scope.row)"
              >
                添加
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div v-if="selectedItems.length > 0" class="batch-actions">
          <el-form-item label="项目描述">
            <el-input 
              v-model="form.itemDescription" 
              type="textarea" 
              :rows="3" 
              placeholder="为所选项目添加描述（可选）"
            />
          </el-form-item>
          
          <el-form-item>
            <el-button 
              type="primary" 
              :loading="submitting"
              @click="addSelectedItems"
            >
              批量添加 ({{ selectedItems.length }})
            </el-button>
          </el-form-item>
        </div>
      </div>
    </el-form>
    
    <!-- 当前套餐的健康项目 -->
    <div v-if="form.planId && planItems.length > 0" class="current-items">
      <el-divider content-position="left">当前健康项目</el-divider>
      
      <el-table 
        :data="planItems" 
        style="width: 100%" 
        border
      >
        <el-table-column prop="healthItemId" label="ID" width="80" />
        <el-table-column prop="itemName" label="项目名称" />
        <el-table-column prop="itemDescription" label="描述" show-overflow-tooltip />
        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button 
              type="danger" 
              link
              @click="removeItem(scope.row)"
            >
              移除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 添加单个项目对话框 -->
    <el-dialog v-model="singleAddDialogVisible" title="添加健康项目" width="500px">
      <el-form v-if="currentItem" :model="singleItemForm" label-width="120px">
        <el-form-item label="项目名称">
          <el-input v-model="currentItem.ItemName" disabled />
        </el-form-item>
        <el-form-item label="项目描述">
          <el-input 
            v-model="singleItemForm.description" 
            type="textarea" 
            :rows="3" 
            placeholder="请输入项目描述（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="singleAddDialogVisible = false">取消</el-button>
          <el-button 
            type="primary"
            :loading="submitting"
            @click="confirmAddSingleItem"
          >添加</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 确认移除对话框 -->
    <el-dialog v-model="removeConfirmVisible" title="确认移除" width="400px">
      <p>确定要从套餐中移除项目 "{{ currentItem?.ItemName }}" 吗？</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="removeConfirmVisible = false">取消</el-button>
          <el-button 
            type="danger"
            :loading="submitting"
            @click="confirmRemoveItem"
          >确认移除</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const props = defineProps({
  institutionId: {
    type: Number,
    default: 0
  }
})

const emit = defineEmits(['updated'])

// 状态变量
const loading = ref(false)
const submitting = ref(false)
const searchKeyword = ref('')
const plans = ref([])
const availableItems = ref([])
const selectedItems = ref([])
const planItems = ref([])
const singleAddDialogVisible = ref(false)
const removeConfirmVisible = ref(false)
const currentItem = ref(null)

// 表单数据
const form = reactive({
  planId: null,
  itemDescription: ''
})

// 单个项目表单
const singleItemForm = reactive({
  description: ''
})

// 初始化
onMounted(async () => {
  await fetchInstitutionPlans()
})

// 监听套餐变化
watch(() => form.planId, async (newVal) => {
  if (newVal) {
    await fetchPlanItems()
    await searchHealthItems()
  } else {
    planItems.value = []
    availableItems.value = []
  }
})

// 获取机构的套餐列表
const fetchInstitutionPlans = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    let url = '/api/plans'
    
    // 如果有指定机构ID，则只获取该机构的套餐
    if (props.institutionId) {
      url = `/api/institutions/${props.institutionId}/plans`
    }
    
    const response = await axios.get(url, {
      headers: { Authorization: token }
    })
    
    if (response.data && Array.isArray(response.data.plans)) {
      plans.value = response.data.plans.map(plan => ({
        id: plan.ID || plan.id,
        name: plan.PlanName || plan.plan_name,
        price: plan.PlanPrice || plan.plan_price,
        institutionId: plan.RelationInstitutionID || plan.institution_id
      }))
    }
  } catch (error) {
    console.error('获取套餐列表失败:', error)
    ElMessage.error('获取套餐列表失败')
  } finally {
    loading.value = false
  }
}

// 获取套餐的健康项目
const fetchPlanItems = async () => {
  if (!form.planId) return
  
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/healthitem-plan/plans/${form.planId}/items`, {
      headers: { Authorization: token }
    })
    
    if (response.data && Array.isArray(response.data.items)) {
      planItems.value = response.data.items.map(item => ({
        id: item.id,
        healthItemId: item.health_item_id,
        itemName: item.item_name,
        itemDescription: item.item_description || ''
      }))
    }
  } catch (error) {
    console.error('获取套餐健康项目失败:', error)
    ElMessage.error('获取套餐健康项目失败')
  } finally {
    loading.value = false
  }
}

// 搜索健康项目
const searchHealthItems = async () => {
  if (!form.planId) return
  
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get('/api/healthitems', {
      headers: { Authorization: token },
      params: { keyword: searchKeyword.value }
    })
    
    if (response.data && Array.isArray(response.data.items)) {
      // 过滤掉已经在套餐中的项目
      const existingItemIds = planItems.value.map(item => item.healthItemId)
      availableItems.value = response.data.items.filter(item => 
        !existingItemIds.includes(item.ID)
      )
    }
  } catch (error) {
    console.error('搜索健康项目失败:', error)
    ElMessage.error('搜索健康项目失败')
  } finally {
    loading.value = false
  }
}

// 处理套餐变更
const onPlanChange = () => {
  selectedItems.value = []
  form.itemDescription = ''
}

// 处理多选
const handleSelectionChange = (selection) => {
  selectedItems.value = selection
}

// 添加单个项目（打开对话框）
const addSingleItem = (item) => {
  currentItem.value = item
  singleItemForm.description = ''
  singleAddDialogVisible.value = true
}

// 确认添加单个项目
const confirmAddSingleItem = async () => {
  if (!currentItem.value || !form.planId) return
  
  submitting.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.post(`/api/healthitem-plan/plans/${form.planId}/item`, {
      health_item_id: currentItem.value.ID,
      item_description: singleItemForm.description
    }, {
      headers: { Authorization: token }
    })
    
    ElMessage.success('成功添加健康项目')
    singleAddDialogVisible.value = false
    
    // 刷新数据
    await fetchPlanItems()
    await searchHealthItems()
    emit('updated')
  } catch (error) {
    console.error('添加健康项目失败:', error)
    ElMessage.error('添加健康项目失败: ' + (error.response?.data?.error || '未知错误'))
  } finally {
    submitting.value = false
  }
}

// 批量添加选中的项目
const addSelectedItems = async () => {
  if (selectedItems.value.length === 0 || !form.planId) {
    ElMessage.warning('请至少选择一个健康项目')
    return
  }
  
  submitting.value = true
  try {
    const token = localStorage.getItem('jwt')
    const itemsToAdd = selectedItems.value.map(item => ({
      health_item_id: item.ID,
      item_description: form.itemDescription
    }))
    
    await axios.post(`/api/healthitem-plan/plans/${form.planId}/items`, itemsToAdd, {
      headers: { Authorization: token }
    })
    
    ElMessage.success('成功添加健康检查项目')
    selectedItems.value = []
    form.itemDescription = ''
    
    // 刷新数据
    await fetchPlanItems()
    await searchHealthItems()
    emit('updated')
  } catch (error) {
    console.error('添加健康检查项目失败:', error)
    ElMessage.error('添加健康检查项目失败: ' + (error.response?.data?.error || '未知错误'))
  } finally {
    submitting.value = false
  }
}

// 移除项目（打开确认对话框）
const removeItem = (item) => {
  currentItem.value = {
    ID: item.healthItemId,
    ItemName: item.itemName
  }
  removeConfirmVisible.value = true
}

// 确认移除项目
const confirmRemoveItem = async () => {
  if (!currentItem.value || !form.planId) return
  
  submitting.value = true
  try {
    const token = localStorage.getItem('jwt')
    await axios.delete(`/api/institutions/plan/item`, {
      headers: { Authorization: token },
      data: {
        plan_id: form.planId,
        item_id: currentItem.value.ID
      }
    })
    
    ElMessage.success('项目已从套餐中移除')
    removeConfirmVisible.value = false
    
    // 刷新数据
    await fetchPlanItems()
    await searchHealthItems()
    emit('updated')
  } catch (error) {
    console.error('移除项目失败:', error)
    ElMessage.error('移除项目失败: ' + (error.response?.data?.error || '未知错误'))
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.plan-item-add-form {
  margin-bottom: 20px;
}

.items-section {
  margin-top: 20px;
}

.search-bar {
  margin-bottom: 15px;
}

.batch-actions {
  margin-top: 15px;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.current-items {
  margin-top: 30px;
}

.loading-state, .empty-state {
  padding: 20px 0;
}
</style>
