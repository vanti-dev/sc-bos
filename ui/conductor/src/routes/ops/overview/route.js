import {useUiConfigStore} from '@/stores/ui-config';

export default [
  {
    path: 'overview',
    components: {
      default: () => import('@/routes/ops/OpsHome.vue'),
      sidebar: () => import('@/routes/ops/overview/OpsSideBar.vue')
    },
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
      }
    },
    children: [
      {
        name: 'dynamic-areas-overview',
        path: ':pathMatch(.+)*', // Captures all segments after /building/
        component: () => import('@/routes/ops/overview/OpsDashPage.vue'),
        // Splits segments into an array and passes it as a prop so the component can use it to find the active item
        props: route => ({pathSegments: route.params['pathMatch']}),
        meta: {
          authentication: {
            rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
          }
        }
      }
    ],
    beforeEnter: async (to, from, next) => {
      if (to.path !== basePath) {
        return next(); // skip processing if we already have a path to go to
      }
      // else find a leaf page
      const uiConfig = useUiConfigStore();
      await uiConfig.loadConfig();

      const pages = uiConfig.getOrDefault('ops.pages');
      if (!pages) {
        // legacy behaviour:
        return next(to.path + '/building');
      }
      return next(to.path + '/' + encodeURIComponent(pages[0].path ?? pages[0].title));
    }
  }
];

const basePath = '/ops/overview';

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
 * @property {string?} path - A path segment locating this page in the parent, defaults to title.
 * @property {Object} widgets - Object containing various trait flags.
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
