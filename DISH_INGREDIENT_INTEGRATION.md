# 菜品食材关联功能实现总结

## 功能需求

在菜品管理页面，创建新菜品和编辑菜品时，可以选择该菜品使用了哪些食材，建立菜品和食材的多对多关联关系。

## 实现方案

### 1. 数据模型设计

#### 现有模型
- `Dish`: 菜品模型
- `Ingredient`: 食材模型  
- `DishIngredient`: 菜品食材关联模型（已存在）

#### 关联关系
```go
type Dish struct {
    // ... 其他字段
    Ingredients []DishIngredient `json:"ingredients,omitempty" gorm:"foreignKey:DishID"`
}

type DishIngredient struct {
    DishID       uint    `json:"dish_id"`
    IngredientID uint    `json:"ingredient_id"`
    Quantity     float64 `json:"quantity"`
    
    Dish       *Dish       `json:"dish,omitempty"`
    Ingredient *Ingredient `json:"ingredient,omitempty"`
}
```

### 2. 后端API实现

#### 扩展Repository接口
```go
type DishRepository interface {
    CreateWithIngredients(ctx context.Context, dish *models.Dish, ingredients []DishIngredientRequest) error
    UpdateWithIngredients(ctx context.Context, dish *models.Dish, ingredients []DishIngredientRequest) error
    // ... 其他方法
}

type DishIngredientRequest struct {
    IngredientID uint    `json:"ingredient_id"`
    Quantity     float64 `json:"quantity"`
}
```

#### 更新Handler
- 修改 `CreateDishRequest` 和 `UpdateDishRequest` 结构体，添加 `Ingredients` 字段
- 更新 `Create` 和 `Update` 方法，使用事务处理菜品和食材关联
- 更新 `GetByID` 和 `List` 方法，预加载食材信息

#### MySQL实现
```go
func (r *MySQLDishRepository) CreateWithIngredients(ctx context.Context, dish *models.Dish, ingredients []repositories.DishIngredientRequest) error {
    tx := r.db.WithContext(ctx).Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 创建菜品
    if err := tx.Create(dish).Error; err != nil {
        tx.Rollback()
        return err
    }

    // 创建食材关联
    for _, ingredient := range ingredients {
        dishIngredient := &models.DishIngredient{
            DishID:       dish.ID,
            IngredientID: ingredient.IngredientID,
            Quantity:     ingredient.Quantity,
        }
        if err := tx.Create(dishIngredient).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    return tx.Commit().Error
}
```

### 3. 前端实现

#### 更新菜品表单
- 添加食材选择区域
- 支持动态添加/删除食材
- 每个食材包含：食材选择、数量输入

#### 表单结构
```vue
<el-form-item label="食材" prop="ingredients">
  <div class="ingredients-section">
    <div v-for="(ingredient, index) in dishForm.ingredients" :key="index" class="ingredient-item">
      <el-row :gutter="10">
        <el-col :span="12">
          <el-select v-model="ingredient.ingredient_id" placeholder="选择食材">
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
          />
        </el-col>
        <el-col :span="4">
          <el-button type="danger" size="small" @click="removeIngredient(index)">
            删除
          </el-button>
        </el-col>
      </el-row>
    </div>
    <el-button type="primary" size="small" @click="addIngredient">
      添加食材
    </el-button>
  </div>
</el-form-item>
```

#### 数据管理
```javascript
const dishForm = reactive({
  name: '',
  description: '',
  price: 0,
  category_id: '',
  image_url: '',
  cooking_link: '',
  ingredients: []  // 新增食材数组
})

const addIngredient = () => {
  dishForm.ingredients.push({
    ingredient_id: '',
    quantity: 1
  })
}

const removeIngredient = (index) => {
  dishForm.ingredients.splice(index, 1)
}
```

### 4. API测试

#### 创建带食材的菜品
```bash
curl -X POST "http://192.168.2.184:8080/api/dishes" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "宫保鸡丁",
    "description": "经典川菜，酸甜微辣",
    "price": 32.00,
    "category_id": 1,
    "image_url": "https://example.com/kungpao-chicken.jpg",
    "cooking_link": "https://example.com/recipe",
    "ingredients": [
      {"ingredient_id": 1, "quantity": 0.5},
      {"ingredient_id": 3, "quantity": 0.3}
    ]
  }'
```

#### 获取菜品详情（包含食材）
```bash
curl -X GET "http://192.168.2.184:8080/api/dishes/10" \
  -H "Content-Type: application/json"
```

### 5. 功能特点

#### 数据一致性
- 使用数据库事务确保菜品和食材关联的原子性操作
- 更新时先删除旧关联，再创建新关联

#### 用户体验
- 动态添加/删除食材
- 食材选择下拉框显示名称和单位
- 数量支持小数点，精确到0.01

#### 数据完整性
- 创建和更新时验证食材ID的有效性
- 数量必须大于0
- 支持空食材列表（菜品可以不使用食材）

### 6. 数据库结构

#### dish_ingredients 表
```sql
CREATE TABLE dish_ingredients (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  dish_id BIGINT UNSIGNED NOT NULL,
  ingredient_id BIGINT UNSIGNED NOT NULL,
  quantity DECIMAL(10,2) NOT NULL,
  FOREIGN KEY (dish_id) REFERENCES dishes(id),
  FOREIGN KEY (ingredient_id) REFERENCES ingredients(id)
);
```

### 7. 后续优化建议

1. **食材库存管理**: 根据菜品食材关联计算库存消耗
2. **成本计算**: 根据食材价格和数量自动计算菜品成本
3. **营养分析**: 基于食材营养成分计算菜品营养价值
4. **批量操作**: 支持批量设置菜品食材
5. **食材推荐**: 根据常用食材推荐相关菜品

## 总结

成功实现了菜品和食材的多对多关联功能，包括：

- ✅ 后端API支持创建和更新菜品时关联食材
- ✅ 前端表单支持动态添加/删除食材
- ✅ 数据库事务确保数据一致性
- ✅ 预加载食材信息，优化查询性能
- ✅ 完整的CRUD操作支持

现在用户可以在菜品管理页面为每个菜品指定使用的食材和数量，实现了菜品和食材的完整关联。 