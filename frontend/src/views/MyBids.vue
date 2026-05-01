<template>
  <div class="my-bids-page">
    <div class="page-header">
      <h1>我的投标</h1>
    </div>

    <div v-if="loading" class="text-center">
      <p>加载中...</p>
    </div>

    <div v-else-if="bids.length === 0" class="text-center">
      <p>暂无投标记录</p>
      <router-link to="/tasks" class="btn btn-primary mt-4">浏览任务</router-link>
    </div>

    <div v-else class="bids-list">
      <div v-for="bid in bids" :key="bid.id" class="bid-item card">
        <div class="bid-header">
          <h3>
            <router-link :to="`/tasks/${bid.task_id}`">{{ bid.task?.title }}</router-link>
          </h3>
          <span :class="`badge badge-${bid.status}`">{{ getStatusText(bid.status) }}</span>
        </div>
        
        <div class="bid-info">
          <span class="bid-price">我的报价: ¥{{ bid.proposed_price }}</span>
          <span class="task-budget">任务预算: ¥{{ bid.task?.budget }}</span>
          <span class="task-employer">发布者: {{ bid.task?.employer?.username }}</span>
        </div>
        
        <p class="bid-proposal">{{ bid.proposal }}</p>
        
        <div class="bid-footer">
          <span class="bid-date">投标时间: {{ formatDate(bid.created_at) }}</span>
          <router-link :to="`/tasks/${bid.task_id}`" class="btn btn-primary btn-sm">查看任务</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { bidApi } from '../services/api'

const bids = ref([])
const loading = ref(false)

async function fetchBids() {
  loading.value = true
  try {
    const response = await bidApi.getMy()
    bids.value = response.data
  } catch (error) {
    console.error('Failed to fetch bids:', error)
  } finally {
    loading.value = false
  }
}

function getStatusText(status) {
  const statusMap = {
    'pending': '待处理',
    'accepted': '已接受',
    'rejected': '已拒绝'
  }
  return statusMap[status] || status
}

function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchBids()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.bids-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.bid-item {
  transition: box-shadow 0.2s;
}

.bid-item:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.bid-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.bid-header h3 {
  margin: 0;
}

.bid-header h3 a {
  color: #333;
  text-decoration: none;
}

.bid-header h3 a:hover {
  color: #667eea;
}

.bid-info {
  display: flex;
  gap: 2rem;
  color: #666;
  font-size: 0.875rem;
  margin-bottom: 1rem;
}

.bid-price {
  color: #28a745;
  font-weight: 500;
}

.bid-proposal {
  color: #444;
  line-height: 1.6;
  margin-bottom: 1rem;
  padding: 0.75rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.bid-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #999;
  font-size: 0.875rem;
}
</style>
