<template>
  <div class="institution-detail-container">
    <div class="page-header-actions">
      <el-page-header @back="goBack" title="返回机构列表">
        <template #content>
          <span v-if="institution">{{ institution.institution_name }}</span>
          <span v-else>机构详情</span>
        </template>
      </el-page-header>
      <el-button type="primary" plain icon="House" @click="goToDashboard">返回主页</el-button>
    </div>

    <el-card v-if="loading" class="loading-card">
      <el-skeleton :rows="20" animated />
    </el-card>

    <div v-else-if="!institution" class="empty-state">
      <el-empty description="未找到机构信息" />
    </div>

    <div v-else class="institution-content">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="基本信息" name="info">
          <div class="full-height-layout">
            <!-- 左侧边栏：机构信息 -->
            <div class="sidebar">
              <div class="sidebar-header">
                <h2>机构信息</h2>
                <el-tag v-if="isAdmin || canManage" :type="getStatusType(institution.status)" effect="plain">
                  {{ getStatusText(institution.status) }}
                </el-tag>
              </div>
              <div class="info-section">
                <div v-if="canManage" class="admin-actions">
                  <el-button type="danger" size="small" @click="confirmDeleteInstitution">删除机构</el-button>
                </div>
                <div class="info-row">
                  <span class="info-label">机构名称：</span>
                  <span class="info-value">{{ institution.institution_name }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">地址：</span>
                  <span class="info-value">{{ institution.institution_address }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">联系电话：</span>
                  <span class="info-value">{{ institution.institution_phone || '未提供' }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">资质认证：</span>
                  <span class="info-value">{{ institution.institution_qualification }}</span>
                </div>
              </div>
            </div>
            
            <!-- 右侧主内容：体检套餐列表 -->
            <div class="main-content">
              <div class="content-header">
                <h2>体检套餐</h2>
                <p v-if="packages.length > 0" class="package-count">共 {{ packages.length }} 个套餐</p>
              </div>
              
              <div v-if="packages.length > 0" class="packages-list">
                <el-collapse accordion>
                  <el-collapse-item v-for="(pkg, index) in packages" :key="index" :name="index">
                    <template #title>
                      <div class="package-title">
                        <span class="title-text">{{ pkg.name }}</span>
                        <div class="package-title-info">
                          <span v-if="isPackageSelected(pkg.id)" class="package-selected-badge">
                            <el-tag type="success" size="small" effect="light">已选择</el-tag>
                          </span>
                          <span class="price-tag">¥{{ pkg.price }}</span>
                        </div>
                      </div>
                    </template>
                    <div class="package-content">
                      <div class="package-description">
                        <p>{{ pkg.description }}</p>
                      </div>
                      <div class="package-details">
                        <p><strong>适用人群：</strong> {{ pkg.suitableFor }}</p>
                        <p><strong>检查项目：</strong> {{ pkg.items }}</p>
                      </div>
                      <div class="package-actions">
                        <el-button 
                          type="primary" 
                          @click="selectPackage(pkg)" 
                          :disabled="isPackageSelected(pkg.id)">
                          {{ isPackageSelected(pkg.id) ? '已选择此套餐' : '选择此套餐' }}
                        </el-button>
                      </div>
                    </div>
                  </el-collapse-item>
                </el-collapse>
              </div>
              <div v-else class="empty-packages">
                <el-empty description="该机构暂无体检套餐信息" />
                <div class="empty-packages-hint">
                  <p>机构可能尚未设置体检套餐，请稍后再试或选择其他机构</p>
                </div>
              </div>
            </div>
          </div>
        </el-tab-pane>
        
        <el-tab-pane v-if="canManage" label="机构管理" name="manage">
          <institution-update-form :institution-id="institution.ID" />
          
          <el-divider>套餐管理</el-divider>
          
          <package-manager :institution-id="institution.ID" />
        </el-tab-pane>
      </el-tabs>
      
      <!-- 显示已选套餐的对话框 -->
      <el-dialog
        v-model="dialogVisible"
        title="确认选择套餐"
        width="50%"
        :close-on-click-modal="false"
        :show-close="true"
        destroy-on-close
      >
        <div v-if="selectedPackage" class="selected-package-info">
          <h3>{{ selectedPackage.name }}</h3>
          <p class="package-description">{{ selectedPackage.description || '暂无描述' }}</p>
          <div class="details-list">
            <p><strong>适用人群：</strong> {{ selectedPackage.suitableFor || '适用所有人群' }}</p>
            <p><strong>检查项目：</strong> {{ selectedPackage.items || '暂无详细信息' }}</p>
            <p><strong>价格：</strong> {{ selectedPackage.price }} 元</p>
          </div>
          <div class="package-notice">
            <p>选择此套餐后，您可以到机构进行体检，或联系机构预约时间。</p>
          </div>
        </div>
        <div v-else class="no-package-selected">
          <el-empty description="未找到套餐信息" />
        </div>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="dialogVisible = false">取消</el-button>
            <el-button type="primary" :disabled="!selectedPackage" @click="confirmPackageSelection">
              确认选择
            </el-button>
          </span>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'
import InstitutionUpdateForm from '@/components/InstitutionUpdateForm.vue'
import PackageManager from '@/components/PackageManager.vue'

interface Institution {
  ID: number
  institution_name: string
  institution_address: string
  institution_qualification: string
  institution_phone?: string
  examination_package: string
  status: number
  UserID?: number
}

interface Package {
  id?: number
  name: string
  description: string
  suitableFor: string
  items: string
  price: number
}

const route = useRoute()
const router = useRouter()
const institution = ref<Institution | null>(null)
const packages = ref<Package[]>([])
const loading = ref(true)
const isAdmin = ref(false)
const dialogVisible = ref(false)
const selectedPackage = ref<Package | null>(null)
const activeTab = ref('info')
const currentUserId = localStorage.getItem('uid') ? parseInt(localStorage.getItem('uid') || '0') : 0
const userSelectedPackages = ref<number[]>([]) // 存储用户已选择的套餐ID
const isSubmitting = ref(false) // 标记是否正在提交选择

// 判断当前用户是否可以管理该机构
const canManage = computed(() => {
  if (!institution.value) return false
  return isAdmin.value || (institution.value.UserID === currentUserId)
})

const fetchInstitutionDetail = async () => {
  try {
    loading.value = true
    const institutionId = route.params.id
    const token = localStorage.getItem('jwt')
    
    if (!token) {
      ElMessage.error('用户未登录，请先登录')
      router.push('/login')
      return
    }
    
    if (!institutionId) {
      ElMessage.error('无效的机构ID')
      router.push('/institutions')
      return
    }
    
    // Fetch institution details
    const response = await axios.get(`/api/institutions/${institutionId}`, {
      headers: { Authorization: `${token}` }
    })
    
    if (response.data && response.data.institution) {
      institution.value = response.data.institution
      isAdmin.value = response.data.isAdmin || false
      console.log('获取到机构详情:', institution.value)
      
      // Fetch packages separately
      await fetchPackages(institutionId as string)
    } else {
      ElMessage.error('获取机构详情失败：数据格式不正确')
      console.error('返回的机构数据格式不正确:', response.data)
    }
  } catch (error) {
    console.error('获取机构详情出错:', error)
    ElMessage.error('获取机构详情失败，请检查网络连接')
  } finally {
    loading.value = false
  }
}

const fetchPackages = async (institutionId: string) => {
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/institutions/${institutionId}/plans`, {
      headers: { Authorization: `${token}` }
    })
    
    console.log('Packages response:', response.data)
    
    // Parse the packages from the response
    if (response.data.plans && response.data.plans.length > 0) {
      const plans = response.data.plans
      const items = response.data.items || []
      
      // Transform the data into a more usable format for the frontend
      packages.value = plans.map((plan: any) => {
        const planItems = items.filter((item: any) => item.plan_id === plan.id)
        const itemNames = planItems.map((item: any) => item.item_description).join(', ')
        
        return {
          id: plan.id,
          name: plan.name,
          description: plan.description || planItems[0]?.item_description || '',
          suitableFor: plan.suitable_for || '适用所有人群', // Default value as it's not provided by API
          items: itemNames || plan.items || '',
          price: plan.price || 0
        }
      })
    } else {
      packages.value = []
      console.log('No packages found')
    }
  } catch (error) {
    console.error('Failed to fetch institution packages:', error)
    ElMessage.error('获取体检套餐信息失败')
    packages.value = []
  }
}

const selectPackage = (pkg: Package): void => {
  console.log('选择的套餐:', pkg as Package)
  if (!pkg || !pkg.id) {
    ElMessage.warning('套餐信息不完整，无法选择')
    return
  }
  
  selectedPackage.value = pkg
  dialogVisible.value = true
}

const confirmPackageSelection = async () => {
  try {
    if (!selectedPackage.value || !institution.value) {
      ElMessage.error('套餐或机构信息不完整，无法提交')
      return
    }
    
    isSubmitting.value = true
    
    const token = localStorage.getItem('jwt')
    
    await axios.post('/api/users/packages', {
      institution_id: institution.value.ID,
      plan_id: selectedPackage.value.id
    }, {
      headers: { Authorization: `${token}` }
    })
    
    // 将选中的套餐ID添加到已选列表中
    if (selectedPackage.value.id && !userSelectedPackages.value.includes(selectedPackage.value.id)) {
      userSelectedPackages.value.push(selectedPackage.value.id)
    }
    
    ElMessage.success(`您已成功选择 ${selectedPackage.value.name} 套餐`)
    dialogVisible.value = false
    
    // 可以在这里导航到预约页面或其他后续流程
    // router.push('/appointments')
  } catch (error) {
    console.error('选择套餐时出错:', error)
    ElMessage.error('选择套餐失败，请稍后再试')
  } finally {
    isSubmitting.value = false
  }
}

// 确认删除机构
const confirmDeleteInstitution = () => {
  ElMessageBox.confirm(
    '确定要删除此机构吗？此操作将永久删除机构及其所有套餐信息，不可恢复。',
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    deleteInstitution()
  }).catch(() => {
    // 取消删除
    ElMessage({
      type: 'info',
      message: '已取消删除'
    })
  })
}

// 删除机构
const deleteInstitution = async () => {
  try {
    const token = localStorage.getItem('jwt')
    
    if (!token || !institution.value || !institution.value.ID) {
      ElMessage.error('参数无效，无法删除机构')
      return
    }
    
    await axios.delete(`/api/institutions/${institution.value.ID}`, {
      headers: { Authorization: `${token}` }
    })
    
    ElMessage.success('机构已成功删除')
    
    // 删除成功后返回机构列表页
    router.push('/institutions')
    
  } catch (error) {
    console.error('删除机构失败:', error)
    ElMessage.error('删除机构失败，请稍后重试')
  }
}

const getStatusText = (status: number): string => {
  switch (status) {
    case 0: return '待审核'
    case 1: return '已批准'
    case 2: return '已拒绝'
    default: return '未知状态'
  }
}

const getStatusType = (status: number): string => {
  switch (status) {
    case 0: return 'warning'
    case 1: return 'success'
    case 2: return 'danger'
    default: return 'info'
  }
}

const goBack = () => {
  // 返回机构列表页面，使用 router.back() 可以更好地保持导航历史
  router.back()
}

const goToDashboard = () => {
  // 直接跳转到主页
  router.push('/dashboard')
}

onMounted(() => {
  fetchInstitutionDetail()
  fetchUserSelectedPackages() // 获取用户已选择的套餐
})

// 获取用户已选择套餐的函数
const fetchUserSelectedPackages = async () => {
  try {
    const token = localStorage.getItem('jwt')
    const userId = localStorage.getItem('uid')
    
    if (!token || !userId) {
      return
    }
    
    const response = await axios.get(`/api/users/${userId}/packages`, {
      headers: { Authorization: `${token}` }
    })

    console.log('获取用户已选择套餐:', response.data)
    if (response.data && response.data.user_packages) {
      userSelectedPackages.value = response.data.user_packages.map((pkg: { plan_id: number }) => pkg.plan_id)
    }
  } catch (error) {
    console.error('获取用户已选择套餐失败:', error)
  }
}

// 检查套餐是否已被用户选择
const isPackageSelected = (packageId: number | undefined): boolean => {
  if (!packageId) return false
  return userSelectedPackages.value.includes(packageId)
}

// 显示套餐选择状态的标签
const getPackageStatusTag = (pkg: Package): string => {
  if (isPackageSelected(pkg.id)) {
    return "已选择"
  }
  return ""
}
</script>

<style scoped>
.institution-detail-container {
  width: 100vw;
  height: 100vh;
  padding: 20px;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
}

.page-header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.institution-content {
  flex: 1;
  width: 100%;
  margin-top: 20px;
  display: flex;
  flex-direction: column;
}

.loading-card {
  width: 100%;
  min-height: 400px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 全高布局 */
.full-height-layout {
  display: flex;
  width: 100%;
  height: calc(100vh - 150px);
  margin-top: 10px;
}

/* 左侧边栏 */
.sidebar {
  flex: 0 0 300px;
  background-color: #f9f9f9;
  border-right: 1px solid #eaeaea;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 15px;
  border-bottom: 1px solid #eaeaea;
  margin-bottom: 15px;
}

/* 右侧主内容区域 */
.main-content {
  flex: 1;
  padding: 0 20px 20px 20px;
  overflow-y: auto;
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 0;
  border-bottom: 1px solid #eaeaea;
  margin-bottom: 20px;
}

.package-count {
  color: #909399;
  font-size: 14px;
}

.info-section {
  padding: 10px 0;
}

.admin-actions {
  margin-bottom: 15px;
  display: flex;
  justify-content: flex-end;
}

.info-row {
  margin-bottom: 15px;
  display: flex;
  flex-direction: column;
}

.info-label {
  font-weight: bold;
  margin-bottom: 5px;
  color: #606266;
}

.info-value {
  color: #303133;
}

.packages-list {
  margin-top: 15px;
}

.package-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.title-text {
  font-weight: 500;
  font-size: 16px;
}

.price-tag {
  font-size: 16px;
  color: #ff6b6b;
  font-weight: bold;
}

.package-title-info {
  display: flex;
  align-items: center;
  gap: 10px;
}

.package-selected-badge {
  margin-right: 5px;
}

.package-content {
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.package-description {
  margin-bottom: 15px;
  color: #606266;
  line-height: 1.5;
}

.package-details {
  margin-bottom: 20px;
  background-color: #ffffff;
  border-radius: 4px;
  padding: 15px;
  border-left: 3px solid #409EFF;
}

.package-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 10px;
}

.empty-packages,
.empty-state {
  margin: 30px 0;
  text-align: center;
  padding: 30px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.empty-packages-hint {
  margin-top: 15px;
  color: #909399;
}

.selected-package-info {
  padding: 20px;
}

.package-description {
  margin: 10px 0;
  color: #606266;
  line-height: 1.6;
}

.details-list {
  margin-top: 15px;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.package-notice {
  margin-top: 15px;
  padding: 15px;
  background-color: #ecf8ff;
  border-radius: 4px;
  color: #409EFF;
  border-left: 3px solid #409EFF;
}

.no-package-selected {
  padding: 30px;
  text-align: center;
  background-color: #f9f9f9;
  border-radius: 4px;
}

/* 响应式设计 */
@media screen and (max-width: 768px) {
  .full-height-layout {
    flex-direction: column;
    height: auto;
  }
  
  .sidebar {
    flex: none;
    width: 100%;
    border-right: none;
    border-bottom: 1px solid #eaeaea;
    margin-bottom: 20px;
  }
  
  .main-content {
    padding: 0;
  }
}
</style>