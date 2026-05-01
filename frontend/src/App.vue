<template>
  <div id="app">
    <nav class="navbar">
      <router-link to="/" class="nav-brand">招聘平台</router-link>
      <div class="nav-links">
        <router-link to="/" class="nav-link">首页</router-link>
        <router-link to="/tasks" class="nav-link">任务列表</router-link>
        
        <template v-if="userStore.isLoggedIn">
          <template v-if="userStore.isEmployer || userStore.isAdmin">
            <router-link to="/tasks/create" class="nav-link">发布任务</router-link>
            <router-link to="/my-tasks" class="nav-link">我的任务</router-link>
          </template>
          
          <template v-if="userStore.isEmployee || userStore.isAdmin">
            <router-link to="/my-bids" class="nav-link">我的投标</router-link>
            <router-link to="/favorites" class="nav-link">我的收藏</router-link>
          </template>
          
          <router-link to="/my-reviews" class="nav-link">我的评价</router-link>
          
          <template v-if="userStore.isAdmin">
            <router-link to="/admin" class="nav-link admin-link">管理后台</router-link>
          </template>
          
          <span class="nav-user">欢迎, {{ userStore.user?.username }}</span>
          <button @click="logout" class="nav-btn">退出</button>
        </template>
        
        <template v-else>
          <router-link to="/login" class="nav-link">登录</router-link>
          <router-link to="/register" class="nav-link nav-btn-primary">注册</router-link>
        </template>
      </div>
    </nav>
    
    <main class="main-content">
      <router-view />
    </main>
    
    <footer class="footer">
      <p>&copy; 2026 招聘平台. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup>
import { useUserStore } from './stores/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

function logout() {
  userStore.logout()
  router.push('/')
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background-color: #f5f5f5;
  color: #333;
}

#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.nav-brand {
  color: white;
  font-size: 1.5rem;
  font-weight: bold;
  text-decoration: none;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-link:hover {
  background-color: rgba(255, 255, 255, 0.1);
}

.nav-btn {
  background-color: transparent;
  border: 1px solid white;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
}

.nav-btn:hover {
  background-color: white;
  color: #667eea;
}

.nav-btn-primary {
  background-color: white;
  color: #667eea;
  font-weight: bold;
}

.nav-btn-primary:hover {
  background-color: #f0f0f0;
}

.nav-user {
  color: white;
  font-weight: 500;
}

.admin-link {
  background-color: #ff6b6b;
  color: white;
}

.admin-link:hover {
  background-color: #ee5a5a;
}

.main-content {
  flex: 1;
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.footer {
  background-color: #333;
  color: white;
  text-align: center;
  padding: 1.5rem;
}

.card {
  background: white;
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  margin-bottom: 1rem;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #333;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 1rem;
  transition: all 0.3s;
  text-decoration: none;
  display: inline-block;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.btn-primary:hover {
  opacity: 0.9;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #5a6268;
}

.btn-success {
  background-color: #28a745;
  color: white;
}

.btn-success:hover {
  background-color: #218838;
}

.btn-danger {
  background-color: #dc3545;
  color: white;
}

.btn-danger:hover {
  background-color: #c82333;
}

.btn-sm {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: #333;
}

.form-input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.form-input:focus {
  outline: none;
  border-color: #667eea;
}

.form-textarea {
  min-height: 120px;
  resize: vertical;
}

.alert {
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.alert-success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.alert-danger {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
}

.badge-open {
  background-color: #d4edda;
  color: #155724;
}

.badge-in_progress {
  background-color: #fff3cd;
  color: #856404;
}

.badge-completed {
  background-color: #cce5ff;
  color: #004085;
}

.badge-cancelled {
  background-color: #f8d7da;
  color: #721c24;
}

.badge-pending {
  background-color: #fff3cd;
  color: #856404;
}

.badge-accepted {
  background-color: #d4edda;
  color: #155724;
}

.badge-rejected {
  background-color: #f8d7da;
  color: #721c24;
}

.text-center {
  text-align: center;
}

.mt-2 { margin-top: 0.5rem; }
.mt-4 { margin-top: 1rem; }
.mb-2 { margin-bottom: 0.5rem; }
.mb-4 { margin-bottom: 1rem; }
</style>
