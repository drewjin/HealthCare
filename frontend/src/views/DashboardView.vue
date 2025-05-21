<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElContainer, ElAside, ElMain, ElHeader, ElMenu, ElMenuItem, ElCard, ElAvatar, ElRow, ElCol, ElButton, ElDivider, ElForm, ElFormItem, ElInput, ElMessage, ElDialog, ElTable, ElTableColumn, ElBadge, ElTimeline, ElTimelineItem, ElSelect, ElOption, ElCollapse, ElCollapseItem, ElSkeleton, ElEmpty, ElUpload } from 'element-plus'
import { User, List, Setting, UserFilled, HomeFilled, Tools, PictureFilled, Download } from '@element-plus/icons-vue'
import InstitutionForm from '@/components/InstitutionForm.vue'
import AdminReview from '@/components/AdminReview.vue'
import HealthItemForm from '@/components/HealthItemForm.vue'
import axios from 'axios'
import { use } from 'echarts/core';
import { BarChart } from 'echarts/charts';
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components';
import { CanvasRenderer } from 'echarts/renderers';
import VChart from 'vue-echarts';

// 健康数据相关变量
const healthItemsData = ref(null); // 存储原始健康项目数据
const parsedHealthItems = ref({}); // 解析后的健康项目键值对

const router = useRouter()
const uid = ref(localStorage.getItem('uid'))
const userType = ref(1) // Default to normal user
const activeMenu = ref('1')
const userInfo = ref({
  username: '',
  name: '',
  email: '',
  phone: '',
  gender: '',
  birthday: '',
  address: '',
  user_type: 1
})
const updateInfo = ref(false)
const showResetPwdDialog = ref(false)
const resetPwdForm = ref({
  prevPassword: '',
  newPassword: '',
  newPasswordConfirm: ''
})
const hasPendingRequests = ref(false)
const familyRequests = ref<Array<{
  id: number,
  requester: string,
  name: string,
  relationship: string,
  created_at: string
}>>([])
const confirmedFamilyMembers = ref<Array<{
  id: string;
  user_id: string | number;
  username: string;
  name: string;
  relationship: string;
}>>([])

const showCreateFamilyDialog = ref(false)
const familyRequestForm = ref({
  id: '',
  relative_username: '',
  relationship: ''
})

const relationshipOptions = [
  { value: '父亲', label: '父亲' },
  { value: '母亲', label: '母亲' },
  { value: '儿子', label: '儿子' },
  { value: '女儿', label: '女儿' },
  { value: '配偶', label: '配偶' },
]

const getUserTypeLabel = (type: number) => {
  switch (type) {
    case 1:
      return '普通用户'
    case 2:
      return '管理员'
    case 3:
      return '机构用户'
    default:
      return '未知类型'
  }
}

const genderOptions = [
  { value: 'M', label: '男' },
  { value: 'F', label: '女' }
]

interface InstitutionItem {
  ID: number
  institution_name: string
  institution_address: string
  institution_qualification: string
  examination_package: string
  status: number
}
const institutions = ref<InstitutionItem[]>([])
const loadingInstitutions = ref(false)
const activeInstitution = ref<number | string>('')

const fetchInstitutions = async () => {
  try {
    loadingInstitutions.value = true
    const rawToken = localStorage.getItem('jwt')
    if (!rawToken) {
      ElMessage.error('未登录或令牌丢失，无法获取机构列表')
      return
    }
    // Ensure single '' prefix
    const authToken = rawToken.startsWith('') ? rawToken : `${rawToken}`
    const response = await axios.get('/api/institutions', {
      headers: { Authorization: authToken }
    })
    
    // 处理不同的响应格式：可能是数组或单个对象
    console.log('获取到的机构数据:', response.data)
    
    try {
      if (Array.isArray(response.data)) {
        // 如果是数组，过滤并进行类型强制转换
        const validInstitutions = response.data
          .filter(inst => inst && typeof inst === 'object' && 'ID' in inst)
          .map(inst => inst as InstitutionItem);
        
        institutions.value = validInstitutions;
        console.log('获取到机构列表（数组）：', institutions.value.length)
      } else if (response.data && typeof response.data === 'object') {
        // 如果是单个对象，将其转换为数组
        if ('ID' in response.data) {
          institutions.value = [response.data as InstitutionItem];
          console.log('获取到单个机构，已转换为数组')
        } else {
          // 尝试从对象中提取机构
          try {
            const values = Object.values(response.data);
            const validInstitutions = values
              .filter(item => item && typeof item === 'object' && 'ID' in item)
              .map(item => item as InstitutionItem);
              
            if (validInstitutions.length > 0) {
              institutions.value = validInstitutions;
              console.log('从对象中提取机构数组：', validInstitutions.length)
            } else {
              throw new Error('未找到有效的机构数据');
            }
          } catch (extractError) {
            console.error('从对象提取机构列表失败:', extractError);
            institutions.value = [];
            ElMessage.warning('获取机构列表格式不正确');
          }
        }
      } else {
        // 无法识别的格式
        throw new Error('响应数据格式不正确');
      }
    } catch (formatError) {
      console.error('处理机构数据出错:', formatError);
      institutions.value = [];
      ElMessage.warning('获取机构列表格式不正确');
    }
  } catch (e) {
    console.error('获取机构列表出错：', e)
    institutions.value = []
    ElMessage.error('获取机构列表失败，请检查网络连接')
  } finally {
    loadingInstitutions.value = false
  }
}

