import {closeResource, newResourceValue} from '@/api/resource';
import {pullDemand} from '@/api/sc/traits/electric';
import {computed, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {PullDemandRequest.AsObject} [props.request]
 * @param {boolean} [props.paused]
 * @return {{
 *  demandValue: import('vue').UnwrapNestedRefs<
 *    ResourceValue<ElectricDemand.AsObject, proto.smartcore.traits.PullDemandResponse>
 *  >
 * }}
 */
export default function(props) {
  const demandValue = reactive(
      /** @type {ResourceValue<ElectricDemand.AsObject, PullDemandResponse>} */
      newResourceValue()
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
          closeResource(demandValue);
        }

        if (!newPaused && (oldPaused || !reqEqual)) {
          closeResource(demandValue);
          pullDemand(newReq, demandValue);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  onUnmounted(() => {
    closeResource(demandValue);
  });

  return {
    demandValue
  };
}
