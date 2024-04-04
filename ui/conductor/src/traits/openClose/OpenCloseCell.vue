<template>
  <StatusAlert v-if="props.streamError" icon="mdi-cancel" :resource="props.streamError"/>

  <v-tooltip v-else left>
    <template #activator="{ on }">
      <v-icon :class="doorState.class" right size="20" v-on="on">{{ doorState.icon }}</v-icon>
    </template>
    <span class="text-capitalize">{{ doorState?.text }}</span>
  </v-tooltip>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useOpenClosePositions} from '@/traits/openClose/openClose.js';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object, // of OpenClosePositions.AsObject
    default: () => {}
  },
  loading: {
    type: Boolean,
    default: false
  },
  streamError: {
    type: Object,
    default: null
  }
});

const {openStr, openIcon, openClass} = useOpenClosePositions(() => props.value);
const doorState = computed(() => {
  return {
    icon: openIcon.value,
    class: openClass.value,
    text: openStr.value
  };
});
</script>

<style scoped>
.open, .moving {
  color: var(--v-success-base);
}

.closed {
  color: var(--v-warning-base);
}

.unknown {
  color: grey;
}
</style>
