<template>
  <div>
    <slot v-bind="statusLogValue"/>
  </div>
</template>

<script setup>
import {closeResource, newResourceValue} from '@/api/resource';
import {pullCurrentStatus} from '@/api/sc/traits/status';
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
    type: Object, // of type PullCurrentStatusRequest.AsObject
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const statusLogValue = reactive(/** @type {ResourceValue<StatusLog.AsObject>} */ newResourceValue());
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
        closeResource(statusLogValue);
      }

      if (!newPaused && (oldPaused || !reqEqual)) {
        closeResource(statusLogValue);
        pullCurrentStatus(newReq, statusLogValue);
      }
    },
    {immediate: true, deep: true, flush: 'sync'}
);

// UI error handling
const errorStore = useErrorStore();
const unwatchErrors = [];
onMounted(() => {
  unwatchErrors.push(
      errorStore.registerValue(statusLogValue)
  );
});
onUnmounted(() => {
  closeResource(statusLogValue);
  unwatchErrors.forEach(unwatch => unwatch());
});

</script>
