<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ElContainer, ElAside, ElMain, ElHeader, ElMenu, ElMenuItem, ElCard, ElAvatar, ElRow, ElCol, ElButton, ElDivider, ElForm, ElFormItem, ElInput, ElMessage, ElDialog, ElTable, ElTableColumn, ElBadge, ElTimeline, ElTimelineItem, ElSelect, ElOption, ElCollapse, ElCollapseItem, ElSkeleton, ElEmpty, ElUpload } from 'element-plus'
import { User, List, Setting, UserFilled, HomeFilled, Tools, PictureFilled } from '@element-plus/icons-vue'
import InstitutionForm from '@/components/InstitutionForm.vue'
import AdminReview from '@/components/AdminReview.vue'
import axios from 'axios'

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
  username: string,
  name: string,
  relationship: string
}>>([])

const showCreateFamilyDialog = ref(false)
const familyRequestForm = ref({
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

      // 获取已确认的家庭关系
      const confirmedResponse = await fetch(`http://localhost:3000/api/family/confirmed/${uid.value}`, {
        method: 'GET',
        headers: { Authorization: token },
      })
      if (confirmedResponse.ok) {
        const data = await confirmedResponse.json()
        confirmedFamilyMembers.value = data || []
      }
    }
  } catch (error) {
    console.error(error)
    router.push('/login')
  }
})

const handleLogout = () => {
  localStorage.removeItem('jwt')
  localStorage.removeItem('uid')
  router.push('/login')
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
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '处理请求失败')
  }
}

const handleCreateFamilyRequest = async () => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) throw new Error('未登录')

    const response = await fetch(`http://localhost:3000/api/family/request/${uid.value}`, {
      method: 'POST',
      headers: { 
        'Content-Type': 'application/json',
        Authorization: token 
      },
      body: JSON.stringify(familyRequestForm.value)
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error)
    }

    ElMessage.success('请求已发送')
    showCreateFamilyDialog.value = false
    familyRequestForm.value = { relative_username: '', relationship: '' }
  } catch (error) {
    console.error(error)
    ElMessage.error(error instanceof Error ? error.message : '发送请求失败')
  }
}

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
                    <span>{{ userInfo.username }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">邮箱：</span>
                    <span>{{ userInfo.email }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">电话：</span>
                    <span>{{ userInfo.phone }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">性别：</span>
                    <span>{{ userInfo.gender === 'M' ? '男' : '女' }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">生日：</span>
                    <span>{{ userInfo.birthday }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">地址：</span>
                    <span>{{ userInfo.address }}</span>
                  </div>
                  <ElDivider />
                  <div class="info-item">
                    <span class="label">用户类型：</span>
                    <span>{{ getUserTypeLabel(userInfo.user_type) }}</span>
                  </div>
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
        </div>          <!-- 普通用户其他面板 -->
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
                  <el-button type="primary" plain @click="router.push('/health-records')">查看详细健康记录</el-button>
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
                      <ElButton type="primary" @click="showCreateFamilyDialog = true">添加成员</ElButton>
                    </div>
                  </template>
                  <div v-if="confirmedFamilyMembers.length > 0">
                    <ElTable :data="confirmedFamilyMembers" style="width: 100%">
                      <ElTableColumn prop="username" label="用户名" />
                      <ElTableColumn prop="name" label="姓名" />
                      <ElTableColumn prop="relationship" label="关系" />
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
              title="创建家庭关系请求"
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
  width: 80px;
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
</style>