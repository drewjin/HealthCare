<template>
  <div class="plan-health-items-container">
    <h2>套餐健康项目管理</h2>
    
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="10" animated />
    </div>
    
    <template v-else>
      <div class="plan-info">
        <el-descriptions title="套餐信息" border>
          <el-descriptions-item label="套餐ID">{{ planId }}</el-descriptions-item>
          <el-descriptions-item label="套餐名称">{{ planDetails.name }}</el-descriptions-item>
          <el-descriptions-item label="机构名称">{{ planDetails.institutionName }}</el-descriptions-item>
        </el-descriptions>
      </div>
      
      <div class="health-items-section">
        <div class="section-header">
          <h3>健康检查项目</h3>
          <el-button type="primary" @click="showAddItemDialog">添加检查项目</el-button>
        </div>
        
        <el-table v-if="planItems.length > 0" :data="planItems" style="width: 100%" border>
          <el-table-column prop="itemName" label="项目名称" />
          <el-table-column prop="itemDescription" label="项目描述" />
          <el-table-column label="操作" width="220">
            <template #default="scope">
              <el-button size="small" type="primary" @click="editItem(scope.row)">
                编辑描述
              </el-button>
              <el-button size="small" type="danger" @click="removeItem(scope.row)">
                移除
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <div v-else class="empty-state">
          <el-empty description="暂无检查项目" />
        </div>
      </div>
      
      <!-- 添加项目对话框 -->
      <el-dialog v-model="addItemDialogVisible" title="添加健康检查项目" width="500px">
        <el-form :model="addItemForm" label-width="120px">
          <el-form-item label="项目搜索">
            <el-input v-model="searchKeyword" placeholder="输入关键词搜索健康项目" @input="searchHealthItems" />
          </el-form-item>
          
          <el-table 
            v-if="searchResults.length > 0" 
            :data="searchResults" 
            style="width: 100%" 
            border 
            height="300px"
            @selection-change="handleSelectionChange"
          >
            <el-table-column type="selection" width="55" />
            <el-table-column prop="ID" label="ID" width="80" />
            <el-table-column prop="ItemName" label="项目名称" />
          </el-table>
          
          <el-form-item v-if="selectedItems.length > 0" label="项目描述">
            <el-input 
              v-model="addItemForm.description" 
              type="textarea" 
              :rows="3" 
              placeholder="请输入对所选项目的描述（可选）"
            ></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="addItemDialogVisible = false">取消</el-button>
            <el-button 
              type="primary" 
              :disabled="selectedItems.length === 0"
              :loading="submitting"
              @click="addSelectedItemsToPlan"
            >添加所选项目</el-button>
          </span>
        </template>
      </el-dialog>
      
      <!-- 编辑项目描述对话框 -->
      <el-dialog v-model="editItemDialogVisible" title="编辑项目描述" width="500px">
        <el-form v-if="currentItem" :model="editItemForm" label-width="120px">
          <el-form-item label="项目名称">
            <el-input v-model="currentItem.itemName" disabled />
          </el-form-item>
          <el-form-item label="项目描述">
            <el-input 
              v-model="editItemForm.description" 
              type="textarea" 
              :rows="3" 
              placeholder="请输入项目描述"
            ></el-input>
          </el-form-item>
        </el-form>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="editItemDialogVisible = false">取消</el-button>
            <el-button 
              type="primary"
              :loading="submitting"
              @click="updateItemDescription"
            >保存</el-button>
          </span>
        </template>
      </el-dialog>
      
      <!-- 确认移除对话框 -->
      <el-dialog v-model="removeConfirmVisible" title="确认移除" width="400px">
        <p>确定要从套餐中移除项目 "{{ currentItem?.itemName }}" 吗？</p>
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
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()

// 套餐ID
const planId = computed(() => route.params.id || route.query.id)

// 加载状态
const loading = ref(true)
const submitting = ref(false)

// 套餐详情
const planDetails = reactive({
  name: '',
  institutionName: '',
  institutionId: 0,
  price: 0,
  description: ''
})

// 套餐关联的健康项目
interface PlanItem {
  id: number
  healthItemId: number
  itemName: string
  itemDescription: string
}

const planItems = ref<PlanItem[]>([])

// 对话框状态
const addItemDialogVisible = ref(false)
const editItemDialogVisible = ref(false)
const removeConfirmVisible = ref(false)

// 当前操作的项目
const currentItem = ref<PlanItem | null>(null)

// 添加项目表单
const addItemForm = reactive({
  description: ''
})

// 编辑项目表单
const editItemForm = reactive({
  description: ''
})

// 搜索相关
const searchKeyword = ref('')
const searchResults = ref<any[]>([])
const selectedItems = ref<any[]>([])

// 初始化加载
onMounted(async () => {
  await fetchPlanItems()
})

