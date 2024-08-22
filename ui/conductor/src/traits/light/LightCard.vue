<template>
  <v-card elevation="0" tile>
    <!-- Brightness -->
    <v-list tile class="ma-0 pa-0">
      <v-list-subheader class="text-title-caps-large text-neutral-lighten-3">Lighting</v-list-subheader>
      <v-list-item class="py-1">
        <v-list-item-title class="text-body-small text-capitalize">Brightness</v-list-item-title>
        <template #append>
          <v-list-item-subtitle class="text-capitalize text-body-1">{{ levelStr }}</v-list-item-subtitle>
        </template>
      </v-list-item>
    </v-list>
    <v-progress-linear
        height="34"
        class="mx-4 my-2"
        :model-value="level"
        bg-color="neutral-lighten-1"
        bg-opacity="1"
        color="accent"/>
    <v-card-actions class="px-4">
      <v-btn
          v-for="control in brightnessControl.left"
          color="neutral-lighten-1"
          :disabled="control.disabled"
          elevation="0"
          :key="control.label"
          size="small"
          @click="control.onClick">
        {{ control.label }}
      </v-btn>
      <v-spacer/>
      <v-btn
          v-for="control in brightnessControl.right"
          color="neutral-lighten-1"
          :disabled="control.disabled"
          elevation="0"
          :key="control.label"
          size="small"
          @click="control.onClick">
        {{ control.label }}
      </v-btn>
    </v-card-actions>

    <!-- Presets -->
    <v-container v-if="presets.length > 0" class="pa-4">
      <v-list tile class="ma-0 mt-2 pa-0">
        <v-list-item class="py-1 pl-0">
          <v-list-item-title class="text-body-small text-capitalize">Presets</v-list-item-title>
        </v-list-item>
      </v-list>
      <v-btn
          v-for="preset in presets"
          block
          class="py-1 mx-0 mt-1 mb-2 preset"
          :color="getColor(preset.title, currentPresetTitle)"
          elevation="0"
          :key="preset.name"
          size="small"
          width="100%"
          max-width="575"
          @click="updateBrightnessPreset(preset)">
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
import {
  useBrightness,
  useDescribeBrightness,
  usePullBrightness,
  useUpdateBrightness
} from '@/traits/light/light.js';
import {computed} from 'vue';

const {blockActions} = useAuthSetup();
const props = defineProps({
  name: {
    type: String,
    default: ''
  }
});

const {value, loading: pullLoading} = usePullBrightness(() => props.name);
const {response: support, loading: supportLoading} = useDescribeBrightness(() => props.name);
const {updateBrightness, loading: updateLoading} = useUpdateBrightness(() => props.name);
const {levelStr, level, presets, currentPresetTitle} = useBrightness(value, support);

const loading = computed(() => pullLoading.value || supportLoading.value || updateLoading.value);

/**
 * @param {string} title
 * @param {string} currentPresetTitle
 * @return {string}
 */
function getColor(title, currentPresetTitle) {
  return title === currentPresetTitle ? 'primary' : 'neutral-lighten-1';
}

/**
 * Update the brightness level.
 *
 * @param {number} level
 * @return {Promise<Brightness.AsObject>}
 */
function updateBrightnessLevel(level) {
  return updateBrightness(level);
}

/**
 * Update the brightness preset.
 *
 * @param {LightPreset.AsObject} preset
 * @return {Promise<Brightness.AsObject>}
 */
function updateBrightnessPreset(preset) {
  return updateBrightness({preset: preset});
}

const brightnessControl = computed(() => {
  return {
    left: [
      {
        disabled: blockActions.value,
        label: 'On',
        onClick: () => updateBrightnessLevel(100)
      },
      {
        disabled: blockActions.value,
        label: 'Off',
        onClick: () => updateBrightnessLevel(0)
      }
    ],
    right: [
      {
        disabled: blockActions.value || level.value >= 100,
        label: 'Up',
        onClick: () => updateBrightnessLevel(level.value + 1)
      },
      {
        disabled: blockActions.value || level.value <= 0,
        label: 'Down',
        onClick: () => updateBrightnessLevel(level.value - 1)
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

.preset :deep(.v-btn__content) {
  max-width: 100%;
}
</style>
