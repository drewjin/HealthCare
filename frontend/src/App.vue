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
  <router-view v-slot="{ Component, route }">
    <template v-if="route.meta.keepAlive">
      <keep-alive>
        <component :is="Component" :key="route.fullPath" />
      </keep-alive>
    </template>
    <template v-else>
      <component :is="Component" :key="route.fullPath" />
    </template>
  </router-view>
</template>

<style>
#app {
  height: 100vh;
  margin: 0;
  padding: 0;
}
</style>