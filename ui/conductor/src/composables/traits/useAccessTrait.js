import {closeResource, newResourceValue} from '@/api/resource';
import {pullAccessAttempts} from '@/api/sc/traits/access';
import {computed, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {PullAccessAttemptsRequest.AsObject} [props.request]
 * @param {boolean} [props.paused]
 * @return {{accessAttemptValue: ResourceValue<AccessAttempt.AsObject, PullAccessAttemptsResponse>}}
 */
export default function(props) {
  const accessAttemptValue = reactive(
      /** @type {ResourceValue<AccessAttempt.AsObject, PullAccessAttemptsResponse>} */ newResourceValue());

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
          closeResource(accessAttemptValue);
        }

        if (!newPaused && (oldPaused || !reqEqual)) {
          closeResource(accessAttemptValue);
          pullAccessAttempts(newReq, accessAttemptValue);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  onUnmounted(() => {
    closeResource(accessAttemptValue);
  });

  return {
    accessAttemptValue
  };
}
