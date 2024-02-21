import {closeResource, newResourceValue} from '@/api/resource';
import {pullOpenClosePositions} from '@/api/sc/traits/open-close';
import {computed, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {PullOpenClosePositionsRequest.AsObject} [props.request]
 * @param {boolean} [props.paused]
 * @return {{
 *  openClosedValue: import('vue').UnwrapNestedRefs<ResourceValue<OpenClose.AsObject, PullOpenClosePositionsResponse>>
 * }}
 */
export default function(props) {
  const openClosedValue = reactive(
      /** @type {ResourceValue<OpenClose.AsObject, PullOpenClosePositionsResponse>} */ newResourceValue()
  );

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
          closeResource(openClosedValue);
        }

        if (!newPaused && (oldPaused || !reqEqual)) {
          closeResource(openClosedValue);
          pullOpenClosePositions(newReq, openClosedValue);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  onUnmounted(() => {
    closeResource(openClosedValue);
  });

  return {
    openClosedValue
  };
}
