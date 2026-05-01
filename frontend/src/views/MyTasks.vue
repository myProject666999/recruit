<template>
  <div class="my-tasks-page">
    <div class="page-header">
      <h1>我的任务</h1>
      <router-link to="/tasks/create" class="btn btn-primary">发布新任务</router-link>
    </div>

    <div v-if="loading" class="text-center">
      <p>加载中...</p>
    </div>

    <div v-else-if="tasks.length === 0" class="text-center">
      <p>暂无发布的任务</p>
      <router-link to="/tasks/create" class="btn btn-primary mt-4">发布第一个任务</router-link>
    </div>

    <div v-else class="tasks-list">
      <div v-for="task in tasks" :key="task.id" class="task-item card">
        <div class="task-header">
          <h3>
            <router-link :to="`/tasks/${task.id}`">{{ task.title }}</router-link>
          </h3>
          <span :class="`badge badge-${task.status}`">{{ getStatusText(task.status) }}</span>
        </div>
        
        <p class="task-description">{{ truncateText(task.description, 150) }}</p>
        
        <div class="task-info">
          <span class="task-budget">预算: ¥{{ task.budget }}</span>
          <span class="task-bids">投标数: {{ task.bids?.length || 0 }}</span>
          <span class="task-date">{{ formatDate(task.created_at) }}</span>
        </div>
        
        <div v-if="task.winner" class="task-winner">
          <strong>中标者:</strong> {{ task.winner.username }}
        </div>
        
        <div class="task-actions" v-if="task.status === 'open'">
          <router-link :to="`/tasks/${task.id}`" class="btn btn-primary btn-sm">查看详情</router-link>
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

async function fetchTasks() {
  loading.value = true
  try {
    const response = await taskApi.getMy()
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

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.task-item {
  transition: box-shadow 0.2s;
}

.task-item:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.task-header h3 {
  margin: 0;
}

.task-header h3 a {
  color: #333;
  text-decoration: none;
}

.task-header h3 a:hover {
  color: #667eea;
}

.task-description {
  color: #666;
  margin-bottom: 1rem;
  line-height: 1.5;
}

.task-info {
  display: flex;
  gap: 2rem;
  color: #999;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.task-budget {
  color: #28a745;
  font-weight: 500;
}

.task-winner {
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 4px;
  margin-bottom: 1rem;
  color: #495057;
}
</style>
