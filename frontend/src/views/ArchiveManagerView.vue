<template>
  <div class="archive-items-container">
    <el-page-header @back="goBack" title="返回">
      <template #content>
        <span class="page-title">档案管理</span>
      </template>
    </el-page-header>

    <div class="archive-items-section">
      <el-card class="box-card">
        <template #header>
          <div class="card-header">
            <span>档案列表</span>
            <div class="button-group">
              <el-button type="success" @click="editItem()">创建档案</el-button>
              <el-button type="primary" @click="refreshItems">刷新列表</el-button>
            </div>
          </div>
        </template>

        <el-table v-loading="loading" :data="archiveItems" style="width: 100%">
          <el-table-column prop="ID" label="ID" width="80" />
          <el-table-column prop="UserId" label="用户名">
            <template #default="scope">
              {{ formatName(scope.row.UserId) }}
            </template>
          </el-table-column>
          <el-table-column prop="UserHealthInfo" label="用户指标" />
          <el-table-column label="操作" width="300">
            <template #default="scope">
              <el-button size="small" type="primary" @click="editItem(scope.row)">编辑</el-button>
              <el-button size="small" type="danger" @click="deleteItem(scope.row.ID)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </div>

    <!-- 编辑档案对话框 -->
    <el-dialog v-model="editDialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="editForm" label-width="100px">
        <el-form-item label="用户">
          <el-select v-model="editForm.UserId" :disabled="editForm.ID" placeholder="请选择用户">
            <el-option 
              v-for="plan in users" 
              :key="plan.id" 
              :label="plan.name" 
              :value="plan.id" 
            />
          </el-select>
        </el-form-item>
        <el-form-item label="OCR">
          <el-input v-model="editForm.UserHealthInfo" placeholder="请输入" type="textarea" />
          <el-upload
            class="upload-demo"
            accept="image/*"
            :show-file-list="false"
            :on-change="handleFileChange" 
            :auto-upload="false"
          >
            <el-button type="primary">OCR识别</el-button>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="saving" @click="saveItemChanges">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import { en } from 'element-plus/es/locales.mjs'

interface ArchiveItem {
  ID: number
  ItemName: string
  UserId: string | number
  UserHealthInfo: Object
  CreatedAt: string
  UpdatedAt: string
  DeletedAt: string | null
}

const router = useRouter()
const route = useRoute()
const loading = ref(false)
const saving = ref(false)
const archiveItems = ref<ArchiveItem[]>([])
const editDialogVisible = ref(false)

const editForm = reactive({
  ID: '',
  UserId: '',
  UserHealthInfo: ''
})

const users = [{
  id: 1,
  name: 'Qing Xu'
}, {
  id: 3,
  name: 'Yijie Jin'
}, {
  id: 4,
  name: 'Hi'
}]

const goBack = () => {
  router.back()
}

const formatName = (id: number) => {
  if (!id) return '-'
  return users.find(item => item.id === id).name
}  
const formatDate = (dateString: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString()
}  

const fetchArchiveItems = async () => {
  loading.value = true
  try {
    const token = localStorage.getItem('jwt')
    const response = await axios.get('/api/healthitems/all', {
      headers: { Authorization: `${token}` }
    })
    
    const items = response.data.items || response.data.health_items || []
    
    archiveItems.value = items.map((item: any) => ({
      ID: item.ID || item.id,
      ItemName: item.ItemName || item.item_name,
      UserId: item.UserId || item.user_id,
      UserHealthInfo: item.UserHealthInfo || item.user_health_info,
      CreatedAt: item.CreatedAt || item.created_at || '',
      UpdatedAt: item.UpdatedAt || item.updated_at || '',
      DeletedAt: item.DeletedAt || item.deleted_at || null
    }))
    
    console.log('Fetched health items:', archiveItems.value)
  } catch (error) {
    console.error('Failed to fetch health items:', error)
    ElMessage.error('获取档案列表失败')
  } finally {
    loading.value = false
  }
}

const dialogTitle = ref('')
// 编辑
const editItem = (item: ArchiveItem) => {
  if(item) {
    dialogTitle.value = '编辑档案'
    editForm.ID = item.ID
    editForm.UserId = item.UserId
    editForm.UserHealthInfo = item.UserHealthInfo
  } else {
    dialogTitle.value = '新增档案'
    editForm.ID = null
    editForm.UserId = null
    editForm.UserHealthInfo = null
  }
  editDialogVisible.value = true
}

// 保存
const saveItemChanges = async () => {
  saving.value = true
  try {
    const token = localStorage.getItem('jwt')
    // 修改
    if (editForm.ID) {
      console.log('editForm:', editForm)
      await axios.post(`/api/users/update_use_health_item`, {
        id: editForm.ID,
        user_health_info: editForm.UserHealthInfo
      }, {
        headers: { 
          Authorization: `${token}`,
          'Content-Type': 'application/json; charset=utf-8'
        }
      })
      ElMessage.success('档案修改成功')
    } else {
      await axios.post(`/api/users/create_health_item`, {
        user_id: editForm.UserId,
        user_health_info: editForm.UserHealthInfo
      }, {
        headers: { 
          Authorization: `${token}`,
          'Content-Type': 'application/json; charset=utf-8'
        }
      })
      ElMessage.success('档案新增成功')
    }
    editDialogVisible.value = false

    // 更新列表
    await fetchArchiveItems()
  } catch (error: any) {
    console.error('Failed to update health item:', error)
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('更新档案失败')
    }
  } finally {
    saving.value = false
  }
}

// 删除
const deleteItem = async (key: string) => {
  try {
    const token = localStorage.getItem('jwt')
    await axios.get(`/api/users/${key}/del_health_item`, {
      headers: { Authorization: token }
    })
    ElMessage.success('档案删除成功')
    // 更新列表
    await fetchArchiveItems()
  } catch (error: any) {
    console.error('Failed to update health item:', error)
    if (error.response && error.response.data && error.response.data.error) {
      ElMessage.error(error.response.data.error)
    } else {
      ElMessage.error('档案删除失败')
    }
  } finally {
    saving.value = false
  }
}

// 刷新列表
const refreshItems = () => {
  fetchArchiveItems()
}

onMounted(() => {
  fetchArchiveItems()
})

const selectedFile = ref<File | null>(null)

const handleFileChange = (uploadFile: any) => {
  console.log('File selected (handleFileChange):', uploadFile);
  if (uploadFile && uploadFile.raw) {
    selectedFile.value = uploadFile.raw;
    console.log('selectedFile.value assigned:', selectedFile.value);
    submitOcr()
  } else {
    console.log('No file raw data found in uploadFile:', uploadFile);
    // Clear previous selection if the new selection is invalid
    selectedFile.value = null;
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

    const rawOcrResults: OcrResult[] = response.data.result || [];
    console.log('Raw OCR results from backend:', JSON.stringify(rawOcrResults)); // Log raw results

    editForm.UserHealthInfo = rawOcrResults.map(entry => {
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
      // Format: "体检 | 血糖：200"
      return `体检 | ${newName}：${newValue}`;
    });

  } catch (e) {
    ElMessage.error('OCR 识别失败，请重试')
    console.error('OCR 识别失败', e)
  }
}
</script>

<style scoped>
.archive-items-container {
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: bold;
}

.archive-items-section {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.button-group {
  display: flex;
  gap: 10px;
}

.empty-data {
  padding: 20px;
  text-align: center;
  color: #909399;
}

.mt-4 {
  margin-top: 16px;
}
</style>
