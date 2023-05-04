import { createRouter, createWebHistory } from 'vue-router'
import AddService from '../components/AddService.vue'
import Services from '../components/Services.vue'
import Statistics from '../components/Statistics.vue'
import Settings from '../components/Settings.vue'
import ServiceSettings from '../components/ServiceSettings.vue'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
    },
    {
      path: '/settings',
      name: 'settings',
      component: Settings,
    },
    {
      path: '/statistics',
      name: 'statistics',
      component: Statistics,
    },
    {
      path: '/services',
      name: 'services',
      component: Services,
    },
    {
      path: '/services/new',
      name: 'service.new',
      component: AddService,
    },
    {
      path: '/services/:name',
      name: "service.settings",
      component: ServiceSettings,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue')
    }
  ]
})

export default router
