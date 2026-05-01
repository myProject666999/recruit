<template>
  <div class="my-reviews-page">
    <div class="page-header">
      <h1>我的评价</h1>
    </div>

    <div v-if="loading" class="text-center">
      <p>加载中...</p>
    </div>

    <div v-else-if="reviews.length === 0" class="text-center">
      <p>暂无评价记录</p>
    </div>

    <div v-else class="reviews-list">
      <div v-for="review in reviews" :key="review.id" class="review-item card">
        <div class="review-header">
          <div class="review-task">
            <strong>任务:</strong> 
            <router-link :to="`/tasks/${review.task_id}`">{{ review.task?.title }}</router-link>
          </div>
          <div class="review-rating">
            <span v-for="i in 5" :key="i" class="star" :class="{ filled: i <= review.rating }">★</span>
          </div>
        </div>
        
        <div class="review-info">
          <span class="review-employer">评价者: {{ review.employer?.username }}</span>
          <span class="review-date">{{ formatDate(review.created_at) }}</span>
        </div>
        
        <p class="review-comment">{{ review.comment }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { reviewApi } from '../services/api'

const reviews = ref([])
const loading = ref(false)

async function fetchReviews() {
  loading.value = true
  try {
    const response = await reviewApi.getMy()
    reviews.value = response.data
  } catch (error) {
    console.error('Failed to fetch reviews:', error)
  } finally {
    loading.value = false
  }
}

function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchReviews()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.reviews-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.review-item {
  transition: box-shadow 0.2s;
}

.review-item:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.review-task a {
  color: #667eea;
  text-decoration: none;
}

.review-task a:hover {
  text-decoration: underline;
}

.review-rating {
  display: flex;
  gap: 0.25rem;
}

.star {
  font-size: 1.25rem;
  color: #ddd;
}

.star.filled {
  color: #ffc107;
}

.review-info {
  display: flex;
  gap: 2rem;
  color: #999;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.review-comment {
  color: #444;
  line-height: 1.6;
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 4px;
}
</style>
