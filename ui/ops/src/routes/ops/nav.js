import {useUiConfigStore} from '@/stores/uiConfig.js';
import {computed} from 'vue';

/**
 * @typedef {Object} NavItem
 * @property {string} title - The title of the menu item
 * @property {string} icon - The icon to display for the menu item
 * @property {{path: string}} link - The route link for the menu item
 * @property {string} [badgeType] - Optional badge type for the menu item
 * @property {function} [enabled] - Optional function to determine if the item is enabled
 */

/**
 * Menu Items
 * This is the main list of items
 *
 * @type {NavItem[]}
 */
export const navItems = [
  {
    title: 'Notifications',
    icon: 'mdi-bell-outline',
    link: {path: '/ops/notifications'},
    badgeType: 'unacknowledgedAlertCount'
  },
  {
    title: 'Air Quality',
    icon: 'mdi-air-filter',
    link: {path: '/ops/air-quality'},
    badgeType: null
  },
  {
    title: 'Emergency Lighting',
    icon: 'mdi-alarm-light-outline',
    link: {path: '/ops/emergency-lighting'},
    badgeType: null
  },
  {
    title: 'Security',
    icon: 'mdi-shield-key',
    link: {path: '/ops/security'},
    badgeType: null
  },
  {
    title: 'Security Events',
    icon: 'mdi-shield-alert',
    link: {path: '/ops/security-events'},
    badgeType: null
  },
  {
    title: 'Waste Records',
    icon: 'mdi-recycle',
    link: {path: '/ops/waste'},
    badgeType: null,
  }
];

/**
 * Check if a route is enabled.
 *
 * @param {NavItem} item
 * @return {boolean}
 */
export function isRouteEnabled(item) {
  const uiConfig = useUiConfigStore();
  if (!uiConfig.pathEnabled(item.link.path)) return false;
  if (typeof item.enabled === 'function') return item.enabled();
  return true;
}

/**
 * Get a computed reference to the enabled navigation items.
 *
 * @return {ComputedRef<NavItem[]>}
 */
export function useEnabledNavItems() {
  return computed(() => navItems.filter(item => isRouteEnabled(item)));
}
