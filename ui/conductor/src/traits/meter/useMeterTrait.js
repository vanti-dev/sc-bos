import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {describeMeterReading, pullMeterReading} from '@/api/sc/traits/meter';
import {onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

/**
 *
 * @param {Object} props
 * @param {string} props.name
 * @param {boolean} [props.paused]
 * @return {{
 *  meterReadings: import('vue').UnwrapNestedRefs<
 *    ResourceValue<MeterReading.AsObject, proto.smartcore.bos.PullMeterReadingsResponse>
 *    >,
 *    meterReadingInfo: UnwrapNestedRefs<ActionTracker<MeterReadingSupport.AsObject>>
 * }}
 */
export default function(props) {
  const meterReadings = reactive(
      /** @type {ResourceValue<MeterReading.AsObject, PullMeterReadingsResponse>} */
      newResourceValue()
  );
  const meterReadingInfo = reactive(
      /** @type {ActionTracker<MeterReadingSupport.AsObject>} */
      newActionTracker()
  );

  watch(
      [() => props.name, () => props.paused],
      ([newName, newPaused], [oldName, oldPaused]) => {
        const nameEqual = deepEqual(newName, oldName);
        if (newPaused === oldPaused && nameEqual) return;

        if (newPaused) {
          closeResource(meterReadings);
        }

        if (!newPaused && (oldPaused || !nameEqual)) {
          closeResource(meterReadings);
          pullMeterReading({name: newName}, meterReadings); // pulls in unit value
          describeMeterReading({name: newName}, meterReadingInfo); // pulls in unit type
        }
      },
      {immediate: true, deep: true, flush: 'sync'}
  );

  onUnmounted(() => {
    closeResource(meterReadings);
  });

  return {
    meterReadings,
    meterReadingInfo
  };
}
