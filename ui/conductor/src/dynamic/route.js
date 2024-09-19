import {useSidebarStore} from '@/stores/sidebar.js';
import {serviceName} from '@/util/gateway.js';
import {computed, toValue} from 'vue';
import {useRoute} from 'vue-router';

/**
 * Converts a plural word to a singular word.
 * It's not very clever, so only use it for words that end in 's' or 'es'.
 *
 * @param {string} name
 * @return {string}
 */
function toSingular(name) {
  return name.replace(/s$/, '');
}

/**
 * Capitalise the first letter of each word.
 *
 * @param {string} name
 * @return {string}
 */
function toTitleCase(name) {
  return name.replace(/\b./, str => str.toUpperCase());
}

/**
 * Create a route that can show a single service page of the given category.
 *
 * @param {string} category - one of ServiceNames
 * @param {string?} pathPrefix
 * @return {import('vue-router').RouteRecordRaw}
 */
export function useServiceRoute(category, pathPrefix = '/') {
  const single = toSingular(category);
  return {
    name: single,
    path: pathPrefix + single,
    children: useServiceRoutes(category),
    meta: {
      authentication: {
        rolesRequired: ['superAdmin', 'admin', 'commissioner', 'operator', 'viewer']
      },
      title: toTitleCase(single)
    }
  };
}

/**
 * Create routes that can show a single service page of the given category.
 *
 * @param {string} category - one of ServiceNames
 * @return {import('vue-router').RouteRecordRaw[]}
 */
export function useServiceRoutes(category) {
  return [{
    name: toSingular(category) + '-name-id',
    path: ':name/:id',
    component: () => import('@/components/pages/ServiceJsonEditor.vue'),
    props: route => {
      return {
        name: serviceName(route.params.name, category),
        id: route.params.id
      };
    }
  }, {
    name: toSingular(category) + '-id',
    path: ':id',
    component: () => import('@/components/pages/ServiceJsonEditor.vue'),
    props: route => {
      return {
        name: category,
        id: route.params.id
      };
    }
  }];
}

/**
 * Create a link (suitable for `to` in a `router-link`) to a service page that will show the given service.
 *
 * @param {MaybeRefOrGetter<string | undefined>} category - One of ServiceNames
 * @param {MaybeRefOrGetter<string | undefined>} name - The name of the device hosting the service
 * @param {MaybeRefOrGetter<string | undefined>} id - The id of the service
 * @return {{hasLink: Ref<boolean>, to: Ref<undefined|import('vue-router').RouteLocationRaw>}}
 */
export function useServiceRouterLink(category, name, id) {
  const hasLink = computed(() => Boolean(toValue(category) && toValue(id)));
  const to = computed(() => {
    if (!hasLink.value) return undefined;
    if (toValue(name)) {
      return {
        name: toSingular(toValue(category)) + '-name-id',
        params: {name: toValue(name), id: toValue(id)}
      };
    } else {
      return {
        name: toSingular(toValue(category)) + '-id',
        params: {id: toValue(id)}
      };
    }
  });
  return {
    hasLink,
    to
  };
}

/**
 * Like useServiceRouteLink, but uses the context from the current route and sidebar store.
 *
 * @return {{hasLink: Ref<boolean>, to: Ref<import('vue-router').RouteLocationRaw|undefined>}}
 */
export function useSidebarServiceRouterLink() {
  const route = useRoute();
  const sidebar = useSidebarStore();
  return useServiceRouterLink(() => route.meta?.editRoutePrefix, () => sidebar.data?.nodeName, () => sidebar.data?.service?.id);
}
