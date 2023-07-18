<template>
  <div>
    <slot :resource="meterReadings" :info="meterReadingInfo"/>
  </div>
</template>

<script setup>
import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {describeMeterReading, pullMeterReading} from '@/api/sc/traits/meter';
import {useErrorStore} from '@/components/ui-error/error';
import {onMounted, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  paused: {
    type: Boolean,
    default: false
  }
});

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

// UI error handling
const errorStore = useErrorStore();
const unwatchErrors = [];
onMounted(() => {
  unwatchErrors.push(
      errorStore.registerValue(meterReadings)
  );
});
onUnmounted(() => {
  closeResource(meterReadings);
  unwatchErrors.forEach(unwatch => unwatch());
});
</script>
