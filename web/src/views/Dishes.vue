<template>
  <div class="dishes-page">
    <!-- 权限检查 -->
    <div v-if="!authStore.isRoot" class="permission-denied">
      <el-result
        icon="warning"
        title="权限不足"
        sub-title="只有 root 用户才能访问菜品管理功能"
      >
        <template #extra>
          <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
        </template>
      </el-result>
    </div>

    <!-- 菜品管理内容 -->
    <div v-else>
      <div class="page-header">
        <h1>菜品管理</h1>
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          添加菜品
        </el-button>
      </div>

      <!-- 搜索和筛选 -->
      <div class="search-section">
        <el-row :gutter="20">
          <el-col :span="8">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索菜品名称..."
              @keyup.enter="handleSearch"
            >
              <template #append>
                <el-button @click="handleSearch">
                  <el-icon><Search /></el-icon>
                </el-button>
              </template>
            </el-input>
          </el-col>
          <el-col :span="6">
            <el-select v-model="selectedCategory" placeholder="选择分类" clearable @change="handleSearch">
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-col>
          <el-col :span="4">
            <el-button @click="resetSearch">重置</el-button>
          </el-col>
        </el-row>
      </div>

      <!-- 菜品列表 -->
      <el-table
        :data="dishesStore.dishes"
        v-loading="dishesStore.loading"
        style="width: 100%"
        class="dishes-table"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="图片" width="120">
          <template #default="{ row }">
            <el-image
              :src="row.image_url || '/placeholder-dish.jpg'"
              style="width: 80px; height: 60px"
              fit="cover"
            >
              <template #error>
                <div class="image-placeholder">
                  <el-icon><Picture /></el-icon>
                </div>
              </template>
            </el-image>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="菜品名称" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="price" label="价格" width="100">
          <template #default="{ row }">
            ¥{{ row.price }}
          </template>
        </el-table-column>
        <el-table-column prop="category.name" label="分类" width="120" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="editDish(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteDish(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
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

      <!-- 创建/编辑菜品对话框 -->
      <el-dialog
        v-model="showCreateDialog"
        :title="editingDish ? '编辑菜品' : '添加菜品'"
        width="600px"
      >
        <el-form
          ref="dishFormRef"
          :model="dishForm"
          :rules="dishRules"
          label-width="100px"
        >
          <el-form-item label="菜品名称" prop="name">
            <el-input v-model="dishForm.name" placeholder="请输入菜品名称" />
          </el-form-item>
          
          <el-form-item label="描述" prop="description">
            <el-input
              v-model="dishForm.description"
              type="textarea"
              :rows="3"
              placeholder="请输入菜品描述"
            />
          </el-form-item>
          
          <el-form-item label="价格" prop="price">
            <el-input-number
              v-model="dishForm.price"
              :precision="2"
              :step="0.1"
              :min="0"
              style="width: 100%"
            />
          </el-form-item>
          
          <el-form-item label="分类" prop="category_id">
            <el-select v-model="dishForm.category_id" placeholder="选择分类" style="width: 100%">
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>
          
          <el-form-item label="图片链接" prop="image_url">
            <el-input v-model="dishForm.image_url" placeholder="请输入图片链接" />
          </el-form-item>
          
          <el-form-item label="制作链接" prop="cooking_link">
            <el-input v-model="dishForm.cooking_link" placeholder="请输入制作方法链接" />
          </el-form-item>
          
          <el-form-item label="食材" prop="ingredients">
            <div class="ingredients-section">
              <div v-for="(ingredient, index) in dishForm.ingredients" :key="index" class="ingredient-item">
                <el-row :gutter="10">
                  <el-col :span="12">
                    <el-select 
                      v-model="ingredient.ingredient_id" 
                      placeholder="选择食材"
                      style="width: 100%"
                    >
                      <el-option
                        v-for="ing in ingredientsStore.ingredients"
                        :key="ing.id"
                        :label="`${ing.name} (${ing.unit})`"
                        :value="ing.id"
                      />
                    </el-select>
                  </el-col>
                  <el-col :span="8">
                    <el-input-number
                      v-model="ingredient.quantity"
                      :precision="2"
                      :step="0.1"
                      :min="0"
                      placeholder="数量"
                      style="width: 100%"
                    />
                  </el-col>
                  <el-col :span="4">
                    <el-button 
                      type="danger" 
                      size="small" 
                      @click="removeIngredient(index)"
                    >
                      删除
                    </el-button>
                  </el-col>
                </el-row>
              </div>
              <el-button 
                type="primary" 
                size="small" 
                @click="addIngredient"
                style="margin-top: 10px;"
              >
                添加食材
              </el-button>
            </div>
          </el-form-item>
        </el-form>
        
        <template #footer>
          <span class="dialog-footer">
            <el-button @click="showCreateDialog = false">取消</el-button>
            <el-button type="primary" @click="saveDish" :loading="saving">
              {{ editingDish ? '更新' : '创建' }}
            </el-button>
          </span>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useDishesStore } from '@/stores/dishes'
