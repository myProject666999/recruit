<template>
  <div class="admin-dashboard">
    <div class="page-header">
      <h1>管理后台</h1>
    </div>

    <div class="dashboard-stats">
      <div class="stat-card card">
        <h3>用户管理</h3>
        <p>管理平台所有用户</p>
        <router-link to="/admin/users" class="btn btn-primary">进入管理</router-link>
      </div>
      
      <div class="stat-card card">
        <h3>任务管理</h3>
        <p>管理平台所有任务</p>
        <router-link to="/admin/tasks" class="btn btn-primary">进入管理</router-link>
      </div>
    </div>

    <div class="quick-info">
      <div class="info-card card">
        <h3>管理员操作指南</h3>
        <ul>
          <li><strong>用户管理:</strong> 查看、编辑和删除用户账户</li>
          <li><strong>任务管理:</strong> 查看所有任务、修改任务状态、删除任务</li>
          <li><strong>权限说明:</strong> 管理员可以操作所有用户和任务</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.dashboard-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.stat-card {
  text-align: center;
  padding: 2rem;
}

.stat-card h3 {
  margin-bottom: 0.5rem;
  color: #333;
}

.stat-card p {
  color: #666;
  margin-bottom: 1.5rem;
}

.quick-info {
  display: grid;
  gap: 2rem;
}

.info-card h3 {
  margin-bottom: 1rem;
  color: #333;
}

.info-card ul {
  list-style: none;
  padding: 0;
}

.info-card li {
  padding: 0.5rem 0;
  color: #666;
  border-bottom: 1px solid #eee;
}

.info-card li:last-child {
  border-bottom: none;
}
</style>
