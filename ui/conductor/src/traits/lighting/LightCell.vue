<template>
  <StatusAlert v-if="error" icon="mdi-lightbulb-outline" :resource="error"/>

  <span v-else class="d-flex flex-row flex-nowrap">
    <v-tooltip bottom>
      <template #activator="{ on, attrs }">
        <span v-on="on" v-bind="attrs" class="d-flex flex-row">
          <span class="text-caption" style="min-width: 4ex">{{ brightnessLevelString }}</span>
          <v-icon right :color="brightnessLevelNumber > 0 ? 'yellow' : 'white' " size="20">
            {{ lightIcon }}
          </v-icon>
        </span>
      </template>
      <span>Lighting</span>
    </v-tooltip>
  </span>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import useLightingTrait from '@/traits/lighting/useLightingTrait.js';

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  paused: {
    type: Boolean,
    default: false
  }
});
const {
  brightnessLevelString,
  brightnessLevelNumber,
  lightIcon,
  error
} = useLightingTrait(() => props.name, () => props.paused);
</script>
