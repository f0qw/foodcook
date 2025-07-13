<template>
  <div class="home">
    <!-- 移动端导航栏 -->
    <el-header class="header mobile-header" v-if="isMobile">
      <div class="mobile-header-content">
        <div class="mobile-logo">
          <h2>🍽️ FoodCook</h2>
        </div>
        <div class="mobile-actions">
          <el-button 
            v-if="authStore.isAuthenticated"
            type="primary" 
            size="small"
            @click="showCartDialog = true"
            :disabled="cartStore.totalItems === 0"
          >
            🛒 {{ cartStore.totalItems }}
          </el-button>
          <el-button 
            v-else
            type="primary" 
            size="small"
            @click="$router.push('/login')"
          >
            登录
          </el-button>
        </div>
      </div>
    </el-header>

    <!-- 桌面端导航栏 -->
    <el-header class="header" v-else>
      <div class="header-content">
        <div class="logo">
          <h1>🍽️ FoodCook</h1>
        </div>
        <div class="nav">
          <el-menu mode="horizontal" :router="true" :default-active="$route.path">
            <el-menu-item index="/">首页</el-menu-item>
            <el-menu-item index="/dishes" v-if="authStore.isAuthenticated && authStore.isRoot">菜品管理</el-menu-item>
            <el-menu-item index="/ingredients" v-if="authStore.isAuthenticated && authStore.isRoot">食材管理</el-menu-item>
            <el-menu-item index="/meal-records" v-if="authStore.isAuthenticated">用餐记录</el-menu-item>
          </el-menu>
        </div>
        <div class="user-actions">
          <template v-if="authStore.isAuthenticated">
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                {{ authStore.user?.username }}
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" @click="$router.push('/login')">登录</el-button>
            <el-button @click="$router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </el-header>

    <!-- 主要内容 -->
    <div class="main-content" :class="{ 'mobile-content': isMobile }">
      <!-- 搜索区域 -->
      <div class="search-section">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索菜品..."
          class="search-input"
          :class="{ 'mobile-input': isMobile }"
          @keyup.enter="handleSearch"
        >
          <template #append>
            <el-button @click="handleSearch" :class="{ 'mobile-btn': isMobile }">
              <el-icon><Search /></el-icon>
            </el-button>
          </template>
        </el-input>
      </div>

      <!-- 菜品展示区域 -->
      <div class="dishes-section">
        <div class="section-header">
          <h2>今日推荐</h2>
          <div class="header-actions" v-if="authStore.isAuthenticated && !isMobile">
            <el-badge :value="cartStore.totalItems" :hidden="cartStore.totalItems === 0" class="cart-badge">
              <el-button type="info" @click="showCartDialog = true" :disabled="cartStore.totalItems === 0">
                <el-icon><ShoppingCart /></el-icon>
                购物车 ({{ cartStore.totalItems }})
              </el-button>
            </el-badge>
            <el-button type="primary" @click="showOrderDialog = true" :disabled="cartStore.totalItems === 0">
              创建用餐记录
            </el-button>
          </div>
        </div>

        <!-- 移动端菜品网格 -->
        <div v-if="isMobile" class="mobile-grid" v-loading="dishesStore.loading">
          <div 
            v-for="dish in dishesStore.dishes" 
            :key="dish.id" 
            class="mobile-dish-card"
            @click="showDishDetail(dish)"
          >
            <div class="mobile-dish-image">
              <el-image
                :src="dish.image_url || '/placeholder-dish.jpg'"
                fit="contain"
              >
                <template #error>
                  <div class="mobile-image-placeholder">
                    <el-icon><Picture /></el-icon>
                  </div>
                </template>
              </el-image>
            </div>
            <div class="mobile-dish-info">
              <h4>{{ dish.name }}</h4>
              <p class="mobile-dish-description">{{ dish.description }}</p>
              <div class="mobile-dish-meta">
                <span class="mobile-price">¥{{ dish.price }}</span>
                <span class="mobile-category" v-if="dish.category">{{ dish.category.name }}</span>
              </div>
              <el-button 
                v-if="authStore.isAuthenticated"
                :type="cartStore.isInCart(dish.id) ? 'danger' : 'primary'"
                size="small"
                class="mobile-cart-btn"
                @click.stop="toggleCartItem(dish)"
              >
                <el-icon v-if="cartStore.isInCart(dish.id)"><Remove /></el-icon>
                <el-icon v-else><Plus /></el-icon>
                {{ cartStore.isInCart(dish.id) ? '移除' : '加入' }}
              </el-button>
            </div>
          </div>
        </div>

        <!-- 桌面端菜品网格 -->
        <el-row v-else :gutter="20" v-loading="dishesStore.loading">
          <el-col 
            v-for="dish in dishesStore.dishes" 
            :key="dish.id" 
            :xs="24" 
            :sm="12" 
            :md="8" 
            :lg="6"
            class="dish-col"
          >
            <el-card class="dish-card">
              <div class="dish-card-content" @click="showDishDetail(dish)">
                <div class="dish-image">
                  <el-image
                    :src="dish.image_url || '/placeholder-dish.jpg'"
                    fit="contain"
                  >
                    <template #error>
                      <div class="image-placeholder">
                        <el-icon><Picture /></el-icon>
                      </div>
                    </template>
                  </el-image>
                </div>
                <div class="dish-info">
                  <h3>{{ dish.name }}</h3>
                  <p class="dish-description">{{ dish.description }}</p>
                  <div class="dish-meta">
                    <span class="price">¥{{ dish.price }}</span>
                    <span class="category" v-if="dish.category">{{ dish.category.name }}</span>
                  </div>
                </div>
                <div class="dish-actions">
                  <el-button 
                    v-if="authStore.isAuthenticated"
                    :type="cartStore.isInCart(dish.id) ? 'danger' : 'primary'"
                    size="small"
                    @click.stop="toggleCartItem(dish)"
                  >
                    <el-icon v-if="cartStore.isInCart(dish.id)"><Remove /></el-icon>
                    <el-icon v-else><Plus /></el-icon>
                    {{ cartStore.isInCart(dish.id) ? '移除' : '加入购物车' }}
                  </el-button>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 分页 -->
        <div class="pagination-wrapper" v-if="!isMobile">
          <el-pagination
            v-model:current-page="dishesStore.currentPage"
            v-model:page-size="dishesStore.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="dishesStore.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>

        <!-- 移动端加载更多 -->
        <div v-if="isMobile && dishesStore.hasMore" class="mobile-load-more">
          <el-button 
            type="primary" 
            :loading="dishesStore.loading"
            @click="loadMore"
            class="mobile-btn"
          >
            加载更多
          </el-button>
        </div>
      </div>
    </div>

    <!-- 移动端底部导航 -->
    <div v-if="isMobile && authStore.isAuthenticated" class="mobile-bottom-nav">
      <div class="bottom-nav-content">
        <div class="nav-item" @click="$router.push('/')">
          <el-icon><House /></el-icon>
          <span>首页</span>
        </div>
        <div class="nav-item" @click="showCartDialog = true">
          <el-icon><ShoppingCart /></el-icon>
          <span>购物车</span>
          <el-badge v-if="cartStore.totalItems > 0" :value="cartStore.totalItems" class="nav-badge" />
        </div>
        <div class="nav-item" @click="$router.push('/meal-records')">
          <el-icon><Document /></el-icon>
          <span>记录</span>
        </div>
        <div class="nav-item" @click="showUserMenu = true">
          <el-icon><User /></el-icon>
          <span>我的</span>
        </div>
      </div>
    </div>

    <!-- 菜品详情对话框 -->
    <el-dialog
      v-model="dishDetailVisible"
      title="菜品详情"
      :width="isMobile ? '90%' : '600px'"
      :before-close="handleCloseDishDetail"
    >
      <div v-if="selectedDish" class="dish-detail">
        <div class="dish-detail-image">
          <el-image
            :src="selectedDish.image_url || '/placeholder-dish.jpg'"
            fit="contain"
            style="width: 100%; height: 200px;"
          />
        </div>
        <div class="dish-detail-info">
          <h3>{{ selectedDish.name }}</h3>
          <p class="description">{{ selectedDish.description }}</p>
          <div class="meta">
            <p><strong>价格:</strong> ¥{{ selectedDish.price }}</p>
            <p><strong>分类:</strong> {{ selectedDish.category?.name || '未分类' }}</p>
            <p v-if="selectedDish.cooking_link">
              <strong>制作链接:</strong> 
              <a :href="selectedDish.cooking_link" target="_blank">查看制作方法</a>
            </p>
          </div>
          <div v-if="selectedDish.ingredients && selectedDish.ingredients.length > 0" class="ingredients">
            <h4>所需食材:</h4>
            <el-tag
              v-for="ingredient in selectedDish.ingredients"
              :key="ingredient.id"
              class="ingredient-tag"
            >
              {{ ingredient.ingredient.name }} {{ ingredient.quantity }}{{ ingredient.ingredient.unit }}
            </el-tag>
          </div>
        </div>
      </div>
    </el-dialog>

    <!-- 购物车对话框 -->
    <el-dialog
      v-model="showCartDialog"
      title="购物车"
      :width="isMobile ? '90%' : '600px'"
    >
      <div v-if="cartStore.totalItems === 0" class="empty-cart">
        <el-empty description="购物车是空的" />
      </div>
      <div v-else class="cart-content">
        <div class="cart-items">
          <div v-for="item in cartStore.cartItems" :key="item.id" class="cart-item">
            <div class="cart-item-info">
              <el-image
                :src="item.image_url || '/placeholder-dish.jpg'"
                style="width: 60px; height: 60px; border-radius: 8px;"
                fit="contain"
              />
              <div class="cart-item-details">
                <h4>{{ item.name }}</h4>
                <p class="cart-item-price">¥{{ item.price }}</p>
              </div>
            </div>
            <div class="cart-item-actions">
              <el-button-group>
                <el-button size="small" @click="cartStore.decreaseQuantity(item.id)">-</el-button>
                <el-button size="small" disabled>{{ item.quantity }}</el-button>
                <el-button size="small" @click="cartStore.increaseQuantity(item.id)">+</el-button>
              </el-button-group>
              <el-button 
                type="danger" 
                size="small" 
                @click="cartStore.removeFromCart(item.id)"
              >
                删除
              </el-button>
            </div>
          </div>
        </div>
        <div class="cart-summary">
          <div class="cart-total">
            <span>总计: ¥{{ cartStore.totalPrice }}</span>
          </div>
          <el-button type="primary" @click="showOrderDialog = true" class="mobile-btn">
            创建用餐记录
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 创建用餐记录对话框 -->
    <el-dialog
      v-model="showOrderDialog"
      title="创建用餐记录"
      :width="isMobile ? '90%' : '600px'"
    >
      <el-form :model="orderForm" label-width="80px">
        <el-form-item label="总价">
          <el-input v-model="orderForm.totalPrice" disabled />
        </el-form-item>
        <el-form-item label="用餐感想">
          <el-input 
            v-model="orderForm.thoughts" 
            type="textarea" 
            :rows="3"
            placeholder="记录一下今天的用餐体验..."
          />
        </el-form-item>
        <el-form-item label="图片链接">
          <el-input 
            v-model="orderForm.imageUrl" 
            placeholder="可选：添加用餐照片链接"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showOrderDialog = false">取消</el-button>
        <el-button type="primary" @click="createMealRecord" :loading="creatingOrder">
          创建记录
        </el-button>
      </template>
    </el-dialog>

    <!-- 移动端用户菜单 -->
    <el-drawer
      v-model="showUserMenu"
      title="个人中心"
      direction="rtl"
      size="80%"
    >
      <div class="mobile-user-menu">
        <div class="user-info">
          <el-avatar :size="60" icon="UserFilled" />
          <h3>{{ authStore.user?.username }}</h3>
          <p>{{ authStore.user?.email }}</p>
        </div>
        <el-menu>
          <el-menu-item @click="$router.push('/meal-records')">
            <el-icon><Document /></el-icon>
            <span>我的用餐记录</span>
          </el-menu-item>
          <el-menu-item v-if="authStore.isRoot" @click="$router.push('/dishes')">
            <el-icon><Food /></el-icon>
            <span>菜品管理</span>
          </el-menu-item>
          <el-menu-item v-if="authStore.isRoot" @click="$router.push('/ingredients')">
            <el-icon><Basket /></el-icon>
            <span>食材管理</span>
          </el-menu-item>
          <el-menu-item @click="handleLogout">
            <el-icon><SwitchButton /></el-icon>
            <span>退出登录</span>
          </el-menu-item>
        </el-menu>
      </div>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import { useDishesStore } from '../stores/dishes'
