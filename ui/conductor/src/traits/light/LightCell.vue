<template>
  <status-alert v-if="error" icon="mdi-lightbulb-outline" :resource="error"/>

  <span v-else>
    <v-tooltip location="bottom">
      <template #activator="{ props: _props }">
        <span v-bind="_props" class="d-flex flex-row">
          <span class="text-caption" style="min-width: 4ex">{{ levelStr }}</span>
          <light-icon :level="level" class="ml-2" size="20"/>
        </span>
      </template>
      <span>Lighting</span>
    </v-tooltip>
  </span>
</template>

<script setup>
import StatusAlert from '@/components/StatusAlert.vue';
import {useBrightness, usePullBrightness} from '@/traits/light/light.js';
import LightIcon from '@/traits/light/LightIcon.vue';

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
const {value, streamError: error} = usePullBrightness(() => props.name, () => props.paused);
const {levelStr, level} = useBrightness(value);
</script>
