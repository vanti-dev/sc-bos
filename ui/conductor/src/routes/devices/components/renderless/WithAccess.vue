<template>
  <div>
    <slot :resource="accessAttemptValue"/>
  </div>
</template>

<script setup>
import {closeResource, newResourceValue} from '@/api/resource';
import {pullAccessAttempts} from '@/api/sc/traits/access';
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
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const accessAttemptValue = reactive(
    /** @type {ResourceValue<AccessAttempt.AsObject, PullAccessAttemptsResponse>} */ newResourceValue());

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
        closeResource(accessAttemptValue);
      }

      if (!newPaused && (oldPaused || !reqEqual)) {
        closeResource(accessAttemptValue);
        pullAccessAttempts(newReq, accessAttemptValue);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

onUnmounted(() => {
  closeResource(accessAttemptValue);
});
</script>
