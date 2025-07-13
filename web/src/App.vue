<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
// 主应用组件
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
  background-color: #f5f5f5;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

#app {
  min-height: 100vh;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.card-shadow {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.text-center {
  text-align: center;
}

.mb-20 {
  margin-bottom: 20px;
}

.mt-20 {
  margin-top: 20px;
}

.flex {
  display: flex;
}

.flex-center {
  display: flex;
  align-items: center;
  justify-content: center;
}

.flex-between {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* 移动端适配 */
@media (max-width: 768px) {
  .page-container {
    padding: 10px;
  }
  
  .mb-20 {
    margin-bottom: 15px;
  }
  
  .mt-20 {
    margin-top: 15px;
  }
  
  /* 移动端卡片样式 */
  .mobile-card {
    border-radius: 12px;
    margin-bottom: 15px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  /* 移动端按钮样式 */
  .mobile-btn {
    height: 44px;
    border-radius: 22px;
    font-size: 16px;
  }
  
  /* 移动端输入框样式 */
  .mobile-input {
    height: 44px;
    border-radius: 22px;
  }
  
  /* 移动端导航栏 */
  .mobile-header {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    z-index: 1000;
    background: white;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }
  
  /* 移动端底部导航 */
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
  
  /* 移动端内容区域 */
  .mobile-content {
    padding-top: 60px;
    padding-bottom: 80px;
  }
  
  /* 移动端网格布局 */
  .mobile-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 10px;
    padding: 10px;
  }
  
  /* 移动端瀑布流布局 */
  .mobile-waterfall {
    column-count: 2;
    column-gap: 10px;
    padding: 10px;
  }
  
  .mobile-waterfall-item {
    break-inside: avoid;
    margin-bottom: 10px;
  }
}

/* 超小屏幕适配 */
@media (max-width: 480px) {
  .mobile-grid {
    grid-template-columns: 1fr;
  }
  
  .mobile-waterfall {
    column-count: 1;
  }
  
  .page-container {
    padding: 5px;
  }
}

/* 触摸优化 */
@media (hover: none) and (pointer: coarse) {
  .touch-target {
    min-height: 44px;
    min-width: 44px;
  }
  
  .touch-spacing {
    margin: 8px 0;
  }
}
</style> 