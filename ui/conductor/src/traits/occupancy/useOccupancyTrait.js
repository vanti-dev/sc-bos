import {closeResource, newResourceValue} from '@/api/resource';
import {pullOccupancy} from '@/api/sc/traits/occupancy';
import {onUnmounted, reactive, watch} from 'vue';

/**
 *
 * @param {{
 *  paused: boolean,
 *  name: string
 *}} props
 *@return {{
 *  occupancyValue: ResourceValue<Occupancy.AsObject, Occupancy>
 *}}
 * }
 */
export default function(props) {
  const occupancyValue = reactive(
      /** @type {ResourceValue<Occupancy.AsObject, Occupancy>} */
      newResourceValue()
  );

  // Watch
  // Depending on paused state/device name, we close/open data stream(s)
  watch(
      [() => props.paused, () => props.name],
      ([newPaused, newName], [oldPaused, oldName]) => {
        if (newPaused === oldPaused && newName === oldName) return;

        if (newPaused) {
          closeResource(occupancyValue);
        }

        if (!newPaused && (oldPaused || newName !== oldName)) {
          closeResource(occupancyValue);
          pullOccupancy({name: newName}, occupancyValue);
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  onUnmounted(() => {
    closeResource(occupancyValue);
  });

  return {
    occupancyValue
  };
}
