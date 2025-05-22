<template>
  <div class="health-view">
    <div class="header">
      <h2>亲友健康数据</h2>
      <el-button type="primary" icon="ArrowLeft" @click="goBack">返回</el-button>
    </div>
    
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="4" animated />
    </div>
    <template v-else>
      <div v-if="familyMemberName" class="member-info">
        <h3>{{ familyMemberName }} 的健康指标</h3>
      </div>
      <v-chart v-if="chartOption" :option="chartOption" autoresize style="height: 400px;" />
      <el-empty v-else description="暂无健康数据" />
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'

// 注册必要的 ECharts 组件
use([
  BarChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
  CanvasRenderer
])

const route = useRoute()
const router = useRouter()
const loading = ref(true)
const chartOption = ref<any>(null)
const chartData = ref<any[]>([])
const familyMemberName = ref('')
const familyMemberId = ref<string | number>('')

// 返回上一页
const goBack = () => {
  router.back()
}

// 获取亲友健康数据
const fetchHealthItems = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    const id = route.params.id as string
    familyMemberId.value = id
    
    // 先获取亲友的姓名信息
    try {
      const userResponse = await axios.get(`/api/users/${id}/profile`, {
        headers: { Authorization: `${token}` }
      })
      if (userResponse.data && userResponse.data.data) {
        familyMemberName.value = userResponse.data.data.name || '亲友'
      }
    } catch (e) {
      console.error('获取用户信息失败:', e)
      familyMemberName.value = '亲友'
    }
    
    console.log('Fetching health items for family member ID:', id)
    const response = await axios.get(`/api/healthitems/byid/${id}`, {
      headers: { Authorization: `${token}` }
    })

    if (response.data.items && response.data.items.length > 0) {
      const allMetricsMap: Record<string, number> = {}

      response.data.items.forEach((item: any) => {
        const healthInfo = item.user_health_info || ''
        if (typeof healthInfo !== 'string') {
          console.warn('Unexpected healthInfo type:', typeof healthInfo, healthInfo)
          return
        }
        
        const metrics = healthInfo.split(/[;,]/)
          .filter(Boolean)
          .map((str: string) => {
            const match = str.match(/([^:：]+)[：:]?(\d+)/)
            if (match) {
              return { name: match[1].trim(), value: Number(match[2]) }
            }
            return null
          })
          .filter(Boolean) as { name: string, value: number }[]

        metrics.forEach(({ name, value }) => {
          if (allMetricsMap[name]) {
            allMetricsMap[name] += value
          } else {
            allMetricsMap[name] = value
          }
        })
      })

      const averageHealthData: Record<string, number> = {
        '身高': 170,
        '体重': 65,
        '血压': 120,
        '血糖': 5.2,
        '心率': 75,
        '体温': 36.5,
        '血氧': 98,
        '肺活量': 4000,
        '胆固醇': 4.5,
        '尿酸': 360
      }

      const aggregatedMetrics = Object.entries(allMetricsMap).map(([name, value]) => ({
        name,
        userValue: value,
        averageValue: averageHealthData[name] || Math.round(value * (0.8 + Math.random() * 0.4))
      }))

      chartData.value = aggregatedMetrics

      const categories = aggregatedMetrics.map(item => item.name)
      const userValues = aggregatedMetrics.map(item => item.userValue)
      const averageValues = aggregatedMetrics.map(item => item.averageValue)

      const unitMap: Record<string, string> = {
        '身高': 'cm',
        '体重': 'kg',
        '血压': 'mmHg',
        '血糖': 'mmol/L',
        '心率': '次/分',
        '体温': '°C',
        '血氧': '%',
        '肺活量': 'mL',
        '胆固醇': 'mmol/L',
        '尿酸': 'μmol/L'
      }

      chartOption.value = {
        title: {
          text: '健康指标对比图',
          left: 'center'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: { type: 'shadow' },
          formatter: function(params: any[]) {
            const name = params[0].name
            const unit = unitMap[name] || ''
            let result = `<div style="font-weight:bold;">${name}</div>`
            params.forEach((param: any) => {
              result += `<div>${param.marker} ${param.seriesName}: ${param.value}${unit}</div>`
            })
            return result
          }
        },
        legend: {
          data: ['您的数值', '平均水平'],
          top: 'bottom'
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '10%',
          containLabel: true
        },
        xAxis: {
          type: 'category',
          data: categories,
          axisLabel: { rotate: 30 }
        },
        yAxis: {
          type: 'value',
          name: '数值'
        },
        series: [
          {
            name: '您的数值',
            type: 'bar',
            data: userValues,
            itemStyle: { color: '#67c23a' }
          },
          {
            name: '平均水平',
            type: 'bar',
            data: averageValues,
            itemStyle: { color: '#409eff' }
          }
        ]
      }
    } else {
      chartData.value = []
      chartOption.value = null
    }
  } catch (error) {
    console.error('Failed to fetch health items:', error)
    ElMessage.error('获取健康数据失败')
    chartData.value = []
    chartOption.value = null
  } finally {
    loading.value = false
  }
}

onMounted(fetchHealthItems)
watch(() => route.params.id, fetchHealthItems)
</script>

<style scoped>
.health-view {
  margin-top: 20px;
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.member-info {
  margin-bottom: 20px;
  color: #606266;
}

.loading-state {
  padding: 40px 0;
  display: flex;
  justify-content: center;
}
</style>
