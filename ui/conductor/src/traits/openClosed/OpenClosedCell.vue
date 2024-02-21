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
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object,
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

const doorState = computed(() => {
  if (!props.value) return {icon: 'mdi-door', class: 'unknown', text: ''};

  return props.value?.statesList[0].openPercent === 0 ?
      {icon: 'mdi-door-closed', class: 'closed', text: 'Closed'} :
      props.value?.statesList[0].openPercent === 100 ?
          {icon: 'mdi-door-open', class: 'open', text: 'Open'} :
          {
            icon: 'mdi-door',
            class: 'moving',
            text: '' + props.openClosePercentage?.value.statesList[0].openPercent + '%'
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
