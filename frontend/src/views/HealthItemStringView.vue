<template>
  <div class="health-item-string-view">
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="3" animated />
    </div>
    <div v-else-if="error" class="error-state">
      <el-alert :title="error" type="error" />
    </div>
    <div v-else>
      <div v-if="Object.keys(parsedItems).length === 0" class="empty-state">
        <el-empty description="暂无健康指标数据" />
      </div>
      <div v-else>
        <el-table :data="tableData" style="width: 100%" v-loading="loading">
          <el-table-column prop="key" label="项目名称" width="180" />
          <el-table-column prop="value" label="检测值" />
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

const props = defineProps({
  itemString: {
    type: String,
    default: ''
  },
  itemId: {
    type: Number,
    default: 0
  },
  autoFetch: {
    type: Boolean,
    default: false
  }
})

const loading = ref(false)
const error = ref('')
const parsedItems = ref<Record<string, string>>({})

// 转换为表格数据格式
const tableData = computed(() => {
  return Object.entries(parsedItems.value).map(([key, value]) => ({
    key,
    value: value || '无数据'
  }))
})

// 解析健康项目字符串
const parseHealthItemString = (str: string): Record<string, string> => {
  const result: Record<string, string> = {}
  if (!str) return result
  
  // 增强解析，支持多种分隔符（逗号和分号）
  const pairs = str.split(/[;,]/)
  for (const pair of pairs) {
    // 匹配各种格式的健康指标，如"指标名:值"或"指标名：值"
    const match = pair.match(/([^:：]+)[：:]?(.*)/);
    if (match && match[1]) {
      const key = match[1].trim();
      const value = match[2] ? match[2].trim() : '';
      if (key) {
        result[key] = value;
      }
    }
  }
  
  return result
}

// 从字符串解析健康项目
const parseFromString = () => {
  parsedItems.value = parseHealthItemString(props.itemString)
}

// 从API获取解析后的健康项目
const fetchFromAPI = async () => {
  console.log('HealthItemStringView - fetchFromAPI - itemId:', props.itemId)
  
  if (!props.itemId) {
    console.log('HealthItemStringView - fetchFromAPI - 缺少健康项目ID')
    error.value = '缺少健康项目ID'
    return
  }
  
  loading.value = true
  error.value = ''
  
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      console.log('HealthItemStringView - fetchFromAPI - 未登录或会话已过期')
      error.value = '未登录或会话已过期'
      return
    }
    
    console.log(`HealthItemStringView - fetchFromAPI - 发送请求: /api/healthitem-manager/values/${props.itemId}`)
    const response = await axios.get(`/api/healthitem-manager/values/${props.itemId}`, {
      headers: { Authorization: token }
    })
    
    console.log('HealthItemStringView - fetchFromAPI - 响应数据:', response.data)
    parsedItems.value = response.data.values || {}
  } catch (err: any) {
    console.error('获取健康项目解析结果失败:', err)
    if (err.response) {
      console.error('错误响应状态:', err.response.status)
      console.error('错误响应数据:', err.response.data)
    }
    error.value = '获取数据失败'
  } finally {
    loading.value = false
  }
}

// 初始化
onMounted(() => {
  if (props.autoFetch && props.itemId) {
    fetchFromAPI()
  } else if (props.itemString) {
    parseFromString()
  }
})

// 监听属性变化
watch(() => props.itemString, (newVal) => {
  if (newVal && !props.autoFetch) {
    parseFromString()
  }
})

watch(() => props.itemId, (newVal) => {
  if (newVal && props.autoFetch) {
    fetchFromAPI()
  }
})

// 暴露方法给父组件
defineExpose({
  refresh: () => {
    if (props.autoFetch && props.itemId) {
      fetchFromAPI()
    } else if (props.itemString) {
      parseFromString()
    }
  }
})
</script>

<style scoped>
.health-item-string-view {
  margin: 10px 0;
}

.loading-state, .error-state, .empty-state {
  padding: 20px;
  text-align: center;
}
</style>