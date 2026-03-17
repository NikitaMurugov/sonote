import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/auth/LoginView.vue'),
      meta: { guest: true },
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/auth/RegisterView.vue'),
      meta: { guest: true },
    },
    {
      path: '/',
      component: () => import('@/components/layout/AppShell.vue'),
      meta: { auth: true },
      children: [
        {
          path: '',
          redirect: '/w',
        },
        {
          path: 'w',
          name: 'workspace-redirect',
          component: () => import('@/views/app/WorkspaceRedirect.vue'),
        },
        {
          path: 'w/:wsSlug',
          name: 'dashboard',
          component: () => import('@/views/app/DashboardView.vue'),
        },
        {
          path: 'w/:wsSlug/note/:noteId',
          name: 'note-editor',
          component: () => import('@/views/app/NoteEditorView.vue'),
        },
        {
          path: 'w/:wsSlug/graph',
          name: 'graph',
          component: () => import('@/views/app/GraphView.vue'),
        },
        {
          path: 'w/:wsSlug/search',
          name: 'search',
          component: () => import('@/views/app/SearchView.vue'),
        },
        {
          path: 'w/:wsSlug/settings',
          name: 'settings',
          component: () => import('@/views/app/SettingsView.vue'),
        },
        {
          path: 'profile',
          name: 'profile',
          component: () => import('@/views/app/ProfileView.vue'),
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/views/NotFoundView.vue'),
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.meta.auth && !authStore.isAuthenticated) {
    next('/login')
  } else if (to.meta.guest && authStore.isAuthenticated) {
    next('/')
  } else {
    next()
  }
})

export default router
