import {ServiceNames} from '@/api/ui/services.js';
import SidebarPage from '@/components/pages/SidebarPage.vue';
import {serviceName} from '@/util/gateway.js';

export default [{
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
}, {
  name: 'automation',
  path: '/automation',
  children: [{
    name: 'automation-name-id',
    path: ':name/:id',
    component: () => import('@/components/pages/ServiceJsonEditor.vue'),
    props: route => {
      return {
        name: serviceName(route.params.name, ServiceNames.Automations),
        id: route.params.id
      };
    }
  }, {
    name: 'automation-id',
    path: ':id',
    component: () => import('@/components/pages/ServiceJsonEditor.vue'),
    props: route => {
      return {
        name: ServiceNames.Automations,
        id: route.params.id
      };
    }
  }],
  meta: {
    authentication: {
      rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
    },
    title: 'Automation'
  }
}];
