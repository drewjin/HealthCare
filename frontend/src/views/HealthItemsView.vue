<template>
  <div class="health-items-container">
    <el-page-header @back="goBack" title="返回">
      <template #content>
        <span class="page-title">{{ planId ? `管理套餐#${planId}的健康检查项目` : '管理健康检查项目' }}</span>
      </template>
    </el-page-header>

    <div class="health-items-section">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>健康检查项目列表</span>
            <div class="button-group">
              <el-button type="success" @click="createNewItem">创建检查项目</el-button>
              <el-button type="primary" @click="refreshItems">刷新列表</el-button>
            </div>
          </div>
        </template>

        <el-table v-loading="loading" :data="healthItems" style="width: 100%">
          <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="ItemName" label="项目名称" />
          <el-table-column label="创建时间" width="180">
            <template #default="scope">
              {{ formatDate(scope.row.CreatedAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="300">
            <template #default="scope">
              <el-button size="small" @click="viewItem(scope.row)">查看详情</el-button>
              <el-button size="small" type="primary" @click="editItem(scope.row)">编辑</el-button>
              <el-button size="small" type="success" @click="viewParsedItem(scope.row)">解析查看</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 项目详情对话框 -->
    <el-dialog v-model="itemDetailsVisible" title="项目详情" width="600px">
      <div v-if="selectedItem">
        <h3>{{ selectedItem.ItemName }}</h3>
        <p><strong>项目ID:</strong> {{ selectedItem.ID }}</p>
        <p><strong>创建时间:</strong> {{ formatDate(selectedItem.CreatedAt) }}</p>
        <p><strong>更新时间:</strong> {{ formatDate(selectedItem.UpdatedAt) }}</p>

        <h4 class="mt-4">关联套餐:</h4>
        <el-table v-if="selectedItemPlans && selectedItemPlans.length > 0" :data="selectedItemPlans" style="width: 100%">
          <el-table-column prop="PlanName" label="套餐名称" />
          <el-table-column prop="ItemDescription" label="在套餐中的描述" />
          <el-table-column prop="PlanPrice" label="套餐价格" width="120">
            <template #default="scope">
              ¥{{ scope.row.PlanPrice.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="120">
            <template #default="scope">
              <el-button size="small" type="primary" @click="editPlanItemDescription(scope.row)">
                编辑描述
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="empty-data">该项目尚未关联任何套餐</div>
      </div>
    </el-dialog>

    <!-- 解析查看对话框 -->
    <el-dialog v-model="parseViewVisible" title="健康项目解析视图" width="600px">
      <div v-if="selectedItem">
        <health-item-string-view 
          :item-id="selectedItem.ID" 
          :auto-fetch="true" 
          ref="itemStringView"
        />
      </div>
    </el-dialog>

    <!-- 编辑项目对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑健康检查项目" width="500px">
      <el-form v-if="selectedItem" :model="editForm" label-width="100px">
        <el-form-item label="项目名称">
          <el-input v-model="editForm.ItemName" placeholder="请输入项目名称" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="saving" @click="saveItemChanges">保存</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 编辑套餐中项目描述对话框 -->
    <el-dialog v-model="editPlanItemDialogVisible" title="编辑套餐中的项目描述" width="500px">
      <el-form v-if="selectedPlanItem" :model="planItemEditForm" label-width="100px">
        <el-form-item label="套餐名称">
          <el-input v-model="selectedPlanItem.PlanName" disabled />
        </el-form-item>
        <el-form-item label="项目名称">
          <el-input v-model="selectedItem.ItemName" disabled />  
        </el-form-item>
        <el-form-item label="项目描述">
          <el-input
            v-model="planItemEditForm.ItemDescription"
            type="textarea"
            rows="3"
            placeholder="请输入该套餐中项目的描述"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editPlanItemDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="saving" @click="savePlanItemChanges">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import HealthItemStringView from './HealthItemStringView.vue'

interface HealthItem {
  ID: number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
  ItemName: string
}

interface PlanInfo {
  PlanID: number
  PlanName: string
  ItemDescription: string
  PlanPrice: number
}

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const saving = ref(false)
const healthItems = ref<HealthItem[]>([])
const selectedItem = ref<HealthItem | null>(null)
const selectedItemPlans = ref<PlanInfo[]>([])
const selectedPlanItem = ref<PlanInfo | null>(null)
const itemDetailsVisible = ref(false)
const editDialogVisible = ref(false)
const editPlanItemDialogVisible = ref(false)
const parseViewVisible = ref(false)
const planId = ref<number | null>(null) // 从URL获取的套餐ID
const itemStringView = ref<InstanceType<typeof HealthItemStringView> | null>(null)

const editForm = reactive({
  ItemName: ''
})

const planItemEditForm = reactive({
  ItemDescription: ''
})

const goBack = () => {
  router.back()
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString()
}  // 获取所有健康检查项目
const fetchHealthItems = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get('/api/healthitems', {
      headers: { Authorization: `${token}` }
    })
    
    // 处理返回的健康项目数据 - 可能在items或health_items字段中
    const items = response.data.items || response.data.health_items || []
    
    // 标准化数据结构
    healthItems.value = items.map((item: any) => ({
      ID: item.ID || item.id,
      ItemName: item.ItemName || item.item_name,
      CreatedAt: item.CreatedAt || item.created_at || '',
      UpdatedAt: item.UpdatedAt || item.updated_at || '',
      DeletedAt: item.DeletedAt || item.deleted_at || null
    }))
    
    console.log('Fetched health items:', healthItems.value)
  } catch (error) {
    console.error('Failed to fetch health items:', error)
    ElMessage.error('获取健康检查项目列表失败')
  } finally {
    loading.value = false
  }
}

// 查看项目详情
const viewItem = async (item: HealthItem) => {
  selectedItem.value = item
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/healthitems/${item.ID}`, {
      headers: { Authorization: `${token}` }
    })
    selectedItemPlans.value = response.data.plans
    itemDetailsVisible.value = true
  } catch (error) {
    console.error('Failed to fetch item details:', error)
    ElMessage.error('获取项目详情失败')
  }
}

// 编辑项目
const editItem = (item: HealthItem) => {
  selectedItem.value = item
  editForm.ItemName = item.ItemName
  editDialogVisible.value = true
}

// 保存项目修改
const saveItemChanges = async () => {
  if (!selectedItem.value) return
  
  if (!editForm.ItemName.trim()) {
    ElMessage.warning('项目名称不能为空')
    return
  }
  
  saving.value = true
  try {
    const token = localStorage.getItem('jwt')
    await axios.patch(`/api/healthitems/${selectedItem.value.ID}`, {
      item_name: editForm.ItemName
    }, {
      headers: { Authorization: `${token}` }
    })
    
    ElMessage.success('项目更新成功')
    editDialogVisible.value = false
    
    // 更新列表
    await fetchHealthItems()
  } catch (error: any) {
    console.error('Failed to update health item:', error)
    
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('更新项目失败')
    }
  } finally {
    saving.value = false
  }
}

// 编辑套餐中的项目描述
const editPlanItemDescription = (planItem: PlanInfo) => {
  selectedPlanItem.value = planItem
  planItemEditForm.ItemDescription = planItem.ItemDescription
  editPlanItemDialogVisible.value = true
}

// 保存套餐项目描述修改
const savePlanItemChanges = async () => {
  if (!selectedPlanItem.value || !selectedItem.value) return
  
  saving.value = true
  try {
    const token = localStorage.getItem('jwt')
    await axios.patch('/api/healthitems/plan-item', {
      plan_id: selectedPlanItem.value.PlanID,
      item_id: selectedItem.value.ID,
      item_description: planItemEditForm.ItemDescription
    }, {
      headers: { Authorization: `${token}` }
    })
    
    ElMessage.success('项目描述更新成功')
    editPlanItemDialogVisible.value = false
    
    // 刷新项目详情
    if (selectedItem.value) {
      await viewItem(selectedItem.value)
    }
  } catch (error) {
    console.error('Failed to update plan item description:', error)
    ElMessage.error('更新项目描述失败')
  } finally {
    saving.value = false
  }
}

// 查看解析后的项目
const viewParsedItem = (item: HealthItem) => {
  selectedItem.value = item
  parseViewVisible.value = true
  // 确保在下一个更新周期刷新，防止在对话框未完全打开时调用方法
  setTimeout(() => {
    if (itemStringView.value) {
      itemStringView.value.refresh()
    }
  }, 100)
}

// 刷新列表
const refreshItems = () => {
  fetchHealthItems()
}

// 创建新的检查项目
const createNewItem = () => {
  router.push('/health-item-manager')
}

onMounted(() => {
  // 检查URL参数中是否有plan_id
  const planIdParam = route.query.plan_id as string
  if (planIdParam) {
    planId.value = parseInt(planIdParam)
    // 只加载指定套餐的健康项目
    fetchHealthItemsByPlanId(planId.value)
  } else {
    // 加载所有健康项目
    fetchHealthItems()
  }
})

// 获取指定套餐的健康检查项目
const fetchHealthItemsByPlanId = async (id: number) => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/userview/plan?plan_id=${id}`, {
      headers: { Authorization: `${token}` }
    })
    
    if (response.data && response.data.plan_items && Array.isArray(response.data.plan_items)) {
      // 转换为健康项目格式
      const planItems = response.data.plan_items.map((item: any) => {
        return {
          ID: item.item_id || 0,
          ItemName: item.item_name || '',
          CreatedAt: '',  // 这些字段可能在API中没有
          UpdatedAt: '',
          DeletedAt: null
        } as HealthItem
      })
      
      healthItems.value = planItems
    } else {
      healthItems.value = []
      ElMessage.warning('未找到该套餐的健康检查项目')
    }
  } catch (error) {
    console.error('Failed to fetch health items for plan:', error)
    ElMessage.error('获取套餐项目失败')
    healthItems.value = []
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.health-items-container {
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: bold;
}

.health-items-section {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.button-group {
  display: flex;
  gap: 10px;
}

.empty-data {
  padding: 20px;
  text-align: center;
  color: #909399;
}

.mt-4 {
  margin-top: 16px;
}
</style>
