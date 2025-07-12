-- 修复旧的用餐记录数据
-- 为没有菜品关联的用餐记录添加默认菜品

-- 查找没有菜品关联的用餐记录
SELECT id, user_id, total_price, created_at 
FROM meal_records 
WHERE id NOT IN (SELECT DISTINCT meal_record_id FROM meal_record_dishes);

-- 为这些用餐记录添加默认菜品关联（假设菜品ID=1和2）
INSERT INTO meal_record_dishes (meal_record_id, dish_id, quantity) 
SELECT id, 1, 1 
FROM meal_records 
WHERE id NOT IN (SELECT DISTINCT meal_record_id FROM meal_record_dishes);

INSERT INTO meal_record_dishes (meal_record_id, dish_id, quantity) 
SELECT id, 2, 1 
FROM meal_records 
WHERE id NOT IN (SELECT DISTINCT meal_record_id FROM meal_record_dishes);

-- 验证修复结果
SELECT 
    mr.id,
    mr.user_id,
    mr.total_price,
    mr.created_at,
    COUNT(mrd.dish_id) as dish_count
FROM meal_records mr
LEFT JOIN meal_record_dishes mrd ON mr.id = mrd.meal_record_id
GROUP BY mr.id
ORDER BY mr.id; 