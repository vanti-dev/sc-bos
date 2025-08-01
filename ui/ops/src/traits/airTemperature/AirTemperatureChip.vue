<template>
  <div class="d-flex align-center justify-center overlap" :class="[layoutClass, variantClass]">
    <v-chip v-if="hasTemp || hasSetPoint" v-bind="chipAttrs">
      <span v-if="hasTemp" v-tooltip="'Current temperature'">{{ tempStr }}</span>
      <v-icon class="" v-if="hasTemp && hasSetPoint">mdi-chevron-right</v-icon>
      <span v-if="hasSetPoint" v-tooltip="'Target temperature (set point)'">{{ setPointStr }}</span>
    </v-chip>
    <v-avatar class="icon-container" v-bind="avatarAttrs" v-tooltip="avatarTooltip">
      <air-temperature-icon :size="props.size - 6" :icon="iconStr" :value="currentTempNum" v-model:comfort="comfort"/>
    </v-avatar>
  </div>
</template>

<script setup>
import {useAirTemperatureValues} from '@/traits/airTemperature/airTemperature.js';
import AirTemperatureIcon from '@/traits/airTemperature/AirTemperatureIcon.vue';
import {roundTo} from '@/util/number.js';
import {computed, ref} from 'vue';

const props = defineProps({
  setPoint: {
    type: [Number, String],
    default: null
  },
  currentTemp: {
    type: [Number, String],
    default: null
  },
  size: {
    type: [Number, String],
    default: 40
  },
  variant: {
    type: String,
    default: 'outlined'
  },
  layout: {
    type: String,
    default: 'start'
  },
  color: {
    type: String,
    default: ''
  }
});

const currentTempNum = computed(() => props.currentTemp);
const setPointNum = computed(() => props.setPoint);
const {
  hasTemp,
  tempStr,
  hasSetPoint,
  setPointStr
} = useAirTemperatureValues(currentTempNum, setPointNum);
const canTellDirection = computed(() => hasTemp.value && hasSetPoint.value);

const isHeating = computed(() => roundTo(props.currentTemp, 1) < roundTo(props.setPoint, 1));
const isCooling = computed(() => roundTo(props.currentTemp, 1) > roundTo(props.setPoint, 1));
const iconStr = computed(() => {
  if (!canTellDirection.value) return 'mdi-thermometer';
  if (isHeating.value) return 'mdi-fire';
  if (isCooling.value) return 'mdi-snowflake';
  return 'mdi-stop-circle';
});

// layout and sizing for the chip
const chipSize = computed(() => {
  const s = +props.size;
  if (s < 32) return 'x-small';
  if (s < 44) return 'small';
  if (s < 56) return 'default';
  if (s < 68) return 'large';
  return 'x-large';
});
const sizeVar = computed(() => {
  return props.size + 'px';
});
const layoutClass = computed(() => `air-temperature-chip--layout-${props.layout ?? 'start'}`);
const variantClass = computed(() => `air-temperature-chip--variant-${props.variant ?? 'outlined'}`);

const comfort = ref('normal');

const avatarAttrs = computed(() => {
  const attrs = {
    color: props.color,
    size: props.size,
    variant: props.variant
  };
  if (props.variant.startsWith('outlined')) {
    attrs.variant = 'outlined';
  }
  return attrs;
});
const avatarTooltip = computed(() => {
  let str = '';
  switch (comfort.value) {
    case 'cold':
      str = 'Temperature too cold';
      break;
    case 'normal':
      str = 'Temperature just right';
      break;
    case 'warm':
      str = 'Temperature too warm';
      break;
    case 'hot':
      str = 'Temperature very hot';
      break;
    default:
      str = 'Unknown';
      break;
  }
  if (isHeating.value) {
    str += ', is heating';
  } else if (isCooling.value) {
    str += ', is cooling';
  }
  return str;
});
const chipAttrs = computed(() => {
  const attrs = {
    size: chipSize.value,
    variant: props.variant,
    color: props.color
  };
  if (props.variant.startsWith('outlined')) {
    attrs.variant = 'outlined';
  }
  return attrs;
});

</script>

<style scoped lang="scss">
.overlap {
  --size: v-bind(sizeVar);
  --r: calc(var(--size) / 2);
}

.v-chip {
  mask-image: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100" fill="none"><circle r="50" cx="50" cy="50" fill="black"/></svg>'), linear-gradient(#fff, #fff);
  mask-size: auto var(--size);
  mask-repeat: no-repeat;
  mask-composite: exclude;

  font-size: calc(var(--size) * 0.4);
  height: auto;
  padding-block: .15em;

  overflow: visible;
  min-width: min-content;

  > * {
    // small adjustment to make the text appear more central vertically
    margin-top: .15em;
  }
}

.air-temperature-chip {
  &--layout {
    &-start, &-left {
      flex-direction: row-reverse;

      .v-chip {
        mask-position: calc(-1 * var(--r)) center, center;
        padding-left: calc(var(--r) + .6em);
        margin-left: calc(var(--r) * -1);
        border-left-color: transparent;
        border-bottom-left-radius: 0;
        border-top-left-radius: 0;
      }
    }

    &-end, &-right {
      flex-direction: row;
      justify-content: start;

      .v-chip {
        mask-position: calc(100% + var(--r)) center, center;
        padding-right: calc(var(--r) + .6em);
        margin-right: calc(var(--r) * -1);
        border-right-color: transparent;
        border-bottom-right-radius: 0;
        border-top-right-radius: 0;
      }
    }

    &-top {
      flex-direction: column-reverse;

      .v-chip {
        mask-position: center calc(-1 * var(--size) + .6em), center;
        margin-top: -.6em;
      }
    }

    &-bottom {
      flex-direction: column;

      .v-chip {
        mask-position: center calc(100% + var(--size) - .6em), center;
        margin-bottom: -.6em;
      }
    }
  }

  &--variant {
    &-outlined-filled {
      .v-chip, .v-avatar {
        color: rgb(var(--v-theme-on-surface));
        border-color: rgb(var(--v-theme-on-surface));
        background-color: rgb(var(--v-theme-surface));
      }
    }

    &-outlined-inverted {
      .v-chip, .v-avatar {
        color: rgb(var(--v-theme-surface));
        border-color: rgb(var(--v-theme-surface));
        background-color: rgb(var(--v-theme-on-surface));
      }
    }
  }
}
</style>
