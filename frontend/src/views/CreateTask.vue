<template>
  <div class="create-task-page">
    <div class="create-task-card">
      <h2 class="card-title">发布新任务</h2>
      
      <div v-if="error" class="alert alert-danger">{{ error }}</div>
      <div v-if="success" class="alert alert-success">{{ success }}</div>
      
      <form @submit.prevent="handleSubmit">
        <div class="form-group">
          <label class="form-label">任务标题</label>
          <input 
            type="text" 
            v-model="form.title" 
            class="form-input" 
            placeholder="请输入任务标题"
            required
          />
        </div>
        
        <div class="form-group">
          <label class="form-label">任务描述</label>
          <textarea 
            v-model="form.description" 
            class="form-input form-textarea" 
            placeholder="请详细描述任务要求..."
            required
          ></textarea>
        </div>
        
        <div class="form-group">
          <label class="form-label">预算 (¥)</label>
          <input 
            type="number" 
            v-model="form.budget" 
            class="form-input" 
            placeholder="请输入预算金额"
            min="0"
            step="0.01"
            required
          />
        </div>
        
        <div class="form-actions">
          <router-link to="/my-tasks" class="btn btn-secondary">取消</router-link>
          <button type="submit" class="btn btn-primary" :disabled="loading">
            {{ loading ? '发布中...' : '发布任务' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { taskApi } from '../services/api'

const router = useRouter()

const form = ref({
  title: '',
  description: '',
  budget: 0
})

const error = ref('')
const success = ref('')
const loading = ref(false)

async function handleSubmit() {
  error.value = ''
  success.value = ''
  loading.value = true
  
  try {
    await taskApi.create(form.value)
    success.value = '任务发布成功！'
    setTimeout(() => {
      router.push('/my-tasks')
    }, 1000)
  } catch (err) {
    error.value = err.response?.data?.error || '发布失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.create-task-page {
  display: flex;
  justify-content: center;
}

.create-task-card {
  width: 100%;
  max-width: 600px;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}
</style>
