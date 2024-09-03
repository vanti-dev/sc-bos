import SidebarPage from '@/components/page-layout/SidebarPage.vue';

export default {
  name: 'automations',
  path: '/automations',
  redirect: '/automations/all',
  components: {
    default: SidebarPage,
    nav: () => import('./AutomationNav.vue')
  },
  children: [
    {
      path: ':type',
      components: {
        default: () => import('./components/AutomationList.vue'),
        sidebar: () => import('./components/AutomationSideBar.vue')
      },
      props: {
        default: true
      },
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    }
  ],
  meta: {
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    },
    title: 'Automations'
  }
};
