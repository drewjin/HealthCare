<template>
  <div class="family-health-view">
    <div v-if="loading" class="loading-state">
      <el-skeleton :rows="3" animated />
    </div>
    <div v-else-if="error" class="error-state">
      <el-alert :title="error" type="error" />
    </div>
    <div v-else-if="items.length === 0" class="empty-state">
      <el-empty description="暂无健康指标数据" />
    </div>
    <div v-else>
      <div class="health-item-description">
        <p style="color: #333; margin-bottom: 15px;">以下是 {{ familyMemberName }} 的健康指标数据，与平均水平对比：</p>
      </div>
      <el-form label-position="left" label-width="120px">
        <el-form-item v-for="item in items" :key="item.name" :label="item.name">
          <div class="item-value">
            <div style="color: #303133; font-weight: bold; font-size: 16px;">
              {{ item.userValue !== undefined ? item.userValue : (item.value !== undefined ? item.value : '-') }}
              <span class="unit-text" v-if="getUnit(item.name)">{{ getUnit(item.name) }}</span>
            </div>
            <div v-if="item.averageValue" class="average-value">
              平均水平: {{ item.averageValue }} {{ getUnit(item.name) }}
            </div>
          </div>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import axios from 'axios'

interface HealthItem {
  name: string;
  value?: number | string;
  userValue?: number | string;
  averageValue?: number | string;
}

const props = defineProps({
  familyMemberId: {
    type: [String, Number],
    required: true
  },
  familyMemberName: {
    type: String,
    default: '亲友'
  }
})

const loading = ref(false)
const error = ref('')
const items = ref<HealthItem[]>([])

// 根据指标名称获取单位
const getUnit = (name: string): string => {
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
  
  return unitMap[name] || ''
}

// 加载亲友健康数据
const loadFamilyHealthData = async () => {
  if (!props.familyMemberId) {
    error.value = '亲友ID不能为空'
    return
  }
  
  try {
    loading.value = true
    const token = localStorage.getItem('jwt')
    if (!token) {
      error.value = '未登录或会话已过期'
      return
    }
    
    // 获取亲友的健康数据，使用亲友的ID，而不是当前用户的ID
    const response = await axios.get(`/api/healthitems/byid/${props.familyMemberId}`, {
      headers: { Authorization: token }
    })
    
    if (response.data.items && response.data.items.length > 0) {
      console.log('健康数据:', response.data.items)
      
      // 汇总所有健康指标
      const allMetricsMap: Record<string, number> = {};
      
      response.data.items.forEach((item: any) => {
        const healthInfo = item.user_health_info;
        
        // 使用所有可能的分隔符处理
        const metrics = healthInfo.split(/[;,]/)
          .filter(Boolean)
          .map((str: string) => {
            const match = str.match(/([^:：]+)[：:]?(\d+)/);
            if (match) {
              return {
                name: match[1].trim(),
                value: Number(match[2])
              };
            }
            return null;
          })
          .filter(Boolean) as { name: string, value: number }[];
          
        // 合并同名指标的值（可以自定义策略：求和、取平均等）
        metrics.forEach(({ name, value }) => {
          if (allMetricsMap[name]) {
            allMetricsMap[name] += value; // 简单求和，可按需改为平均值
          } else {
            allMetricsMap[name] = value;
          }
        });
      });
      
      // 生成平均值参考数据
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
      };
      
      // 将用户数据与平均数据组合
      const healthItems = Object.entries(allMetricsMap).map(([name, value]) => {
        return {
          name,
          userValue: value,
          averageValue: averageHealthData[name as keyof typeof averageHealthData] || Math.round(value * (0.8 + Math.random() * 0.4)) // 如果没有预设值，随机生成一个平均值
        };
      });
      
      console.log('亲友健康指标数据:', healthItems);
      items.value = healthItems;
    } else {
      items.value = [];
      console.log('未找到亲友健康数据');
    }
  } catch (err) {
    console.error('加载健康数据失败:', err)
    error.value = '加载健康数据失败，请重试'
    items.value = [];
  } finally {
    loading.value = false
  }
}

// 监听家庭成员ID变化
watch(() => props.familyMemberId, () => {
  if (props.familyMemberId) {
    loadFamilyHealthData()
  }
}, { immediate: true })

// 初始化
onMounted(() => {
  if (props.familyMemberId) {
    loadFamilyHealthData()
  }
})
</script>

<style scoped>
.family-health-view {
  max-width: 800px;
  margin: 0 auto;
}

.loading-state, .error-state, .empty-state {
  padding: 20px;
  text-align: center;
}

.item-value {
  padding: 8px 0;
  color: #606266;
  font-size: 16px;
}

.unit-text {
  color: #909399;
  font-size: 0.85em;
  margin-left: 4px;
}

.average-value {
  margin-top: 4px;
  color: #5470c6;
  font-size: 13px;
}

.health-item-description {
  background-color: #f5f7fa;
  padding: 10px 15px;
  border-radius: 4px;
  margin-bottom: 20px;
  border-left: 4px solid #409eff;
}
</style>
