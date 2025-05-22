<template>
  <div class="user-packages-container">
    <h2>机构用户套餐管理</h2>
    
    <!-- 搜索区域 -->
    <div class="search-box">
      <el-form :model="searchForm" :inline="true">
        <el-form-item label="用户姓名">
          <el-input v-model="searchForm.user_name" placeholder="输入用户姓名搜索"></el-input>
        </el-form-item>
        <el-form-item label="套餐名称">
          <el-input v-model="searchForm.plan_name" placeholder="输入套餐名称搜索"></el-input>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="选择状态" clearable>
            <el-option label="待完成" :value="0" />
            <el-option label="已完成" :value="1" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="searchPackages">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <!-- 数据表格 -->
    <el-table
      v-loading="loading"
      :data="userPackages"
      border
      stripe
      style="width: 100%"
    >
      <el-table-column prop="user_id" label="用户ID" width="80" />
      <el-table-column prop="user_name" label="用户姓名" />
      <el-table-column prop="plan_id" label="套餐ID" width="80" />
      <el-table-column prop="plan_name" label="套餐名称" />
      <el-table-column label="完成状态" width="150">
        <template #default="scope">
          <el-progress 
            :percentage="calculateProgress(scope.row)" 
            :status="scope.row.status === 1 ? 'success' : undefined"
          />
          <div class="progress-text">
            {{ scope.row.completed_items }}/{{ scope.row.total_items }} 项
          </div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 1 ? 'success' : 'warning'">
            {{ scope.row.status === 1 ? '已完成' : '待完成' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200">
        <template #default="scope">
          <el-button 
            type="primary" 
            size="small"
            @click="editUserData(scope.row)"
          >
            {{ scope.row.status === 0 ? '填写数据' : '查看数据' }}
          </el-button>
          <el-button 
            v-if="scope.row.status === 0 && calculateProgress(scope.row) === 100"
            type="success" 
            size="small"
            @click="updatePackageStatus(scope.row, 1)"
          >
            标记完成
          </el-button>
          <el-button 
            v-if="scope.row.status === 1"
            type="warning" 
            size="small"
            @click="updatePackageStatus(scope.row, 0)"
          >
            重新开放
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    
    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        background
        layout="total, sizes, prev, pager, next, jumper"
        :total="totalItems"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'

interface UserPackage {
  id: number
  user_id: number
  plan_id: number
  institution_id: number
  status: number
  user_name: string
  plan_name: string
  completed_items: number
  total_items: number
}

const router = useRouter()
const loading = ref(false)
const userPackages = ref<UserPackage[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const totalItems = ref(0)

const searchForm = reactive({
  user_name: '',
  plan_name: '',
  status: null as number | null
})

// 获取用户套餐列表
const fetchUserPackages = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('未登录或令牌丢失')
      router.push('/login')
      return
    }
    
    // 构建查询参数
    const params: Record<string, string> = {}
    
    if (searchForm.user_name) {
      params.user_name = searchForm.user_name
    }
    
    if (searchForm.plan_name) {
      params.plan_name = searchForm.plan_name
    }
    
    if (searchForm.status !== null) {
      params.status = searchForm.status.toString()
    }
    
    const response = await axios.get('/api/institution/user-packages', {
      params,
      headers: { Authorization: token }
    })
    
    if (response.data && Array.isArray(response.data.packages)) {
      userPackages.value = response.data.packages
      totalItems.value = response.data.packages.length
    } else {
      userPackages.value = []
      totalItems.value = 0
    }
  } catch (error: any) {
    console.error('获取用户套餐列表失败:', error)
    ElMessage.error(error.response?.data?.error || '获取用户套餐列表失败')
  } finally {
    loading.value = false
  }
}

// 计算进度百分比
const calculateProgress = (userPackage: UserPackage): number => {
  if (userPackage.total_items === 0) return 0
  return Math.round((userPackage.completed_items / userPackage.total_items) * 100)
}

// 更新套餐状态
const updatePackageStatus = async (userPackage: UserPackage, newStatus: number) => {
  try {
    const token = localStorage.getItem('jwt')
    if (!token) {
      ElMessage.error('未登录或令牌丢失')
      router.push('/login')
      return
    }
    
    const response = await axios.patch(
      `/api/user-packages/${userPackage.user_id}/${userPackage.plan_id}/status`,
      { status: newStatus },
      { headers: { Authorization: token } }
    )
    
    if (response.status === 200) {
      ElMessage.success('套餐状态更新成功')
      // 更新本地数据
      const packageIndex = userPackages.value.findIndex(pkg => 
        pkg.user_id === userPackage.user_id && pkg.plan_id === userPackage.plan_id
      )
      if (packageIndex !== -1) {
        userPackages.value[packageIndex].status = newStatus
      }
    } else {
      throw new Error('更新套餐状态失败')
    }
  } catch (error: any) {
    console.error('更新套餐状态失败:', error)
    ElMessage.error(error.response?.data?.error || '更新套餐状态失败')
  }
}

// 搜索用户套餐
const searchPackages = () => {
  currentPage.value = 1
  fetchUserPackages()
}

// 重置搜索条件
const resetSearch = () => {
  searchForm.user_name = ''
  searchForm.plan_name = ''
  searchForm.status = null
  searchPackages()
}

// 编辑用户体检数据
const editUserData = (userPackage: UserPackage) => {
  router.push({
    name: 'add-user-data-detail',  // 改为add-user-data-detail以确保正确路由
    params: {
      customer_id: userPackage.user_id.toString(),
      plan_id: userPackage.plan_id.toString()
    }
  })
}

// 每页大小变化
const handleSizeChange = (size: number) => {
  pageSize.value = size
  fetchUserPackages()
}

// 页码变化
const handleCurrentChange = (page: number) => {
  currentPage.value = page
  fetchUserPackages()
}

onMounted(() => {
  fetchUserPackages()
})
</script>

<style scoped>
.user-packages-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-box {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.progress-text {
  font-size: 12px;
  color: #606266;
  text-align: center;
  margin-top: 5px;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
