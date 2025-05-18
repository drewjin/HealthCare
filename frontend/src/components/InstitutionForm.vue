<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage, ElTabs, ElTabPane } from 'element-plus'
import PackageManager from './PackageManager.vue'

const institutionForm = ref({
  institution_name: '',
  institution_address: '',
  institution_qualification: '',
  examination_package: ''
})

const emit = defineEmits(['submit-success'])
const activeTab = ref('basic')
const institutionId = ref(0)
const isSubmitted = ref(false)

// 检查用户是否已经提交过机构信息
const checkExistingInstitution = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const uid = localStorage.getItem('uid')
    if (!token || !uid) return

    const response = await fetch(`/api/user/${uid}/institution`, {
      method: 'GET',
      headers: {
        Authorization: `${token}`
      }
    })

    if (response.ok) {
      const data = await response.json()
      // 如果用户已经有机构信息，则填充表单
      if (data && data.ID) {
        institutionId.value = data.ID
        institutionForm.value.institution_name = data.institution_name || ''
        institutionForm.value.institution_address = data.institution_address || ''
        institutionForm.value.institution_qualification = data.institution_qualification || ''
        institutionForm.value.examination_package = data.examination_package || ''
        isSubmitted.value = true
      }
    }
  } catch (error) {
    console.error('Failed to fetch institution info:', error)
  }
}

const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const uid = localStorage.getItem('uid')
    if (!token || !uid) throw new Error('未登录')

    // 验证必填字段
    if (!institutionForm.value.institution_name || !institutionForm.value.institution_address) {
      ElMessage.warning('请填写机构名称和地址')
      return
    }

    console.info('Submitting institution form:', institutionForm.value)
    const response = await fetch(`/api/institutions/${uid}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `${token}`
      },
      body: JSON.stringify({
        institution_name: institutionForm.value.institution_name,
        institution_address: institutionForm.value.institution_address,
        institution_qualification: institutionForm.value.institution_qualification,
        examination_package: institutionForm.value.examination_package
      })
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || '提交失败')
    }

    // 获取新创建的机构ID
    await checkExistingInstitution()
    
    ElMessage.success('机构信息提交成功，等待管理员审核')
    isSubmitted.value = true
    emit('submit-success')
    
    // 提交成功后切换到套餐管理标签页
    activeTab.value = 'packages'
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '提交失败')
  }
}

onMounted(() => {
  checkExistingInstitution()
})
</script>

<template>
  <div class="institution-form-container">
    <ElTabs v-model="activeTab" class="form-tabs">
      <ElTabPane label="基本信息" name="basic">
        <ElForm :model="institutionForm" label-width="120px">
          <ElFormItem label="机构名称">
            <ElInput v-model="institutionForm.institution_name" :disabled="isSubmitted" />
          </ElFormItem>
          <ElFormItem label="机构地址">
            <ElInput v-model="institutionForm.institution_address" :disabled="isSubmitted" />
          </ElFormItem>
          <ElFormItem label="资质证明">
            <ElInput 
              v-model="institutionForm.institution_qualification" 
              type="textarea" 
              rows="3"
              placeholder="请输入机构资质证明信息"
              :disabled="isSubmitted"
            />
          </ElFormItem>
          <ElFormItem v-if="!isSubmitted">
            <ElButton type="primary" @click="handleSubmit">提交</ElButton>
          </ElFormItem>
          <div v-else class="submitted-message">
            <p>您的机构信息已提交，请前往"套餐管理"标签页设置体检套餐</p>
          </div>
        </ElForm>
      </ElTabPane>
      
      <ElTabPane label="套餐管理" name="packages" :disabled="!isSubmitted">
        <div v-if="isSubmitted && institutionId">
          <PackageManager :institutionId="institutionId" />
        </div>
        <div v-else class="no-institution-message">
          <p>请先在"基本信息"标签页提交机构信息</p>
        </div>
      </ElTabPane>
    </ElTabs>
  </div>
</template>

<style scoped>
.institution-form-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.form-tabs {
  margin-top: 20px;
}

.submitted-message,
.no-institution-message {
  text-align: center;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 4px;
  margin-top: 20px;
}

.submitted-message {
  color: #67c23a;
}

.no-institution-message {
  color: #e6a23c;
}
</style>