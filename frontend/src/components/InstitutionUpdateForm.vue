<template>
  <div class="institution-update-form">
    <el-form :model="updateForm" label-width="120px">
      <el-form-item label="机构名称">
        <el-input v-model="updateForm.institution_name" placeholder="输入机构名称" />
      </el-form-item>
      <el-form-item label="联系电话">
        <el-input v-model="updateForm.institution_phone" placeholder="输入联系电话" />
      </el-form-item>
      <el-form-item label="机构地址">
        <el-input v-model="updateForm.institution_address" placeholder="输入机构地址" />
      </el-form-item>
      <el-form-item label="资质证明">
        <el-input 
          v-model="updateForm.institution_qualification" 
          type="textarea" 
          rows="3"
          placeholder="请输入机构资质证明信息"
        />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="updating" @click="handleUpdate">更新信息</el-button>
      </el-form-item>
    </el-form>
    
    <el-alert
      v-if="updateSuccess"
      title="机构信息更新成功"
      type="success"
      :closable="false"
      show-icon
      class="alert-message"
    />
    
    <el-alert
      v-if="updateError"
      :title="errorMessage"
      type="error"
      :closable="false"
      show-icon
      class="alert-message"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  institutionId: number | string
}>()

const updateForm = ref({
  institution_name: '',
  institution_phone: '',
  institution_address: '',
  institution_qualification: ''
})

const updating = ref(false)
const updateSuccess = ref(false)
const updateError = ref(false)
const errorMessage = ref('')

// 获取机构原始信息
const fetchInstitutionInfo = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/institutions/${props.institutionId}`, {
      headers: { Authorization: `${token}` }
    })
    
    const institution = response.data.institution
    
    // 填充表单
    updateForm.value.institution_name = institution.institution_name || ''
    updateForm.value.institution_phone = institution.institution_phone || ''
    updateForm.value.institution_address = institution.institution_address || ''
    updateForm.value.institution_qualification = institution.institution_qualification || ''
  } catch (error) {
    console.error('Failed to fetch institution info:', error)
    ElMessage.error('获取机构信息失败')
  }
}

const handleUpdate = async () => {
  resetAlerts()
  
  // 构建更新数据，只包含非空的字段
  const updateData: Record<string, string> = {}
  
  if (updateForm.value.institution_name) {
    updateData.institution_name = updateForm.value.institution_name
  }
  
  if (updateForm.value.institution_phone) {
    updateData.institution_phone = updateForm.value.institution_phone
  }
  
  if (updateForm.value.institution_address) {
    updateData.institution_address = updateForm.value.institution_address
  }
  
  if (updateForm.value.institution_qualification) {
    updateData.institution_qualification = updateForm.value.institution_qualification
  }
  
  // 如果没有任何更新数据，提示用户
  if (Object.keys(updateData).length === 0) {
    ElMessage.warning('请至少修改一项信息')
    return
  }
  
  updating.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.patch(`/api/institutions/${props.institutionId}/update`, updateData, {
      headers: { Authorization: `${token}` }
    })
    
    // 检查响应内容
    if (response.data && response.data.institution) {
      // 更新本地表单数据以反映最新的服务器状态
      updateForm.value.institution_name = response.data.institution.institution_name || ''
      updateForm.value.institution_phone = response.data.institution.institution_phone || ''
      updateForm.value.institution_address = response.data.institution.institution_address || ''
      updateForm.value.institution_qualification = response.data.institution.institution_qualification || ''
    }
    
    updateSuccess.value = true
    ElMessage.success('机构信息更新成功')
  } catch (error: any) {
    console.error('Failed to update institution:', error)
    updateError.value = true
    
    // 提取具体错误信息
    if (error.response && error.response.data && error.response.data.error) {
      errorMessage.value = error.response.data.error
    } else {
      errorMessage.value = '更新机构信息失败'
    }
  } finally {
    updating.value = false
    
    // 3秒后隐藏成功消息
    if (updateSuccess.value) {
      setTimeout(() => {
        updateSuccess.value = false
      }, 3000)
    }
  }
}

const resetAlerts = () => {
  updateSuccess.value = false
  updateError.value = false
  errorMessage.value = ''
}

onMounted(() => {
  fetchInstitutionInfo()
})
</script>

<style scoped>
.institution-update-form {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px 0;
}

.alert-message {
  margin-top: 20px;
}
</style>
