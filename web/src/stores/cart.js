import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useCartStore = defineStore('cart', () => {
  // 购物车中的菜品
  const cartItems = ref([])

  // 计算属性
  const totalItems = computed(() => {
    return cartItems.value.length
  })

  const totalPrice = computed(() => {
    return cartItems.value.reduce((sum, item) => sum + item.price, 0)
  })

  const dishIds = computed(() => {
    return cartItems.value.map(item => item.id)
  })

  // 添加菜品到购物车
  const addToCart = (dish) => {
    const existingItem = cartItems.value.find(item => item.id === dish.id)
    if (!existingItem) {
      cartItems.value.push({
        id: dish.id,
        name: dish.name,
        price: dish.price,
        image_url: dish.image_url,
        description: dish.description,
        category: dish.category
      })
    }
  }

  // 从购物车移除菜品
  const removeFromCart = (dishId) => {
    const index = cartItems.value.findIndex(item => item.id === dishId)
    if (index > -1) {
      cartItems.value.splice(index, 1)
    }
  }

  // 清空购物车
  const clearCart = () => {
    cartItems.value = []
  }

  // 检查菜品是否在购物车中
  const isInCart = (dishId) => {
    return cartItems.value.some(item => item.id === dishId)
  }

  return {
    cartItems,
    totalItems,
    totalPrice,
    dishIds,
    addToCart,
    removeFromCart,
    clearCart,
    isInCart
  }
}) 