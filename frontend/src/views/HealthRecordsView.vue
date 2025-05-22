<template>
  <div class="health-records-container">
    <el-button type="primary" icon="ArrowLeft" @click="goDashboard" class="back-button">
      返回 Dashboard
    </el-button>
    <div class="cards">
      <el-card v-for="rec in records" :key="rec.plan_id" class="record-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span>{{ rec.institution_name }} - {{ rec.plan_name }}</span>
            <span class="status-badge" :class="rec.status === 1 ? 'completed' : 'pending'">
              {{ rec.status === 1 ? '已完成' : '待完成' }}
            </span>
          </div>
        </template>
        <div class="progress-info">
          <span class="item-count">共 {{ rec.item_count }} 项，已完成 {{ rec.completed_count }} 项</span>
          <el-progress 
            :percentage="calculateProgress(rec)" 
            :status="rec.status === 1 ? 'success' : undefined"
          />
        </div>
        <el-table v-if="hasItems(rec)" :data="parseJson(rec.items)" stripe border>
          <el-table-column prop="item_name" label="项目名称" />
          <el-table-column prop="item_description" label="项目描述" />
          <el-table-column prop="item_value" label="检测结果" />
        </el-table>
        <div v-else class="no-items">暂无体检项目</div>
        <div class="footer">
          <el-button type="text" @click="viewPlanItems(rec.plan_id)">查看套餐详情</el-button>
        </div>
        <el-button type="success" @click="addItem(rec.plan_id)" style="margin-bottom: 20px;">创建评论</el-button>
        <el-table :data="rec.commentItems">
          <el-table-column prop="Commentary" label="评论" />
          <el-table-column label="创建时间" width="180">
            <template #default="scope">
              {{ formatDate(scope.row.CreatedAt) }}
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
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'


interface RecordItem {
  plan_id: number
  institution_id: number
  institution_name: string
  plan_name: string
  status: number // 0: pending, 1: completed
  items: string | null
  item_count: number
  completed_count: number
  commentItems: any
}

const records = ref<RecordItem[]>([])
const router = useRouter()

const parseJson = (str: string | null): Array<{ item_name: string; item_description: string; item_value: string }> => {
  if (!str) return []
  try {
    const parsed = JSON.parse(str)
    return Array.isArray(parsed) ? parsed : []
  } catch {
    return []
  }
}

const hasItems = (rec: RecordItem): boolean => {
  const arr = parseJson(rec.items)
  return arr.length > 0
}

const calculateProgress = (rec: RecordItem): number => {
  if (rec.item_count === 0) return 0
  return Math.round((rec.completed_count / rec.item_count) * 100)
}

const goDashboard = (): void => {
  router.push('/dashboard')
}

const viewPlanItems = (planId: number): void => {
  router.push({ name: 'plan-items', query: { plan_id: planId } })
}

const uid = ref(localStorage.getItem('uid'))
onMounted(async () => {
  await getUserView()
})

const getUserView = async () => {
  const response = await axios.get('/api/userview/')
  records.value = response.data.records
  records.value.map(async(item) => {
    item.commentItems = await fetchCommentItems(item)
  })
}
const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString()
}  

const fetchCommentItems = async (item: any) => {
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/commentary/get/plan/${item.plan_id}`, {
      headers: { Authorization: `${token}` }
    })
    const items = response.data.commentaries || []
    const commentItems = items.map((item: any) => ({
      ID: item.ID || item.id,
      UserId: item.UserId || item.user_id,
      PlanId: item.PlanId || item.plan_id,
      Commentary: item.Commentary || item.commentary,
      CreatedAt: item.CreatedAt || item.created_at || ''
    }))
    return commentItems
  } catch (error) {
    console.error('Failed to fetch health items:', error)
  } finally {
  }
}

const editDialogVisible = ref(false)
const editForm = reactive({
  ID: null,
  UserId: null,
  PlanId: null,
  Commentary: ''
})
// 新增
const addItem = (planId) => {
  editForm.UserId = Number(uid.value)
  editForm.PlanId = planId
  editForm.Commentary = ''
  editDialogVisible.value = true
}

// 保存
const saveItemChanges = async () => {
  try {
    const token = localStorage.getItem('jwt')
    await axios.post('/api/commentary/add', {
      user_id: editForm.UserId,
      plan_id: editForm.PlanId,
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
    await getUserView()
  } catch (error: any) {
    console.error('Failed to update health item:', error)
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('评论新增失败')
    }
  }
}
</script>

<style scoped>
.health-records-container {
  padding: 20px;
}
.back-button {
  margin-bottom: 20px;
}
.cards {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}
.record-card {
  width: calc(50% - 20px);
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-badge {
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  color: white;
}

.status-badge.completed {
  background-color: #67C23A;
}

.status-badge.pending {
  background-color: #E6A23C;
}

.progress-info {
  margin: 10px 0 15px 0;
}

.item-count {
  color: #606266;
  font-size: 14px;
  display: block;
  margin-bottom: 5px;
}
.no-items {
  padding: 20px;
  color: #909399;
  text-align: center;
}
.footer {
  text-align: right;
  margin-top: 10px;
}
</style>