watch(activeMenu, (val) => {
  if (val === '5') fetchInstitutions()
})

onMounted(async () => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('未登录')

    // 获取用户信息
    const response = await fetch(`http://localhost:3000/api/users/${uid.value}/profile`, {
      method: 'GET',
      headers: { Authorization: token },
    })
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || '获取用户信息失败')
    }
    const userData = (await response.json()).data
    userInfo.value = userData
    userType.value = userData.user_type

    // Only fetch family-related data for normal users
    if (userType.value === 1) {
      // 获取待处理的家庭关系请求
      const requestsResponse = await fetch(`http://localhost:3000/api/family/pending/${uid.value}`, {
        method: 'GET',
        headers: { Authorization: token },
      })
      if (requestsResponse.ok) {
        const data = await requestsResponse.json()
        familyRequests.value = data || []
        hasPendingRequests.value = familyRequests.value && familyRequests.value.length > 0
      }
      await fetchFamilyItems()
      await fetchHealthItems()

    }
  } catch (error) {
    console.error(error)
    router.push('/login')
  }
})

const fetchFamilyItems = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/family/confirmed/${uid.value}`, {
      headers: { Authorization: `${token}` }
    })
    confirmedFamilyMembers.value = response.data || []
    console.log('Fetched family members:', confirmedFamilyMembers.value)
  } catch (error) {
    console.error('Failed to fetch family members:', error)
    ElMessage.error('获取家庭列表失败')
  }
}

const fetchHealthItems = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/healthitems/byid/${uid.value}`, {
      headers: { Authorization: `${token}` }
    })

    if (response.data.items && response.data.items.length > 0) {
      console.log('Health items data:', response.data.items)

      // 汇总所有健康指标
      const allMetricsMap: Record<string, number> = {};

      response.data.items.forEach(item => {
        const healthInfo = item.user_health_info;

        // 使用所有可能的分隔符处理
        const metrics = healthInfo.split(/[;,]/)
          .filter(Boolean)
          .map(str => {
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
      const averageHealthData = {
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
      const aggregatedMetrics = Object.entries(allMetricsMap).map(([name, value]) => {
        return {
          name,
          userValue: value,
          averageValue: averageHealthData[name] || Math.round(value * (0.8 + Math.random() * 0.4)) // 如果没有预设值，随机生成一个平均值
        };
      });

      console.log('Aggregated health metrics:', aggregatedMetrics);
      chartData.value = aggregatedMetrics;

      // 生成图表所需的数据格式
      const categories = aggregatedMetrics.map(item => item.name);
      const userValues = aggregatedMetrics.map(item => item.userValue);
      const averageValues = aggregatedMetrics.map(item => item.averageValue);

      // 单位映射
      const unitMap = {
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
      };

      chartOption.value = {
        title: {
          text: '健康指标对比图',
          left: 'center'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'shadow'
          },
          formatter: function(params) {
            const name = params[0].name;
            const unit = unitMap[name] || '';
            let result = `<div style="font-weight:bold;color:#333;margin-bottom:5px;">${name}</div>`;
            
            params.forEach(param => {
              const value = param.value;
              const seriesName = param.seriesName;
              const marker = param.marker;
              const color = param.color;
              result += `<div style="display:flex;align-items:center;margin:3px 0;">
                ${marker} <span style="color:${color};font-weight:bold;">${seriesName}: ${value}${unit}</span>
              </div>`;
            });
            
            return result;
          },
          backgroundColor: 'rgba(255, 255, 255, 0.9)',
          borderColor: '#ccc',
          borderWidth: 1,
          textStyle: {
            color: '#333'
          }
        },
        legend: {
          data: ['您的数值', '平均水平'],
          top: 'bottom',
          textStyle: {
            color: '#333',
            fontSize: 12
          }
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
          axisLabel: {
            interval: 0,
            rotate: 30,
            color: '#333',  // Dark text for better visibility
            fontSize: 12
          }
        },
        yAxis: {
          type: 'value',
          name: '数值',
          nameTextStyle: {
            color: '#333'  // Dark text for better visibility
          },
          axisLabel: {
            formatter: function(value: number) {
              return value.toString();
            },
            color: '#333'  // Dark text for better visibility
          }
        },
        series: [
          {
            name: '您的数值',
            type: 'bar',
            data: userValues,
            itemStyle: {
              color: '#67c23a'  // Bright green - more visible
            }
          },
          {
            name: '平均水平',
            type: 'bar',
            data: averageValues,
            itemStyle: {
              color: '#409eff'  // Bright blue - more visible
            }
          }
        ]
      };
    } else {
      console.log('No health items found');
      chartData.value = [];
    }
  } catch (error) {
    console.error('Failed to fetch health items:', error);
    ElMessage.error('获取健康数据失败');
    chartData.value = [];
  }
}


const handleLogout = () => {
  localStorage.removeItem('jwt')
  localStorage.removeItem('uid')
  router.push('/login')
}

// 查看亲友健康记录
const viewFamilyHealth = (familyMember: any) => {
  console.log('查看亲友健康数据:', familyMember)
  if (familyMember && familyMember.user_id) {
    router.push(`/family-health/${familyMember.user_id}`)
  } else {
    ElMessage.warning('无法获取亲友ID，请重试')
  }
}

const handleGoToProfile = () => {
  activeMenu.value = '1'
}

const handleResetPwd = async () => {
  try {
    if (resetPwdForm.value.newPassword !== resetPwdForm.value.newPasswordConfirm) {
      ElMessage.error('新密码不一致')
      return
    }

    const response = await fetch(`http://localhost:3000/api/users/${uid.value}/reset_pwd`, {
      method: 'PUT',
      headers: { 
        'Content-Type': 'application/json',
        'Authorization': localStorage.getItem('jwt') || ''
      },
      body: JSON.stringify({
        prev_password: resetPwdForm.value.prevPassword,
        new_password: resetPwdForm.value.newPassword,
        new_password_confirm: resetPwdForm.value.newPasswordConfirm
      })
    })

    const data = await response.json()
    if (!response.ok) throw new Error(data.error)

    ElMessage.success('密码修改成功，请重新登录')
    handleLogout()
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '密码修改失败')
  }
}

const handleUpdateInfo = async () => {
  try {
    const token = localStorage.getItem('jwt')
    await axios.post(`/api/users/${uid.value}/profile`, {
      username: userInfo.value.username,
      name: userInfo.value.name,
      gender: userInfo.value.gender,
      birthday: userInfo.value.birthday,
      phone: userInfo.value.phone,
      email: userInfo.value.email,
      address: userInfo.value.address,
    }, {
      headers: { 
        Authorization: `${token}`,
        'Content-Type': 'application/json; charset=utf-8'
      }
    })
    ElMessage.success('基本信息修改成功')
  } catch (error) {
    console.error('Failed to fetch health items:', error)
    ElMessage.error('基本信息修改失败')
  }
  updateInfo.value = false
}

const handleFamilyRequest = async (requestId: number, accept: boolean) => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('未登录')

    const response = await fetch(`http://localhost:3000/api/family/handle/${uid.value}/${requestId}`, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        Authorization: token 
      },
      body: JSON.stringify({ accept })
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error)
    }

    // 更新请求列表
    familyRequests.value = familyRequests.value.filter(req => req.id !== requestId)
    hasPendingRequests.value = familyRequests.value.length > 0
    ElMessage.success(accept ? '已接受请求' : '已拒绝请求')
    fetchFamilyItems()
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '处理请求失败')
  }
}

