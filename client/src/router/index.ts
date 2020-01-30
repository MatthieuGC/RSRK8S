import Vue from 'vue'
import Router from 'vue-router'
import Platforms from '@/components/Platform/index.vue'
import Tasks from '@/components/Task/index.vue'
import NewJob from '@/components/Job/new.vue'
import Index from '@/components/index.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Index,
    },
    {
      path: '/platforms',
      name: 'platforms',
      component: Platforms,
    },
    {
      path: '/:platform/tasks',
      name: 'tasks',
      component: Tasks,
    },
    {
      path: '/:platform/:task/jobs/new',
      name: 'new_job',
      component: NewJob,
    },
  ],
})
