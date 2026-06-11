import { createRouter, createWebHistory } from 'vue-router'

import CategoryList from '../views/CategoryList.vue'
import Dashboard from '../views/Dashboard.vue'
import DishList from '../views/DishList.vue'
import Login from '../views/Login.vue'
import OrderDetail from '../views/OrderDetail.vue'
import OrderList from '../views/OrderList.vue'
import UserWhitelist from '../views/UserWhitelist.vue'

// router 管理端页面路由。
export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: Login,
    },
    {
      path: '/',
      component: Dashboard,
      meta: {
        requiresAuth: true,
      },
      children: [
        {
          path: '',
          redirect: '/categories',
        },
        {
          path: 'categories',
          name: 'categories',
          component: CategoryList,
        },
        {
          path: 'dishes',
          name: 'dishes',
          component: DishList,
        },
        {
          path: 'orders',
          name: 'orders',
          component: OrderList,
        },
        {
          path: 'orders/:id',
          name: 'order-detail',
          component: OrderDetail,
          props: true,
        },
        {
          path: 'users',
          name: 'users',
          component: UserWhitelist,
        },
      ],
    },
  ],
})

// 路由守卫，未登录用户跳转登录页。
router.beforeEach((to) => {
  const token = localStorage.getItem('admin_token')
  if (to.meta.requiresAuth && !token) {
    return {
      path: '/login',
      query: {
        redirect: to.fullPath,
      },
    }
  }
  if (to.path === '/login' && token) {
    return '/'
  }
  return true
})
