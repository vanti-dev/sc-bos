/**
 * @param {import('vue-router').Route} route
 * @return {string|undefined}
 */
export function routeTitle(route) {
  for (let i = route.matched.length - 1; i >= 0; i--) {
    const r = route.matched[i];
    let title = r.meta?.['title'];
    if (typeof title === 'function') title = title();
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
  return [route];
}

/**
 * Finds active items in the building structure based on a list of titles.
 * It iterates through the provided titles, searching for corresponding items in the building structure.
 * If an item matching a title is found, it's added to the result. The search continues into the children of
 * the found item for the next title. If any title in the sequence is not found, the search stops and returns
 * the items found up to that point.
 *
 * If all titles are found, the last item in the result is the active item and is returned.
 * If any title is not found, the active item is null and is returned.
 *
 * Example:
 * const overviewChildren = [
 *    {
 *      title: 'Area 1',
 *      children: [
 *        {
 *          title: 'Area 2',
 *          children: [
 *            {
 *              title: 'Area 3'
 *            }
 *          ]
 *        }
 *      ]
 *    }
 *  ];
 *
 * // Returns the item with title 'Area 3'
 * findActiveItem(overviewChildren, ['Area 1', 'Area 2', 'Area 3']);
 *
 * // Returns the item with title 'Area 2'
 * findActiveItem(overviewChildren, ['Area 1', 'Area 2']);
 *
 * // Returns null
 * findActiveItem(overviewChildren, ['Area 1', 'Area 2', 'Area 4']);
 *
 * @param {OverviewChild[]} children - Array of building children, each following the structure of OverviewChild.
 * @param {string[]} pathSegments - Array of titles to find in sequence.
 * @return {OverviewChild[]|null} - Array of found items in the order of the titles provided. If a title is not found,
 *                             the array includes items up to the last found title.
 */
export const findActiveItem = (children, pathSegments) => {
  let currentItems = children;
  const result = [];

  // Iterate through the titles, searching for the corresponding item in the current items
  for (const segment of pathSegments) {
    const _segment = decodeURIComponent(segment);
    const foundItem = currentItems.find(
        item => _segment === item.path || _segment === item.title);

    if (!foundItem) break;
    result.push(foundItem);
    currentItems = foundItem.children || [];
  }

  // If the result array is the same length as the pathSegments array, then all titles were found
  // and the last item in the result array is the active item,
  // otherwise the active item is null
  return result.length === pathSegments.length ? result[result.length - 1] : null;
};
