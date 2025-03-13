import {timestampToDate} from '@/api/convpb.js';
import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {describeMeterReading, listMeterReadingHistory, pullMeterReading} from '@/api/sc/traits/meter';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {isNullOrUndef} from '@/util/types.js';
import {computed, effectScope, onScopeDispose, reactive, ref, toRefs, toValue, watch} from 'vue';

/**
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/meter_pb').MeterReading} MeterReading
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/meter_pb').MeterReadingSupport} MeterReadingSupport
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/meter_pb').PullMeterReadingsRequest} PullMeterReadingsRequest
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/meter_pb').PullMeterReadingsResponse} PullMeterReadingsResponse
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/meter_pb').DescribeMeterReadingRequest} DescribeMeterReadingRequest
 * @typedef {import('vue').ComputedRef} ComputedRef
 * @typedef {import('vue').Ref} Ref
 * @typedef {import('vue').ToRefs} ToRefs
 * @typedef {import('@/api/resource').ResourceValue} ResourceValue
 * @typedef {import('@/api/resource').ActionTracker} ActionTracker
 */

/**
 * @param {MaybeRefOrGetter<string|PullMeterReadingsRequest.AsObject>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>=} paused - Whether to pause the data stream
 * @return {ToRefs<ResourceValue<MeterReading.AsObject, PullMeterReadingsResponse>>}
 */
export function usePullMeterReading(query, paused = false) {
  const resource = reactive(
      /** @type {ResourceValue<MeterReading.AsObject, PullMeterReadingsResponse>} */
      newResourceValue()
  );
  onScopeDispose(() => closeResource(resource));

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullMeterReading(req, resource);
        return () => closeResource(resource);
      }
  );

  return toRefs(resource);
}

/**
 * @param {MaybeRefOrGetter<string|DescribeMeterReadingRequest.AsObject>} query - The name of the device or a query
 *   object
 * @return {ToRefs<ActionTracker<MeterReadingSupport.AsObject>>}
 */
export function useDescribeMeterReading(query) {
  const tracker = reactive(
      /** @type {ActionTracker<MeterReadingSupport.AsObject>} */
      newActionTracker()
  );

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      false,
      (req) => {
        describeMeterReading(req, tracker)
            .catch(() => {}); // errors are tracked by tracker
        return () => closeResource(tracker);
      }
  );

  return toRefs(tracker);
}

/**
 * Converts a usage and optional unit to a string with appropriate precision.
 *
 * @param {number | null | undefined} usage
 * @param {string} [unit]
 * @return {string}
 */
export function usageToString(usage, unit = '') {
  const usageStr = (() => {
    if (isNullOrUndef(usage)) return '-';
    if (Math.abs(usage) < 100) return usage.toPrecision(2);
    return usage.toLocaleString(undefined, {maximumFractionDigits: 0});
  })();
  if (unit) {
    return `${usageStr} ${unit}`;
  }
  return usageStr;
}

/**
 * @param {MaybeRefOrGetter<MeterReading.AsObject|null>} value
 * @param {MaybeRefOrGetter<MeterReadingSupport.AsObject|null>=} support
 * @return {{
 *   unit: ComputedRef<string>,
 *   usage: ComputedRef<number | undefined>,
 *   usageStr: ComputedRef<string>,
 *   usageAndUnit: ComputedRef<string>,
 *   table: ComputedRef<Array<{label:string, unit:string, value:string}>>
 * }}
 */
export function useMeterReading(value, support = null) {
  const _v = computed(() => toValue(value));
  const _s = computed(() => toValue(support));

  const unit = computed(() => {
    return _s.value?.usageUnit ?? '';
  });

  const usage = computed(() => {
    return _v.value?.usage;
  });
  const usageStr = computed(() => usageToString(usage.value));
  const usageAndUnit = computed(() => {
    let val = usageStr.value;
    if (unit.value) {
      val += ` ${unit.value}`;
    }
    return val;
  });

  const table = computed(() => {
    return [{
      label: 'Usage',
      unit: unit.value,
      value: usageStr.value
    }];
  });

  return {
    unit,
    usage,
    usageStr,
    usageAndUnit,
    table
  };
}

/**
 * Returns a ref containing the estimated meter reading at time t.
 * Name must implement MeterHistory trait aspect.
 * Results can be very inaccurate if the history records for the meter are sparse around time t.
 *
 * @param {import('vue').MaybeRefOrGetter<string>} name
 * @param {import('vue').MaybeRefOrGetter<Date>} t
 * @return {import('vue').ComputedRef<null | MeterReading.AsObject>}
 */
