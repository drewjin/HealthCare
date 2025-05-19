import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue')
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/institutions',
      name: 'institutions',
      component: () => import('@/views/InstitutionsView.vue'),
      meta: { requiresAuth: true, keepAlive: false }
    },
    {
      path: '/institutions/:id',
      name: 'institution-detail',
      component: () => import('@/views/InstitutionDetailView.vue'),
      meta: { requiresAuth: true, keepAlive: false }
    },
    {
      path: '/institution-manage',
      name: 'institution-manage',
      component: () => import('@/views/InstitutionManageView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/health-records',
      name: 'health-records',
      component: () => import('@/views/HealthRecordsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/plan-items',
      name: 'plan-items',
      component: () => import('@/views/PlanItemsView.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

router.beforeEach((to, _, next) => {
  const token = localStorage.getItem('jwt')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router