import { useCartStore } from '../stores/cart'
import { useMealRecordsStore } from '../stores/mealRecords'

const router = useRouter()
const authStore = useAuthStore()
const dishesStore = useDishesStore()
const cartStore = useCartStore()
const mealRecordsStore = useMealRecordsStore()

// 响应式数据
const searchKeyword = ref('')
const dishDetailVisible = ref(false)
const showCartDialog = ref(false)
const showOrderDialog = ref(false)
const showUserMenu = ref(false)
const selectedDish = ref(null)
const creatingOrder = ref(false)
const isMobile = ref(false)

// 订单表单
const orderForm = ref({
  totalPrice: 0,
  thoughts: '',
  imageUrl: ''
})

// 检测移动端
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

// 监听窗口大小变化
onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  loadDishes()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

// 加载菜品
const loadDishes = async () => {
  await dishesStore.fetchDishes()
}

// 搜索菜品
const handleSearch = () => {
  dishesStore.searchDishes(searchKeyword.value)
}

// 显示菜品详情
const showDishDetail = (dish) => {
  selectedDish.value = dish
  dishDetailVisible.value = true
}

// 关闭菜品详情
const handleCloseDishDetail = () => {
  dishDetailVisible.value = false
  selectedDish.value = null
}

// 切换购物车项目
const toggleCartItem = (dish) => {
  if (cartStore.isInCart(dish.id)) {
    cartStore.removeFromCart(dish.id)
    ElMessage.success('已从购物车移除')
  } else {
    cartStore.addToCart(dish)
    ElMessage.success('已加入购物车')
  }
}

