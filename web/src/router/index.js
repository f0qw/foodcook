import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/Home.vue'),
    meta: { title: '首页' }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/Register.vue'),
    meta: { title: '注册' }
  },
  {
    path: '/dishes',
    name: 'Dishes',
    component: () => import('@/views/Dishes.vue'),
    meta: { title: '菜品管理', requiresAuth: true, requiresRoot: true }
  },
  {
    path: '/ingredients',
    name: 'Ingredients',
    component: () => import('@/views/Ingredients.vue'),
    meta: { title: '食材管理', requiresAuth: true, requiresRoot: true }
  },
  {
    path: '/meal-records',
    name: 'MealRecords',
    component: () => import('@/views/MealRecords.vue'),
    meta: { title: '用餐记录', requiresAuth: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/Profile.vue'),
    meta: { title: '个人中心', requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - FoodCook` : 'FoodCook'
  
  // 检查是否需要认证
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/login')
    return
  }
  
  // 检查是否需要 root 权限
  if (to.meta.requiresRoot && !authStore.isRoot) {
    next('/')
    return
  }
  
  next()
})

export default router 