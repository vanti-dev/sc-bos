import {newResourceValue} from '@/api/resource';
import {pullDemand} from '@/api/sc/traits/electric';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue.js';
import {computed, reactive} from 'vue';

/**
 * @param {MaybeRefOrGetter<string|PullDemandRequest.AsObject>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @return {{
 *  electricDemandValue: import('vue').ResourceValue<ElectricDemand.AsObject, PullDemandResponse>,
 *  electricDemandRealPowerNumber: import('vue').ComputedRef<number>,
 *  electricDemandRealPowerObject: import('vue').ComputedRef<{label: string, value: string, unit: string}>,
 *  electricDemandApparentPowerNumber: import('vue').ComputedRef<number>,
 *  electricDemandApparentPowerObject: import('vue').ComputedRef<{label: string, value: string, unit: string}>,
 *  electricDemandReactivePowerNumber: import('vue').ComputedRef<number>,
 *  electricDemandReactivePowerObject: import('vue').ComputedRef<{label: string, value: string, unit: string}>,
 *  electricDemandPowerFactorNumber: import('vue').ComputedRef<number>,
 *  electricDemandPowerFactorObject: import('vue').ComputedRef<{label: string, value: string, unit: undefined}>,
 *  error: import('vue').ComputedRef<ResourceError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const electricDemandValue = reactive(
      /** @type {ResourceValue<ElectricDemand.AsObject, PullDemandResponse>} */
      newResourceValue()
  );

  const queryObject = computed(() => toQueryObject(query));

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullDemand(req, electricDemandValue);
        return electricDemandValue;
      });


  /** @type {import('vue').ComputedRef<number>} electricDemandRealPower */
  const electricDemandRealPowerNumber = computed(() => electricDemandValue.value?.realPower);

  /**
   * @type {
   *   import('vue').ComputedRef<{label: string, value: string, unit: string}>
   * } electricDemandRealPowerObject
   */
  const electricDemandRealPowerObject = computed(() => {
    return {
      label: 'Real Power',
      value: (electricDemandRealPowerNumber.value ? electricDemandRealPowerNumber.value / 1000 : 0).toFixed(3),
      unit: 'kW'
    };
  });

  /** @type {import('vue').ComputedRef<number>} electricDemandApparentPowerNumber */
  const electricDemandApparentPowerNumber = computed(() => electricDemandValue.value?.apparentPower);

  /**
   * @type {
   *   import('vue').ComputedRef<{label: string, value: string, unit: string}>
   * } electricDemandApparentPowerObject
   */
  const electricDemandApparentPowerObject = computed(() => {
    return {
      label: 'Apparent Power',
      value: (electricDemandApparentPowerNumber.value ? electricDemandApparentPowerNumber.value / 1000 : 0).toFixed(3),
      unit: 'kVA'
    };
  });

  /** @type {import('vue').ComputedRef<number>} electricDemandReactivePowerNumber */
  const electricDemandReactivePowerNumber = computed(() => electricDemandValue.value?.reactivePower);

  /**
   * @type {
   *   import('vue').ComputedRef<{label: string, value: string, unit: string}>
   * } electricDemandReactivePowerObject
   */
  const electricDemandReactivePowerObject = computed(() => {
    return {
      label: 'Reactive Power',
      value: (electricDemandReactivePowerNumber.value ? electricDemandReactivePowerNumber.value / 1000 : 0).toFixed(3),
      unit: 'kVAr'
    };
  });

  /** @type {import('vue').ComputedRef<number>} electricDemandPowerFactorNumber */
  const electricDemandPowerFactorNumber = computed(() => electricDemandValue.value?.powerFactor);

  /**
   * @type {
   *   import('vue').ComputedRef<{label: string, value: string, unit: undefined}>
   * } electricDemandPowerFactorObject
   */
  const electricDemandPowerFactorObject = computed(() => {
    return {
      label: 'Power Factor',
      value: electricDemandPowerFactorNumber.value?.toFixed(2),
      unit: undefined
    };
  });


  /** @type {import('vue').ComputedRef<ResourceError>} */
  const error = computed(() => electricDemandValue.streamError);


  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => electricDemandValue.loading);

  return {
    electricDemandValue,

    electricDemandRealPowerNumber,
    electricDemandRealPowerObject,
    electricDemandApparentPowerNumber,
    electricDemandApparentPowerObject,
    electricDemandReactivePowerNumber,
    electricDemandReactivePowerObject,
    electricDemandPowerFactorNumber,
    electricDemandPowerFactorObject,

    error,
    loading
  };
}
