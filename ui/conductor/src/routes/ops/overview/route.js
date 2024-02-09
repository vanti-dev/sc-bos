import {useUiConfigStore} from '@/stores/ui-config';
import {findActiveItem} from '@/util/router.js';

export default [
  {
    path: 'overview',
    component: () => import('@/routes/ops/OpsHome.vue'),
    redirect: 'overview/building',
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
      }
    },
    children: [
      {
        name: 'building-overview',
        path: 'building', // Directly map to the building path
        component: () => import('@/routes/ops/overview/pages/BuildingOverview.vue'),
        meta: {
          authentication: {
            rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
          }
        }
      },
      {
        name: 'dynamic-areas-overview',
        path: 'building/:pathMatch(.+)*', // Captures all segments after /building/
        component: () => import('@/routes/ops/overview/pages/DynamicAreasOverview.vue'),
        // Splits segments into an array and passes it as a prop so the component can use it to find the active item
        props: route => ({pathSegments: modifiedPath(route.params.pathMatch).split('/')}),
        meta: {
          authentication: {
            rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
          }
        }
      }
    ],
    beforeEnter: async (to, from, next) => {
      const uiConfig = useUiConfigStore();
      await uiConfig.loadConfig();

      // Get the building children from the app config
      /**
       * @typedef {Object} EnvironmentTrait
       * @property {boolean|string} indoor - Flag or specific value for indoor environment.
       * @property {boolean|string} outdoor - Flag or specific value for outdoor environment.
       */

      /**
       * @typedef {Object} OverviewChild
       * @property {boolean} disabled - Indicates if the item is disabled.
       * @property {string} icon - Icon identifier, e.g., 'mdi-select-all'.
       * @property {string} shortTitle - A short title for the item - for the mini sized navigation.
       * @property {string} title - The full title of the item.
       * @property {Object} traits - Object containing various trait flags.
       * @property {boolean} widgets.showAirQuality - Flag to show air quality.
       * @property {boolean} widgets.showEmergencyLighting - Flag to show emergency lighting.
       * @property {boolean} widgets.showEnergyConsumption - Flag to show energy consumption.
       * @property {boolean|EnvironmentTrait} widgets.showEnvironment - Flag to show environment,
       *                            or an object detailing indoor and outdoor environment widgets.
       * @property {boolean} widgets.showNotifications - Flag to show notifications.
       * @property {boolean|string} widgets.showOccupancy -
       Flag to show occupancy, can be a string for specific occupancy.
       * @property {boolean} widgets.showPower - Flag to show power.
       * @property {OverviewChild[]} [children] - Optional array of children, each following the same structure.
       */
      const overviewChildren = uiConfig.config?.ops?.overview?.children || [];

      // Split the modified path into segments and remove empty segments then return an array of the segments
      const currentPathSegments = modifiedPath(to.path).split('/').filter(segment => segment);

      // Find active item based on the current path segments
      const activeItem = findActiveItem(overviewChildren, currentPathSegments);

      // If the path is '/ops/overview/building' and there is no active item, redirect to the building overview
      const overviewPath = to.path === basePath;

      if (!overviewPath && !activeItem) {
        next(basePath);
      } else {
        next();
      }
    }
  }
];

// Remove the specific beginning '/ops/overview/building' from the path, if it exists
const basePath = '/ops/overview/building';
const modifiedPath = (path) => path.startsWith(basePath) ? path.slice(basePath.length) : path;
