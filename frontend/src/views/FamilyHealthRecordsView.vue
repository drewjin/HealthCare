<template>
  <div class="family-health-records-container">
    <el-page-header @back="goBack" title="返回">
      <template #content>
        <span class="page-title">{{ familyMember ? `${familyMember.name}的健康记录` : '亲友健康记录' }}</span>
      </template>
    </el-page-header>
    
    <el-card class="health-record-card" shadow="hover">
      <div v-if="loading" class="loading-state">
        <el-skeleton :rows="5" animated />
      </div>
      <div v-else-if="error" class="error-state">
        <el-alert :title="error" type="error" />
      </div>
      <template v-else>
        <div class="member-info">
          <h3>{{ familyMember ? familyMember.name : '亲友' }} ({{ familyMember ? familyMember.relationship : '' }})</h3>
          <p>用户名: {{ familyMember ? familyMember.username : '' }}</p>
        </div>
        
        <el-divider />
        
        <div class="health-data-section">
          <FamilyHealthView 
            :family-member-id="memberId" 
            :family-member-name="familyMember ? familyMember.name : '亲友'"
          />
        </div>
      </template>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import FamilyHealthView from '@/components/FamilyHealthView.vue'
import axios from 'axios'

const route = useRoute()
const router = useRouter()

const loading = ref(true)
const error = ref('')
const memberId = ref<string | number>('')
const familyMember = ref<{
  id: string | number;
  user_id: string | number;
  username: string;
  name: string;
  relationship: string;
}>()

// 返回上一页
const goBack = () => {
  router.push('/dashboard')
}

// 获取亲友信息
const fetchFamilyMemberInfo = async () => {
  const id = route.params.id as string
  
  if (!id) {
    error.value = '亲友ID不能为空'
    loading.value = false
    return
  }
  
  memberId.value = id
  
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      error.value = '未登录或会话已过期'
      loading.value = false
      return
    }
    
    const userId = localStorage.getItem('uid')
    
    // 获取所有已确认的家庭关系
    const response = await axios.get(`/api/family/confirmed/${userId}`, {
      headers: { Authorization: token }
    })
    
    // 查找当前查看的亲友信息
    const member = response.data.find((member: any) => member.user_id.toString() === id)
    
    if (member) {
      familyMember.value = {
        id: member.id,
        user_id: member.user_id,
        username: member.username,
        name: member.name,
        relationship: member.relationship
      }
    } else {
      error.value = '未找到该亲友信息'
    }
  } catch (err) {
    console.error('获取亲友信息失败:', err)
    error.value = '获取亲友信息失败，请重试'
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchFamilyMemberInfo()
})
</script>

<style scoped>
.family-health-records-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: bold;
}

.health-record-card {
  margin-top: 20px;
}

.loading-state, .error-state {
  padding: 30px;
  text-align: center;
}

.member-info {
  margin-bottom: 15px;
}

.member-info h3 {
  color: #303133;
  margin-bottom: 10px;
}

.member-info p {
  color: #606266;
  margin: 5px 0;
}

.health-data-section {
  margin-top: 20px;
}
</style>
