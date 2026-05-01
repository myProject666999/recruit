<template>
  <div class="favorites-page">
    <div class="page-header">
      <h1>我的收藏</h1>
    </div>

    <div v-if="loading" class="text-center">
      <p>加载中...</p>
    </div>

    <div v-else-if="favorites.length === 0" class="text-center">
      <p>暂无收藏的任务</p>
      <router-link to="/tasks" class="btn btn-primary mt-4">浏览任务</router-link>
    </div>

    <div v-else class="favorites-list">
      <div v-for="favorite in favorites" :key="favorite.id" class="favorite-item card">
        <div class="favorite-header">
          <h3>
            <router-link :to="`/tasks/${favorite.task_id}`">{{ favorite.task?.title }}</router-link>
          </h3>
          <span :class="`badge badge-${favorite.task?.status}`">{{ getStatusText(favorite.task?.status) }}</span>
        </div>
        
        <p class="task-description">{{ truncateText(favorite.task?.description, 150) }}</p>
        
        <div class="task-info">
          <span class="task-budget">预算: ¥{{ favorite.task?.budget }}</span>
          <span class="task-employer">发布者: {{ favorite.task?.employer?.username }}</span>
          <span class="favorite-date">收藏时间: {{ formatDate(favorite.created_at) }}</span>
        </div>
        
        <div class="task-actions">
          <router-link :to="`/tasks/${favorite.task_id}`" class="btn btn-primary btn-sm">查看详情</router-link>
          <button @click="removeFavorite(favorite.task_id)" class="btn btn-secondary btn-sm">取消收藏</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { favoriteApi } from '../services/api'

const favorites = ref([])
const loading = ref(false)

async function fetchFavorites() {
  loading.value = true
  try {
    const response = await favoriteApi.getMy()
    favorites.value = response.data
  } catch (error) {
    console.error('Failed to fetch favorites:', error)
  } finally {
    loading.value = false
  }
}

async function removeFavorite(taskId) {
  try {
    await favoriteApi.toggle(taskId)
    fetchFavorites()
  } catch (error) {
    console.error('Failed to remove favorite:', error)
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
  fetchFavorites()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.favorites-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.favorite-item {
  transition: box-shadow 0.2s;
}

.favorite-item:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.favorite-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.favorite-header h3 {
  margin: 0;
}

.favorite-header h3 a {
  color: #333;
  text-decoration: none;
}

.favorite-header h3 a:hover {
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

.task-actions {
  display: flex;
  gap: 0.5rem;
}
</style>
