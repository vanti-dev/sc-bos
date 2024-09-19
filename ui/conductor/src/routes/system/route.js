import {ServiceNames} from '@/api/ui/services.js';
import SidebarPage from '@/components/pages/SidebarPage.vue';
import {useServiceRoutes} from '@/dynamic/route.js';
import {useUiConfigStore} from '@/stores/uiConfig.js';

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
        editRoutePrefix: ServiceNames.Drivers,
        authentication: {
          rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
        }
      }
    },
    {
      path: 'zones',
      components: {
        default: () => import('./components/pages/ZonesList.vue'),
        sidebar: () => import('./components/ServicesSideBar.vue')
      },
      meta: {
        editRoutePrefix: ServiceNames.Zones,
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
  children: useServiceRoutes(ServiceNames.Drivers),
  meta: {
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    },
    title: 'Driver'
  }
}, {
  name: 'system-zone',
  path: '/system/zone',
  children: useServiceRoutes(ServiceNames.Zones),
  meta: {
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    },
    title: 'Zone'
  }
}
];
