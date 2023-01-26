import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';
import thirdParty from '@/routes/auth/third-party/route.js';
import {route} from '@/util/router.js';

export default [
  {
    name: 'auth',
    path: '/auth',
    redirect: '/auth/users',
    components: {
      default: ChildOnlyPage,
      nav: () => import('./AdminNav.vue')
    },
    children: [
      {path: 'users', component: () => import('./users/Users.vue')}
    ],
    meta: {
      title: 'Auth'
    }
  },
  ...route(thirdParty)
];