import { useIngredientsStore } from '@/stores/ingredients'
import { useAuthStore } from '@/stores/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const dishesStore = useDishesStore()
const ingredientsStore = useIngredientsStore()
const authStore = useAuthStore()

// 响应式数据
const searchKeyword = ref('')
const selectedCategory = ref('')
const showCreateDialog = ref(false)
const editingDish = ref(null)
const saving = ref(false)
const dishFormRef = ref()

// 模拟分类数据
const categories = ref([
  { id: 1, name: '川菜' },
  { id: 2, name: '粤菜' },
  { id: 3, name: '湘菜' },
  { id: 4, name: '鲁菜' },
  { id: 5, name: '苏菜' },
  { id: 6, name: '浙菜' },
  { id: 7, name: '闽菜' },
  { id: 8, name: '徽菜' }
])

const dishForm = reactive({
  name: '',
  description: '',
  price: 0,
  category_id: '',
  image_url: '',
  cooking_link: '',
  ingredients: []
})

const dishRules = {
  name: [
    { required: true, message: '请输入菜品名称', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' }
  ],
  category_id: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ]
}

// 方法
const handleSearch = () => {
  const params = {}
  if (searchKeyword.value) {
    params.q = searchKeyword.value
  }
  if (selectedCategory.value) {
    params.category_id = selectedCategory.value
  }
  dishesStore.getDishes(params)
}

const resetSearch = () => {
  searchKeyword.value = ''
  selectedCategory.value = ''
  dishesStore.getDishes()
}

const handleSizeChange = (size) => {
  dishesStore.setPageSize(size)
  dishesStore.getDishes()
}

const handleCurrentChange = (page) => {
  dishesStore.setPage(page)
  dishesStore.getDishes()
}

const editDish = (dish) => {
  editingDish.value = dish
  Object.assign(dishForm, {
    name: dish.name,
    description: dish.description,
    price: dish.price,
    category_id: dish.category_id,
    image_url: dish.image_url,
    cooking_link: dish.cooking_link,
    ingredients: dish.ingredients ? dish.ingredients.map(ing => ({
      ingredient_id: ing.ingredient_id,
      quantity: ing.quantity
    })) : []
  })
  showCreateDialog.value = true
}

const saveDish = async () => {
  if (!dishFormRef.value) return

  try {
    await dishFormRef.value.validate()
    saving.value = true

    if (editingDish.value) {
      await dishesStore.updateDish(editingDish.value.id, dishForm)
    } else {
      await dishesStore.createDish(dishForm)
    }

    showCreateDialog.value = false
    resetForm()
  } catch (error) {
    console.error('保存菜品失败:', error)
  } finally {
    saving.value = false
  }
}

const deleteDish = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这个菜品吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await dishesStore.deleteDish(id)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除菜品失败:', error)
    }
  }
}

const resetForm = () => {
  editingDish.value = null
  Object.assign(dishForm, {
    name: '',
    description: '',
    price: 0,
    category_id: '',
    image_url: '',
    cooking_link: '',
    ingredients: []
  })
  if (dishFormRef.value) {
    dishFormRef.value.resetFields()
  }
}

const addIngredient = () => {
  dishForm.ingredients.push({
    ingredient_id: '',
    quantity: 1
  })
}

const removeIngredient = (index) => {
  dishForm.ingredients.splice(index, 1)
}

// 监听对话框关闭
const handleDialogClose = () => {
  resetForm()
}

// 生命周期
onMounted(() => {
  dishesStore.getDishes()
  ingredientsStore.getIngredients()
})
</script>

<style scoped>
.dishes-page {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h1 {
  margin: 0;
  color: #333;
}

.search-section {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.dishes-table {
  margin-bottom: 20px;
}

.image-placeholder {
  width: 80px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f5f5;
  color: #999;
}

.ingredients-section {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  padding: 15px;
  background-color: #fafafa;
}

.ingredient-item {
  margin-bottom: 10px;
  padding: 10px;
  background-color: white;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}

.ingredient-item:last-child {
  margin-bottom: 0;
}

.pagination-wrapper {
  text-align: center;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.pagination-wrapper {
  text-align: center;
  margin-top: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.permission-denied {
  margin-top: 20px;
  padding: 20px;
  background-color: #fffbe6;
  border: 1px solid #ffe58f;
  border-radius: 8px;
  text-align: center;
}

.permission-denied .el-result__title {
  color: #faad14;
}

.permission-denied .el-result__subtitle {
  color: #faad14;
}
</style> 