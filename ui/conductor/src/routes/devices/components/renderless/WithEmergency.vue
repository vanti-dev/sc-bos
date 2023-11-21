<template>
  <div>
    <slot :resource="emergencyValue"/>
  </div>
</template>

<script setup>
import {closeResource, newResourceValue} from '@/api/resource';
import {pullEmergency} from '@/api/sc/traits/emergency';
import {computed, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullEmergencyRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const emergencyValue = reactive(
    /** @type {ResourceValue<Emergency.AsObject, PullEmergencyResponse>} */ newResourceValue());
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
        closeResource(emergencyValue);
      }

      if (!newPaused && (oldPaused || !reqEqual)) {
        closeResource(emergencyValue);
        pullEmergency(newReq, emergencyValue);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

onUnmounted(() => {
  closeResource(emergencyValue);
});
</script>
