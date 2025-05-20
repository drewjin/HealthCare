<template>
  <div class="health-item-manager">
    <el-card class="health-item-form">
      <template #header>
        <div class="card-header">
          <span>{{ isEditing ? '编辑检查项目' : '创建检查项目模板' }}</span>
        </div>
      </template>
      
      <!-- 创建模板表单 -->
      <div v-if="!isEditing" class="create-template-section">
        <el-form :model="createForm" label-width="120px">
          <el-form-item label="检查项目列表">
            <el-input
              v-model="createForm.healthItems"
              type="textarea"
              placeholder="输入格式：项目1, 项目2, 项目3"
              :rows="4"
            />
            <div class="item-hint">多个项目请用逗号分隔</div>
          </el-form-item>
          
          <el-form-item>
            <el-button type="primary" :loading="loading" @click="createTemplate">创建模板</el-button>
          </el-form-item>
        </el-form>
        
        <!-- 预览区域 -->
        <div v-if="templatePreview" class="preview-section">
          <h4>模板预览</h4>
          <div class="preview-content">{{ templatePreview }}</div>
        </div>
      </div>
      
      <!-- 编辑检查项目表单 -->
      <div v-else class="edit-values-section">
        <el-form :model="editForm" label-width="120px">
          <el-form-item label="检查项目字符串">
            <el-input
              v-model="editForm.itemString"
              type="textarea"
              placeholder="检查项目字符串"
              :rows="2"
              :disabled="true"
            />
          </el-form-item>
          
          <!-- 现有项目列表 -->
          <div class="existing-items">
            <h4>当前检查项目</h4>
            <el-table :data="parsedItems" style="width: 100%">
              <el-table-column prop="key" label="项目名称" width="180" />
              <el-table-column prop="value" label="数值" />
              <el-table-column fixed="right" label="操作" width="150">
                <template #default="scope">
                  <el-button 
                    link 
                    type="primary" 
                    @click="editItem(scope.row.key)"
                  >
                    编辑
                  </el-button>
                  <el-button 
                    link 
                    type="danger" 
                    @click="deleteItem(scope.row.key)"
                  >
                    删除
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
          
          <!-- 添加/编辑项目值表单 -->
          <div v-if="showItemForm" class="item-form">
            <h4>{{ itemFormMode === 'add' ? '添加检查项目数值' : '编辑检查项目数值' }}</h4>
            <el-form :model="itemForm" label-width="100px">
              <el-form-item label="项目名称">
                <el-input v-model="itemForm.key" :disabled="itemFormMode === 'edit'" />
              </el-form-item>
              <el-form-item label="数值">
                <el-input v-model="itemForm.value" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="saveItemValue">保存</el-button>
                <el-button @click="cancelItemEdit">取消</el-button>
              </el-form-item>
            </el-form>
          </div>
          
          <el-form-item v-else>
            <el-button type="primary" @click="showAddItem">添加项目</el-button>
            <el-button type="success" :loading="loading" @click="saveChanges">保存全部更改</el-button>
            <el-button @click="cancelEdit">返回</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watchEffect } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const props = defineProps({
  initialItemId: {
    type: Number,
    default: 0
  },
  initialItemString: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:itemString', 'created', 'updated', 'cancel'])

// 状态变量
const isEditing = ref(!!props.initialItemId)
const loading = ref(false)
const templatePreview = ref('')
const showItemForm = ref(false)
const itemFormMode = ref<'add' | 'edit'>('add')
const itemsToDelete = ref<string[]>([])

// 创建模板表单
const createForm = reactive({
  healthItems: ''
})

// 编辑表单
const editForm = reactive({
  itemId: props.initialItemId,
  itemString: props.initialItemString,
})

// 项目编辑表单
const itemForm = reactive({
  key: '',
  value: ''
})

// 计算属性：解析项目字符串为键值对数组
const parsedItems = computed(() => {
  const items: { key: string, value: string }[] = []
  
  if (!editForm.itemString) return items
  
  // 简单解析，可以替换为与后端一致的正则表达式
  const pairs = editForm.itemString.split(',')
  for (const pair of pairs) {
    const [key, value] = pair.split(':')
    if (key) {
      items.push({
        key: key.trim(),
        value: value ? value.trim() : ''
      })
    }
  }
  
  return items
})

// 创建模板
const createTemplate = async () => {
  if (!createForm.healthItems) {
    ElMessage.warning('请输入检查项目')
    return
  }
  
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    // 先创建模板字符串
    const templateResponse = await axios.post('/api/healthitem-manager/template', {
      health_items: createForm.healthItems
    }, {
      headers: { Authorization: token }
    })
    
    templatePreview.value = templateResponse.data.template
    
    // 然后保存到数据库
    const saveResponse = await axios.post('/api/healthitem-manager/save-template', {
      item_name: templateResponse.data.template
    }, {
      headers: { Authorization: token }
    })
    
    ElMessage.success('检查项目模板创建并保存成功')
    console.log('Template saved:', saveResponse.data)
    emit('created', templateResponse.data.template)
  } catch (error) {
    console.error('创建检查项目模板失败:', error)
    ElMessage.error('创建失败，请重试')
  } finally {
    loading.value = false
  }
}

// 显示添加项目表单
const showAddItem = () => {
  itemForm.key = ''
  itemForm.value = ''
  itemFormMode.value = 'add'
  showItemForm.value = true
}

// 编辑项目
const editItem = (key: string) => {
  const item = parsedItems.value.find(item => item.key === key)
  if (item) {
    itemForm.key = item.key
    itemForm.value = item.value
    itemFormMode.value = 'edit'
    showItemForm.value = true
  }
}

// 删除项目
const deleteItem = (key: string) => {
  itemsToDelete.value.push(key)
  ElMessage.success(`项目 "${key}" 标记为删除`)
}

// 保存项目值
const saveItemValue = () => {
  if (!itemForm.key.trim()) {
    ElMessage.warning('项目名称不能为空')
    return
  }
  
  // 更新本地的检查项目字符串
  const items = [...parsedItems.value]
  const existingIndex = items.findIndex(item => item.key === itemForm.key)
  
  if (existingIndex >= 0) {
    // 更新已有项目
    items[existingIndex].value = itemForm.value
  } else {
    // 添加新项目
    items.push({ key: itemForm.key, value: itemForm.value })
  }
  
  // 重建检查项目字符串
  updateItemString(items)
  
  showItemForm.value = false
  ElMessage.success('项目已更新，请记得保存全部更改')
}

// 取消项目编辑
const cancelItemEdit = () => {
  showItemForm.value = false
}

// 更新检查项目字符串
const updateItemString = (items: { key: string, value: string }[]) => {
  const pairs = items.map(item => `${item.key}:${item.value}`)
  editForm.itemString = pairs.join(', ')
}

// 保存所有更改
const saveChanges = async () => {
  if (!editForm.itemId) {
    ElMessage.error('缺少检查项目ID')
    return
  }
  
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    
    // 准备更新数据
    const updates: Record<string, string> = {}
    parsedItems.value.forEach(item => {
      updates[item.key] = item.value
    })
    
    const response = await axios.put('/api/healthitem-manager/values', {
      item_id: editForm.itemId,
      item_string: editForm.itemString,
      updates: updates,
      delete_keys: itemsToDelete.value
    }, {
      headers: { Authorization: token }
    })
    
    // 更新成功，更新组件状态
    editForm.itemString = response.data.item_string
    itemsToDelete.value = []
    
    ElMessage.success('检查项目更新成功')
    emit('updated', response.data.item_string)
  } catch (error) {
    console.error('更新检查项目失败:', error)
    ElMessage.error('更新失败，请重试')
  } finally {
    loading.value = false
  }
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  emit('cancel')
}

// 初始化
watchEffect(() => {
  if (props.initialItemId) {
    isEditing.value = true
    editForm.itemId = props.initialItemId
    editForm.itemString = props.initialItemString
  } else {
    isEditing.value = false
  }
})
</script>

<style scoped>
.health-item-manager {
  margin: 20px;
}

.health-item-form {
  max-width: 800px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-hint {
  font-size: 12px;
  color: #999;
  margin-top: 5px;
}

.preview-section {
  margin-top: 20px;
  padding: 15px;
  background-color: #f8f8f8;
  border-radius: 4px;
}

.preview-content {
  font-family: monospace;
  padding: 10px;
  background-color: #fff;
  border: 1px solid #eee;
  border-radius: 4px;
}

.existing-items {
  margin: 20px 0;
}

.item-form {
  margin: 20px 0;
  padding: 15px;
  background-color: #f8f8f8;
  border-radius: 4px;
}
</style>
