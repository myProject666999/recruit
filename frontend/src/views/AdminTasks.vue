<template>
  <div class="admin-tasks-page">
    <div class="page-header">
      <h1>任务管理</h1>
      <div class="filters">
        <select v-model="filterStatus" @change="fetchTasks" class="form-input" style="width: auto;">
          <option value="">全部状态</option>
          <option value="open">进行中</option>
          <option value="in_progress">已中标</option>
          <option value="completed">已完成</option>
          <option value="cancelled">已取消</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="text-center">
      <p>加载中...</p>
    </div>

    <div v-else-if="tasks.length === 0" class="text-center">
      <p>暂无任务</p>
    </div>

    <div v-else class="tasks-list">
      <div v-for="task in tasks" :key="task.id" class="task-item card">
        <div class="task-header">
          <h3>
            <router-link :to="`/tasks/${task.id}`">{{ task.title }}</router-link>
          </h3>
          <span :class="`badge badge-${task.status}`">{{ getStatusText(task.status) }}</span>
        </div>
        
        <div class="task-info">
          <div class="info-row">
            <span><strong>ID:</strong> {{ task.id }}</span>
            <span><strong>预算:</strong> ¥{{ task.budget }}</span>
          </div>
          <div class="info-row">
            <span><strong>发布者:</strong> {{ task.employer?.username }}</span>
            <span><strong>发布时间:</strong> {{ formatDate(task.created_at) }}</span>
          </div>
          <div v-if="task.winner" class="info-row">
            <span><strong>中标者:</strong> {{ task.winner.username }}</span>
          </div>
        </div>
        
        <div class="task-actions">
          <div class="status-selector">
            <label>修改状态:</label>
            <select 
              v-model="editStatus[task.id]" 
              @change="updateTaskStatus(task)"
              class="form-input"
              style="width: auto;"
            >
              <option value="open">进行中</option>
              <option value="in_progress">已中标</option>
              <option value="completed">已完成</option>
              <option value="cancelled">已取消</option>
            </select>
          </div>
          
          <button 
            @click="deleteTask(task)" 
            class="btn btn-danger btn-sm"
          >
            删除任务
          </button>
          
          <router-link :to="`/tasks/${task.id}`" class="btn btn-secondary btn-sm">
            查看详情
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { adminApi } from '../services/api'

const tasks = ref([])
const loading = ref(false)
const filterStatus = ref('')
const editStatus = reactive({})

async function fetchTasks() {
  loading.value = true
  try {
    const params = {}
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    const response = await adminApi.getTasks(params)
    tasks.value = response.data
    
    // 初始化编辑状态
    tasks.value.forEach(task => {
      editStatus[task.id] = task.status
    })
  } catch (error) {
    console.error('Failed to fetch tasks:', error)
  } finally {
    loading.value = false
  }
}

async function updateTaskStatus(task) {
  if (!confirm(`确定要将任务 "${task.title}" 的状态更改为 "${getStatusText(editStatus[task.id])}" 吗？`)) {
    editStatus[task.id] = task.status
    return
  }
  
  try {
    await adminApi.updateTaskStatus(task.id, { status: editStatus[task.id] })
    task.status = editStatus[task.id]
  } catch (error) {
    console.error('Failed to update task:', error)
    alert('更新失败')
    editStatus[task.id] = task.status
  }
}

async function deleteTask(task) {
  if (!confirm(`确定要删除任务 "${task.title}" 吗？此操作不可恢复。`)) {
    return
  }
  
  try {
    await adminApi.deleteTask(task.id)
    fetchTasks()
  } catch (error) {
    console.error('Failed to delete task:', error)
    alert('删除失败')
  }
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

.tasks-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
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
  margin-bottom: 1rem;
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

.task-info {
  margin-bottom: 1rem;
  color: #666;
}

.info-row {
  display: flex;
  gap: 2rem;
  margin-bottom: 0.5rem;
}

.task-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

.status-selector {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-selector label {
  font-weight: 500;
  color: #333;
}
</style>
