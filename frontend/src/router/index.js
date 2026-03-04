import { createRouter, createWebHashHistory } from 'vue-router'
import Test from '../views/Test.vue'
import WizardLayout from '../views/wizard/Layout.vue'
import ConnectDevice from '../views/wizard/ConnectDevice.vue'
import Login from '../views/wizard/Login.vue'
import AutoMount from '../views/wizard/AutoMount.vue'
import LoginSuccess from '../views/wizard/LoginSuccess.vue'
import MainLayout from '../views/main/Layout.vue'
import Dashboard from '../views/main/Dashboard.vue'
import Storage from '../views/main/Storage.vue'
import ZeroTier from '../views/main/ZeroTier.vue'
import Settings from '../views/main/Settings.vue'

const routes = [
  {
    path: '/test',
    component: Test
  },
  {
    path: '/',
    redirect: '/wizard'
  },
  {
    path: '/wizard',
    component: WizardLayout,
    children: [
      {
        path: '',
        redirect: '/wizard/connect'
      },
      {
        path: 'connect',
        name: 'WizardConnect',
        component: ConnectDevice
      },
      {
        path: 'login',
        name: 'WizardLogin',
        component: Login
      },
      {
        path: 'auto-mount',
        name: 'WizardAutoMount',
        component: AutoMount
      },
      {
        path: 'success',
        name: 'WizardSuccess',
        component: LoginSuccess
      }
    ]
  },
  {
    path: '/main',
    component: MainLayout,
    children: [
      {
        path: '',
        redirect: '/main/dashboard'
      },
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: Dashboard
      },
      {
        path: 'storage',
        name: 'Storage',
        component: Storage
      },
      {
        path: 'zerotier',
        name: 'ZeroTier',
        component: ZeroTier
      },
      {
        path: 'settings',
        name: 'Settings',
        component: Settings
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

// 路由守卫：检查登录状态
router.beforeEach((to, from, next) => {
  const isLoggedIn = localStorage.getItem('rnctl_logged_in') === 'true'

  if (to.path.startsWith('/main') && !isLoggedIn) {
    // 未登录访问主界面，重定向到向导
    next('/wizard/connect')
  } else if (to.path.startsWith('/wizard') && isLoggedIn) {
    // 已登录访问向导，重定向到主界面
    next('/main/dashboard')
  } else {
    next()
  }
})

export default router