// 创建用餐记录
const createMealRecord = async () => {
  if (cartStore.totalItems === 0) {
    ElMessage.warning('购物车为空')
    return
  }

  creatingOrder.value = true
  try {
    const dishIds = cartStore.cartItems.map(item => item.id)

    await mealRecordsStore.createMealRecord({
      dish_ids: dishIds,
      thoughts: orderForm.value.thoughts,
      image_url: orderForm.value.imageUrl
    })

    cartStore.clearCart()
    showOrderDialog.value = false
    orderForm.value = { totalPrice: 0, thoughts: '', imageUrl: '' }
    ElMessage.success('用餐记录创建成功')
  } catch (error) {
    ElMessage.error('创建用餐记录失败')
  } finally {
    creatingOrder.value = false
  }
}

// 加载更多
const loadMore = async () => {
  await dishesStore.loadMore()
}

// 分页处理
const handleSizeChange = (size) => {
  dishesStore.setPageSize(size)
  loadDishes()
}

const handleCurrentChange = (page) => {
  dishesStore.setCurrentPage(page)
  loadDishes()
}

// 用户操作
const handleCommand = (command) => {
  if (command === 'logout') {
    handleLogout()
  } else if (command === 'profile') {
    // 处理个人中心
  }
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
  ElMessage.success('已退出登录')
}
</script>

