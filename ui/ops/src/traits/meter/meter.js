import {timestampToDate} from '@/api/convpb.js';
import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {describeMeterReading, listMeterReadingHistory, pullMeterReading} from '@/api/sc/traits/meter';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {computed, onScopeDispose, reactive, toRefs, toValue} from 'vue';

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
    return _s.value?.unit ?? '';
  });

  const usage = computed(() => {
    return _v.value?.usage;
  });
  const usageStr = computed(() => {
    return usage.value?.toFixed(2) ?? '-';
  });
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
 * @param {MaybeRefOrGetter<string|ListMeterReadingHistoryRequest.AsObject>} req - The request
 * @param {MaybeRefOrGetter<boolean>=} paused - Whether to pause the data stream
 * @return {Promise<ListMeterReadingHistoryResponse.AsObject>}
 */
export function useListMeterReadingHistory(req, paused = false) {
  const tracker = reactive(
      /** @type {ActionTracker<MeterReadingRecord.AsObject>} */
      newActionTracker()
  );
  return listMeterReadingHistory(req, tracker);
}

/**
 * @param {MeterReadingRecord | MeterReadingRecord.AsObject} obj
 * @return {MeterReadingRecord.AsObject & {recordTime: Date|undefined}}
 */
export function meterReadingRecordToObject(obj) {
  if (!obj) return undefined;
  if (typeof obj.toObject === 'function') obj = obj.toObject();
  if (obj.recordTime) obj.recordTime = timestampToDate(obj.recordTime);
  return obj;
}
