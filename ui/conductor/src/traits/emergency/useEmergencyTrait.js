import {newResourceValue} from '@/api/resource';
import {pullEmergency} from '@/api/sc/traits/emergency';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue.js';
import {Emergency} from '@smart-core-os/sc-api-grpc-web/traits/emergency_pb';
import {computed, reactive} from 'vue';


/**
 * @template T
 * @param {MaybeRefOrGetter<T>} query
 * @param {MaybeRefOrGetter<boolean>} paused
 * @return {{
 *  emergencyValue: ResourceValue<Emergency.AsObject, PullEmergencyResponse>,
 *  emergencyColorClass: import('vue').ComputedRef<string>,
 *  emergencyIconString: import('vue').ComputedRef<string>,
 *  emergencyTooltipString: import('vue').ComputedRef<string>,
 *  error: import('vue').ComputedRef<ResourceError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const emergencyValue = reactive(
      /** @type {ResourceValue<Emergency.AsObject, PullEmergencyResponse>} */
      newResourceValue()
  );

  const queryObject = computed(() => toQueryObject(query));

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullEmergency(req, emergencyValue);
        return emergencyValue;
      });

  // ---------------- Emergency Display ---------------- //

  /** @type {import('vue').ComputedRef<string>} emergencyColorClass */
  const emergencyColorClass = computed(() => {
    const val = emergencyValue.value?.level;
    const drill = emergencyValue.value?.drill;
    switch (val) {
      default:
      case Emergency.Level.OK:
        return '';
      case Emergency.Level.WARNING:
        return 'warning--text';
      case Emergency.Level.EMERGENCY:
        if (drill) {
          return 'info--text';
        }
        return 'error--text';
    }
  });

  /** @type {import('vue').ComputedRef<string>} emergencyIconString */
  const emergencyIconString = computed(() => {
    const val = emergencyValue.value?.level;
    switch (val) {
      default:
      case Emergency.Level.OK:
        return 'mdi-smoke-detector-outline';
      case Emergency.Level.WARNING:
        return 'mdi-smoke-detector';
      case Emergency.Level.EMERGENCY:
        return 'mdi-smoke-detector-alert';
    }
  });

  /** @type {import('vue').ComputedRef<string>} emergencyTooltipString */
  const emergencyTooltipString = computed(() => {
    // todo: work out a better message based on current state
    return 'Emergency status';
  });

  // ---------------- Error ---------------- //

  /** @type {import('vue').ComputedRef<ResourceError>} */
  const error = computed(() => emergencyValue.streamError);

  // ---------------- Loading ---------------- //

  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => emergencyValue.loading);

  return {
    emergencyValue,

    emergencyColorClass,
    emergencyIconString,
    emergencyTooltipString,

    error,
    loading
  };
}
