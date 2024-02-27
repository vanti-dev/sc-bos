<template>
  <v-card elevation="0" tile>
    <!-- Brightness -->
    <v-list tile class="ma-0 pa-0">
      <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Lighting</v-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">Brightness</v-list-item-title>
        <v-list-item-subtitle class="text-capitalize">{{ brightnessLevelString }}</v-list-item-subtitle>
      </v-list-item>
    </v-list>
    <v-progress-linear
        height="34"
        class="mx-4 my-2"
        :value="brightnessLevelNumber"
        background-color="neutral lighten-1"
        color="accent"/>
    <v-card-actions class="px-4">
      <v-btn
          v-for="control in brightnessControl.left"
          color="neutral lighten-1"
          :disabled="control.disabled"
          elevation="0"
          :key="control.label"
          small
          @click="control.onClick">
        {{ control.label }}
      </v-btn>
      <v-spacer/>
      <v-btn
          v-for="control in brightnessControl.right"
          color="neutral lighten-1"
          :disabled="control.disabled"
          elevation="0"
          :key="control.label"
          small
          @click="control.onClick">
        {{ control.label }}
      </v-btn>
    </v-card-actions>

    <!-- Presets -->
    <v-container v-if="brightnessPresets.length > 0" class="pa-0 px-4">
      <v-list tile class="ma-0 mt-2 pa-0">
        <v-list-item class="py-1 pl-0">
          <v-list-item-title class="text-body-small text-capitalize">Presets</v-list-item-title>
        </v-list-item>
      </v-list>
      <v-btn
          v-for="preset in brightnessPresets"
          block
          class="py-1 mx-0 mt-1 mb-2 preset"
          color="neutral lighten-1"
          elevation="0"
          :key="preset.name"
          small
          width="100%"
          max-width="575"
          @click="updateBrightness(toPresetObject(preset))">
        <span class="text-truncate">
          {{ preset.title ? preset.title : preset.name }}
        </span>
      </v-btn>
    </v-container>
    <v-progress-linear color="primary" indeterminate :active="loading"/>
  </v-card>
</template>

<script setup>
import useAuthSetup from '@/composables/useAuthSetup';
import useLightingTrait from '@/traits/lighting/useLightingTrait.js';
import {computed} from 'vue';

const {blockActions} = useAuthSetup();
const props = defineProps({
  name: {
    type: String,
    default: ''
  }
});
const {
  brightnessLevelString,
  brightnessLevelNumber,
  brightnessPresets,
  toLevelPercentObject,
  toPresetObject,
  updateBrightness,
  loading
} = useLightingTrait(() => props.name, false);

const brightnessControl = computed(() => {
  return {
    left: [
      {
        disabled: blockActions.value,
        label: 'On',
        onClick: () => updateBrightness(toLevelPercentObject(100))
      },
      {
        disabled: blockActions.value,
        label: 'Off',
        onClick: () => updateBrightness(toLevelPercentObject(0))
      }
    ],
    right: [
      {
        disabled: blockActions.value || brightnessLevelNumber.value >= 100,
        label: 'Up',
        onClick: () => updateBrightness(toLevelPercentObject(brightnessLevelNumber.value + 1))
      },
      {
        disabled: blockActions.value || brightnessLevelNumber.value <= 0,
        label: 'Down',
        onClick: () => updateBrightness(toLevelPercentObject(brightnessLevelNumber.value - 1))
      }
    ]
  };
});
</script>

<style lang="scss" scoped>
.v-list-item {
  min-height: auto;
}

.v-progress-linear {
  width: auto;
}

.preset ::v-deep(.v-btn__content) {
  max-width: 100%;
}
</style>
