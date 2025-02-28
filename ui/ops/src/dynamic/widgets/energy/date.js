import {DAY, HOUR} from '@/components/now.js';
import {isNamedPeriod, previousNamedPeriod, usePastDates, usePeriod} from '@/composables/time.js';
import {
  eachDayOfInterval,
  eachHourOfInterval,
  eachMinuteOfInterval,
  eachMonthOfInterval,
  eachYearOfInterval
} from 'date-fns';
import {computed, toValue} from 'vue';

/**
 * @param {import('vue').MaybeRefOrGetter<keyof useStartOf | string | number | Date>} start
 * @param {import('vue').MaybeRefOrGetter<keyof useStartOf | string | number | Date>} end
 * @param {import('vue').MaybeRefOrGetter<null | undefined | number | string>} offset
 * @return {{
 *   edges: import('vue').ComputedRef<Date[]>,
 *   pastEdges: import('vue').ComputedRef<Date[]>,
 *   tickUnit: import('vue').ComputedRef<string>,
 *   startDate: import('vue').Ref<Date|null>,
 *   endDate: import('vue').Ref<Date|null>
 * }}
 */
export function useDateScale(start, end, offset) {
  const {start: startDate, end: endDate} = usePeriod(start, end, offset);

  // Returns a named period that should be used for calculating how many ticks between startDate and endDate we should use.
  // For example, if start and end date encompass a single day, then the ticks should be in hours.
  const tickUnit = computed(() => {
    if (isNamedPeriod(toValue(start)) && toValue(start) === toValue(end)) {
      return previousNamedPeriod[toValue(start)] ?? 'year';
    }
    const diff = endDate.value.getTime() - startDate.value.getTime();
    if (diff <= 12 * HOUR) return 'minute';
    if (diff <= 3 * DAY) return 'hour';
    if (diff <= 2 * 30 * DAY) return 'day';
    if (diff <= 5 * 365 * DAY) return 'month';
    return 'year';
  })

  const edges = computed(() => {
    const start = startDate.value;
    const end = endDate.value;
    if (!start || !end || start >= end) return [];

    // we want to maintain no more than about 60 edges, so we increase the gaps between edges as the range increases
    switch (tickUnit.value) {
      case 'minute':
        return eachMinuteOfInterval({start, end}, {step: 5});
      case 'hour':
        return eachHourOfInterval({start, end});
      case 'day':
        return eachDayOfInterval({start, end});
      case 'month':
        return eachMonthOfInterval({start, end});
      case 'year':
      default: {
        const years = Math.ceil((end.getTime() - start.getTime()) / (365 * DAY));
        return eachYearOfInterval({start, end}, {step: Math.ceil(years / 60)});
      }
    }
  });

  // work out which edges are in the past, these are the ones we want to fetch data for
  const pastEdges = usePastDates(edges);

  return {
    edges, pastEdges,
    tickUnit,
    startDate, endDate
  }
}