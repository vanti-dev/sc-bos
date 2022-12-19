import {onBeforeUnmount, onMounted, ref} from 'vue';

export const MILLISECOND = 1;
export const SECOND = 1000 * MILLISECOND;
export const MINUTE = 60 * SECOND;
export const HOUR = 60 * MINUTE;
export const DAY = 24 * HOUR;

/**
 *
 * @param resolution
 */
export function useNow(resolution = MINUTE) {
  const now = ref(new Date());

  /**
   * @param {Date} t
   */
  function nextDelay(t) {
    const ms = t.getTime();
    // note: if ms is exactly on a resolution boundary then instead of returning 0 we should wait a full hop
    return (ms % resolution) || resolution;
  }

  let handle = 0;

  /**
   *
   * @param t
   */
  function updateNowWhenNeeded(t) {
    const delay = nextDelay(t);
    clearTimeout(handle);
    handle = setTimeout(() => {
      now.value = new Date();
      updateNowWhenNeeded(now.value);
    }, delay);
  }

  onMounted(() => {
    updateNowWhenNeeded(now.value);
  });
  onBeforeUnmount(() => {
    clearTimeout(handle);
  });

  return {now};
}
