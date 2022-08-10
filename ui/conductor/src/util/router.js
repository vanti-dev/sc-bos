// See https://github.com/vuejs/vue-router/issues/3760#issuecomment-1191774443

// if using 2.6.x import from @vue/composition-api
import {computed, getCurrentInstance, reactive} from 'vue'

export function useRouter() {
  return getCurrentInstance().proxy.$router
}

export function useRoute() {
  const currentRoute = computed(() => getCurrentInstance().proxy.$route)

  const protoRoute = Object.keys(currentRoute.value).reduce(
      (acc, key) => {
        acc[key] = computed(() => currentRoute.value[key])
        return acc
      },
      {}
  )

  return reactive(protoRoute)
}
