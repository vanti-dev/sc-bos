import ChildOnlyPage from '@/components/page-layout/ChildOnlyPage.vue';
import thirdParty from '@/routes/auth/third-party/route.js';
import {route} from '@/util/router.js';
import {useAppConfigStore} from '@/stores/app-config';

export default [
  {
    name: 'auth',
    path: '/auth',
    components: {
      default: ChildOnlyPage,
      nav: () => import('./AuthNav.vue')
    },
    children: [
      {
        path: 'users',
        component: () => import('./users/Users.vue'),
        meta: {
          authentication: {
            rolesRequired: ['superAdmin', 'admin', 'viewer']
          },
          title: 'Users'
        }
      }
    ],
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'operator', 'viewer']
      },
      title: 'Auth'
    },
    beforeEnter: async (to, from, next) => {
      const appConfig = useAppConfigStore();
      if (to.path === '/auth') {
        if (appConfig.pathEnabled('/auth/users')) {
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
