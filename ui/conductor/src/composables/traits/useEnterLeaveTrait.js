import {closeResource, newResourceValue} from '@/api/resource';
import {pullEnterLeaveEvents} from '@/api/sc/traits/enter-leave';
import {computed, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {PullEnterLeaveEventsRequest.AsObject} [props.request]
 * @param {boolean} [props.paused]
 * @return {
 *  {enterLeaveEventValue: import('vue').UnwrapNestedRefs<
 *    ResourceValue<EnterLeaveEvent.AsObject, proto.smartcore.traits.PullEnterLeaveEventsResponse>
 *  >}
 * }
 */
export default function(props) {
  const enterLeaveEventValue = reactive(
      /** @type {ResourceValue<EnterLeaveEvent.AsObject, PullEnterLeaveEventsResponse>} */ newResourceValue());
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
          closeResource(enterLeaveEventValue);
        }

        if (!newPaused && (oldPaused || !reqEqual)) {
          closeResource(enterLeaveEventValue);
          pullEnterLeaveEvents(newReq, enterLeaveEventValue);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  onUnmounted(() => {
    closeResource(enterLeaveEventValue);
  });

  return {
    enterLeaveEventValue
  };
}
