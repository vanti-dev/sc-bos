/**
 * @param {import('vue-router').Route} route
 * @returns {string|undefined}
 */
export function routeTitle(route) {
  for (let i = route.matched.length - 1; i >= 0; i--) {
    const r = route.matched[i];
    const title = r.meta?.['title'];
    if (title) return title;
  }
}

/**
 * @param {import('vue-router').RouteConfig | import('vue-router').RouteConfig[]} route
 * @return {import('vue-router').RouteConfig[]}
 */
export function route(route) {
  if (Array.isArray(route)) {
    return route;
  }
  return [/** @type {import('vue-router').RouteConfig}*/ route];
}
