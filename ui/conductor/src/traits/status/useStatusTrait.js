import {closeResource, newResourceValue} from '@/api/resource';
import {pullCurrentStatus} from '@/api/sc/traits/status';
import {computed, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {PullCurrentStatusRequest.AsObject} [props.request]
 * @param {boolean} [props.paused]
 * @return {{statusLogValue: import('vue').UnwrapNestedRefs<ResourceValue<StatusLog.AsObject>>}}
 */
export default function(props) {
  const statusLogValue = reactive(/** @type {ResourceValue<StatusLog.AsObject>} */ newResourceValue());
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
          closeResource(statusLogValue);
        }

        if (!newPaused && (oldPaused || !reqEqual)) {
          closeResource(statusLogValue);
          pullCurrentStatus(newReq, statusLogValue);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  onUnmounted(() => {
    closeResource(statusLogValue);
  });

  return {
    statusLogValue
  };
}
