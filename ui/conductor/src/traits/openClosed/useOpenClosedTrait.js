import {newResourceValue} from '@/api/resource';
import {pullOpenClosePositions} from '@/api/sc/traits/open-close';
import {toQueryObject, watchResource} from '@/util/traits';
import {toValue} from '@/util/vue';
import {computed, reactive} from 'vue';

/**
 *
 * @param {MaybeRefOrGetter<string>} query
 * @param {MaybeRefOrGetter<boolean>} paused
 * @return {{
 *  openClosedValue: import('vue').UnwrapNestedRefs<ResourceValue<OpenClose.AsObject, PullOpenClosePositionsResponse>>,
 *  openClosedPercentage: import('vue').ComputedRef<string>,
 *  openClosedDoorState: import('vue').ComputedRef<{icon: string, text: string, class: string}>,
 *  error: import('vue').ComputedRef<ResourceError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const openClosedValue = reactive(
      /** @type {ResourceValue<OpenClose.AsObject, PullOpenClosePositionsResponse>} */ newResourceValue()
  );

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullOpenClosePositions(req, openClosedValue);
        return openClosedValue;
      }
  );

  // ---------------- Open/Close ---------------- //
  /**
   * Returns if the door is open/closed or how wide open (in percentage) as a string
   *
   * @type {import('vue').ComputedRef<string>} openClosedPercentage
   */
  const openClosedPercentage = computed(() => {
    if (!openClosedValue.value) return 'Unknown';

    const openPercent = openClosedValue.value?.statesList[0].openPercent;

    if (openPercent === 0) {
      return 'Closed';
    } else if (openPercent === 100) {
      return 'Open';
    } else {
      return openPercent + '%';
    }
  });

  /**
   * Returns the state of the door as an object with an icon, text, and class for styling
   *
   * @type {import('vue').ComputedRef<{icon: string, text: string, class: string}>} openClosedDoorState
   */
  const openClosedDoorState = computed(() => {
    if (openClosedPercentage.value === 'Unknown') {
      return {
        class: 'unknown',
        icon: 'mdi-door',
        text: ''
      };
    } else if (openClosedPercentage.value === 'Closed') {
      return {
        class: 'closed',
        icon: 'mdi-door-closed',
        text: 'Closed'
      };
    } else if (openClosedPercentage.value === 'Open') {
      return {
        class: 'open',
        icon: 'mdi-door-open',
        text: 'Open'
      };
    } else {
      return {
        class: 'moving',
        icon: 'mdi-door',
        text: 'Door ' + openClosedPercentage.value + ' open'
      };
    }
  });


  // ---------------- Error ---------------- //
  /** @type {import('vue').ComputedRef<ResourceError>} */
  const error = computed(() => openClosedValue.streamError);

  // ---------------- Loading ---------------- //
  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => openClosedValue.loading);


  return {
    // Resource
    openClosedValue,

    // Computed
    openClosedPercentage,
    openClosedDoorState,

    // State
    error,
    loading
  };
}