<style scoped>
.home {
  min-height: 100vh;
  background-color: #f5f5f5;
}

/* 桌面端样式 */
.header {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.logo h1 {
  margin: 0;
  color: #409EFF;
  font-size: 24px;
}

.nav {
  flex: 1;
  margin: 0 40px;
}

.user-actions {
  display: flex;
  gap: 10px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 5px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.user-info:hover {
  background-color: #f5f5f5;
}

.main-content {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.search-section {
  margin-bottom: 30px;
}

.search-input {
  max-width: 400px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h2 {
  margin: 0;
  color: #333;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.dish-col {
  margin-bottom: 20px;
}

.dish-card {
  height: 100%;
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}

.dish-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.dish-card-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.dish-image {
  height: 200px;
  overflow: hidden;
  border-radius: 8px;
  margin-bottom: 15px;
}

.dish-image .el-image {
  width: 100%;
  height: 100%;
}

.image-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background-color: #f5f5f5;
  color: #999;
  font-size: 24px;
}

.dish-info {
  flex: 1;
  margin-bottom: 15px;
}

.dish-info h3 {
  margin: 0 0 8px 0;
  color: #333;
  font-size: 18px;
}

.dish-description {
  color: #666;
  font-size: 14px;
  line-height: 1.4;
  margin-bottom: 10px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.dish-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.price {
  color: #f56c6c;
  font-size: 18px;
  font-weight: bold;
}

.category {
  color: #409EFF;
  font-size: 12px;
  background-color: #ecf5ff;
  padding: 2px 8px;
  border-radius: 12px;
}

.dish-actions {
  display: flex;
  justify-content: center;
}

.pagination-wrapper {
  margin-top: 30px;
  display: flex;
  justify-content: center;
}

/* 移动端样式 */
.mobile-header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 15px;
  height: 100%;
}

.mobile-logo h2 {
  margin: 0;
  color: #409EFF;
  font-size: 20px;
}

.mobile-actions {
  display: flex;
  gap: 10px;
}

.mobile-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  padding: 10px 0;
}

.mobile-dish-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: transform 0.2s;
}

.mobile-dish-card:active {
  transform: scale(0.98);
}

.mobile-dish-image {
  height: 120px;
  overflow: hidden;
}

.mobile-dish-image .el-image {
  width: 100%;
  height: 100%;
}

.mobile-image-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  background-color: #f5f5f5;
  color: #999;
  font-size: 20px;
}

