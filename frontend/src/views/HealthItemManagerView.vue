<template>
  <div class="health-item-management-view">
    <el-card class="main-card">
      <template #header>
        <div class="card-header">
          <h2>检查项目管理</h2>
        </div>
      </template>
      
      <div v-if="loading" class="loading-container">
        <el-skeleton :rows="10" animated />
      </div>
      
      <div v-else-if="selectedItem">
        <div class="action-bar">
          <el-button @click="backToList" icon="ArrowLeft">返回列表</el-button>
        </div>
        
        <health-item-manager
          :initial-item-id="selectedItem.id"
          :initial-item-string="selectedItem.itemName"
          @updated="onItemUpdated"
          @cancel="backToList"
        />
      </div>
      
      <div v-else>
        <!-- 项目列表 -->
        <div class="list-actions">
          <el-button type="primary" @click="createNewItem">创建新检查项目</el-button>
          <el-input 
            v-model="searchQuery" 
            placeholder="搜索检查项目" 
            clearable 
            style="width: 300px; margin-left: 16px;" 
          />
        </div>
        
        <el-table 
          v-loading="loading" 
          :data="filteredHealthItems" 
          style="width: 100%; margin-top: 20px;"
        >
          <el-table-column prop="id" label="ID" width="80">
            <template #default="scope">
              <!-- 使用调试输出显示ID值和类型 -->
              <div>
                {{ scope.row.id || scope.row.ID || '无ID' }}
                <div style="font-size: 10px; color: #999">
                  (类型: {{ typeof (scope.row.id || scope.row.ID) }})
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="itemName" label="检查项目" />
          <el-table-column label="操作" width="250">
            <template #default="scope">
              <el-button 
                link 
                type="primary" 
                @click="viewItem(scope.row)"
              >
                查看/编辑
              </el-button>
              <el-button 
                link 
                type="primary" 
                @click="parseItem(scope.row)"
              >
                解析
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        
        <!-- 创建新项目 -->
        <el-dialog v-model="showCreateDialog" title="创建新检查项目" width="700px">
          <health-item-manager
            @created="onItemCreated"
            @cancel="showCreateDialog = false"
          />
        </el-dialog>
        
        <!-- 解析结果对话框 -->
        <el-dialog v-model="showParseDialog" title="检查项目解析结果" width="700px">
          <div v-if="parsedItemData">
            <p><strong>项目字符串:</strong> {{ parsedItemData.item_string }}</p>
            <el-divider />
            <h4>解析结果:</h4>
            <el-table :data="parsedItemList" style="width: 100%">
              <el-table-column prop="key" label="项目" />
              <el-table-column prop="value" label="值" />
            </el-table>
          </div>
        </el-dialog>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import HealthItemManager from '@/components/HealthItemManager.vue'

// 状态变量
const loading = ref(false)
const healthItems = ref<any[]>([])
const selectedItem = ref<any>(null)
const showCreateDialog = ref(false)
const showParseDialog = ref(false)
const parsedItemData = ref<any>(null)
const searchQuery = ref('')

// 过滤后的健康项目列表
const filteredHealthItems = computed(() => {
  if (!searchQuery.value) return healthItems.value
  
  const query = searchQuery.value.toLowerCase()
  return healthItems.value.filter(item => {
    // 考虑 itemName, ItemName, id 和 ID 这几种可能的属性格式
    const name = (item.itemName || item.ItemName || '').toLowerCase()
    const id = (item.id || item.ID || '').toString()
    return name.includes(query) || id.includes(query)
  })
})

// 解析后的项目列表
const parsedItemList = computed(() => {
  if (!parsedItemData.value || !parsedItemData.value.values) return []
  
  return Object.entries(parsedItemData.value.values).map(([key, value]) => ({
    key,
    value
  }))
})

// 获取所有健康项目
const fetchHealthItems = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('未登录或会话已过期')
      return
    }
    
    // 固定使用正确的URL路径
    const response = await axios.get('/api/healthitems', {
      headers: { Authorization: token }
    })
    
    // 显示API返回的原始数据结构
    console.log('健康项目API响应数据结构:', response.data)
    
    // 检查返回的项目数据 - 可能在items或health_items字段中
    const items = response.data.items || response.data.health_items || []
    console.log('健康项目列表:', items)
    
    // 标准化数据格式
    const normalizedItems = items.map((item: any) => ({
      id: item.ID || item.id,
      itemName: item.ItemName || item.item_name
    }))
    
    healthItems.value = normalizedItems
  } catch (error: any) {
    console.error('获取健康项目列表失败:', error)
    ElMessage.error('获取健康项目列表失败，请检查API路径')
    // 查看错误详情以帮助调试
    if (error.response) {
      console.error('错误响应:', error.response.status, error.response.data)
    }
  } finally {
    loading.value = false
  }
}

// 查看/编辑项目
const viewItem = (item: any) => {
  selectedItem.value = item
}

// 返回列表
const backToList = () => {
  selectedItem.value = null
  fetchHealthItems() // 刷新列表
}

// 创建新项目
const createNewItem = () => {
  showCreateDialog.value = true
}

// 解析项目
const parseItem = async (item: any) => {
  // 添加调试输出，查看输入的item对象
  console.log('解析项目 - 输入项目对象:', item)
  
  if (!item) {
    console.log('解析项目失败 - 项目对象为空')
    ElMessage.error('无效的检查项目ID')
    return
  }
  
  // 检查item的属性名，可能是ID而不是id
  console.log('项目对象的属性:', Object.keys(item))
  
  // 检查对象中是否有ID或id属性
  const itemId = item.ID ?? item.id ?? null
  console.log('使用的项目ID:', itemId, '类型:', typeof itemId)
  
  if (itemId === null || itemId === undefined) {
    console.log('解析项目失败 - 项目ID为空')
    ElMessage.error('无效的检查项目ID')
    return
  }
  
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('未登录或会话已过期')
      return
    }
    
    // 确保ID是数字类型
    console.log(`准备发送请求到: /api/healthitem-manager/values/${itemId}`)
    const response = await axios.get(`/api/healthitem-manager/values/${itemId}`, {
      headers: { Authorization: token }
    })
    
    console.log('解析项目成功，响应数据:', response.data)
    parsedItemData.value = response.data
    showParseDialog.value = true
  } catch (error: any) {
    console.error('解析健康项目失败:', error)
    if (error.response) {
      console.error('错误响应状态:', error.response.status)
      console.error('错误响应数据:', error.response.data)
    }
    ElMessage.error('解析失败，请重试')
  }
}

// 项目创建成功回调
const onItemCreated = (_: string) => {
  showCreateDialog.value = false
  ElMessage.success('检查项目创建成功')
  fetchHealthItems() // 刷新列表
}

// 项目更新成功回调
const onItemUpdated = (updatedItemString: string) => {
  ElMessage.success('检查项目更新成功')
  if (selectedItem.value) {
    selectedItem.value.itemName = updatedItemString
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchHealthItems()
})
</script>

<style scoped>
.health-item-management-view {
  padding: 20px;
}

.main-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.loading-container {
  padding: 20px;
}

.action-bar {
  margin-bottom: 20px;
}

.list-actions {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
}
</style>
