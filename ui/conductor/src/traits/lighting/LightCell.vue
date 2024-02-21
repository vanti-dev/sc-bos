<template>
  <StatusAlert v-if="props.streamError" icon="mdi-lightbulb-outline" :resource="props.streamError"/>

  <span v-else class="d-flex flex-row flex-nowrap">
    <v-tooltip bottom>
      <template #activator="{ on, attrs }">
        <span v-on="on" v-bind="attrs" class="d-flex flex-row">
          <span class="text-caption" style="min-width: 4ex">{{ brightnessStr }}</span>
          <v-icon right :color="brightness > 0 ? 'yellow' : 'white' " size="20">
            {{ lightingIcon }}
          </v-icon>
        </span>
      </template>
      <span>Lighting</span>
    </v-tooltip>
  </span>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => {
    }
  },
  streamError: {
    type: Object,
    default: null
  }
});

//
//
// Computed
const brightness = computed(() => props.value?.levelPercent);
const brightnessStr = computed(() => {
  if (brightness.value === 0) {
    return 'Off';
  } else if (brightness.value === 100) {
    return 'Max';
  } else if (brightness.value > 0 && brightness.value < 100) {
    return `${brightness.value.toFixed(0)}%`;
  }

  return '';
});

const lightingIcon = computed(() => {
  if (brightness.value === 0) return 'mdi-lightbulb-outline';
  if (brightness.value > 0) return 'mdi-lightbulb-on';
  return '';
});
</script>