// const handleCreateFamilyRequest = async () => {
//   try {
//     const token = localStorage.getItem('jwt')
//     if (!token) throw new Error('未登录')

//     const response = await fetch(`http://localhost:3000/api/family/request/${uid.value}`, {
//       method: 'POST',
//       headers: { 
//         'Content-Type': 'application/json',
//         Authorization: token 
//       },
//       body: JSON.stringify(familyRequestForm.value)
//     })

//     if (!response.ok) {
//       const errorData = await response.json()
//       throw new Error(errorData.error)
//     }

//     ElMessage.success('请求已发送')
//     showCreateFamilyDialog.value = false
//     familyRequestForm.value = { relative_username: '', relationship: '' }
//   } catch (error) {
//     console.error(error)
//     ElMessage.error(error instanceof Error ? error.message : '发送请求失败')
//   }
// }

// OCR related variables and functions
const selectedFile = ref<File | null>(null)
const imageUrl = ref<string>('')
const ocrResults = ref<{item_name: string; item_value: string}[]>([])

const handleFileChange = (uploadFile: any) => {
  console.log('File selected (handleFileChange):', uploadFile);
  if (uploadFile && uploadFile.raw) {
    selectedFile.value = uploadFile.raw;
    console.log('selectedFile.value assigned:', selectedFile.value);
    const reader = new FileReader();
    reader.readAsDataURL(uploadFile.raw);
    reader.onload = () => {
      imageUrl.value = reader.result as string;
      console.log('Image URL set:', imageUrl.value);
    };
    reader.onerror = (error) => {
      console.error('FileReader error:', error);
    };
  } else {
    console.log('No file raw data found in uploadFile:', uploadFile);
    selectedFile.value = null;
    imageUrl.value = '';
  }
};

