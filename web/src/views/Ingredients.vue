<template>
  <div class="ingredients-page">
    <div class="page-header">
      <h1>食材管理</h1>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加食材
      </el-button>
    </div>

    <!-- 搜索和筛选 -->
    <div class="search-section">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索食材名称..."
            @keyup.enter="handleSearch"
          >
            <template #append>
              <el-button @click="handleSearch">
                <el-icon><Search /></el-icon>
              </el-button>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-button @click="resetSearch">重置</el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 食材列表 -->
    <el-table
      :data="ingredients"
      v-loading="loading"
      style="width: 100%"
      class="ingredients-table"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="食材名称" />
      <el-table-column prop="unit" label="单位" width="100" />
      <el-table-column prop="price" label="单价" width="120">
        <template #default="{ row }">
          ¥{{ row.price }}/{{ row.unit }}
        </template>
      </el-table-column>
      <el-table-column prop="stock" label="库存" width="100" />
      <el-table-column prop="description" label="描述" show-overflow-tooltip />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="editIngredient(row)">编辑</el-button>
          <el-button size="small" type="danger" @click="deleteIngredient(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 创建/编辑食材对话框 -->
    <el-dialog
      v-model="showCreateDialog"
      :title="editingIngredient ? '编辑食材' : '添加食材'"
      width="500px"
    >
      <el-form
        ref="ingredientFormRef"
        :model="ingredientForm"
        :rules="ingredientRules"
        label-width="100px"
      >
        <el-form-item label="食材名称" prop="name">
          <el-input v-model="ingredientForm.name" placeholder="请输入食材名称" />
        </el-form-item>
        
        <el-form-item label="单位" prop="unit">
          <el-select v-model="ingredientForm.unit" placeholder="选择单位" style="width: 100%">
            <el-option label="克" value="克" />
            <el-option label="千克" value="千克" />
            <el-option label="个" value="个" />
            <el-option label="包" value="包" />
            <el-option label="瓶" value="瓶" />
            <el-option label="袋" value="袋" />
            <el-option label="盒" value="盒" />
            <el-option label="斤" value="斤" />
            <el-option label="两" value="两" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="单价" prop="price">
          <el-input-number
            v-model="ingredientForm.price"
            :precision="2"
            :step="0.1"
            :min="0"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="库存" prop="stock">
          <el-input-number
            v-model="ingredientForm.stock"
            :precision="2"
            :step="0.1"
            :min="0"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="ingredientForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入食材描述"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateDialog = false">取消</el-button>
          <el-button type="primary" @click="saveIngredient" :loading="saving">
            {{ editingIngredient ? '更新' : '创建' }}
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ingredientsAPI } from '@/utils/api'
import { ElMessage, ElMessageBox } from 'element-plus'

// 响应式数据
const searchKeyword = ref('')
const showCreateDialog = ref(false)
const editingIngredient = ref(null)
const saving = ref(false)
const loading = ref(false)
const ingredientFormRef = ref()

const ingredients = ref([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const ingredientForm = reactive({
  name: '',
  unit: '',
  price: 0,
  stock: 0,
  description: ''
})

const ingredientRules = {
  name: [
    { required: true, message: '请输入食材名称', trigger: 'blur' }
  ],
  unit: [
    { required: true, message: '请选择单位', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入单价', trigger: 'blur' }
  ],
  stock: [
    { required: true, message: '请输入库存', trigger: 'blur' }
  ]
}

// 方法
const getIngredients = async (params = {}) => {
  loading.value = true
  try {
    const response = await ingredientsAPI.getList({
      offset: (currentPage.value - 1) * pageSize.value,
      limit: pageSize.value,
      ...params
    })
    ingredients.value = response.data
    total.value = response.total
  } catch (error) {
    console.error('获取食材列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  const params = {}
  if (searchKeyword.value) {
    params.q = searchKeyword.value
  }
  getIngredients(params)
}

const resetSearch = () => {
  searchKeyword.value = ''
  getIngredients()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  getIngredients()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  getIngredients()
}

const editIngredient = (ingredient) => {
  editingIngredient.value = ingredient
  Object.assign(ingredientForm, {
    name: ingredient.name,
    unit: ingredient.unit,
    price: ingredient.price,
    stock: ingredient.stock,
    description: ingredient.description
  })
  showCreateDialog.value = true
}

const saveIngredient = async () => {
  if (!ingredientFormRef.value) return

  try {
    await ingredientFormRef.value.validate()
    saving.value = true

    if (editingIngredient.value) {
      await ingredientsAPI.update(editingIngredient.value.id, ingredientForm)
      ElMessage.success('食材更新成功')
    } else {
      await ingredientsAPI.create(ingredientForm)
      ElMessage.success('食材创建成功')
    }

    showCreateDialog.value = false
    resetForm()
    getIngredients()
  } catch (error) {
    console.error('保存食材失败:', error)
  } finally {
    saving.value = false
  }
}

const deleteIngredient = async (id) => {
  try {
    await ElMessageBox.confirm('确定要删除这个食材吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await ingredientsAPI.delete(id)
    ElMessage.success('食材删除成功')
    getIngredients()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除食材失败:', error)
    }
  }
}

const resetForm = () => {
  editingIngredient.value = null
  Object.assign(ingredientForm, {
    name: '',
    unit: '',
    price: 0,
    stock: 0,
    description: ''
  })
  if (ingredientFormRef.value) {
    ingredientFormRef.value.resetFields()
  }
}

// 生命周期
onMounted(() => {
  getIngredients()
})
</script>

<style scoped>
.ingredients-page {
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

.ingredients-table {
  margin-bottom: 20px;
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
</style> 