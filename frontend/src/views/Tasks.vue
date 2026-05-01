<template>
  <div class="tasks-page">
    <div class="page-header">
      <h1>任务列表</h1>
      <div class="filters">
        <select v-model="filterStatus" @change="fetchTasks" class="form-input" style="width: auto;">
          <option value="">全部状态</option>
          <option value="open">进行中</option>
          <option value="in_progress">已中标</option>
          <option value="completed">已完成</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="text-center">
      <p>加载中...</p>
    </div>

    <div v-else-if="tasks.length === 0" class="text-center">
      <p>暂无任务</p>
    </div>

    <div v-else class="tasks-grid">
      <div v-for="task in tasks" :key="task.id" class="task-card card">
        <h3 class="card-title">
          <router-link :to="`/tasks/${task.id}`">{{ task.title }}</router-link>
        </h3>
        <p class="task-description">{{ truncateText(task.description, 100) }}</p>
        <div class="task-meta">
          <span class="task-budget">¥{{ task.budget }}</span>
          <span :class="`badge badge-${task.status}`">{{ getStatusText(task.status) }}</span>
        </div>
        <div class="task-footer">
          <span class="task-employer">发布者: {{ task.employer?.username || '未知' }}</span>
          <span class="task-date">{{ formatDate(task.created_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { taskApi } from '../services/api'

const tasks = ref([])
const loading = ref(false)
const filterStatus = ref('')

async function fetchTasks() {
  loading.value = true
  try {
    const params = {}
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    const response = await taskApi.getAll(params)
    tasks.value = response.data
  } catch (error) {
    console.error('Failed to fetch tasks:', error)
  } finally {
    loading.value = false
  }
}

function truncateText(text, maxLength) {
  if (!text) return ''
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

function getStatusText(status) {
  const statusMap = {
    'open': '进行中',
    'in_progress': '已中标',
    'completed': '已完成',
    'cancelled': '已取消'
  }
  return statusMap[status] || status
}

function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchTasks()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.filters {
  display: flex;
  gap: 1rem;
}

.tasks-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.task-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.task-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.task-card h3 a {
  color: #333;
  text-decoration: none;
}

.task-card h3 a:hover {
  color: #667eea;
}

.task-description {
  color: #666;
  margin-bottom: 1rem;
  line-height: 1.5;
}

.task-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.task-budget {
  font-size: 1.25rem;
  font-weight: bold;
  color: #28a745;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  color: #999;
  font-size: 0.875rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}
</style>
