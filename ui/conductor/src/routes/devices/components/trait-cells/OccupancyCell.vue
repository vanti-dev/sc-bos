<template>
  <div>
    <p :class="[occupancyState.toLowerCase(), 'ma-0 text-body-2']">
      {{ occupancyState }}
    </p>
    <v-progress-linear color="primary" indeterminate :active="props.loading"/>
  </div>
</template>

<script setup>
import {occupancyStateToString} from '@/api/sc/traits/occupancy';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => {}
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const occupancyState = computed(() => {
  if (props.value) {
    return occupancyStateToString(props.value.state);
  }
  return '';
});
</script>
