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
      <el-card v-for="institution in institutions" :key="institution.ID" class="institution-card">
        <div class="institution-header">
          <h2>{{ institution.institution_name }}</h2>
        </div>
        <div class="institution-info">
          <p><strong>地址:</strong> {{ institution.institution_address }}</p>
          <p><strong>资质:</strong> {{ institution.institution_qualification }}</p>
        </div>
        <div class="institution-actions">
          <el-button type="primary" @click="viewDetails(institution.ID)">查看详情</el-button>
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
    const response = await axios.get('/api/institutions', {
      headers: { Authorization: `Bearer ${token}` }
    })
    institutions.value = response.data
  } catch (error) {
    console.error('Failed to fetch institutions:', error)
    ElMessage.error('获取机构列表失败')
  } finally {
    loading.value = false
  }
}

const viewDetails = (institutionId: number) => {
  router.push(`/institutions/${institutionId}`)
}

onMounted(() => {
  fetchInstitutions()
})
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