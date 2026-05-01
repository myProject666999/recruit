<template>
  <div class="admin-users-page">
    <div class="page-header">
      <h1>用户管理</h1>
      <div class="filters">
        <select v-model="filterRole" @change="fetchUsers" class="form-input" style="width: auto;">
          <option value="">全部角色</option>
          <option value="employer">雇主</option>
          <option value="employee">雇员</option>
          <option value="admin">管理员</option>
        </select>
      </div>
    </div>

    <div v-if="loading" class="text-center">
      <p>加载中...</p>
    </div>

    <div v-else-if="users.length === 0" class="text-center">
      <p>暂无用户</p>
    </div>

    <div v-else class="users-table card">
      <table>
        <thead>
          <tr>
            <th>ID</th>
            <th>用户名</th>
            <th>邮箱</th>
            <th>角色</th>
            <th>注册时间</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.id }}</td>
            <td>{{ user.username }}</td>
            <td>{{ user.email }}</td>
            <td>
              <span :class="`badge badge-${user.role}`">{{ getRoleText(user.role) }}</span>
            </td>
            <td>{{ formatDate(user.created_at) }}</td>
            <td>
              <div class="action-buttons">
                <template v-if="user.role !== 'admin'">
                  <select 
                    v-model="editRoles[user.id]" 
                    @change="updateUserRole(user)"
                    class="form-input"
                    style="width: auto; padding: 0.25rem 0.5rem;"
                  >
                    <option value="employer">雇主</option>
                    <option value="employee">雇员</option>
                  </select>
                  <button 
                    @click="deleteUser(user)" 
                    class="btn btn-danger btn-sm"
                  >
                    删除
                  </button>
                </template>
                <template v-else>
                  <span class="text-muted">管理员</span>
                </template>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { adminApi } from '../services/api'

const users = ref([])
const loading = ref(false)
const filterRole = ref('')
const editRoles = reactive({})

async function fetchUsers() {
  loading.value = true
  try {
    const params = {}
    if (filterRole.value) {
      params.role = filterRole.value
    }
    const response = await adminApi.getUsers(params)
    users.value = response.data
    
    // 初始化编辑角色
    users.value.forEach(user => {
      editRoles[user.id] = user.role
    })
  } catch (error) {
    console.error('Failed to fetch users:', error)
  } finally {
    loading.value = false
  }
}

async function updateUserRole(user) {
  if (!confirm(`确定要将用户 "${user.username}" 的角色更改为 "${getRoleText(editRoles[user.id])}" 吗？`)) {
    editRoles[user.id] = user.role
    return
  }
  
  try {
    await adminApi.updateUser(user.id, { role: editRoles[user.id] })
    user.role = editRoles[user.id]
  } catch (error) {
    console.error('Failed to update user:', error)
    alert('更新失败')
    editRoles[user.id] = user.role
  }
}

async function deleteUser(user) {
  if (!confirm(`确定要删除用户 "${user.username}" 吗？此操作不可恢复。`)) {
    return
  }
  
  try {
    await adminApi.deleteUser(user.id)
    fetchUsers()
  } catch (error) {
    console.error('Failed to delete user:', error)
    alert('删除失败')
  }
}

function getRoleText(role) {
  const roleMap = {
    'employer': '雇主',
    'employee': '雇员',
    'admin': '管理员'
  }
  return roleMap[role] || role
}

function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

onMounted(() => {
  fetchUsers()
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

.users-table {
  padding: 0;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead {
  background: #f8f9fa;
}

th, td {
  padding: 1rem;
  text-align: left;
  border-bottom: 1px solid #eee;
}

th {
  font-weight: 600;
  color: #333;
}

tr:hover {
  background: #f8f9fa;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.text-muted {
  color: #999;
}
</style>
