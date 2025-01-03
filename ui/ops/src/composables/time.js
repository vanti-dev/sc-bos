import {computed, onMounted, onUnmounted, ref, toValue} from 'vue';

/**
 * Returns values that can show time since a given time.
 *
 * @example
 * const {showTimeSince, timeSinceStr} = useTimeSince(updateDate);
 * const tooltipStr = computed(() => {
 *   if (showTimeSince.value) {
 *     return `${stateStr.value} for ${timeSinceStr.value}`;
 *   } else {
 *     return stateStr.value;
 *   }
 * })
 * // tooltipStr output will be something like "Active for 5m"
 *
 * @param {Date} date
 * @return {{
 *   now: Ref<Date>,
 *   showTimeSince: ComputedRef<boolean>,
 *   timeSinceStr: ComputedRef<string>
 * }}
 */
export function useTimeSince(date) {
  const nowHandle = ref(0);
  const now = ref(new Date());
  onUnmounted(() => clearInterval(nowHandle.value));
  onMounted(() => {
    nowHandle.value = setInterval(() => {
      now.value = new Date();
    }, 1000);
  });

  const millisSince = computed(() => {
    return now.value.getTime() - toValue(date)?.getTime();
  });
  const showTimeSince = computed(() => {
    return Boolean(toValue(date) && millisSince.value > 1000);
  });
  const timeSinceStr = computed(() => {
    if (!showTimeSince.value) return '';
    const t = millisSince.value;
    if (t > 1000 * 60 * 60 * 24) {
      const h = Math.floor(t / (1000 * 60 * 60 * 24));
      return `${h}d`;
    } else if (t > 1000 * 60 * 60) {
      const h = Math.floor(t / (1000 * 60 * 60));
      return `${h}h`;
    } else if (t > 1000 * 60) {
      const m = Math.floor(t / (1000 * 60));
      return `${m}m`;
    } else if (t > 1000) {
      const s = Math.floor(t / 1000);
      return `${s}s`;
    } else {
      return '';
    }
  });

  return {
    now,
    showTimeSince,
    timeSinceStr
  };
}
