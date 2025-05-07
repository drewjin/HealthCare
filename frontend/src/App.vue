<script setup lang="ts">
import { RouterView } from 'vue-router'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import { createApp } from 'vue'
import 'element-plus/dist/index.css'
import './style.css'

const app = createApp({})
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}
</script>

<template>
  <RouterView />
</template>

<style>
#app {
  height: 100vh;
  margin: 0;
  padding: 0;
}
</style>