import {ServiceNames} from '@/api/ui/services.js';
import SidebarPage from '@/components/pages/SidebarPage.vue';
import {useServiceRoute} from '@/dynamic/route.js';

export default [
  {
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
          editRoutePrefix: ServiceNames.Automations,
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
  },
  useServiceRoute(ServiceNames.Automations, undefined, {name: 'automations'})];
