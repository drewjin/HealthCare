<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElCard, ElButton, ElTable, ElTableColumn, ElMessage } from 'element-plus'

const pendingInstitutions = ref<Array<{
  ID: number
  institution_name: string
  institution_address: string
  institution_qualification: string
  examination_package: string
  created_at: string
}>>([])

const handleReview = async (id: number, approved: boolean) => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('未登录')

    const response = await fetch(`/api/institutions/${id}/review`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ approved })
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || '操作失败')
    }

    ElMessage.success(approved ? '已通过审核' : '已拒绝申请')
    // 从列表中移除已处理的机构
    pendingInstitutions.value = pendingInstitutions.value.filter(inst => inst.ID !== id)
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '操作失败')
  }
}

const fetchPendingInstitutions = async () => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('未登录')

    const response = await fetch('/api/institutions/pending', {
      headers: { Authorization: `Bearer ${token}` }
    })

    if (!response.ok) {
      const data = await response.json()
      throw new Error(data.error || '获取数据失败')
    }

    const data = await response.json()
    pendingInstitutions.value = data
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '获取数据失败')
  }
}

onMounted(() => {
  fetchPendingInstitutions()
})
</script>

<template>
  <div class="admin-review">
    <ElCard shadow="hover">
      <template #header>
        <div class="card-header">
          <span>待审核机构</span>
          <ElButton type="primary" @click="fetchPendingInstitutions">刷新</ElButton>
        </div>
      </template>

      <ElTable :data="pendingInstitutions" style="width: 100%">
        <ElTableColumn prop="institution_name" label="机构名称" />
        <ElTableColumn prop="institution_address" label="地址" />
        <ElTableColumn prop="institution_qualification" label="资质证明" show-overflow-tooltip />
        <ElTableColumn prop="examination_package" label="体检套餐" show-overflow-tooltip />
        <ElTableColumn prop="created_at" label="申请时间" width="180" />
        <ElTableColumn fixed="right" label="操作" width="200">
          <template #default="{ row }">
            <ElButton
              type="success"
              size="small"
              @click="handleReview(row.ID, true)"
            >
              通过
            </ElButton>
            <ElButton
              type="danger"
              size="small"
              @click="handleReview(row.ID, false)"
            >
              拒绝
            </ElButton>
          </template>
        </ElTableColumn>
      </ElTable>

      <div v-if="pendingInstitutions.length === 0" class="empty-text">
        暂无待审核的机构
      </div>
    </ElCard>
  </div>
</template>

<style scoped>
.admin-review {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-text {
  text-align: center;
  color: #909399;
  padding: 30px 0;
}
</style>