<template>
  <v-card :max-width="width" flat tile>
    <svg
        xmlns="http://www.w3.org/2000/svg"
        xml:space="preserve"
        fill-rule="evenodd"
        stroke-linejoin="round"
        stroke-miterlimit="2"
        clip-rule="evenodd"
        viewBox="-70.5 0 146 146">
      <path
          d="M5 0H0l1 20h3L5 0Z"
          v-for="i in segments"
          :key="i"
          :fill="fillColors[i-1]"
          :transform="transforms[i-1]"/>
    </svg>
    <span>{{ value }}</span>
  </v-card>
</template>

<script setup>
import {computed} from 'vue';

const center = [0, 73];

const props = defineProps({
  value: {
    type: Number,
    default: .52
  },
  min: {
    type: [Number, String],
    default: 0.
  },
  max: {
    type: [Number, String],
    default: 1
  },
  segments: {
    type: [Number, String],
    default: 20
  },
  width: {
    type: [Number, String],
    default: 146
  },
  color: {
    type: String,
    default: '#ff9947'
  }
});

const maxValue = computed(() => parseFloat(props.max));
const minValue = computed(() => parseFloat(props.min));

// how much is each segment worth?
const segValue = computed(() => {
  return (maxValue.value - minValue.value) / props.segments;
});

// list of transforms per segment
const transforms = computed(() => {
  const ts = [];
  for (let i = 0; i<props.segments; i++) {
    const t = [];
    const pos = i / (props.segments - 1);
    const val = minValue.value + i * segValue.value;

    t.push('translate(2.5 0)');

    const rot = -120 + pos * 240;
    t.push('rotate(' + rot + ' ' + center.join(' ') + ')');


    if (val >= props.value) {
      t.push('scale(0.5 0.7)');
    } else if (val < props.value - segValue.value) {
      // do nothing (scale 1 1)
    } else {
      // between this segment and the next - dynamic scale
      const scaleFactor = (props.value - val) / segValue.value;
      const s = [0.5 + scaleFactor * 0.5, 0.7 + scaleFactor * 0.3];
      t.push('scale(' + s.join(' ') + ')');
    }

    t.push('translate(-2.5 0)');

    ts.push(t.join(' '));
  }
  return ts;
});

// list of colors per segment
const fillColors = computed(() => {
  const cols = [];

  for (let i = 0; i<props.segments; i++) {
    const val = minValue.value + i * segValue.value;

    if (val >= props.value) {
      cols.push('#ffffff');
    } else {
      cols.push(props.color);
    }
  }
  return cols;
});


</script>

<style scoped>

</style>
