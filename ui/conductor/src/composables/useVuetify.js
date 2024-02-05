import {getCurrentInstance} from 'vue';

/**
 * Get the Vuetify instance from the current Vue instance
 *
 * @return {import('vuetify').Vuetify} Framework
 */
export default function() {
  const vm = getCurrentInstance();

  return vm?.proxy?.$vuetify;
}
