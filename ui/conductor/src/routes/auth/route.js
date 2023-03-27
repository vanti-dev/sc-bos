import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';
import thirdParty from '@/routes/auth/third-party/route.js';
import {route} from '@/util/router.js';
import {featureEnabled} from '@/routes/config';

export default [
  {
    name: 'auth',
    path: '/auth',
    components: {
      default: ChildOnlyPage,
      nav: () => import('./AuthNav.vue')
    },
    children: [
      {path: 'users', component: () => import('./users/Users.vue')}
    ],
    meta: {
      title: 'Auth'
    },
    beforeEnter: async (to, from, next) => {
      if (to.path === '/auth') {
        if (await featureEnabled('/auth/users')) {
          next('/auth/users');
        } else {
          next('/auth/third-party');
        }
      } else {
        next();
      }
    }
  },
  ...route(thirdParty)
];
