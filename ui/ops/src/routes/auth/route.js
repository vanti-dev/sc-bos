import SidebarPage from '@/components/pages/SidebarPage.vue';
import accounts from '@/routes/auth/accounts/route.js';
import roles from '@/routes/auth/roles/route.js';
import thirdParty from '@/routes/auth/third-party/route.js';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {route} from '@/util/router.js';

export default [
  {
    name: 'auth',
    path: '/auth',
    components: {
      default: SidebarPage,
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
      },
      ...route(thirdParty),
      ...route(accounts),
      ...route(roles),
    ],
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'operator', 'viewer']
      },
      title: 'Access Management'
    },
    beforeEnter: async (to, from, next) => {
      const uiConfig = useUiConfigStore();
      if (to.path === '/auth') {
        if (uiConfig.pathEnabled('/auth/users')) {
          next('/auth/users');
        } else {
          next('/auth/third-party');
        }
      } else {
        next();
      }
    }
  }
];
