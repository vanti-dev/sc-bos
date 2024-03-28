import {newActionTracker, newResourceValue} from '@/api/resource';
import {describeMeterReading, pullMeterReading} from '@/api/sc/traits/meter';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue.js';
import {computed, reactive} from 'vue';

/**
 * @param {MaybeRefOrGetter<string|PullMeterReadingsRequest.AsObject>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @param {{
 *  pullMeterReading?: boolean,
 *  describeMeterReading?: boolean
 * }} [options] - Options to control which API calls to make
 * @return {{
 *  meterValue: ResourceValue<MeterReading.AsObject, PullMeterReadingsResponse>,
 *  meterSupport: ActionTracker<MeterReadingSupport.AsObject>,
 *  meterUnit: import('vue').ComputedRef<string>,
 *  meterReadingNumber: import('vue').ComputedRef<number|string>,
 *  meterReadingObject: import('vue').ComputedRef<{[key: string]: number|string}>,
 *  error: import('vue').ComputedRef<ResourceError|ActionError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused, options) {
  const meterValue = reactive(
      /** @type {ResourceValue<MeterReading.AsObject, PullMeterReadingsResponse>} */
      newResourceValue()
  );

  const meterSupport = reactive(
      /** @type {ActionTracker<MeterReadingSupport.AsObject>} */
      newActionTracker()
  );

  const apiCalls = [];

  // For meter state and updates
  if (options?.pullMeterReading !== false) {
    apiCalls.push((req) => {
      pullMeterReading(req, meterValue);
      return meterValue;
    });
  }
  if (options?.describeMeterReading !== false) {
    apiCalls.push((req) => {
      describeMeterReading(req, meterSupport);
      return meterSupport;
    });
  }

  const queryObject = computed(() => toQueryObject(query));

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      ...apiCalls
  );

  /** @type {import('vue').ComputedRef<string>} */
  const meterUnit = computed(() => {
    return meterSupport.response?.unit || '';
  });

  /** @type {import('vue').ComputedRef<number|string>} */
  const meterReadingNumber = computed(() => {
    return meterValue.value?.usage?.toFixed(2) ?? '-';
  });

  /** @type {import('vue').ComputedRef<{[key: string]: number|string}>} */
  const meterReadingObject = computed(() => {
    return {
      label: 'Usage',
      unit: meterUnit.value,
      value: meterReadingNumber.value
    };
  });

  /** @type {import('vue').ComputedRef<ResourceError|ActionError>} */
  const error = computed(() => {
    return meterValue.streamError || meterSupport.error;
  });

  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => {
    return meterValue.loading || meterSupport.loading;
  });


  return {
    meterValue,
    meterSupport,

    meterUnit,
    meterReadingNumber,
    meterReadingObject,

    error,
    loading
  };
}
