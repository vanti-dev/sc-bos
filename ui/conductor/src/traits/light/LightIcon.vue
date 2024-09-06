<script setup>
import {computed, ref} from 'vue';

const props = defineProps({
  level: {
    type: Number,
    default: 100
  },
  size: {
    type: [Number, String],
    default: 24
  },
  density: {
    type: String,
    default: 'default',
    validator(value) {
      return ['default', 'compact'].includes(value);
    }
  },
  onColor: {
    type: String,
    default: '#FBBF24'
  },
  offColor: {
    type: String,
    default: 'rgba(255, 255, 255, 0.2)'
  }
});
const densityOpts = {
  'default': {
    r: 1 / 3.7,
    lineWidth: 1 / 50,
    pad: 1 / 20,
    numRays: 12,
    rayWidth: 1 / 11
  },
  'compact': {
    r: 1 / 2.9,
    lineWidth: 1 / 15,
    pad: 1 / 35,
    numRays: 8,
    rayWidth: 1 / 8
  }
};
const attrs = computed(() => densityOpts[props.density]);


// attrs for tweaking the icon shape/size
const boxSize = ref(1200); // determines accuracy
const r = computed(() => boxSize.value * attrs.value.r); // outer radius of the circle
const lineWidth = computed(() => boxSize.value * attrs.value.lineWidth); // the ring width
const pad = computed(() => boxSize.value * attrs.value.pad); // space between rays and ring/edge
const numRays = computed(() => attrs.value.numRays);
const rayWidth = computed(() => boxSize.value * attrs.value.rayWidth);

// svg string properties
const sizePx = computed(() => {
  if (typeof props.size === 'string' && props.size.endsWith('px')) return props.size;
  return `${props.size}px`;
});
const viewBox = computed(() => `0 0 ${boxSize.value} ${boxSize.value}`);
const c = computed(() => (boxSize.value / 2).toFixed(3)); // both cx and cy
const lineRadius = computed(() => (r.value - lineWidth.value / 2).toFixed(3));
const fillRadius = computed(() => ((r.value) / 2).toFixed(3));
const fillWidth = computed(() => r.value.toFixed(3));
const maxY = computed(() => boxSize.value / 2 - r.value - pad.value - rayWidth.value/2);
const y1 = computed(() => Math.min(rayWidth.value / 2, maxY.value).toFixed(3));
const y2 = computed(() => maxY.value.toFixed(3));
const rayStrokeWidth = computed(() => rayWidth.value.toFixed(3));
const onColorVal = computed(() => props.onColor);
const offColorVal = computed(() => props.offColor);

// related to the level and how to display it
const levelPercent = computed(() => props.level);
const fillDashOffset = computed(() => (100 - levelPercent.value).toFixed(3));
const numRaysOn = computed(() => {
  return Math.ceil(numRays.value * levelPercent.value / 100);
});
const rayClass = (i) => i < numRaysOn.value ? 'stroke-on' : 'stroke-off';
const rayAngle = (i) => (360 / numRays.value * i).toFixed(1);
</script>

<template>
  <svg :width="sizePx" :height="sizePx" :viewBox="viewBox" fill="none" xmlns="http://www.w3.org/2000/svg">
    <circle :cx="c" :cy="c" :r="lineRadius" fill="none"
            :class="levelPercent < 0.01 ? 'stroke-off': 'stroke-on'" :stroke-width="lineWidth"/>
    <circle :cx="c" :cy="c" fill="none"
            class="stroke-on" :stroke-width="fillWidth" :r="fillRadius"
            :transform="`rotate(-90,${c},${c})`"
            pathLength="99.9" stroke-dasharray="100" :stroke-dashoffset="fillDashOffset"/>
    <line v-for="i in numRays" :key="i"
          :x1="c" :y1="y1" :x2="c" :y2="y2"
          :class="rayClass(i-1)"
          :stroke-width="rayStrokeWidth"
          stroke-linecap="round"
          :transform="`rotate(${rayAngle(i-1)},${c},${c})`"/>
  </svg>
</template>

<style scoped>
circle, line {
  transition: all 0.5s;
  transition-property: stroke, fill, stroke-dashoffset;
}

.fill-on {
  fill: v-bind(onColorVal);
}

.fill-off {
  fill: v-bind(offColorVal);
}

.stroke-on {
  stroke: v-bind(onColorVal);
}

.stroke-off {
  stroke: v-bind(offColorVal);
}
</style>
