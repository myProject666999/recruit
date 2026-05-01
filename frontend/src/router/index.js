import { createRouter, createWebHistory } from 'vue-router'
import { useUserStore } from '../stores/user'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue'),
    meta: { guest: true }
  },
  {
    path: '/tasks',
    name: 'Tasks',
    component: () => import('../views/Tasks.vue')
  },
  {
    path: '/tasks/:id',
    name: 'TaskDetail',
    component: () => import('../views/TaskDetail.vue')
  },
  {
    path: '/tasks/create',
    name: 'CreateTask',
    component: () => import('../views/CreateTask.vue'),
    meta: { requiresAuth: true, role: 'employer' }
  },
  {
    path: '/my-tasks',
    name: 'MyTasks',
    component: () => import('../views/MyTasks.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/my-bids',
    name: 'MyBids',
    component: () => import('../views/MyBids.vue'),
    meta: { requiresAuth: true, role: 'employee' }
  },
  {
    path: '/favorites',
    name: 'Favorites',
    component: () => import('../views/Favorites.vue'),
    meta: { requiresAuth: true, role: 'employee' }
  },
  {
    path: '/my-reviews',
    name: 'MyReviews',
    component: () => import('../views/MyReviews.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/AdminDashboard.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/users',
    name: 'AdminUsers',
    component: () => import('../views/AdminUsers.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/tasks',
    name: 'AdminTasks',
    component: () => import('../views/AdminTasks.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach(async (to, from, next) => {
  const userStore = useUserStore()
  
  // 如果已登录但用户信息可能过期，刷新用户信息
  if (userStore.isLoggedIn && !userStore.user) {
    try {
      await userStore.fetchCurrentUser()
    } catch (e) {
      // token 无效，清除登录状态
      userStore.logout()
    }
  }
  
  // 检查需要登录的路由
  if (to.meta.requiresAuth && !userStore.isLoggedIn) {
    next('/login')
    return
  }
  
  // 检查只能游客访问的路由（登录后不能访问）
  if (to.meta.guest && userStore.isLoggedIn) {
    next('/')
    return
  }
  
  // 检查角色权限
  if (to.meta.role) {
    if (to.meta.role === 'employer' && !userStore.isEmployer && !userStore.isAdmin) {
      next('/')
      return
    }
    if (to.meta.role === 'employee' && !userStore.isEmployee && !userStore.isAdmin) {
      next('/')
      return
    }
    if (to.meta.role === 'admin' && !userStore.isAdmin) {
      next('/')
      return
    }
  }
  
  next()
})

export default router
