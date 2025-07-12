-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- 创建菜品分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建菜品表
CREATE TABLE IF NOT EXISTS dishes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    price DECIMAL(10,2) NOT NULL,
    cooking_link VARCHAR(255),
    category_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);

-- 创建食材表
CREATE TABLE IF NOT EXISTS ingredients (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建菜品食材关联表
CREATE TABLE IF NOT EXISTS dish_ingredients (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    dish_id BIGINT NOT NULL,
    ingredient_id BIGINT NOT NULL,
    quantity DECIMAL(10,2) NOT NULL,
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE,
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE CASCADE
);

-- 创建用餐记录表
CREATE TABLE IF NOT EXISTS meal_records (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    thoughts TEXT,
    image_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 创建用餐记录菜品关联表
CREATE TABLE IF NOT EXISTS meal_record_dishes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    meal_record_id BIGINT NOT NULL,
    dish_id BIGINT NOT NULL,
    quantity INT DEFAULT 1,
    FOREIGN KEY (meal_record_id) REFERENCES meal_records(id) ON DELETE CASCADE,
    FOREIGN KEY (dish_id) REFERENCES dishes(id)
);

-- 插入示例分类数据
INSERT INTO categories (name, description) VALUES
('川菜', '四川菜系，以麻辣著称'),
('粤菜', '广东菜系，清淡鲜美'),
('湘菜', '湖南菜系，香辣可口'),
('鲁菜', '山东菜系，咸鲜为主'),
('苏菜', '江苏菜系，清淡雅致'),
('浙菜', '浙江菜系，清淡爽口'),
('闽菜', '福建菜系，清淡鲜美'),
('徽菜', '安徽菜系，咸鲜为主');

-- 插入示例食材数据
INSERT INTO ingredients (name, price, unit) VALUES
('猪肉', 25.00, '斤'),
('牛肉', 45.00, '斤'),
('鸡肉', 15.00, '斤'),
('鸡蛋', 1.50, '个'),
('土豆', 3.00, '斤'),
('胡萝卜', 2.50, '斤'),
('白菜', 2.00, '斤'),
('青椒', 4.00, '斤'),
('大蒜', 8.00, '斤'),
('生姜', 12.00, '斤'),
('生抽', 15.00, '瓶'),
('老抽', 18.00, '瓶'),
('料酒', 12.00, '瓶'),
('盐', 3.00, '袋'),
('糖', 5.00, '斤'),
('醋', 8.00, '瓶'),
('辣椒', 10.00, '斤'),
('花椒', 25.00, '斤');

-- 插入示例菜品数据
INSERT INTO dishes (name, description, price, cooking_link, category_id) VALUES
('麻婆豆腐', '经典川菜，麻辣鲜香', 28.00, 'https://example.com/recipe1', 1),
('宫保鸡丁', '经典川菜，酸甜微辣', 32.00, 'https://example.com/recipe2', 1),
('白切鸡', '经典粤菜，清淡鲜美', 35.00, 'https://example.com/recipe3', 2),
('红烧肉', '经典家常菜，肥而不腻', 38.00, 'https://example.com/recipe4', 4),
('糖醋里脊', '经典鲁菜，酸甜可口', 30.00, 'https://example.com/recipe5', 4),
('清蒸鲈鱼', '经典粤菜，鲜嫩爽滑', 45.00, 'https://example.com/recipe6', 2),
('剁椒鱼头', '经典湘菜，香辣鲜美', 42.00, 'https://example.com/recipe7', 3),
('东坡肉', '经典苏菜，肥而不腻', 40.00, 'https://example.com/recipe8', 5);

-- 插入菜品食材关联数据
INSERT INTO dish_ingredients (dish_id, ingredient_id, quantity) VALUES
(1, 1, 0.5),   -- 麻婆豆腐：猪肉0.5斤
(1, 11, 0.1),  -- 麻婆豆腐：生抽0.1瓶
(1, 14, 0.05), -- 麻婆豆腐：盐0.05袋
(1, 17, 0.1),  -- 麻婆豆腐：辣椒0.1斤
(1, 18, 0.05), -- 麻婆豆腐：花椒0.05斤
(2, 3, 0.8),   -- 宫保鸡丁：鸡肉0.8斤
(2, 8, 0.3),   -- 宫保鸡丁：青椒0.3斤
(2, 11, 0.1),  -- 宫保鸡丁：生抽0.1瓶
(2, 15, 0.1),  -- 宫保鸡丁：糖0.1斤
(2, 16, 0.05), -- 宫保鸡丁：醋0.05瓶
(3, 3, 1.2),   -- 白切鸡：鸡肉1.2斤
(3, 11, 0.05), -- 白切鸡：生抽0.05瓶
(3, 14, 0.02), -- 白切鸡：盐0.02袋
(4, 1, 1.0),   -- 红烧肉：猪肉1.0斤
(4, 11, 0.1),  -- 红烧肉：生抽0.1瓶
(4, 12, 0.05), -- 红烧肉：老抽0.05瓶
(4, 13, 0.05), -- 红烧肉：料酒0.05瓶
(4, 14, 0.05), -- 红烧肉：盐0.05袋
(4, 15, 0.1);  -- 红烧肉：糖0.1斤 