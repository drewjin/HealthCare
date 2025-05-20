<template>
  <div class="plan-health-item-manager">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>套餐健康项目管理</span>
        </div>
      </template>
      
      <el-tabs v-model="activeTab">
        <el-tab-pane label="添加健康项目" name="add">
          <PlanItemAddForm 
            :institution-id="institutionId" 
            @updated="handleUpdated"
          />
        </el-tab-pane>
        <el-tab-pane label="套餐项目列表" name="list">
          <div class="plan-selector" v-if="plans.length > 0">
            <span class="label">选择套餐：</span>
            <el-select 
              v-model="selectedPlanId" 
              placeholder="请选择套餐" 
              filterable
              @change="loadPlanItems"
            >
              <el-option 
                v-for="plan in plans" 
                :key="plan.id" 
                :label="plan.name" 
                :value="plan.id" 
              />
            </el-select>
          </div>
          
          <div v-if="loading" class="loading-state">
            <el-skeleton :rows="10" animated />
          </div>
          
          <div v-else-if="!selectedPlanId" class="empty-state">
            <el-empty description="请选择一个套餐查看关联的健康项目" />
          </div>
          
          <div v-else-if="planItems.length === 0" class="empty-state">
            <el-empty description="该套餐暂无关联的健康项目" />
          </div>
          
          <el-table v-else :data="planItems" style="width: 100%" border>
            <el-table-column prop="healthItemId" label="ID" width="80" />
            <el-table-column prop="itemName" label="项目名称" />
            <el-table-column prop="itemDescription" label="描述" show-overflow-tooltip />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import PlanItemAddForm from '@/components/PlanItemAddForm.vue'

const route = useRoute()
const activeTab = ref('add')
const loading = ref(false)
const institutionId = ref(0)
const plans = ref([])
const selectedPlanId = ref(null)
const planItems = ref([])

// 初始化
onMounted(async () => {
  // 尝试从路由参数获取机构ID
  const routeInstitutionId = route.params.institutionId || route.query.institutionId
  if (routeInstitutionId) {
    institutionId.value = parseInt(routeInstitutionId as string)
  }
  
  await fetchInstitutionPlans()
})

// 获取机构的套餐列表
const fetchInstitutionPlans = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    let url = '/api/plans'
    
    // 如果有指定机构ID，则只获取该机构的套餐
    if (institutionId.value) {
      url = `/api/institutions/${institutionId.value}/plans`
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
      
      // 如果有套餐，默认选中第一个
      if (plans.value.length > 0) {
        selectedPlanId.value = plans.value[0].id
        await loadPlanItems()
      }
    }
  } catch (error) {
    console.error('获取套餐列表失败:', error)
    ElMessage.error('获取套餐列表失败')
  } finally {
    loading.value = false
  }
}

// 加载套餐的健康项目
const loadPlanItems = async () => {
  if (!selectedPlanId.value) return
  
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/healthitem-plan/plans/${selectedPlanId.value}/items`, {
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

// 处理更新事件
const handleUpdated = async () => {
  if (activeTab.value === 'list' && selectedPlanId.value) {
    await loadPlanItems()
  }
}
</script>

<style scoped>
.plan-health-item-manager {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.plan-selector {
  margin-bottom: 20px;
  display: flex;
  align-items: center;
}

.plan-selector .label {
  margin-right: 10px;
  font-weight: bold;
}

.loading-state, .empty-state {
  padding: 20px 0;
  margin: 20px 0;
}
</style>
