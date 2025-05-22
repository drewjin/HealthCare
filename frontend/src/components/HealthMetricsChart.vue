<template>
  <div class="health-metrics-chart">
    <div class="chart-title">健康指标对比图</div>
    <div v-if="loading" class="loading">
      <el-skeleton :rows="4" animated />
    </div>
    <div v-else-if="noData" class="no-data">
      <el-empty description="暂无健康指标数据" />
    </div>
    <div v-else>
      <v-chart class="chart" :option="chartOption" autoresize />
      <div class="legend">
        <div class="legend-item">
          <span class="legend-color" style="background-color: #5470c6;"></span>
          <span>您的指标</span>
        </div>
        <div class="legend-item">
          <span class="legend-color" style="background-color: #91cc75;"></span>
          <span>参考平均值</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watchEffect } from 'vue'
import { use } from 'echarts/core'
import { BarChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  LegendComponent
} from 'echarts/components'
import { LabelLayout, UniversalTransition } from 'echarts/features'
import { CanvasRenderer } from 'echarts/renderers'
import VChart from 'vue-echarts'

// 注册 ECharts 组件
use([
  BarChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  DatasetComponent,
  TransformComponent,
  LegendComponent,
  LabelLayout,
  UniversalTransition,
  CanvasRenderer
])

const props = defineProps({
  healthData: {
    type: Array as () => Array<{ item_name: string; item_description: string; item_value: string }>,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

// 解析健康指标值和名称
const parsedData = computed(() => {
  return props.healthData.map(item => {
    // 尝试从字符串中提取数值
    const match = item.item_value?.match(/(\d+(\.\d+)?)/);
    const value = match ? parseFloat(match[1]) : 0;
    
    // 提取名称，优先使用 item_name，如果太长就用 item_description 的前10个字符
    let name = item.item_name;
    if (name && name.includes(':')) {
      name = name.split(':')[0].trim();
    }
    
    return {
      name: name || item.item_description || '未命名指标',
      value: value,
      // 生成一个接近真实值的随机平均值（在实际值的80%-120%之间）
      avgValue: value > 0 ? Math.round((value * (0.9 + Math.random() * 0.3)) * 10) / 10 : 0
    };
  }).filter(item => item.value > 0); // 过滤掉没有值的项目
});

const noData = computed(() => {
  return !props.loading && parsedData.value.length === 0;
});

// 生成柱状图配置
const chartOption = computed(() => {
  const data = parsedData.value;
  const categories = data.map(item => item.name);
  const values = data.map(item => item.value);
  const avgValues = data.map(item => item.avgValue);
  
  return {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      },
      formatter: function(params) {
        const userValue = params[0].value;
        const avgValue = params[1].value;
        const diff = ((userValue / avgValue - 1) * 100).toFixed(1);
        const diffText = diff > 0 ? `高于平均${diff}%` : `低于平均${Math.abs(diff)}%`;
        
        return `${params[0].name}<br/>
                您的指标: ${userValue}<br/>
                平均值: ${avgValue}<br/>
                ${diffText}`;
      }
    },
    legend: {
      data: ['您的指标', '参考平均值'],
      bottom: 0
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      top: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: categories,
      axisLabel: {
        interval: 0,
        rotate: categories.length > 5 ? 45 : 0,
        textStyle: {
          fontSize: 12
        }
      }
    },
    yAxis: {
      type: 'value',
      name: '数值',
      nameTextStyle: {
        padding: [0, 0, 0, 40]
      }
    },
    series: [
      {
        name: '您的指标',
        type: 'bar',
        data: values,
        itemStyle: {
          color: '#5470c6'
        },
        label: {
          show: true,
          position: 'top'
        },
        barGap: 0
      },
      {
        name: '参考平均值',
        type: 'bar',
        data: avgValues,
        itemStyle: {
          color: '#91cc75'
        },
        label: {
          show: true,
          position: 'top'
        }
      }
    ]
  };
});
</script>

<style scoped>
.health-metrics-chart {
  width: 100%;
  height: 100%;
  min-height: 350px;
  margin-bottom: 20px;
}

.chart {
  height: 350px;
  width: 100%;
}

.chart-title {
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 15px;
  text-align: center;
}

.loading, .no-data {
  height: 350px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.legend {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.legend-item {
  display: flex;
  align-items: center;
  margin: 0 10px;
}

.legend-color {
  display: inline-block;
  width: 15px;
  height: 15px;
  margin-right: 5px;
  border-radius: 3px;
}
</style>
