import {ServiceNames} from '@/api/ui/services.js';
import SidebarPage from '@/components/pages/SidebarPage.vue';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {serviceName} from '@/util/gateway.js';

export default [{
  name: 'system',
  path: '/system',
  components: {
    default: SidebarPage,
    nav: () => import('./SystemNav.vue')
  },
  children: [
    {
      path: 'drivers',
      components: {
        default: () => import('./components/pages/DriversList.vue'),
        sidebar: () => import('./components/ServicesSideBar.vue')
      },
      meta: {
        editRoutePrefix: 'driver',
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    },
    {
      path: 'features',
      components: {
        default: () => import('./components/pages/FeaturesList.vue'),
        sidebar: () => import('./components/ServicesSideBar.vue')
      },
      meta: {
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    },
    {
      path: 'components',
      components: {
        default: () => import('./components/pages/ComponentsList.vue')
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
    title: 'System'
  },
  beforeEnter: async (to, from, next) => {
    const uiConfig = useUiConfigStore();
    if (to.path === '/system') {
      if (uiConfig.pathEnabled('/system/drivers')) {
        next('/system/drivers');
      } else if (uiConfig.pathEnabled('/system/features')) {
        next('/system/features');
      }
    } else {
      next();
    }
  }
}, {
  name: 'driver',
  path: '/system/driver',
  children: [{
    name: 'driver-name-id',
    path: ':name/:id',
    component: () => import('@/components/pages/ServiceJsonEditor.vue'),
    props: route => {
      return {
        name: serviceName(route.params.name, ServiceNames.Drivers),
        id: route.params.id
      };
    }
  }, {
    name: 'driver-id',
    path: ':id',
    component: () => import('@/components/pages/ServiceJsonEditor.vue'),
    props: route => {
      return {
        name: ServiceNames.Drivers,
        id: route.params.id
      };
    }
  }],
  meta: {
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    },
    title: 'Driver'
  }
}
];
