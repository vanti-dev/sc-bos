import {newResourceValue} from '@/api/resource';
import {pullEnterLeaveEvents} from '@/api/sc/traits/enter-leave';
import {toQueryObject, watchResource} from '@/util/traits';
import {toValue} from '@/util/vue';
import {EnterLeaveEvent} from '@smart-core-os/sc-api-grpc-web/traits/enter_leave_sensor_pb';
import {computed, reactive, ref, watch} from 'vue';

/**
 * @param {MaybeRefOrGetter<string|PullEnterLeaveEventsRequest.AsObject>} query - The name of the device or a query
 *   object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @return {{
 *  enterLeaveValue: ResourceValue<EnterLeaveEvent.AsObject, PullEnterLeaveEventsResponse>,
 *  enterLeaveHasTotals: import('vue').ComputedRef<boolean>,
 *  enterTotal: import('vue').ComputedRef<number>,
 *  leaveTotal: import('vue').ComputedRef<number>,
 *  justEntered: import('vue').Ref<boolean>,
 *  justLeft: import('vue').Ref<boolean>,
 *  error: import('vue').ComputedRef<ResourceError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const enterLeaveValue = reactive(
      /** @type {ResourceValue<EnterLeaveEvent.AsObject, PullEnterLeaveEventsResponse>} */
      newResourceValue()
  );

  const queryObject = computed(() => toQueryObject(query));

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullEnterLeaveEvents(req, enterLeaveValue);
        return enterLeaveValue;
      }
  );

  // ---------------- Enter Leave Values ---------------- //

  /** @type {import('vue').ComputedRef<boolean>} */
  const enterLeaveHasTotals = computed(
      () => enterLeaveValue.value?.enterTotal !== undefined || enterLeaveValue.value?.leaveTotal !== undefined
  );

  /** @type {import('vue').ComputedRef<number>} */
  const enterTotal = computed(() => enterLeaveValue.value?.enterTotal || 0);

  /** @type {import('vue').ComputedRef<number>} */
  const leaveTotal = computed(() => enterLeaveValue.value?.leaveTotal || 0);

  /** @type {import('vue').Ref<number>} */
  const enterTimeoutHandle = ref(0);

  /** @type {import('vue').Ref<number>} */
  const leaveTimeoutHandle = ref(0);

  /** @type {import('vue').Ref<boolean>} */
  const justEntered = ref(false);

  /** @type {import('vue').Ref<boolean>} */
  const justLeft = ref(false);

  // Watch for changes in the enter/leave values and set the justEntered and justLeft flags accordingly
  watch(() => enterLeaveValue.value, (newVal, oldVal) => {
    if (!oldVal || !newVal) {
      justEntered.value = false;
      justLeft.value = false;
      clearTimeout(enterTimeoutHandle.value);
      clearTimeout(leaveTimeoutHandle.value);
      return;
    }

    if (newVal.direction === EnterLeaveEvent.Direction.ENTER) {
      justEntered.value = true;
      clearTimeout(enterTimeoutHandle.value);
      enterTimeoutHandle.value = setTimeout(() => {
        justEntered.value = false;
      }, props.showChangeDuration);
    }
    if (newVal.direction === EnterLeaveEvent.Direction.LEAVE) {
      justLeft.value = true;
      clearTimeout(leaveTimeoutHandle.value);
      leaveTimeoutHandle.value = setTimeout(() => {
        justLeft.value = false;
      }, props.showChangeDuration);
    }
  }, {deep: true});

  // ---------------- Error ---------------- //
  /** @type {import('vue').ComputedRef<ResourceError>} */
  const error = computed(() => {
    return enterLeaveValue.streamError;
  });

  // ---------------- Loading ---------------- //
  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => {
    return enterLeaveValue.loading;
  });

  return {
    enterLeaveValue,

    enterLeaveHasTotals,
    enterTotal,
    leaveTotal,

    justEntered,
    justLeft,

    error,
    loading
  };
}
