// See https://github.com/vuejs/vue-router/issues/3760#issuecomment-1191774443

// if using 2.6.x import from @vue/composition-api
import {computed, getCurrentInstance, reactive} from 'vue'

export function useRouter() {
  return getCurrentInstance().proxy.$router
}

export function useRoute() {
  const currentRoute = computed(() => getCurrentInstance().proxy.$route)

  const protoRoute = /** @type {import('vue-router').Route} */ Object.keys(currentRoute.value).reduce(
      (acc, key) => {
        acc[key] = computed(() => currentRoute.value[key])
        return acc
      },
      {}
  )

  return reactive(protoRoute)
}

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
