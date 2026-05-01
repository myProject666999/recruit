<template>
  <div class="login-page">
    <div class="login-card">
      <h2 class="card-title">登录</h2>
      
      <div v-if="error" class="alert alert-danger">{{ error }}</div>
      <div v-if="success" class="alert alert-success">{{ success }}</div>
      
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label class="form-label">邮箱</label>
          <input 
            type="email" 
            v-model="form.email" 
            class="form-input" 
            placeholder="请输入邮箱"
            required
          />
        </div>
        
        <div class="form-group">
          <label class="form-label">密码</label>
          <input 
            type="password" 
            v-model="form.password" 
            class="form-input" 
            placeholder="请输入密码"
            required
          />
        </div>
        
        <button type="submit" class="btn btn-primary w-full" :disabled="loading">
          {{ loading ? '登录中...' : '登录' }}
        </button>
      </form>
      
      <div class="mt-4 text-center">
        还没有账号？<router-link to="/register">立即注册</router-link>
      </div>
      
      <div class="mt-2 text-center" style="font-size: 0.875rem; color: #666;">
        <p>管理员测试账号：admin@example.com / admin123</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '../stores/user'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  email: '',
  password: ''
})

const error = ref('')
const success = ref('')
const loading = ref(false)

async function handleLogin() {
  error.value = ''
  success.value = ''
  loading.value = true
  
  try {
    await userStore.login(form.value)
    success.value = '登录成功！'
    setTimeout(() => {
      router.push('/')
    }, 1000)
  } catch (err) {
    error.value = err.response?.data?.error || '登录失败，请检查邮箱和密码'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.login-card {
  width: 100%;
  max-width: 400px;
}

.w-full {
  width: 100%;
}
</style>
