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
    },
    {
      path: '/plan-health-items/:id',
      name: 'plan-health-items',
      component: () => import('@/components/PlanHealthItems.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/user-health-data',
      name: 'user-health-data',
      component: () => import('@/views/UserHealthDataView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/health-items',
      name: 'health-items',
      component: () => import('@/views/HealthItemsView.vue'),
      meta: { requiresAuth: true, userTypes: [3, 2] } // Only institution users and admins
    },
    {
      path: '/health-item-manager',
      name: 'health-item-manager',
      component: () => import('@/views/HealthItemManagerView.vue'),
      meta: { requiresAuth: true, userTypes: [3, 2] } // Only institution users and admins
    },
    {
      path: '/plan-health-item-manager',
      name: 'plan-health-item-manager',
      component: () => import('@/views/PlanHealthItemManagerView.vue'),
      meta: { requiresAuth: true, userTypes: [3, 2] } // Only institution users and admins
    },
    {
      path: '/add-user-data',
      name: 'add-user-data',
      component: () => import('@/views/AddUserDataView.vue'),
      meta: { requiresAuth: true, userTypes: [3, 2] } // Only institution users and admins
    },
    {
      path: '/add-user-data/:customer_id/:plan_id',
      name: 'add-user-data-detail',
      component: () => import('@/components/AddUserHealthData.vue'),
      meta: { requiresAuth: true, userTypes: [3, 2] } // Only institution users and admins
    },
    {
      path: '/user-packages',
      name: 'user-packages',
      component: () => import('@/views/UserPackageManagerView.vue'),
      meta: { requiresAuth: true, userTypes: [3, 2] } // Only institution users and admins
    },
    // OCR function is now integrated into the dashboard
    // {
    //   path: '/ocr',
    //   name: 'ocr',
    //   component: () => import('@/views/OcrView.vue'),
    //   meta: { requiresAuth: true }
    // }
  ]
})

router.beforeEach((to, _, next) => {
  const token = localStorage.getItem('jwt')
  if (to.meta.requiresAuth && !token) {
    next('/login')
  } else if (to.meta.userTypes && Array.isArray(to.meta.userTypes)) {
    // Check user type if required
    const userType = parseInt(localStorage.getItem('userType') || '0')
    if (!to.meta.userTypes.includes(userType)) {
      next('/dashboard')
    } else {
      next()
    }
  } else {
    next()
  }
})

export default router