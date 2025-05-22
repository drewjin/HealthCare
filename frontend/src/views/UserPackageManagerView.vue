<template>
  <div class="user-package-manager-view">
    <el-page-header @back="goBack" title="返回">
      <template #content>
        <span class="page-title">机构用户套餐管理</span>
      </template>
    </el-page-header>
    
    <div class="manager-container">
      <user-package-manager />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import UserPackageManager from '@/components/UserPackageManager.vue'

const router = useRouter()

const goBack = () => {
  router.push('/dashboard')
}
</script>

<style scoped>
.user-package-manager-view {
  padding: 20px;
}

.page-title {
  font-size: 18px;
  font-weight: bold;
}

.manager-container {
  margin-top: 20px;
}
</style>
