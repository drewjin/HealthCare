<template>
  <div>
    <h1>套餐体检项目</h1>
    <table v-if="planItems.length">
      <thead>
        <tr>
          <th>套餐名称</th>
          <th>项目名称</th>
          <th>项目值</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(item, index) in planItems" :key="index">
          <td>{{ item.plan_name }}</td>
          <td>{{ item.item_name }}</td>
          <td>{{ item.item_value }}</td>
        </tr>
      </tbody>
    </table>
    <p v-else>暂无体检项目</p>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'

interface PlanItem {
  plan_name: string
  item_name: string
  item_value: string
}

const planItems = ref<PlanItem[]>([])
const route = useRoute()

onMounted(async () => {
  const planId = Number(route.query.plan_id)
  const response = await axios.get('/api/userview/plan', { params: { plan_id: planId } })
  planItems.value = response.data.plan_items
})
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1em;
}
th, td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}
th {
  background-color: #f2f2f2;
}
</style>
