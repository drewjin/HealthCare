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
            <span class="item-count">共 {{ parseJson(rec.items).length }} 项</span>
          </div>
        </template>
        <el-table v-if="hasItems(rec)" :data="parseJson(rec.items)" stripe border>
          <el-table-column prop="item_name" label="项目名称" />
          <el-table-column prop="item_value" label="项目值" />
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
  institution_name: string
  plan_name: string
  items: string | null
}

const records = ref<RecordItem[]>([])
const router = useRouter()

const parseJson = (str: string | null): Array<{ item_name: string; item_value: string }> => {
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
.item-count {
  color: #909399;
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
