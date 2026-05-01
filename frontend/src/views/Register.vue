<template>
  <div class="register-page">
    <div class="register-card">
      <h2 class="card-title">注册</h2>
      
      <div v-if="error" class="alert alert-danger">{{ error }}</div>
      <div v-if="success" class="alert alert-success">{{ success }}</div>
      
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label class="form-label">用户名</label>
          <input 
            type="text" 
            v-model="form.username" 
            class="form-input" 
            placeholder="请输入用户名"
            required
          />
        </div>
        
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
            placeholder="请输入密码（至少6位）"
            required
            minlength="6"
          />
        </div>
        
        <div class="form-group">
          <label class="form-label">角色</label>
          <select v-model="form.role" class="form-input" required>
            <option value="employer">雇主</option>
            <option value="employee">雇员</option>
          </select>
          <p class="form-hint">雇主可以发布任务，雇员可以投标任务</p>
        </div>
        
        <button type="submit" class="btn btn-primary w-full" :disabled="loading">
          {{ loading ? '注册中...' : '注册' }}
        </button>
      </form>
      
      <div class="mt-4 text-center">
        已有账号？<router-link to="/login">立即登录</router-link>
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
  username: '',
  email: '',
  password: '',
  role: 'employer'
})

const error = ref('')
const success = ref('')
const loading = ref(false)

async function handleRegister() {
  error.value = ''
  success.value = ''
  loading.value = true
  
  try {
    await userStore.register(form.value)
    success.value = '注册成功！'
    setTimeout(() => {
      router.push('/')
    }, 1000)
  } catch (err) {
    error.value = err.response?.data?.error || '注册失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 60vh;
}

.register-card {
  width: 100%;
  max-width: 400px;
}

.w-full {
  width: 100%;
}

.form-hint {
  font-size: 0.875rem;
  color: #666;
  margin-top: 0.25rem;
}
</style>
