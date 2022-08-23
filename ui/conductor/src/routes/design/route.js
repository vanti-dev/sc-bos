import ChildOnlyPage from '@/components/ChildOnlyPage.vue';

export default {
  name: 'design',
  path: '/design',
  redirect: '/design/buildings',
  components: {
    default: ChildOnlyPage,
    sections: () => import('./DesignSections.vue')
  },
  children: [
    {path: 'buildings', component: () => import('./DesignBuildings.vue')},
    {path: 'spaces', component: () => import('./DesignSpaces.vue')},
    {path: 'plan', component: () => import('./DesignPlan.vue')},
  ],
  meta: {
    title: 'Designer'
  }
}
