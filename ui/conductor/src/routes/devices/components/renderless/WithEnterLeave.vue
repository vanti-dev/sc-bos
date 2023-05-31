<template>
  <div>
    <slot :value="enterLeaveEventValue"/>
  </div>
</template>

<script setup>
import {closeResource, newResourceValue} from '@/api/resource';
import {pullEnterLeaveEvents} from '@/api/sc/traits/enter-leave';
import {useErrorStore} from '@/components/ui-error/error';
import {computed, onMounted, onUnmounted, reactive, watch} from 'vue';
import {deepEqual} from 'vuetify/src/util/helpers';

const props = defineProps({
  // unique name of the device
  name: {
    type: String,
    default: ''
  },
  request: {
    type: Object, // of type PullEnterLeaveEventsRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const enterLeaveEventValue = reactive(
    /** @type {ResourceValue<EnterLeaveEvent.AsObject, PullEnterLeaveEventsResponse>} */ newResourceValue());
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
        closeResource(enterLeaveEventValue);
      }

      if (!newPaused && (oldPaused || !reqEqual)) {
        closeResource(enterLeaveEventValue);
        pullEnterLeaveEvents(newReq, enterLeaveEventValue);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

// UI error handling
const errorStore = useErrorStore();
const unwatchErrors = [];
onMounted(() => {
  unwatchErrors.push(
      errorStore.registerValue(enterLeaveEventValue)
  );
});
onUnmounted(() => {
  closeResource(enterLeaveEventValue);
  unwatchErrors.forEach(unwatch => unwatch());
});
</script>
