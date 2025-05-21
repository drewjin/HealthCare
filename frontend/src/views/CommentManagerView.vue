<template>
  <div class="comment-items-container">
    <el-page-header @back="goBack" title="返回">
      <template #content>
        <span class="page-title">评论管理</span>
      </template>
    </el-page-header>

    <div class="comment-items-section">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>评论列表</span>
            <div class="button-group">
              <el-button type="success" @click="addItem()">创建评论</el-button>
              <el-button type="primary" @click="refreshItems">刷新列表</el-button>
            </div>
          </div>
        </template>

        <el-table v-loading="loading" :data="commentItems" style="width: 100%">
          <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="Commentary" label="评论" />
          <el-table-column label="创建时间" width="180">
            <template #default="scope">
              {{ formatDate(scope.row.CreatedAt) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="300">
            <template #default="scope">
              <el-button size="small" type="danger" @click="deleteItem(scope.row.ID)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <el-dialog v-model="editDialogVisible" title="新增评论" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="评论">
          <el-input v-model="editForm.Commentary" placeholder="请输入" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="saving" @click="saveItemChanges">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'

interface CommentItem {
  ID: number
  UserId: string | number
  PlanId: string | number
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
}

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const saving = ref(false)
const commentItems = ref<CommentItem[]>([])
const editDialogVisible = ref(false)

const editForm = reactive({
  ID: '',
  Commentary: ''
})

const goBack = () => {
  router.back()
}

const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString()
}  

const fetchCommentItems = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get('/api/commentary/get/user_list', {
      headers: { Authorization: `${token}` }
    })
    
    const items = response.data.commentaries || []
    
    commentItems.value = items.map((item: any) => ({
      ID: item.ID || item.id,
      UserId: item.UserId || item.user_id,
      PlanId: item.PlanId || item.plan_id,
      Commentary: item.Commentary || item.commentary,
      CreatedAt: item.CreatedAt || item.created_at || '',
    }))
    
    console.log('Fetched health items:', commentItems.value)
  } catch (error) {
    console.error('Failed to fetch health items:', error)
    ElMessage.error('获取评论列表失败')
  } finally {
    loading.value = false
  }
}

// 新增
const addItem = () => {
  editForm.ID = null
  editForm.UserId = null
  editForm.PlanId = null
  editForm.Commentary = null
  editDialogVisible.value = true
}

// 保存
const saveItemChanges = async () => {
  saving.value = true
  try {
    const token = localStorage.getItem('jwt')
    await axios.post('/api/commentary/add', {
      user_id: 1,
      plan_id: 1,
      commentary: editForm.Commentary
    }, {
      headers: { 
        Authorization: `${token}`,
        'Content-Type': 'application/json; charset=utf-8'
      }
    })
    ElMessage.success('评论新增成功')
    editDialogVisible.value = false

    // 更新列表
    await fetchCommentItems()
  } catch (error: any) {
    console.error('Failed to update health item:', error)
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('评论新增失败')
    }
  } finally {
    saving.value = false
  }
}

// 删除
const deleteItem = async (key: string) => {
  try {
    const token = localStorage.getItem('jwt')
    await axios.delete(`/api/commentary/delete/${key}`, {
      headers: { Authorization: token }
    })
    ElMessage.success('评论删除成功')
    // 更新列表
    await fetchCommentItems()
  } catch (error: any) {
    console.error('Failed to update health item:', error)
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('评论删除失败')
    }
  } finally {
    saving.value = false
  }
}

// 刷新列表
const refreshItems = () => {
  fetchCommentItems()
}

onMounted(() => {
  fetchCommentItems()
})
</script>

<style scoped>
.comment-items-container {
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: bold;
}

.comment-items-section {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.button-group {
  display: flex;
  gap: 10px;
}

.empty-data {
  padding: 20px;
  text-align: center;
  color: #909399;
}

.mt-4 {
  margin-top: 16px;
}
</style>
