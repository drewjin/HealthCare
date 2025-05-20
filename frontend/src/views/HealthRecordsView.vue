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
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

interface RecordItem {
  plan_id: number
  institution_id: number
  institution_name: string
  plan_name: string
  status: number // 0: pending, 1: completed
  items: string | null
  item_count: number
  completed_count: number
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

onMounted(async () => {
  const response = await axios.get('/api/userview/')
  records.value = response.data.records
})
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
