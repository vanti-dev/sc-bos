<template>
  <div>
    <slot :value="slotValue"/>
  </div>
</template>

<script setup>
import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';
import {closeResource, newActionTracker, newResourceValue} from '@/api/resource';
import {pullBrightness, updateBrightness} from '@/api/sc/traits/light';
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

const errorStore = useErrorStore();

const lightValue = reactive(
    /** @type {ResourceValue<Brightness.AsObject, Brightness>} */
    newResourceValue());

const updateValue = reactive(
    /** @type {ActionTracker<Brightness.AsObject>}  */
    newActionTracker());

//
//
// Computed
const brightness = computed(() => {
  if (lightValue && lightValue.value) {
    return Math.round(lightValue.value.levelPercent);
  }
  return '';
});

const slotValue = computed(() => {
  return {
    lightValue,
    updateValue,
    brightness: brightness.value,
    updateLight
  };
});

//
//
// Methods
/**
 * @param {number} value
 */
function updateLight(value) {
  /* @type {UpdateBrightnessRequest.AsObject} */
  const req = {
    name: props.name,
    brightness: {
      levelPercent: Math.min(100, Math.round(value))
    }
  };
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
        pullBrightness(newName, lightValue);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

//
//
// UI error handling
let unwatchLightError; let unwatchUpdateError;
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