// 获取套餐健康项目
const fetchPlanItems = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('您需要登录才能访问此页面')
      router.push('/login')
      return
    }
    
    // 获取套餐详情
    const planResponse = await axios.get(`/api/plans/${planId.value}`, {
      headers: { Authorization: token }
    })
    
    if (planResponse.data) {
      planDetails.name = planResponse.data.plan_name || ''
      planDetails.price = planResponse.data.plan_price || 0
      planDetails.description = planResponse.data.description || ''
      
      // 获取机构信息
      const institutionId = planResponse.data.institution_id
      planDetails.institutionId = institutionId
      
      const institutionResponse = await axios.get(`/api/institutions/${institutionId}`, {
        headers: { Authorization: token }
      })
      
      if (institutionResponse.data && institutionResponse.data.institution) {
        planDetails.institutionName = institutionResponse.data.institution.institution_name || ''
      }
    }
    
    // 获取套餐项目
    const itemsResponse = await axios.get(`/api/healthitem-plan/plans/${planId.value}/items`, {
      headers: { Authorization: token }
    })
    
    if (itemsResponse.data && Array.isArray(itemsResponse.data.items)) {
      planItems.value = itemsResponse.data.items.map((item: any) => ({
        id: item.id,
        healthItemId: item.health_item_id,
        itemName: item.item_name,
        itemDescription: item.item_description || ''
      }))
    }
  } catch (error: any) {
    console.error('获取套餐健康项目失败:', error)
    ElMessage.error('获取套餐健康项目失败: ' + (error.response?.data?.error || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 显示添加项目对话框
const showAddItemDialog = () => {
  searchKeyword.value = ''
  searchResults.value = []
  selectedItems.value = []
  addItemForm.description = ''
  addItemDialogVisible.value = true
}

// 搜索健康项目
const searchHealthItems = async () => {
  if (!searchKeyword.value.trim()) {
    searchResults.value = []
    return
  }
  
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get('/api/healthitems', {
      headers: { Authorization: token },
      params: { keyword: searchKeyword.value }
    })
    
    if (response.data && Array.isArray(response.data.items)) {
      // 过滤掉已经在套餐中的项目
      const existingItemIds = planItems.value.map(item => item.healthItemId)
      searchResults.value = response.data.items.filter((item: any) => 
        !existingItemIds.includes(item.ID)
      )
    }
  } catch (error) {
    console.error('搜索健康项目失败:', error)
  }
}

// 处理多选
const handleSelectionChange = (selection: any[]) => {
  selectedItems.value = selection
}

// 添加所选项目到套餐
const addSelectedItemsToPlan = async () => {
  if (selectedItems.value.length === 0) {
    ElMessage.warning('请至少选择一项检查项目')
    return
  }
  
  submitting.value = true
  try {
    const token = localStorage.getItem('jwt')
    const itemsToAdd = selectedItems.value.map(item => ({
      health_item_id: item.ID,
      item_description: addItemForm.description
    }))
    
    await axios.post(`/api/healthitem-plan/plans/${planId.value}/items`, itemsToAdd, {
      headers: { Authorization: token }
    })
    
    ElMessage.success('成功添加健康检查项目')
    addItemDialogVisible.value = false
    
    // 重新加载套餐项目
    await fetchPlanItems()
  } catch (error: any) {
    console.error('添加健康检查项目失败:', error)
    ElMessage.error('添加健康检查项目失败: ' + (error.response?.data?.error || '未知错误'))
  } finally {
    submitting.value = false
  }
}

// 编辑项目
const editItem = (item: PlanItem) => {
  currentItem.value = item
  editItemForm.description = item.itemDescription
  editItemDialogVisible.value = true
}

// 更新项目描述
const updateItemDescription = async () => {
  if (!currentItem.value) return
  
  submitting.value = true
  try {
    const token = localStorage.getItem('jwt')
    await axios.patch(`/api/healthitems/plan-item`, {
      plan_id: planId.value,
      item_id: currentItem.value.healthItemId,
      item_description: editItemForm.description
    }, {
      headers: { Authorization: token }
    })
    
    ElMessage.success('项目描述更新成功')
    editItemDialogVisible.value = false
    
    // 更新本地数据
    const index = planItems.value.findIndex(item => item.id === currentItem.value?.id)
    if (index !== -1) {
      planItems.value[index].itemDescription = editItemForm.description
    }
  } catch (error: any) {
    console.error('更新项目描述失败:', error)
    ElMessage.error('更新项目描述失败: ' + (error.response?.data?.error || '未知错误'))
  } finally {
    submitting.value = false
  }
}

// 移除项目
const removeItem = (item: PlanItem) => {
  currentItem.value = item
  removeConfirmVisible.value = true
}

// 确认移除项目
const confirmRemoveItem = async () => {
  if (!currentItem.value) return
  
  submitting.value = true
  try {
    const token = localStorage.getItem('jwt')
    await axios.delete(`/api/institutions/plan/item`, {
      headers: { Authorization: token },
      data: {
        plan_id: planId.value,
        item_id: currentItem.value.healthItemId
      }
    })
    
    ElMessage.success('项目已从套餐中移除')
    removeConfirmVisible.value = false
    
    // 从本地移除
    planItems.value = planItems.value.filter(item => item.id !== currentItem.value?.id)
  } catch (error: any) {
    console.error('移除项目失败:', error)
    ElMessage.error('移除项目失败: ' + (error.response?.data?.error || '未知错误'))
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.plan-health-items-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.loading-state {
  padding: 20px;
}

.plan-info {
  margin-bottom: 30px;
}

.health-items-section {
  margin-top: 30px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.empty-state {
  margin: 40px 0;
  display: flex;
  justify-content: center;
}

.dialog-footer {
  margin-top: 20px;
  text-align: right;
}
</style>
