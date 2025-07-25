<template>
  <div class="register-container">
    <!-- 移动端头部 -->
    <div v-if="isMobile" class="mobile-header">
      <div class="mobile-header-content">
        <el-button @click="$router.go(-1)" type="text" class="back-btn">
          <el-icon><ArrowLeft /></el-icon>
        </el-button>
        <h2>注册</h2>
        <div></div>
      </div>
    </div>

    <div class="register-card" :class="{ 'mobile-register-card': isMobile }">
      <div class="register-header">
        <h1>🍽️ FoodCook</h1>
        <p>创建您的账户，开始美食之旅</p>
      </div>

      <el-form
        ref="registerFormRef"
        :model="registerForm"
        :rules="registerRules"
        label-width="0"
        class="register-form"
      >
        <el-form-item prop="username">
          <el-input
            v-model="registerForm.username"
            placeholder="用户名"
            size="large"
            :class="{ 'mobile-input': isMobile }"
            prefix-icon="User"
          />
        </el-form-item>

        <el-form-item prop="email">
          <el-input
            v-model="registerForm.email"
            placeholder="邮箱"
            size="large"
            :class="{ 'mobile-input': isMobile }"
            prefix-icon="Message"
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="registerForm.password"
            type="password"
            placeholder="密码"
            size="large"
            :class="{ 'mobile-input': isMobile }"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>

        <el-form-item prop="confirmPassword">
          <el-input
            v-model="registerForm.confirmPassword"
            type="password"
            placeholder="确认密码"
            size="large"
            :class="{ 'mobile-input': isMobile }"
            prefix-icon="Lock"
            show-password
            @keyup.enter="handleRegister"
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            class="register-button"
            :class="{ 'mobile-btn': isMobile }"
            :loading="loading"
            @click="handleRegister"
          >
            注册
          </el-button>
        </el-form-item>
      </el-form>

      <div class="register-footer">
        <p>
          已有账户？
          <router-link to="/login" class="login-link">立即登录</router-link>
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

const registerFormRef = ref()
const loading = ref(false)
const isMobile = ref(false)

const registerForm = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value === '') {
    callback(new Error('请再次输入密码'))
  } else if (value !== registerForm.password) {
    callback(new Error('两次输入密码不一致'))
  } else {
    callback()
  }
}

const registerRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户名只能包含字母、数字和下划线', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, validator: validateConfirmPassword, trigger: 'blur' }
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

const handleRegister = async () => {
  if (!registerFormRef.value) return

  try {
    await registerFormRef.value.validate()
    loading.value = true

    const { confirmPassword, ...registerData } = registerForm
    await authStore.register(registerData)
    router.push('/')
  } catch (error) {
    console.error('注册失败:', error)
    ElMessage.error('注册失败，请检查输入信息')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-container {
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

.register-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  padding: 40px;
  width: 100%;
  max-width: 400px;
}

.mobile-register-card {
  margin-top: 60px;
  padding: 30px 20px;
  border-radius: 0;
  box-shadow: none;
  background: transparent;
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
}

.register-header h1 {
  margin: 0 0 10px 0;
  color: #333;
  font-size: 28px;
}

.register-header p {
  margin: 0;
  color: #666;
  font-size: 14px;
}

.register-form {
  margin-bottom: 20px;
}

.register-button {
  width: 100%;
  height: 44px;
  font-size: 16px;
}

.mobile-btn {
  height: 48px;
  border-radius: 24px;
  font-size: 16px;
}

.register-footer {
  text-align: center;
  margin-top: 20px;
}

.register-footer p {
  margin: 0 0 10px 0;
  color: #666;
  font-size: 14px;
}

.login-link {
  color: #409eff;
  text-decoration: none;
  font-weight: 500;
}

.login-link:hover {
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
  .register-container {
    padding: 0;
    background: white;
  }
  
  .register-card {
    box-shadow: none;
    border-radius: 0;
  }
  
  .register-header h1 {
    font-size: 24px;
  }
  
  .register-header p {
    font-size: 13px;
  }
}

@media (max-width: 480px) {
  .mobile-register-card {
    padding: 20px 15px;
  }
  
  .register-header {
    margin-bottom: 25px;
  }
}
</style> 