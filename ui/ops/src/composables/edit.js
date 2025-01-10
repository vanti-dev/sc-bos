import deepEqual from 'fast-deep-equal';
import {computed, ref, watch} from 'vue';

/**
 * @template T
 * @param {Readonly<import('vue').Ref<T>>} remote
 * @param {import('vue').Ref<T>} local - must be writable
 * @return {{
 *   remoteIsChanged: Ref<boolean>,
 *   sync: () => void,
 *   localIsChanged: Ref<boolean>
 * }}
 */
export function useOfflineEdit(remote, local) {
  const remoteIsChanged = ref(false);
  const sync = () => {
    local.value = remote.value;
    remoteIsChanged.value = false;
  };
  const localIsChanged = computed(() => !deepEqual(local.value, remote.value));

  watch(remote, (newValue, oldValue) => {
    if (deepEqual(newValue, oldValue)) return;
    if (!oldValue || !local.value || deepEqual(oldValue, local.value)) {
      local.value = newValue;
      return;
    }
    remoteIsChanged.value = !deepEqual(newValue, local.value);
  }, {immediate: true, deep: true});

  return {
    remoteIsChanged,
    sync,
    localIsChanged
  };
}
