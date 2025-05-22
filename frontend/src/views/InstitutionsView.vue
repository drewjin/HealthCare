<template>
  <div class="institutions-container">
    <h1>体检机构列表</h1>
    <el-card v-if="loading" class="loading-card">
      <el-skeleton :rows="5" animated />
    </el-card>
    
    <div v-else-if="institutions.length === 0" class="empty-state">
      <el-empty description="暂无可用的体检机构" />
    </div>
    
    <div v-else class="institution-list">
      <el-card v-for="(institution, index) in institutions" :key="institution?.ID || index" class="institution-card">
        <div class="institution-header">
          <h2>{{ institution?.institution_name || '未命名机构' }}</h2>
        </div>
        <div class="institution-info">
          <p><strong>地址:</strong> {{ institution?.institution_address || '暂无地址' }}</p>
          <p><strong>资质:</strong> {{ institution?.institution_qualification || '暂无资质信息' }}</p>
        </div>
        <div class="institution-actions">
          <el-button type="primary" @click="institution?.ID ? viewDetails(institution.ID) : ElMessage.warning('机构ID无效，无法查看详情')">查看详情</el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

interface Institution {
  ID: number
  institution_name: string
  institution_address: string
  institution_qualification: string
  examination_package: string
  status: number
}

const router = useRouter()
const institutions = ref<Institution[]>([])
const loading = ref(true)

const fetchInstitutions = async () => {
  try {
    loading.value = true
    const token = localStorage.getItem('jwt')
    
    if (!token) {
      ElMessage.error('用户未登录，请先登录')
      router.push('/login')
      return
    }
    
    const response = await axios.get('/api/institutions', {
      headers: { Authorization: `${token}` }
    })
    
    // 处理不同的响应格式：可能是数组或单个对象
    console.log('获取到的机构数据:', response.data)
    
    try {
      if (Array.isArray(response.data)) {
        // 如果是数组，过滤并进行类型强制转换
        const validInstitutions = response.data
          .filter(inst => inst && typeof inst === 'object' && 'ID' in inst)
          .map(inst => inst as Institution);
        
        institutions.value = validInstitutions;
        console.log('获取到机构列表（数组）：', institutions.value.length)
      } else if (response.data && typeof response.data === 'object') {
        // 如果是单个对象，将其转换为数组
        if ('ID' in response.data) {
          institutions.value = [response.data as Institution];
          console.log('获取到单个机构，已转换为数组')
        } else {
          // 尝试从对象中提取机构
          try {
            const values = Object.values(response.data);
            const validInstitutions = values
              .filter(item => item && typeof item === 'object' && 'ID' in item)
              .map(item => item as Institution);
              
            if (validInstitutions.length > 0) {
              institutions.value = validInstitutions;
              console.log('从对象中提取机构数组：', validInstitutions.length)
            } else {
              throw new Error('未找到有效的机构数据');
            }
          } catch (extractError) {
            console.error('从对象提取机构列表失败:', extractError);
            institutions.value = [];
            ElMessage.warning('获取机构列表格式不正确');
          }
        }
      } else {
        // 无法识别的格式
        throw new Error('响应数据格式不正确');
      }
    } catch (formatError) {
      console.error('处理机构数据出错:', formatError);
      institutions.value = [];
      ElMessage.warning('获取机构列表格式不正确');
    }
  } catch (error) {
    console.error('Failed to fetch institutions:', error)
    ElMessage.error('获取机构列表失败')
  } finally {
    loading.value = false
  }
}

const viewDetails = (institutionId: number) => {
  if (!institutionId) {
    ElMessage.warning('无效的机构ID，无法查看详情')
    return
  }
  router.push(`/institutions/${institutionId}`)
}

// 无论是挂载时还是每次路由进入时都刷新数据
onMounted(() => {
  fetchInstitutions()
})

// 这个页面不再被缓存，因此不需要 onActivated 钩子
</script>

<style scoped>
.institutions-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.institution-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
  margin-top: 20px;
}

.institution-card {
  transition: transform 0.3s;
}

.institution-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
}

.institution-header {
  margin-bottom: 15px;
  border-bottom: 1px solid #eee;
  padding-bottom: 10px;
}

.institution-info {
  margin-bottom: 15px;
}

.institution-actions {
  display: flex;
  justify-content: flex-end;
}

.loading-card {
  padding: 20px;
  margin-top: 20px;
}

.empty-state {
  margin-top: 40px;
  text-align: center;
}
</style>