<template>
  <div class="ocr-container">
    <el-card>
      <template #header>
        <span>OCR 识别</span>
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
      <div v-if="results.length">
        <el-table :data="results" stripe>
          <el-table-column prop="item_name" label="项目名称" />
          <el-table-column prop="item_value" label="识别结果" />
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus';

interface OcrResult { item_name: string; item_value: string }

const selectedFile = ref<File | null>(null)
const imageUrl = ref<string>('')
const results = ref<OcrResult[]>([])

// const beforeUpload = (file: File) => {
//   console.log('Selected file (beforeUpload):', file); 
//   selectedFile.value = file
//   const reader = new FileReader()
//   reader.readAsDataURL(file)
//   reader.onload = () => { imageUrl.value = reader.result as string }
//   return false // 阻止自动上传
// }

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
    // Clear previous selection if the new selection is invalid
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

    const rawOcrResults: OcrResult[] = response.data.result || [];
    console.log('Raw OCR results from backend:', JSON.stringify(rawOcrResults)); // Log raw results

    results.value = rawOcrResults.map(entry => {
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

<style scoped>
.ocr-container { padding: 20px; }
.preview img { max-width: 200px; margin-top: 10px; }
</style>
