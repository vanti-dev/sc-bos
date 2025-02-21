import {
  addDays,
  addHours,
  addMinutes,
  addMonths,
  addWeeks,
  addYears,
  startOfDay,
  startOfHour,
  startOfMinute,
  startOfMonth,
  startOfWeek,
  startOfYear
} from 'date-fns';
import {computed, onMounted, onScopeDispose, onUnmounted, ref, toValue} from 'vue';

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

/**
 * Returns a ref to a date that is updated to represent the start of a period.
 * The startOf function converts a date to the start of a period.
 * The nextOf function converts a date to an equivalent date in the next period.
 *
 * @example
 * const t = useStartOfPeriod(startOfDay, addDays);
 * // t will always be the start of the current day
 *
 * @param {function(Date):Date} startOf - return a new date that is at the start of the period enclosing the passed Date
 * @param {function(Date, number):Date} nextOf - return a new date that is some number of periods after the passed Date
 * @return {import('vue').Ref<Date>}
 */
export function useStartOfPeriod(startOf, nextOf) {
  const tRef = ref(startOf(new Date()));
  let startTimeHandle = 0;
  const updateStartTime = () => {
    const n = new Date();
    const t = startOf(n);
    tRef.value = t;
    const startOfNext = startOf(nextOf(t, 1));
    startTimeHandle = setTimeout(updateStartTime, startOfNext.getTime() - n.getTime());
  }

  updateStartTime();
  onScopeDispose(() => {
    clearTimeout(startTimeHandle);
  });

  return tRef;
}

export const useStartOf = {
  minute: () => useStartOfPeriod(startOfMinute, addMinutes),
  hour: () => useStartOfPeriod(startOfHour, addHours),
  day: () => useStartOfPeriod(startOfDay, addDays),
  week: () => useStartOfPeriod(startOfWeek, addWeeks),
  month: () => useStartOfPeriod(startOfMonth, addMonths),
  year: () => useStartOfPeriod(startOfYear, addYears)
}
