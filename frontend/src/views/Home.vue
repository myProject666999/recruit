<template>
  <div class="home">
    <div class="hero-section">
      <h1>欢迎来到招聘平台</h1>
      <p>连接雇主与人才，创造无限可能</p>
      <div class="hero-actions">
        <router-link to="/tasks" class="btn btn-primary">浏览任务</router-link>
        <router-link v-if="!userStore.isLoggedIn" to="/register" class="btn btn-secondary">立即注册</router-link>
      </div>
    </div>

    <div class="features-section">
      <h2>平台特色</h2>
      <div class="features-grid">
        <div class="feature-card">
          <h3>雇主服务</h3>
          <p>发布任务，寻找合适的人才，选择中标者，评价服务质量</p>
        </div>
        <div class="feature-card">
          <h3>雇员服务</h3>
          <p>浏览任务，提交投标，收藏感兴趣的任务，完成任务获得报酬</p>
        </div>
        <div class="feature-card">
          <h3>安全可靠</h3>
          <p>完善的评价系统，确保交易双方的权益得到保障</p>
        </div>
      </div>
    </div>

    <div class="stats-section" v-if="stats">
      <h2>平台数据</h2>
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-number">{{ stats.tasks || 0 }}</div>
          <div class="stat-label">发布任务</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ stats.users || 0 }}</div>
          <div class="stat-label">注册用户</div>
        </div>
        <div class="stat-card">
          <div class="stat-number">{{ stats.completed || 0 }}</div>
          <div class="stat-label">完成任务</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '../stores/user'
import { taskApi } from '../services/api'

const userStore = useUserStore()
const stats = ref(null)

onMounted(async () => {
  try {
    const response = await taskApi.getAll()
    const tasks = response.data
    stats.value = {
      tasks: tasks.length,
      users: 0,
      completed: tasks.filter(t => t.status === 'completed').length
    }
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
})
</script>

<style scoped>
.home {
  text-align: center;
}

.hero-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 4rem 2rem;
  border-radius: 12px;
  margin-bottom: 3rem;
}

.hero-section h1 {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.hero-section p {
  font-size: 1.25rem;
  margin-bottom: 2rem;
  opacity: 0.9;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.features-section {
  margin-bottom: 3rem;
}

.features-section h2 {
  font-size: 2rem;
  margin-bottom: 2rem;
  color: #333;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.feature-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.feature-card h3 {
  font-size: 1.25rem;
  margin-bottom: 1rem;
  color: #667eea;
}

.feature-card p {
  color: #666;
  line-height: 1.6;
}

.stats-section {
  margin-bottom: 3rem;
}

.stats-section h2 {
  font-size: 2rem;
  margin-bottom: 2rem;
  color: #333;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 2rem;
}

.stat-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.stat-number {
  font-size: 3rem;
  font-weight: bold;
  color: #667eea;
  margin-bottom: 0.5rem;
}

.stat-label {
  color: #666;
  font-size: 1rem;
}
</style>