const submitOcr = async () => {
  console.log('Attempting to submit OCR. Current selectedFile:', selectedFile.value);
  if (!selectedFile.value) {
    ElMessage.error('请先选择图片');
    console.log('No file selected for submission.');
    return;
  }
  const formData = new FormData()
  formData.append('image', selectedFile.value)
  console.log('FormData to be sent. Image appended:', formData.get('image')); 
  try {
    const token = localStorage.getItem('jwt') || ''
    const response = await axios.post('/api/imageocr/solve', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: token
      }
    });

    const rawOcrResults: {item_name: string; item_value: string}[] = response.data.result || [];
    console.log('Raw OCR results from backend:', JSON.stringify(rawOcrResults)); // Log raw results

    ocrResults.value = rawOcrResults.map(entry => {
      const textToSplit = entry.item_value; // This is the string like "血压 124"
      console.log(`Processing entry.item_value: "${textToSplit}"`); 

      const firstSpaceIndex = textToSplit.indexOf(' ');
      let newName = '';
      let newValue = '';

      if (firstSpaceIndex !== -1) {
        newName = textToSplit.substring(0, firstSpaceIndex).trim();
        newValue = textToSplit.substring(firstSpaceIndex + 1).trim();
      } else {
        // If no space, assume the whole string is the item name
        newName = textToSplit.trim();
        newValue = ''; // Or you could set it to something like '-' or 'N/A' to indicate missing value
      }
      console.log(`Parsed to: name="${newName}", value="${newValue}"`); 
      return { item_name: newName, item_value: newValue };
    });

  } catch (e) {
    ElMessage.error('OCR 识别失败，请重试')
    console.error('OCR 识别失败', e)
  }
}

const dialogTitle = ref('')
// 编辑
interface FamilyMember {
  id?: string;
  username?: string;
  relationship?: string;
}

const editItem = (item?: FamilyMember) => {
  if(item) {
    dialogTitle.value = '编辑家庭关系请求'
    familyRequestForm.value.id = item.id || ''
    familyRequestForm.value.relative_username = item.username || ''
    familyRequestForm.value.relationship = item.relationship || ''
  } else {
    dialogTitle.value = '创建家庭关系请求'
    familyRequestForm.value.id = ''
    familyRequestForm.value.relative_username = ''
    familyRequestForm.value.relationship = ''
  }
  showCreateFamilyDialog.value = true
}
// 保存
const handleCreateFamilyRequest = async () => {
  try {
    const token = localStorage.getItem('jwt')
    if(familyRequestForm.value.id) {
      await axios.post(`/api/family/update_family_name`, {
        id: familyRequestForm.value.id,
        relationship: familyRequestForm.value.relationship
      }, {
        headers: { 
          Authorization: `${token}`,
          'Content-Type': 'application/json; charset=utf-8'
        }
      })
      ElMessage.success('家庭关系修改成功')
    } else {
      await axios.post(`/api/family/request/${uid.value}`, {
        relative_username: familyRequestForm.value.relative_username,
        relationship: familyRequestForm.value.relationship
      }, {
        headers: { 
          Authorization: `${token}`,
          'Content-Type': 'application/json'
        }
      })
      ElMessage.success('家庭关系新增成功')
    }

    showCreateFamilyDialog.value = false

    // 更新列表
    await fetchFamilyItems()
  } catch (error: any) {
    console.error('Failed to update health item:', error)
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('家庭关系修改失败')
    }
  }
}
// 删除
const deleteItem = async (key: string) => {
  try {
    const token = localStorage.getItem('jwt')
    await axios.get(`/api/family/del_confirmed/${key}`, {
      headers: { Authorization: token }
    })
    ElMessage.success('家庭关系删除成功')
    // 更新列表
    await fetchFamilyItems()
  } catch (error: any) {
    console.error('Failed to update health item:', error)
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('档案删除失败')
    }
  }
}


// 按需注册必要的组件
use([
  CanvasRenderer,
  BarChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
]);

// 图表配置项
const chartOption = ref();
const chartData = ref([])

