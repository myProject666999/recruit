<template>
  <div class="task-detail-page">
    <div v-if="loading" class="text-center">
      <p>加载中...</p>
    </div>

    <div v-else-if="error" class="alert alert-danger">{{ error }}</div>

    <div v-else class="task-detail">
      <div class="task-main card">
        <div class="task-header">
          <h1>{{ task.title }}</h1>
          <span :class="`badge badge-${task.status}`">{{ getStatusText(task.status) }}</span>
        </div>
        
        <div class="task-meta-info">
          <div class="meta-item">
            <span class="meta-label">预算:</span>
            <span class="meta-value budget">¥{{ task.budget }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">发布者:</span>
            <span class="meta-value">{{ task.employer?.username }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">发布时间:</span>
            <span class="meta-value">{{ formatDate(task.created_at) }}</span>
          </div>
          <div v-if="task.winner" class="meta-item">
            <span class="meta-label">中标者:</span>
            <span class="meta-value">{{ task.winner.username }}</span>
          </div>
        </div>

        <div class="task-description-section">
          <h3>任务描述</h3>
          <p>{{ task.description }}</p>
        </div>

        <div class="task-actions">
          <template v-if="userStore.isLoggedIn">
            <!-- 雇员可以投标和收藏 -->
            <template v-if="userStore.isEmployee || userStore.isAdmin">
              <template v-if="task.status === 'open'">
                <button @click="showBidForm = !showBidForm" class="btn btn-primary">
                  投标任务
                </button>
              </template>
              
              <button @click="toggleFavorite" class="btn btn-secondary" :disabled="favoriteLoading">
                {{ isFavorited ? '取消收藏' : '收藏任务' }}
              </button>
            </template>

            <!-- 雇主（自己的任务）或管理员可以查看投标、选择中标者、完成任务、评价 -->
            <template v-if="(userStore.isEmployer && task.employer_id === userStore.user.id) || userStore.isAdmin">
              <template v-if="task.status === 'in_progress'">
                <button @click="completeTask" class="btn btn-success" :disabled="completing">
                  完成任务
                </button>
              </template>
              
              <template v-if="task.status === 'completed' && !hasReviewed">
                <button @click="showReviewForm = !showReviewForm" class="btn btn-primary">
                  评价雇员
                </button>
              </template>
            </template>

            <!-- 中标者可以完成任务 -->
            <template v-if="task.winner && task.winner.id === userStore.user.id && task.status === 'in_progress'">
              <button @click="completeTask" class="btn btn-success" :disabled="completing">
                完成任务
              </button>
            </template>
          </template>

          <template v-else>
            <router-link to="/login" class="btn btn-primary">登录后投标</router-link>
          </template>
        </div>

        <!-- 投标表单 -->
        <div v-if="showBidForm" class="bid-form card mt-4">
          <h3>提交投标</h3>
          <form @submit.prevent="submitBid">
            <div class="form-group">
              <label class="form-label">报价 (¥)</label>
              <input 
                type="number" 
                v-model="bidForm.proposed_price" 
                class="form-input" 
                min="0"
                step="0.01"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label">投标方案</label>
              <textarea 
                v-model="bidForm.proposal" 
                class="form-input form-textarea" 
                placeholder="请描述您的投标方案..."
                required
              ></textarea>
            </div>
            <div class="form-actions">
              <button type="button" @click="showBidForm = false" class="btn btn-secondary">取消</button>
              <button type="submit" class="btn btn-primary" :disabled="bidding">
                {{ bidding ? '提交中...' : '提交投标' }}
              </button>
            </div>
          </form>
        </div>

        <!-- 评价表单 -->
        <div v-if="showReviewForm" class="review-form card mt-4">
          <h3>评价雇员</h3>
          <form @submit.prevent="submitReview">
            <div class="form-group">
              <label class="form-label">评分</label>
              <div class="rating-input">
                <button 
                  v-for="i in 5" 
                  :key="i"
                  type="button"
                  @click="reviewForm.rating = i"
                  class="star-btn"
                  :class="{ active: i <= reviewForm.rating }"
                >
                  ★
                </button>
              </div>
            </div>
            <div class="form-group">
              <label class="form-label">评价内容</label>
              <textarea 
                v-model="reviewForm.comment" 
                class="form-input form-textarea" 
                placeholder="请输入您的评价..."
                required
              ></textarea>
            </div>
            <div class="form-actions">
              <button type="button" @click="showReviewForm = false" class="btn btn-secondary">取消</button>
              <button type="submit" class="btn btn-primary" :disabled="reviewing">
                {{ reviewing ? '提交中...' : '提交评价' }}
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- 投标列表 -->
      <div v-if="task.bids && task.bids.length > 0" class="bids-section card">
        <h3>投标列表 ({{ task.bids.length }})</h3>
        <div class="bids-list">
          <div v-for="bid in task.bids" :key="bid.id" class="bid-item">
            <div class="bid-header">
              <span class="bidder">{{ bid.employee?.username }}</span>
              <span :class="`badge badge-${bid.status}`">{{ getBidStatusText(bid.status) }}</span>
            </div>
            <div class="bid-info">
              <span class="bid-price">报价: ¥{{ bid.proposed_price }}</span>
              <span class="bid-date">{{ formatDate(bid.created_at) }}</span>
            </div>
            <p class="bid-proposal">{{ bid.proposal }}</p>
            
            <!-- 雇主（自己的任务）或管理员可以接受投标 -->
            <div v-if="((userStore.isEmployer && task.employer_id === userStore.user.id) || userStore.isAdmin) && task.status === 'open' && bid.status === 'pending'" class="bid-actions">
              <button @click="acceptBid(bid.id)" class="btn btn-success btn-sm">
                接受投标
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '../stores/user'
import { taskApi, bidApi, favoriteApi, reviewApi } from '../services/api'

const route = useRoute()
const userStore = useUserStore()

const task = ref({})
const loading = ref(true)
const error = ref('')
const showBidForm = ref(false)
const showReviewForm = ref(false)
const bidding = ref(false)
const reviewing = ref(false)
const favoriteLoading = ref(false)
const completing = ref(false)
const isFavorited = ref(false)
const hasReviewed = ref(false)

const bidForm = ref({
  proposed_price: 0,
  proposal: ''
})

const reviewForm = ref({
  rating: 5,
  comment: ''
})

async function fetchTask() {
  loading.value = true
  error.value = ''
  
  try {
    const taskId = route.params.id
    const response = await taskApi.getById(taskId)
    task.value = response.data
    
    // 检查是否已收藏
    if (userStore.isLoggedIn) {
      try {
        const favResponse = await favoriteApi.check(taskId)
        isFavorited.value = favResponse.data.is_favorited
      } catch (e) {
        console.error('Failed to check favorite:', e)
      }
    }
  } catch (err) {
    error.value = err.response?.data?.error || '加载任务失败'
  } finally {
    loading.value = false
  }
}

async function submitBid() {
  bidding.value = true
  try {
    await bidApi.create(task.value.id, bidForm.value)
    showBidForm.value = false
    bidForm.value = { proposed_price: 0, proposal: '' }
    fetchTask()
  } catch (err) {
    alert(err.response?.data?.error || '投标失败')
  } finally {
    bidding.value = false
  }
}

async function acceptBid(bidId) {
  if (!confirm('确定要接受这个投标吗？')) return
  
  try {
    await bidApi.accept(bidId)
    fetchTask()
  } catch (err) {
    alert(err.response?.data?.error || '接受投标失败')
  }
}

async function toggleFavorite() {
  favoriteLoading.value = true
  try {
    const response = await favoriteApi.toggle(task.value.id)
    isFavorited.value = response.data.is_favorited
  } catch (err) {
    alert(err.response?.data?.error || '操作失败')
  } finally {
    favoriteLoading.value = false
  }
}

async function submitReview() {
  reviewing.value = true
  try {
    await reviewApi.create(task.value.id, reviewForm.value)
    showReviewForm.value = false
    hasReviewed.value = true
    reviewForm.value = { rating: 5, comment: '' }
  } catch (err) {
    alert(err.response?.data?.error || '评价失败')
  } finally {
    reviewing.value = false
  }
}

async function completeTask() {
  if (!confirm('确定要完成这个任务吗？')) return
  
  completing.value = true
  try {
    await taskApi.complete(task.value.id)
    fetchTask()
  } catch (err) {
    alert(err.response?.data?.error || '完成任务失败')
  } finally {
    completing.value = false
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

function getBidStatusText(status) {
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
  fetchTask()
})
</script>

<style scoped>
.task-detail-page {
  max-width: 800px;
  margin: 0 auto;
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.task-header h1 {
  margin: 0;
  font-size: 1.75rem;
}

.task-meta-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 8px;
}

.meta-item {
  display: flex;
  flex-direction: column;
}

.meta-label {
  font-size: 0.875rem;
  color: #666;
  margin-bottom: 0.25rem;
}

.meta-value {
  font-weight: 500;
}

.meta-value.budget {
  font-size: 1.5rem;
  color: #28a745;
  font-weight: bold;
}

.task-description-section {
  margin-bottom: 2rem;
}

.task-description-section h3 {
  margin-bottom: 1rem;
}

.task-description-section p {
  line-height: 1.8;
  color: #444;
}

.task-actions {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.rating-input {
  display: flex;
  gap: 0.5rem;
}

.star-btn {
  background: none;
  border: none;
  font-size: 2rem;
  cursor: pointer;
  color: #ddd;
  padding: 0;
}

.star-btn.active {
  color: #ffc107;
}

.bids-section h3 {
  margin-bottom: 1.5rem;
}

.bids-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.bid-item {
  padding: 1rem;
  border: 1px solid #eee;
  border-radius: 8px;
}

.bid-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.bidder {
  font-weight: 500;
}

.bid-info {
  display: flex;
  gap: 1rem;
  margin-bottom: 0.5rem;
  color: #666;
  font-size: 0.875rem;
}

.bid-price {
  color: #28a745;
  font-weight: 500;
}

.bid-proposal {
  color: #444;
  line-height: 1.5;
}

.bid-actions {
  margin-top: 0.75rem;
}
</style>
