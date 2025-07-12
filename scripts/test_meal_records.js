// 测试用餐记录API
const axios = require('axios');

const API_BASE = 'http://192.168.2.184:8080';

async function testMealRecords() {
  try {
    // 1. 先注册一个用户
    console.log('1. 注册用户...');
    const registerResponse = await axios.post(`${API_BASE}/api/auth/register`, {
      username: 'testuser3',
      password: 'testpass123',
      email: 'test3@example.com'
    });
    
    const token = registerResponse.data.token;
    console.log('注册成功，获取到token:', token.substring(0, 50) + '...');

    // 2. 获取用餐记录列表
    console.log('\n2. 获取用餐记录列表...');
    const recordsResponse = await axios.get(`${API_BASE}/api/meal-records`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    console.log('用餐记录总数:', recordsResponse.data.total);
    console.log('用餐记录数据:');
    recordsResponse.data.data.forEach((record, index) => {
      console.log(`\n记录 ${index + 1}:`);
      console.log(`  ID: ${record.id}`);
      console.log(`  用户ID: ${record.user_id}`);
      console.log(`  总价格: ¥${record.total_price}`);
      console.log(`  感想: ${record.thoughts || '无'}`);
      console.log(`  创建时间: ${record.created_at}`);
      console.log(`  菜品数量: ${record.dishes.length}`);
      
      record.dishes.forEach((dishItem, dishIndex) => {
        console.log(`    菜品 ${dishIndex + 1}: ${dishItem.dish.name} (¥${dishItem.dish.price})`);
      });
    });

    // 3. 创建新的用餐记录
    console.log('\n3. 创建新的用餐记录...');
    const createResponse = await axios.post(`${API_BASE}/api/meal-records`, {
      dish_ids: [1, 3, 5],
      thoughts: '测试用餐记录 - 包含多个菜品',
      image_url: 'https://example.com/test-meal.jpg'
    }, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    console.log('创建成功:', createResponse.data);

    // 4. 再次获取用餐记录列表验证
    console.log('\n4. 验证新创建的记录...');
    const updatedResponse = await axios.get(`${API_BASE}/api/meal-records`, {
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });

    console.log('更新后的用餐记录总数:', updatedResponse.data.total);
    const newRecord = updatedResponse.data.data[0]; // 最新的记录
    console.log('最新记录包含菜品:');
    newRecord.dishes.forEach((dishItem, index) => {
      console.log(`  ${index + 1}. ${dishItem.dish.name} - ¥${dishItem.dish.price} - 数量: ${dishItem.quantity}`);
    });

  } catch (error) {
    console.error('测试失败:', error.response?.data || error.message);
  }
}

testMealRecords(); 