export function useMeterReadingAt(name, t) {
  const usageAtT = ref(/** @type {null|number} */ null);
  const readingAtT = computed(() => {
    if (usageAtT.value === null) return null;
    return /** @type {MeterReading.AsObject} */ {usage: usageAtT.value};
  });

  const fetching = ref(/** @type {null | {t: Date, cancel: () => void}} */ null);
  watch([() => toValue(name), () => toValue(t)], async ([name, t]) => {
    const cancelled = () => {
      if (fetching.value === null) return true;
      return fetching.value.t.getTime() !== t.getTime();
    };
    if (fetching.value !== null) {
      if (fetching.value.t.getTime() === t.getTime()) {
        return; // already working on it
      }
      fetching.value.cancel();
    }

    // if we don't have all the info, don't do anything
    if (isNullOrUndef(name) || isNullOrUndef(t)) {
      fetching.value = null;
      return;
    }

    // note: currently no way to cancel unary RPCs,
    // see: https://github.com/grpc/grpc-web/issues/946
    fetching.value = {t, cancel: () => {}};
    try {
      const [before, after] = await Promise.all([
        getReadingBefore(name, t),
        getReadingOnOrAfter(name, t)
      ]);
      if (cancelled()) return;
      usageAtT.value = interpolateUsage(before, after, t);
    } catch (e) {
      if (!cancelled()) {
        console.warn('Failed to get meter reading at', t, e.message ?? e);
      }
    } finally {
      if (!cancelled()) {
        fetching.value = null;
      }
    }
  }, {immediate: true});

  return readingAtT;
}

/**
 * Returns an of meter readings, one for each of the times in ts.
 * Times in the future will result in a null returned reading at that index.
 *
 * @param {import('vue').MaybeRefOrGetter<string>} name
 * @param {import('vue').MaybeRefOrGetter<Date[]>} ts
 * @return {import('vue').ComputedRef<Array<null | MeterReading.AsObject>>}
 */
export function useMeterReadingsAt(name, ts) {
  /**
   * @typedef {Object} Watcher
   * @property {function(): void} stop
   * @property {import('vue').ComputedRef<null | MeterReading.AsObject>} reading
   */
  /** @type {Record<number, Watcher>} */
  const trackers = {}; // keyed by Date.getTime()
  onScopeDispose(() => Object.values(trackers).forEach(({stop}) => stop()));

  return computed(() => {
    const _ts = toValue(ts);
    const toStop = Object.fromEntries(Object.keys(trackers).map(k => [k, true]));

    for (const t of _ts) {
      const k = t.getTime();
      const tracker = trackers[k];
      if (tracker) {
        delete toStop[k];
        continue;
      }
      const scope = effectScope();
      scope.run(() => {
        const reading = useMeterReadingAt(name, t);
        trackers[k] = {stop: () => scope.stop(), reading};
      });
    }

    for (const k of Object.keys(toStop)) {
      trackers[k].stop();
      delete trackers[k];
    }

    return _ts.map(t => {
      const k = t.getTime();
      return trackers[k]?.reading.value ?? null;
    });
  });
}

/**
 * Returns the newest meter reading record before t, if there is one.
 *
 * @param {string} name
 * @param {Date} t
 * @return {Promise<MeterReadingRecord.AsObject | undefined>}
 */
async function getReadingBefore(name, t) {
  const res = await listMeterReadingHistory({
    name,
    pageSize: 1,
    orderBy: 'recordTime desc',
    period: {endTime: t}
  });
  return res.meterReadingRecordsList?.[0];
};

/**
 * Returns the oldest meter reading record at or after t.
 *
 * @param {string} name
 * @param {Date} t
 * @return {Promise<MeterReadingRecord.AsObject | undefined>}
 */
async function getReadingOnOrAfter(name, t) {
  const res = await listMeterReadingHistory({name, pageSize: 1, period: {startTime: t}});
  return res.meterReadingRecordsList?.[0];
}


/**
 * Returns the estimated meter usage at `at` between the two known readings.
 *
 * @param {MeterReadingRecord.AsObject | undefined} a
 * @param {MeterReadingRecord.AsObject | undefined} b
 * @param {Date} at
 * @return {number | null}
 */
function interpolateUsage(a, b, at) {
  if (a && b) {
    // meters only decrease if they are reset, if they are reset use the later reading
    if (a.meterReading.usage > b.meterReading.usage) return b.meterReading.usage;
    const dt = (at.getTime() - timestampToDate(a.recordTime)) /
        (timestampToDate(b.recordTime) - timestampToDate(a.recordTime));
    return (dt * (b.meterReading.usage - a.meterReading.usage)) + a.meterReading.usage;
  } else if (a) {
    return a.meterReading.usage;
  } else if (b) {
    return b.meterReading.usage
  } else {
    return null;
  }
}