// 导出健康数据为CSV格式
const exportHealthData = () => {
  if (chartData.value.length === 0) {
    ElMessage.warning('没有可导出的健康数据');
    return;
  }
  
  // 获取当前日期作为文件名
  const date = new Date();
  const fileName = `健康数据_${date.getFullYear()}-${date.getMonth()+1}-${date.getDate()}.csv`;
  
  // 单位映射
  const unitMap = {
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
  };
  
  // 创建CSV内容
  let csvContent = '健康指标,您的数值,平均水平,单位\n';
  
  // 使用类型断言确保TypeScript知道每个项目的类型
  interface ChartDataItem {
    name: string;
    userValue?: number | string;
    averageValue?: number | string;
    value?: number | string;
  }
  
  (chartData.value as ChartDataItem[]).forEach(item => {
    const unit = unitMap[item.name as keyof typeof unitMap] || '';
    const userVal = item.userValue !== undefined ? item.userValue : (item.value || '-');
    const avgVal = item.averageValue || '-';
    csvContent += `${item.name},${userVal},${avgVal},${unit}\n`;
  });
  
  // 创建下载链接
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.setAttribute('href', url);
  link.setAttribute('download', fileName);
  link.style.visibility = 'hidden';
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}

</script>

<template>
  <ElContainer class="dashboard-container">
    <ElHeader class="dashboard-header">
      <div class="header-content">
        <div class="header-left">
          <h2>健康管理系统</h2>
          <div class="notification-wrapper">
            <ElButton type="primary" @click="handleGoToProfile" plain>
              <el-icon><User /></el-icon>
              个人信息
            </ElButton>
            <div v-if="hasPendingRequests && userType === 1" class="notification-dot"></div>
          </div>
        </div>
        <ElButton type="danger" @click="handleLogout">退出登录</ElButton>
      </div>
    </ElHeader>
    
    <ElContainer>
      <ElAside width="200px">
        <ElMenu
          default-active="1"
          class="dashboard-menu"
          @select="activeMenu = $event"
        >
          <ElMenuItem index="1">
            <el-icon><User /></el-icon>
            <span>个人资料</span>
          </ElMenuItem>
          
          <!-- 机构用户菜单 -->
          <template v-if="userType === 3">
            <ElMenuItem index="2">
              <el-icon><HomeFilled /></el-icon>
              <span>机构信息</span>
            </ElMenuItem>
            <ElMenuItem index="6">
              <el-icon><List /></el-icon>
              <span>用户套餐管理</span>
            </ElMenuItem>
          </template>

          <!-- 管理员用户菜单 -->
          <template v-if="userType === 2">
            <ElMenuItem index="2">
              <el-icon><Tools /></el-icon>
              <span>机构审核</span>
            </ElMenuItem>
          </template>

          <!-- 普通用户菜单 -->
          <template v-if="userType === 1">
            <ElMenuItem index="2">
              <el-icon><List /></el-icon>
              <span>健康记录</span>
            </ElMenuItem>
            <ElMenuItem index="3">
              <el-icon><Setting /></el-icon>
              <span>系统设置</span>
            </ElMenuItem>
            <ElMenuItem index="4">
              <el-icon><UserFilled /></el-icon>
              <span>家庭关系</span>
            </ElMenuItem>
            <!-- 新增体检机构列表菜单 -->
            <ElMenuItem index="5">
              <el-icon><HomeFilled /></el-icon>
              <span>机构列表</span>
            </ElMenuItem>
            <!-- OCR 识别 菜单 -->
            <ElMenuItem index="7">
              <el-icon><PictureFilled /></el-icon>
              <span>OCR 识别</span>
            </ElMenuItem>
          </template>
          
          <!-- 机构管理菜单项 -->
          <ElMenuItem index="6" v-if="userType === 3 || userType === 2" @click="router.push('/institution-manage')">
            <el-icon><tools /></el-icon>
            <span>机构管理</span>
          </ElMenuItem>
          
          <!-- 套餐健康项目管理菜单项 -->
          <ElMenuItem index="8" v-if="userType === 3 || userType === 2" @click="router.push('/plan-health-item-manager')">
            <el-icon><List /></el-icon>
            <span>套餐项目管理</span>
          </ElMenuItem>
          
          <!-- 添加用户体检数据 -->
          <ElMenuItem index="8" v-if="userType === 3 || userType === 2" @click="router.push('/add-user-data')">
            <el-icon><List /></el-icon>
            <span>添加用户体检数据</span>
          </ElMenuItem>
          
          <!-- 用户套餐管理 -->
          <ElMenuItem index="11" v-if="userType === 3 || userType === 2" @click="router.push('/user-packages')">
            <el-icon><List /></el-icon>
            <span>用户套餐管理</span>
          </ElMenuItem>
          
          <!-- 管理健康检查项目 -->
          <ElMenuItem index="9" v-if="userType === 3 || userType === 2" @click="router.push('/health-items')">
            <el-icon><List /></el-icon>
            <span>健康检查项目管理</span>
          </ElMenuItem>
          
          <!-- 检查项目存储管理 -->
          <ElMenuItem index="10" v-if="userType === 3 || userType === 2" @click="router.push('/health-item-manager')">
            <el-icon><List /></el-icon>
            <span>检查项目存储管理</span>
          </ElMenuItem>

          <!-- 档案管理 -->
          <ElMenuItem index="11" v-if="userType === 3 || userType === 2" @click="router.push('/archive-manager')">
            <el-icon><List /></el-icon>
            <span>档案管理</span>
          </ElMenuItem>

          <!-- 评论管理 -->
          <ElMenuItem index="11" v-if="userType === 3 || userType === 2" @click="router.push('/comment-manager')">
            <el-icon><List /></el-icon>
            <span>评论管理</span>
          </ElMenuItem>
        </ElMenu>
      </ElAside>
      
      <ElMain>
        <!-- 个人资料面板 -->
        <div v-if="activeMenu === '1'">
          <ElRow :gutter="20">
            <ElCol :span="8">
              <ElCard shadow="hover">
                <div class="user-profile-card">
                  <ElAvatar :size="100" src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" />
                  <h3>{{ userInfo.name }}</h3>
                  <p>ID: {{ uid }}</p>
                  <ElButton type="primary" @click="showResetPwdDialog = true">修改密码</ElButton>
                </div>
              </ElCard>
            </ElCol>
            <ElCol :span="16">
              <ElCard shadow="hover">
                <template #header>
                  <div class="card-header">
                    <span>基本信息</span>
                  </div>
                </template>
                <div class="user-info-list">
                  <div class="info-item">
                    <span class="label">用户名：</span>
                    <el-input v-if="updateInfo" v-model="userInfo.username" placeholder="请输入" />
                    <span v-else>{{ userInfo.username }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">邮箱：</span>
                    <el-input v-if="updateInfo" v-model="userInfo.email" placeholder="请输入" />
                    <span v-else>{{ userInfo.email }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">电话：</span>
                    <el-input v-if="updateInfo" v-model="userInfo.phone" placeholder="请输入" />
                    <span v-else>{{ userInfo.phone }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">性别：</span>
                    <el-select v-if="updateInfo" v-model="userInfo.gender" placeholder="请选择">
                      <el-option v-for="item in genderOptions" :key="item.value" :label="item.label" :value="item.value" />
                    </el-select>
                    <span v-else>{{ userInfo.gender === 'M' ? '男' : '女' }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">生日：</span>
                    <el-input v-if="updateInfo" v-model="userInfo.birthday" placeholder="请输入" />
                    <span v-else>{{ userInfo.birthday }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">地址：</span>
                    <el-input v-if="updateInfo" v-model="userInfo.address" placeholder="请输入" />
                    <span v-else>{{ userInfo.address }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">用户类型：</span>
                    <span>{{ getUserTypeLabel(userInfo.user_type) }}</span>
                  </div>
                </div>
                <ElButton v-if="!updateInfo" type="primary" @click="updateInfo = true">修改</ElButton>
                <div v-else>
                  <ElButton @click="updateInfo = false">取消</ElButton>
                  <ElButton type="primary" @click="handleUpdateInfo">保存</ElButton>
                </div>
              </ElCard>
            </ElCol>
          </ElRow>

          <!-- 修改密码对话框 -->
          <ElDialog
            v-model="showResetPwdDialog"
            title="修改密码"
            width="30%"
            :close-on-click-modal="false"
          >
            <ElForm :model="resetPwdForm" label-width="100px">
              <ElFormItem label="当前密码">
                <ElInput v-model="resetPwdForm.prevPassword" type="password" show-password />
              </ElFormItem>
              <ElFormItem label="新密码">
                <ElInput v-model="resetPwdForm.newPassword" type="password" show-password />
              </ElFormItem>
              <ElFormItem label="确认新密码">
                <ElInput v-model="resetPwdForm.newPasswordConfirm" type="password" show-password />
              </ElFormItem>
            </ElForm>
            <template #footer>
              <span class="dialog-footer">
                <ElButton @click="showResetPwdDialog = false">取消</ElButton>
                <ElButton type="primary" @click="handleResetPwd">
                  确认修改
                </ElButton>
              </span>
            </template>
          </ElDialog>
        </div>

        <!-- 机构用户面板 -->
        <div v-if="activeMenu === '2' && userType === 3">
          <ElCard shadow="hover">
            <template #header>
              <div class="card-header">
                <span>机构信息管理</span>
              </div>
            </template>
            <InstitutionForm @submit-success="activeMenu = '1'" />
          </ElCard>
        </div>

        <!-- 管理员面板 -->
        <div v-if="activeMenu === '2' && userType === 2">
          <AdminReview />
        </div>          
        <!-- 普通用户其他面板 -->
        <template v-if="userType === 1">
          <!-- 健康记录面板 -->
          <div v-if="activeMenu === '2'">
            <ElCard shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>健康记录</span>
                  <ElButton type="primary" @click="router.push('/health-records')">查看全部</ElButton>
                </div>
              </template>
              <div class="health-summary">
                <p>这里展示您的健康记录摘要</p>
                <div v-if="chartData.length > 0" class="health-data-section">
                  <!-- 图表展示 -->
                  <div class="chart-section">
                    <v-chart class="chart-container" :option="chartOption" autoresize />
                  </div>
                  
                  <!-- 表单展示 - 使用新组件 -->
                  <div class="form-section">
                    <div class="form-header">
                      <h3>健康指标详情</h3>
                      <el-button type="primary" size="small" @click="exportHealthData" plain>
                        <el-icon><Download /></el-icon> 导出数据
                      </el-button>
                    </div>
                    <HealthItemForm :items="chartData" :editable="false" />
                  </div>
                </div>
                <div v-else class="empty-data">
                  <el-empty description="暂无健康数据" />
                </div>
                <el-button type="primary" plain @click="router.push('/health-records')" style="margin-top: 20px;">查看详细健康记录</el-button>
              </div>
            </ElCard>
          </div>

          <!-- 系统设置面板 -->
          <div v-if="activeMenu === '3'">
            <ElCard shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>系统设置</span>
                </div>
              </template>
              <p>系统设置内容将在此显示</p>
            </ElCard>
          </div>

          <!-- 家庭关系面板 -->
          <div v-if="activeMenu === '4'">
            <ElRow :gutter="20">
              <ElCol :span="12">
                <ElCard shadow="hover">
                  <template #header>
                    <div class="card-header">
                      <span>已有家庭关系</span>
                      <ElButton type="primary" @click="editItem()">添加成员</ElButton>
                    </div>
                  </template>
                  <div v-if="confirmedFamilyMembers.length > 0">
                    <ElTable :data="confirmedFamilyMembers" style="width: 100%">
                      <ElTableColumn prop="username" label="用户名" />
                      <ElTableColumn prop="name" label="姓名" />
                      <ElTableColumn prop="relationship" label="关系" />
                      <ElTableColumn label="操作" width="240">
                        <template #default="scope">
                          <el-button size="small" type="primary" @click="editItem(scope.row)">编辑</el-button>
                          <el-button size="small" type="danger" @click="deleteItem(scope.row.id)">删除</el-button>
                          <el-button size="small" type="success" @click="viewFamilyHealth(scope.row)">查看健康</el-button>
                        </template>
                      </ElTableColumn>
                    </ElTable>
                  </div>
                  <div v-else>
                    <p>暂无已确认的家庭关系。</p>
                  </div>
                </ElCard>
              </ElCol>
              
              <ElCol :span="12">
                <ElCard shadow="hover">
                  <template #header>
                    <div class="card-header">
                      <span>待处理请求</span>
                      <ElBadge :value="familyRequests.length" :hidden="!hasPendingRequests">
                        <ElButton type="primary" plain>查看全部</ElButton>
                      </ElBadge>
                    </div>
                  </template>
                  <div v-if="hasPendingRequests">
                    <ElTimeline>
                      <ElTimelineItem
                        v-for="request in familyRequests"
                        :key="request.id"
                        :timestamp="request.created_at"
                        placement="top"
                      >
                        <ElCard>
                          <h4>来自 {{ request.requester }} ({{ request.name }}) 的请求</h4>
                          <p>关系：{{ request.relationship }}</p>
                          <div class="request-actions">
                            <ElButton type="success" size="small" @click="handleFamilyRequest(request.id, true)">接受</ElButton>
                            <ElButton type="danger" size="small" @click="handleFamilyRequest(request.id, false)">拒绝</ElButton>
                          </div>
                        </ElCard>
                      </ElTimelineItem>
                    </ElTimeline>
                  </div>
                  <div v-else>
                    <p>暂无待处理的家庭关系请求。</p>
                  </div>
                </ElCard>
              </ElCol>
            </ElRow>

            <!-- 创建家庭关系请求对话框 -->
            <ElDialog
              v-model="showCreateFamilyDialog"
              :title="dialogTitle"
              width="30%"
              :close-on-click-modal="false"
            >
              <ElForm :model="familyRequestForm" label-width="100px">
                <ElFormItem label="成员用户名">
                  <ElInput v-model="familyRequestForm.relative_username" />
                </ElFormItem>
                <ElFormItem label="关系">
                  <ElSelect v-model="familyRequestForm.relationship" placeholder="请选择关系">
                    <ElOption
                      v-for="item in relationshipOptions"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value"
                    />
                  </ElSelect>
                </ElFormItem>
              </ElForm>
              <template #footer>
                <span class="dialog-footer">
                  <ElButton @click="showCreateFamilyDialog = false">取消</ElButton>
                  <ElButton type="primary" @click="handleCreateFamilyRequest">
                    发送请求
                  </ElButton>
                </span>
              </template>
            </ElDialog>
          </div>

          <!-- 体检机构列表面板 -->
          <div v-if="activeMenu === '5'">
            <ElCard shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>体检机构列表</span>
                  <ElButton type="primary" size="small" @click="router.push('/institutions')">查看全部</ElButton>
                </div>
              </template>
              <div v-if="loadingInstitutions">
                <el-skeleton :rows="4" animated />
              </div>
              <div v-else-if="institutions.length === 0" class="empty-state">
                <el-empty description="暂无可用的体检机构" />
                <div class="empty-hint">
                  <p>没有找到可用的体检机构，请稍后再试</p>
                </div>
              </div>
              <div v-else>
                <el-collapse v-model="activeInstitution">
                  <el-collapse-item
                    v-for="(inst, index) in institutions"
                    :key="inst?.ID || index"
                    :title="inst?.institution_name || '未命名机构'"
                    :name="inst?.ID || index">
                    <p><strong>地址:</strong> {{ inst?.institution_address || '暂无地址' }}</p>
                    <p><strong>资质:</strong> {{ inst?.institution_qualification || '暂无资质信息' }}</p>
                    <p><strong>套餐:</strong> {{ inst?.examination_package || '暂无套餐信息' }}</p>
                    <div class="institution-actions">
                      <ElButton type="info" size="small" style="margin-right: 10px;" @click="inst?.ID ? router.push(`/institutions/${inst.ID}`) : ElMessage.warning('机构ID无效，无法查看详情')">查看详情</ElButton>
                      <ElButton type="primary" size="small" @click="inst?.ID ? router.push(`/institutions/${inst.ID}`) : ElMessage.warning('机构ID无效，无法选择套餐')">选择套餐</ElButton>
                    </div>
                  </el-collapse-item>
                </el-collapse>
              </div>
            </ElCard>
          </div>

          <!-- OCR 识别面板 -->
          <div v-if="activeMenu === '7'">
            <ElCard shadow="hover">
              <template #header>
                <div class="card-header">
                  <span>OCR 识别</span>
                </div>
              </template>
              <el-form>
                <el-form-item label="选择图片">
                  <el-upload
                    class="upload-demo"
                    accept="image/*"
                    :show-file-list="false"
                    :on-change="handleFileChange" 
                    :auto-upload="false"
                  >
                    <el-button type="primary">选择图片</el-button>
                  </el-upload>
                  <div v-if="imageUrl" class="preview">
                    <img :src="imageUrl" alt="预览" />
                  </div>
                </el-form-item>
                <el-form-item>
                  <el-button type="success" @click="submitOcr" :disabled="!selectedFile">上传并识别</el-button>
                </el-form-item>
              </el-form>
              <el-divider />
              <div v-if="ocrResults.length">
                <el-table :data="ocrResults" stripe>
                  <el-table-column prop="item_name" label="项目名称" />
                  <el-table-column prop="item_value" label="识别结果" />
                </el-table>
              </div>
            </ElCard>
          </div>
        </template>
      </ElMain>
    </ElContainer>
  </ElContainer>
</template>

<style scoped>
.dashboard-container {
  height: 100vh;
}

.dashboard-header {
  background-color: #fff;
  border-bottom: 1px solid #dcdfe6;
  padding: 0 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 20px;
}

.notification-wrapper {
  position: relative;
  display: inline-block;
}

.notification-dot {
  position: absolute;
  top: -5px;
  right: -5px;
  width: 10px;
  height: 10px;
  background-color: #f56c6c;
  border-radius: 50%;
}

.dashboard-menu {
  height: 100%;
  border-right: none;
}

.user-profile-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 20px 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.user-info-list {
  padding: 0 20px;
}

.info-item {
  display: flex;
  padding: 10px 0;
}

.label {
  font-weight: bold;
  width: 120px;
  color: #606266;
}

.el-divider {
  margin: 8px 0;
}

.settings-container {
  max-width: 500px;
  margin: 0 auto;
  padding: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.request-actions {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

:deep(.el-timeline-item__content) {
  width: 100%;
}

.empty-state {
  text-align: center;
  padding: 20px;
}

.empty-hint {
  margin-top: 10px;
  text-align: center;
  color: #909399;
}

.institution-actions {
  margin-top: 10px;
}

/* OCR style */
.preview img {
  max-width: 200px;
  margin-top: 10px;
}
.chart-container {
  width: 100%;
  height: 400px;
  margin: auto;
  background-color: #fff;
  padding: 10px;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.health-data-section {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  margin-top: 20px;
}

.chart-section {
  flex: 1;
  min-width: 400px;
}

.form-section {
  flex: 1;
  min-width: 300px;
  padding: 0 20px;
  border-left: 1px solid #eee;
}

.form-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.form-header h3 {
  margin: 0;
}
</style>