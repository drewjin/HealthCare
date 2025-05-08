<script setup lang="ts">
import { ref } from 'vue'
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage } from 'element-plus'

const institutionForm = ref({
  institution_name: '',
  institution_address: '',
  institution_qualification: '',
  examination_package: ''
})

const emit = defineEmits(['submit-success'])

const handleSubmit = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const uid = localStorage.getItem('uid')
    if (!token || !uid) throw new Error('未登录')

    console.info('Submitting institution form:', institutionForm.value)
    const response = await fetch(`http://localhost:3000/api/institutions/${uid}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: token
      },
      body: JSON.stringify({
        institution_name: institutionForm.value.institution_name,
        institution_address: institutionForm.value.institution_address,
        institution_qualification: institutionForm.value.institution_qualification,
        examination_package: institutionForm.value.examination_package
      })
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || '提交失败')
    }

    ElMessage.success('机构信息提交成功，等待管理员审核')
    emit('submit-success')
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '提交失败')
  }
}
</script>

<template>
  <ElForm :model="institutionForm" label-width="120px">
    <ElFormItem label="机构名称">
      <ElInput v-model="institutionForm.institution_name" />
    </ElFormItem>
    <ElFormItem label="机构地址">
      <ElInput v-model="institutionForm.institution_address" />
    </ElFormItem>
    <ElFormItem label="资质证明">
      <ElInput 
        v-model="institutionForm.institution_qualification" 
        type="textarea" 
        rows="3"
        placeholder="请输入机构资质证明信息"
      />
    </ElFormItem>
    <ElFormItem label="体检套餐">
      <ElInput
        v-model="institutionForm.examination_package"
        type="textarea"
        rows="4"
        placeholder="请输入体检套餐信息"
      />
    </ElFormItem>
    <ElFormItem>
      <ElButton type="primary" @click="handleSubmit">提交</ElButton>
    </ElFormItem>
  </ElForm>
</template>

<style scoped>
.el-form {
  max-width: 600px;
  margin: 0 auto;
}
</style>