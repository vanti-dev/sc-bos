<template>
  <div>
    <slot :resource="lightValue" :update="doUpdateBrightness" :update-tracker="updateValue"/>
  </div>
</template>

<script setup>
import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {pullBrightness, updateBrightness} from '@/api/sc/traits/light';
import {useErrorStore} from '@/components/ui-error/error';
import {onMounted, onUnmounted, reactive, watch} from 'vue';

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

const errorStore = useErrorStore();

const lightValue = reactive(
    /** @type {ResourceValue<Brightness.AsObject, Brightness>} */
    newResourceValue());

const updateValue = reactive(
    /** @type {ActionTracker<Brightness.AsObject>}  */
    newActionTracker());


//
//
// Methods
/**
 * @param {number|Brightness.AsObject|UpdateBrightnessRequest.AsObject} req
 */
function doUpdateBrightness(req) {
  if (typeof req === 'number') {
    req = {levelPercent: Math.min(100, Math.round(req))};
  }
  if (!req.hasOwnProperty('brightness')) {
    req = {brightness: req};
  }
  req.name = props.name;
  updateBrightness(req, updateValue);
}

//
//
// Watch
// Depending on paused state/device name, we close/open data stream(s)
watch(
    [() => props.paused, () => props.name],
    ([newPaused, newName], [oldPaused, oldName]) => {
      // only for LightSensor
      if (newPaused === oldPaused && newName === oldName) return;

      if (newPaused) {
        closeResource(lightValue);
      }

      if (!newPaused && (oldPaused || newName !== oldName)) {
        closeResource(lightValue);
        pullBrightness({name: newName}, lightValue);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

//
//
// UI error handling
let unwatchLightError;
let unwatchUpdateError;
onMounted(() => {
  unwatchLightError = errorStore.registerValue(lightValue);
  unwatchUpdateError = errorStore.registerTracker(updateValue);
});
onUnmounted(() => {
  closeResource(lightValue);
  if (unwatchLightError) unwatchLightError();
  if (unwatchUpdateError) unwatchUpdateError();
});
</script>

<style lang="scss">
.occupied {
  color: var(--v-success-lighten1) !important;
}

.idle {
  color: var(--v-info-base) !important;
}

.unoccupied {
  color: var(--v-warning-base) !important;
}
</style>
