<template>
  <div>
    <slot :resource="meterValue" :type="meterType"/>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, reactive, watch} from 'vue';
import {closeResource, newResourceValue, newActionTracker} from '@/api/resource';
import {pullMeterReading, describeMeterReading} from '@/api/sc/traits/meter';
import {deepEqual} from 'vuetify/src/util/helpers';
import {useErrorStore} from '@/components/ui-error/error';

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

const meterValue = reactive(
    /** @type {ResourceValue<MeterReading.AsObject, MeterReading>} */
    newResourceValue()
);
const meterType = reactive(
    /** @type {ActionTracker<MeterReadingSupport.AsObject, MeterReadingSupport>} */
    newActionTracker()
);

watch(
    [() => props.name, () => props.paused],
    ([newName, newPaused], [oldName, oldPaused]) => {
      const nameEqual = deepEqual(newName, oldName);
      if (newPaused === oldPaused && nameEqual) return;

      if (newPaused) {
        closeResource(meterValue);
      }

      if (!newPaused && (oldPaused || !nameEqual)) {
        closeResource(meterValue);
        pullMeterReading(newName, meterValue); // pulls in unit value
        describeMeterReading(newName, meterType); // pulls in unit type
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

// UI error handling
const errorStore = useErrorStore();
const unwatchErrors = [];
onMounted(() => {
  unwatchErrors.push(
      errorStore.registerValue(meterValue)
  );
});
onUnmounted(() => {
  closeResource(meterValue);
  unwatchErrors.forEach(unwatch => unwatch());
});
</script>