.mobile-dish-info {
  padding: 12px;
}

.mobile-dish-info h4 {
  margin: 0 0 6px 0;
  font-size: 16px;
  color: #333;
}

.mobile-dish-description {
  color: #666;
  font-size: 12px;
  line-height: 1.3;
  margin-bottom: 8px;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.mobile-dish-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.mobile-price {
  color: #f56c6c;
  font-size: 16px;
  font-weight: bold;
}

.mobile-category {
  color: #409EFF;
  font-size: 10px;
  background-color: #ecf5ff;
  padding: 2px 6px;
  border-radius: 8px;
}

.mobile-cart-btn {
  width: 100%;
  height: 32px;
  font-size: 12px;
}

.mobile-load-more {
  text-align: center;
  margin: 20px 0;
}

.mobile-bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  border-top: 1px solid #eee;
  z-index: 1000;
  padding: 8px 0;
}

.bottom-nav-content {
  display: flex;
  justify-content: space-around;
  align-items: center;
}

.nav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px;
  cursor: pointer;
  position: relative;
  color: #666;
  transition: color 0.3s;
}

.nav-item:hover {
  color: #409EFF;
}

.nav-item span {
  font-size: 12px;
}

.nav-badge {
  position: absolute;
  top: 0;
  right: 0;
}

.mobile-user-menu {
  padding: 20px;
}

.mobile-user-menu .user-info {
  text-align: center;
  margin-bottom: 30px;
}

.mobile-user-menu .user-info h3 {
  margin: 10px 0 5px 0;
  color: #333;
}

.mobile-user-menu .user-info p {
  color: #666;
  margin: 0;
}

/* 对话框样式 */
.dish-detail {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.dish-detail-image {
  border-radius: 8px;
  overflow: hidden;
}

.dish-detail-info h3 {
  margin: 0 0 10px 0;
  color: #333;
}

.description {
  color: #666;
  line-height: 1.6;
  margin-bottom: 15px;
}

.meta p {
  margin: 8px 0;
  color: #333;
}

.meta a {
  color: #409EFF;
  text-decoration: none;
}

.ingredients {
  margin-top: 15px;
}

.ingredients h4 {
  margin: 0 0 10px 0;
  color: #333;
}

.ingredient-tag {
  margin: 0 5px 5px 0;
}

.cart-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.cart-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.cart-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 8px;
}

.cart-item-info {
  display: flex;
  align-items: center;
  gap: 15px;
  flex: 1;
}

.cart-item-details h4 {
  margin: 0 0 5px 0;
  color: #333;
}

.cart-item-price {
  color: #f56c6c;
  font-weight: bold;
  margin: 0;
}

.cart-item-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-end;
}

.cart-summary {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.cart-total {
  font-size: 18px;
  font-weight: bold;
  color: #333;
}

.empty-cart {
  text-align: center;
  padding: 40px 0;
}

/* 响应式适配 */
@media (max-width: 768px) {
  .main-content {
    padding: 10px;
  }
  
  .search-section {
    margin-bottom: 20px;
  }
  
  .section-header {
    flex-direction: column;
    gap: 15px;
    align-items: flex-start;
  }
  
  .mobile-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .mobile-grid {
    grid-template-columns: 1fr;
  }
  
  .main-content {
    padding: 5px;
  }
}
</style> 