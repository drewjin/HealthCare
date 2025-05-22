<template>
  <div class="health-item-form">
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="3" animated />
    </div>
    <div v-else-if="error" class="error-state">
      <el-alert :title="error" type="error" />
    </div>
    <div v-else-if="items.length === 0" class="empty-state">
      <el-empty description="暂无健康指标数据" />
    </div>
    <div v-else>
      <div class="health-item-description">
        <p style="color: #333; margin-bottom: 15px;">以下是您的健康指标数据，与平均水平对比：</p>
      </div>
      <el-form :model="formData" label-position="left" label-width="120px">
        <el-form-item v-for="item in items" :key="item.name" :label="item.name">
          <el-input 
            v-if="editable" 
            v-model="formValues[item.name]" 
            :placeholder="'请输入' + item.name"
          />
          <div v-else class="item-value">
            <div style="color: #303133; font-weight: bold; font-size: 16px;">
              {{ item.userValue !== undefined ? item.userValue : (item.value !== undefined ? item.value : '-') }}
              <span class="unit-text" v-if="getUnit(item.name)">{{ getUnit(item.name) }}</span>
            </div>
            <div v-if="item.averageValue" class="average-value">
              平均水平: {{ item.averageValue }} {{ getUnit(item.name) }}
            </div>
          </div>
        </el-form-item>
        
        <el-form-item v-if="editable">
          <el-button type="primary" @click="handleSubmit">保存</el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
// Using ElMessage in the loadHealthData function
import { ElMessage } from 'element-plus'

interface HealthItem {
  name: string;
  value?: number | string;
  userValue?: number | string;
  averageValue?: number | string;
}

const props = defineProps({
  // 健康指标项列表
  items: {
    type: Array as () => HealthItem[],
    default: () => []
  },
  // 是否可编辑
  editable: {
    type: Boolean,
    default: false
  },
  // 是否自动加载数据
  autoLoad: {
    type: Boolean,
    default: false
  },
  // 用户ID，用于自动加载数据
  userId: {
    type: [String, Number],
    default: ''
  }
})

const emit = defineEmits(['submit', 'cancel', 'update'])

const loading = ref(false)
const error = ref('')
const formData = reactive({})
const formValues = reactive<Record<string, number | string>>({})

// 根据指标名称获取单位
const getUnit = (name: string): string => {
  const unitMap: Record<string, string> = {
    '身高': 'cm',
    '体重': 'kg',
    '血压': 'mmHg',
    '血糖': 'mmol/L',
    '心率': '次/分',
    '体温': '°C',
    '血氧': '%',
    '肺活量': 'mL',
    '胆固醇': 'mmol/L',
    '尿酸': 'μmol/L'
  }
  
  return unitMap[name] || ''
}

// 初始化表单值
const initFormValues = () => {
  props.items.forEach(item => {
    formValues[item.name] = item.userValue || item.value || '';
  })
}

// 提交表单
const handleSubmit = () => {
  const updatedItems = props.items.map(item => ({
    name: item.name,
    userValue: formValues[item.name],
    averageValue: item.averageValue
  }))
  
  emit('submit', updatedItems)
}

// 取消编辑
const handleCancel = () => {
  initFormValues()
  emit('cancel')
}

// 加载健康数据
const loadHealthData = async () => {
  if (!props.userId) {
    error.value = '用户ID不能为空'
    return
  }
  
  try {
    loading.value = true
    const token = localStorage.getItem('jwt')
    if (!token) {
      error.value = '未登录或会话已过期'
      return
    }
    
    // 这里可以根据实际API调整请求
    const response = await fetch(`/api/healthitems/byid/${props.userId}`, {
      headers: { Authorization: token }
    })
    
    if (!response.ok) {
      throw new Error('加载健康数据失败')
    }
    
    const data = await response.json()
    if (data.items && data.items.length > 0) {
      // 解析健康指标数据
      const healthInfo = data.items[0].user_health_info
      const parsedItems = healthInfo.split(/[;,]/)
        .filter(Boolean)
        .map((item: string) => {
          const match = item.match(/([^:：]+)[：:]?(\d+)/)
          if (match) {
            return {
              name: match[1].trim(),
              value: Number(match[2])
            }
          }
          return null
        })
        .filter(Boolean)
      
      emit('update', parsedItems)
    }
  } catch (err) {
    console.error('加载健康数据失败:', err)
    error.value = '加载健康数据失败，请重试'
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  initFormValues()
  
  if (props.autoLoad && props.userId) {
    loadHealthData()
  }
})

// 监听项目变化
watch(() => props.items, () => {
  initFormValues()
}, { deep: true })

// 暴露方法给父组件
defineExpose({
  refresh: loadHealthData
})
</script>

<style scoped>
.health-item-form {
  max-width: 800px;
  margin: 0 auto;
}

.loading-state, .error-state, .empty-state {
  padding: 20px;
  text-align: center;
}

.item-value {
  padding: 8px 0;
  color: #606266;
  font-size: 16px;
}

.unit-text {
  color: #909399;
  font-size: 0.85em;
  margin-left: 4px;
}

.average-value {
  margin-top: 4px;
  color: #5470c6;
  font-size: 13px;
}

.health-item-description {
  background-color: #f5f7fa;
  padding: 10px 15px;
  border-radius: 4px;
  margin-bottom: 20px;
  border-left: 4px solid #409eff;
}
</style>
