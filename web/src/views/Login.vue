<template>
  <div class="login-container">
    <!-- 移动端头部 -->
    <div v-if="isMobile" class="mobile-header">
      <div class="mobile-header-content">
        <el-button @click="$router.go(-1)" type="text" class="back-btn">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <h2>登录</h2>
        <div></div>
      </div>
    </div>

    <div class="login-card" :class="{ 'mobile-login-card': isMobile }">
      <div class="login-header">
        <h1>🍽️ FoodCook</h1>
        <p>欢迎回来，请登录您的账户</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        label-width="0"
        class="login-form"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="用户名"
            size="large"
            :class="{ 'mobile-input': isMobile }"
            prefix-icon="User"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :class="{ 'mobile-input': isMobile }"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleLogin"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="login-button"
            :class="{ 'mobile-btn': isMobile }"
            :loading="loading"
            @click="handleLogin"
          >
            登录
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <p>
          还没有账户？
          <router-link to="/register" class="register-link">立即注册</router-link>
        </p>
        <router-link to="/" class="back-home">返回首页</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const router = useRouter()
const authStore = useAuthStore()

const loginFormRef = ref()
const loading = ref(false)
const isMobile = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
}

// 检测移动端
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

// 监听窗口大小变化
onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    await authStore.login(loginForm)
    router.push('/')
  } catch (error) {
    console.error('登录失败:', error)
    ElMessage.error('登录失败，请检查用户名和密码')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

/* 移动端头部 */
.mobile-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  z-index: 1000;
  padding: 10px 0;
}

.mobile-header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 15px;
}

.back-btn {
  font-size: 18px;
}

.mobile-header-content h2 {
  margin: 0;
  color: #333;
  font-size: 18px;
}

.login-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.mobile-login-card {
  margin-top: 60px;
  padding: 30px 20px;
  border-radius: 0;
  box-shadow: none;
  background: transparent;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.login-header h1 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 28px;
}

.login-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.login-form {
  margin-bottom: 20px;
}

.login-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
}

.mobile-btn {
  height: 48px;
  border-radius: 24px;
  font-size: 16px;
}

.login-footer {
  text-align: center;
  margin-top: 20px;
}

.login-footer p {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
}

.register-link {
  color: #409eff;
  text-decoration: none;
  font-weight: 500;
}

.register-link:hover {
  text-decoration: underline;
}

.back-home {
  color: #999;
  text-decoration: none;
  font-size: 12px;
}

.back-home:hover {
  color: #666;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .login-container {
    padding: 0;
    background: white;
  }
  
  .login-card {
    box-shadow: none;
    border-radius: 0;
  }
  
  .login-header h1 {
    font-size: 24px;
  }
  
  .login-header p {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .mobile-login-card {
    padding: 20px 15px;
  }
  
  .login-header {
    margin-bottom: 25px;
  }
}
</style> 