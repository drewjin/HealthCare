<template>
  <div class="institution-detail-container">
    <el-page-header @back="goBack" title="返回机构列表">
      <template #content>
        <span v-if="institution">{{ institution.institution_name }}</span>
        <span v-else>机构详情</span>
      </template>
    </el-page-header>

    <el-card v-if="loading" class="loading-card">
      <el-skeleton :rows="10" animated />
    </el-card>

    <div v-else-if="!institution" class="empty-state">
      <el-empty description="未找到机构信息" />
    </div>

    <div v-else class="institution-content">
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <h2>机构信息</h2>
            <el-tag v-if="isAdmin" :type="getStatusType(institution.status)" effect="plain">
              {{ getStatusText(institution.status) }}
            </el-tag>
          </div>
        </template>
        <div class="info-section">
          <div class="info-row">
            <span class="info-label">机构名称：</span>
            <span class="info-value">{{ institution.institution_name }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">地址：</span>
            <span class="info-value">{{ institution.institution_address }}</span>
          </div>
          <div class="info-row">
            <span class="info-label">资质认证：</span>
            <span class="info-value">{{ institution.institution_qualification }}</span>
          </div>
        </div>
      </el-card>

      <el-card class="package-card">
        <template #header>
          <div class="card-header">
            <h2>体检套餐</h2>
          </div>
        </template>
        <div v-if="packages.length > 0" class="packages-list">
          <el-collapse>
            <el-collapse-item v-for="(pkg, index) in packages" :key="index" :title="pkg.name">
              <div class="package-content">
                <div class="package-description">
                  <p>{{ pkg.description }}</p>
                </div>
                <div class="package-details">
                  <p><strong>适用人群：</strong> {{ pkg.suitableFor }}</p>
                  <p><strong>检查项目：</strong> {{ pkg.items }}</p>
                  <p><strong>价格：</strong> {{ pkg.price }} 元</p>
                </div>
                <div class="package-actions">
                  <el-button type="primary" @click="selectPackage(pkg)">选择此套餐</el-button>
                </div>
              </div>
            </el-collapse-item>
          </el-collapse>
        </div>
        <div v-else class="empty-packages">
          <el-empty description="暂无体检套餐信息" />
        </div>
      </el-card>
      
      <!-- 显示已选套餐的对话框 -->
      <el-dialog
        v-model="dialogVisible"
        title="确认选择套餐"
        width="50%"
      >
        <div v-if="selectedPackage" class="selected-package-info">
          <h3>{{ selectedPackage.name }}</h3>
          <p>{{ selectedPackage.description }}</p>
          <div class="details-list">
            <p><strong>适用人群：</strong> {{ selectedPackage.suitableFor }}</p>
            <p><strong>检查项目：</strong> {{ selectedPackage.items }}</p>
            <p><strong>价格：</strong> {{ selectedPackage.price }} 元</p>
          </div>
        </div>
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="dialogVisible = false">取消</el-button>
            <el-button type="primary" @click="confirmPackageSelection">
              确认选择
            </el-button>
          </span>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

interface Institution {
  ID: number
  institution_name: string
  institution_address: string
  institution_qualification: string
  examination_package: string
  status: number
}

interface Package {
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

const fetchInstitutionDetail = async () => {
  try {
    loading.value = true
    const institutionId = route.params.id
    const token = localStorage.getItem('jwt')
    
    // Fetch institution details
    const response = await axios.get(`/api/institutions/${institutionId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    institution.value = response.data.institution
    isAdmin.value = response.data.isAdmin
    
    // Fetch packages separately
    await fetchPackages(institutionId as string)
  } catch (error) {
    console.error('Failed to fetch institution details:', error)
    ElMessage.error('获取机构详情失败')
  } finally {
    loading.value = false
  }
}

const fetchPackages = async (institutionId: string) => {
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get(`/api/institutions/${institutionId}/packages`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    
    // Parse the packages (assuming packages are stored as JSON string)
    if (response.data.packages) {
      try {
        packages.value = JSON.parse(response.data.packages)
      } catch (e) {
        console.error('Failed to parse examination packages:', e)
        packages.value = []
      }
    }
  } catch (error) {
    console.error('Failed to fetch institution packages:', error)
    ElMessage.error('获取体检套餐信息失败')
    packages.value = []
  }
}

const selectPackage = (pkg: Package) => {
  selectedPackage.value = pkg
  dialogVisible.value = true
}

const confirmPackageSelection = () => {
  if (selectedPackage.value) {
    // 这里可以实现保存选择的套餐到用户的预约或者其他操作
    ElMessage.success(`您已成功选择 ${selectedPackage.value.name} 套餐`)
    dialogVisible.value = false
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
  router.push('/institutions')
}

onMounted(() => {
  fetchInstitutionDetail()
})
</script>

<style scoped>
.institution-detail-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 20px;
}

.loading-card,
.info-card,
.package-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.info-section {
  padding: 10px 0;
}

.info-row {
  margin-bottom: 15px;
  display: flex;
}

.info-label {
  font-weight: bold;
  min-width: 100px;
}

.info-value {
  flex: 1;
}

.packages-list {
  margin-top: 15px;
}

.package-content {
  padding: 10px;
}

.package-description {
  margin-bottom: 15px;
}

.package-details {
  margin-bottom: 15px;
}

.package-actions {
  display: flex;
  justify-content: flex-end;
}

.empty-packages,
.empty-state {
  margin: 30px 0;
  text-align: center;
}

.selected-package-info {
  padding: 15px;
}

.details-list {
  margin-top: 15px;
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 4px;
}
</style>