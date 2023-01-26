import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';

export default {
  name: 'automations',
  path: '/automations',
  components: {
    default: ChildOnlyPage,
    nav: () => import('./AutomationNav.vue')
  },
  children: [
    {
      path: ':type?',
      components: {
        default: () => import('./AutomationConfig.vue')
      },
      props: {
        default: true
      }
    }
  ],
  meta: {
    title: 'Automations'
  }
};
