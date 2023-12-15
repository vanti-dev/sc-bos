import {closeResource, newResourceValue} from '@/api/resource';
import {pullEmergency} from '@/api/sc/traits/emergency';
import {computed, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {PullEmergencyRequest.AsObject} [props.request]
 * @param {boolean} [props.paused]
 * @return {{emergencyValue: import('vue').UnwrapNestedRefs<
 *  ResourceValue<Emergency.AsObject, proto.smartcore.traits.PullEmergencyResponse>
 * >}}
 */
export default function(props) {
  const emergencyValue = reactive(
      /** @type {ResourceValue<Emergency.AsObject, PullEmergencyResponse>} */ newResourceValue());
  const _request = computed(() => {
    if (props.request) {
      return props.request;
    } else {
      return {name: props.name};
    }
  });

  watch(
      [() => _request.value, () => props.paused],
      ([newReq, newPaused], [oldReq, oldPaused]) => {
        const reqEqual = deepEqual(newReq, oldReq);
        if (newPaused === oldPaused && reqEqual) return;

        if (newPaused) {
          closeResource(emergencyValue);
        }

        if (!newPaused && (oldPaused || !reqEqual)) {
          closeResource(emergencyValue);
          pullEmergency(newReq, emergencyValue);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  onUnmounted(() => {
    closeResource(emergencyValue);
  });

  return {
    emergencyValue
  };
}
