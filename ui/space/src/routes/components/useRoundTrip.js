import {SECOND, useNow} from '@/components/now';
import {computed, ref, toValue, watch} from 'vue';

/**
 * Represents a value that we both read and write to a server.
 * In these cases there can be some update loops or other jittery behavior caused by a combination of
 * local+remote+push event state.
 * This composable hopes to resolve these issues.
 *
 * @param {MaybeRefOrGetter<T>} [remoteValue]
 * @param {MaybeRefOrGetter<T>} [localValue]
 * @param {MaybeRefOrGetter<number>} [age]
 * @return {{
 *   value: import('vue').ComputedRef<T>,
 *   remoteUpdateTime: import('vue').Ref<Date>,
 *   localUpdateTime: import('vue').Ref<Date>,
 *   useLocalValue: import('vue').ComputedRef<boolean>,
 *   remoteValue: import('vue').Ref<T>,
 *   localValue: import('vue').Ref<T>
 * }}
 * @template T
 */
export const useRoundTrip = (remoteValue = null, localValue = null, age = 15 * SECOND) => {
  remoteValue = remoteValue ? ref(remoteValue) : ref(null);
  localValue = localValue ? ref(localValue) : ref(null);
  const remoteUpdateTime = ref(/** @type {Date} */null);
  const localUpdateTime = ref(/** @type {Date} */null);
  const {now} = useNow(age);

  watch(() => toValue(remoteValue), (n, o) => {
    if (o !== n) {
      remoteUpdateTime.value = new Date();
    }
  });
  watch(() => toValue(localValue), (n, o) => {
    if (o !== n) {
      localUpdateTime.value = new Date();
    }
  });

  const useLocalValue = computed(() => {
    if (!remoteUpdateTime.value) {
      return true;
    }
    if (!localUpdateTime.value) {
      return false;
    }
    const localTime = localUpdateTime.value.getTime();
    const remoteTime = remoteUpdateTime.value.getTime();
    if (localTime > remoteTime) {
      return true;
    }
    if (remoteTime - localTime > toValue(age)) {
      return false; // we got a response from the server significantly newer than the local value
    }
    const time = now.value.getTime();
    // if enough time has passed, we should use the remote value
    return time - localTime < toValue(age);
  });
  const value = computed(() => {
    if (useLocalValue.value) {
      return toValue(localValue);
    } else {
      return toValue(remoteValue);
    }
  });

  return {
    value,
    remoteUpdateTime,
    localUpdateTime,
    useLocalValue,
    remoteValue,
    localValue
  };
};
