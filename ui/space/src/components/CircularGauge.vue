<template>
  <v-sheet elevation="0" color="transparent" class="gauge">
    <svg xmlns="http://www.w3.org/2000/svg"
         xml:space="preserve"
         viewBox="0 0 100 100"
         :width="props.size"
         :height="props.size">
      <circle class="track" :style="trackStyles" :class="trackClasses"
              cx="50" cy="50" :r="radius"
              :stroke-width="trackStrokeWidth"
              pathLength="100" stroke-dasharray="100 100" :stroke-dashoffset="trackStrokeDashOffset"/>
      <circle class="remainder" :style="remainderStyles" :class="remainderClasses"
              cx="50" cy="50" :r="radius"
              :stroke-width="remainderStrokeWidth"
              pathLength="100" stroke-dasharray="100 100" :stroke-dashoffset="remainderStrokeDashOffset"/>
      <circle class="fill" :style="fillStyles" :class="fillClasses"
              cx="50" cy="50" :r="radius"
              :stroke-width="fillStrokeWidth"
              pathLength="100" stroke-dasharray="100 100" :stroke-dashoffset="fillStrokeDashOffset"/>
    </svg>
  </v-sheet>
</template>

<script setup>
import {computed, toRef} from 'vue';
import {useTextColor} from 'vuetify/lib/composables/color';

const props = defineProps({
  value: {
    type: Number,
    default: 30
  },
  min: {
    type: [Number, String],
    default: 0.
  },
  max: {
    type: [Number, String],
    default: 100
  },
  arcStart: {
    type: [Number, String],
    default: -120,
  },
  arcEnd: {
    type: [Number, String],
    default: 120
  },
  size: {
    type: [Number, String],
    default: 155
  },
  fillColor: {
    type: String,
    default: 'primary'
  },
  trackColor: {
    type: String,
    default: 'currentColor'
  },
  remainderColor: {
    type: String,
    default: 'transparent'
  },
  width: {
    type: [Number, String],
    default: 10
  },
  trackWidth: {
    type: [Number, String],
    default: undefined // defaults to width
  },
  fillWidth: {
    type: [Number, String],
    default: undefined // defaults to width
  },
  remainderWidth: {
    type: [Number, String],
    default: undefined // defaults to width
  },
  innerGap: {
    type: [Number, String],
    default: 5 // in degrees
  }
});

const _min = computed(() => parseFloat(props.min));
const _max = computed(() => parseFloat(props.max));
const _value = computed(() => Math.min(Math.max(parseFloat(props.value), _min.value), _max.value));
const _innerGap = computed(() => parseFloat(props.innerGap));

const trackStrokeWidth = computed(() => parseFloat(props.trackWidth ?? props.width));
const fillStrokeWidth = computed(() => parseFloat(props.fillWidth ?? props.width));
const remainderStrokeWidth = computed(() => parseFloat(props.remainderWidth ?? props.width));
const maxStrokeWidth = computed(() => Math.max(trackStrokeWidth.value, fillStrokeWidth.value, remainderStrokeWidth.value));

const radius = computed(() => 50 - maxStrokeWidth.value / 2);
const rotation = computed(() => parseFloat(props.arcStart) - 90);
const arcLength = computed(() => Math.abs(parseFloat(props.arcEnd) - parseFloat(props.arcStart)));
const arcPathLength = computed(() => arcLength.value / 360 * 100)
const trackStrokeDashOffset = computed(() => 100 - arcPathLength.value);
const fillPercent = computed(() => (_value.value - _min.value) / (_max.value - _min.value));
const fillStrokeDashOffset = computed(() => 100 - arcPathLength.value * fillPercent.value);
const remainderStrokeDashOffset = computed(() => {
  if (fillPercent.value === 0) {
    return 100 - arcPathLength.value; // remove inner gap
  }
  return 100 - (arcPathLength.value * (1 - fillPercent.value) - (_innerGap.value / 360 * 100));
});
const remainderRotation = computed(() => {
  if (fillPercent.value === 0) {
    return rotation.value; // remove inner gap
  }
  return rotation.value + (arcLength.value * fillPercent.value) + _innerGap.value;
});

const {textColorClasses: trackClasses, textColorStyles: trackStyles} = useTextColor(toRef(props, 'trackColor'));
const {textColorClasses: fillClasses, textColorStyles: fillStyles} = useTextColor(toRef(props, 'fillColor'));
const {
  textColorClasses: remainderClasses,
  textColorStyles: remainderStyles
} = useTextColor(toRef(props, 'remainderColor'));
</script>

<style scoped>
.gauge {
  display: grid;
  align-items: center;
  justify-items: center;
}

svg {
  /* get rid of the descendent for baseline alignment */
  vertical-align: top;
}

svg circle {
  transform-origin: center;
  transform: rotate(calc(v-bind(rotation) * 1deg));
  fill: transparent;
  stroke: currentColor;
}

circle.remainder {
  transform: rotate(calc(v-bind(remainderRotation) * 1deg));
}

.track {
  stroke: color-mix(in srgb, currentColor, transparent calc(100% * (1 - var(--v-border-opacity))));
}
</style>
