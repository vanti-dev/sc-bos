<template>
  <div>
    <slot :resource="openClosedValue"/>
  </div>
</template>

<script setup>
import {closeResource, newResourceValue} from '@/api/resource';
import {pullOpenClosePositions} from '@/api/sc/traits/open-close';
import {computed, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullAccessAttemptsRequest.AsObject
    default: () => {}
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const openClosedValue = reactive(
    /** @type {ResourceValue<OpenClose.AsObject, PullOpenClosesResponse>} */ newResourceValue()
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
</script>
