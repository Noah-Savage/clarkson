import { createRouter, createWebHistory } from "vue-router"
import Dashboard from "./views/Dashboard.vue"
import VehicleDetail from "./views/VehicleDetail.vue"
import Login from "./views/Login.vue"
import Register from "./views/Register.vue"

const routes = [
  { path: "/login", component: Login },
  { path: "/register", component: Register },
  { path: "/dashboard", component: Dashboard, meta: { requiresAuth: true } },
  { path: "/vehicle/:id", component: VehicleDetail, meta: { requiresAuth: true } },
  { path: "/", redirect: "/dashboard" },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

router.beforeEach((to, from, next) => {
  const token = localStorage.getItem("token")
  if (to.meta.requiresAuth && !token) {
    next("/login")
  } else {
    next()
  }
})

export default router
