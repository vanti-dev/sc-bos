<template>
  <v-card :width="width" :height="heightComputed" flat tile class="gauge">
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
    <span class="text-h1 value">
      <slot/>
    </span>
    <span class="text-title gauge-title">
      <slot name="title"/>
    </span>
  </v-card>
</template>

<script setup>
import {computed, getCurrentInstance} from 'vue';

const center = [0, 73];

const props = defineProps({
  value: {
    type: Number,
    default: .5
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
    default: 25
  },
  width: {
    type: [Number, String],
    default: 155
  },
  height: {
    type: [Number, String],
    default: null
  },
  color: {
    type: String,
    default: 'primary'
  }
});

const maxValue = computed(() => parseFloat(props.max));
const minValue = computed(() => parseFloat(props.min));

const heightComputed = computed(() => props.height || props.width);

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

    // move 0 point for rotation and scaling to top, centre of path
    t.push('translate(2.5 0)');

    // calculate rotation based on position in sequence
    const rot = -120 + pos * 240;
    t.push('rotate(' + rot + ' ' + center.join(' ') + ')');

    // calculate scale based on value
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

    // move back again so it displays correctly
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
      // todo: allow override of 'off' colour
      cols.push('#ffffff');
    } else {
      if (props.color.startsWith('#')) {
        cols.push(props.color);
      } else {
        const theme = getCurrentInstance().proxy.$vuetify.theme.currentTheme;
        const c = props.color.split(' ');
        const col = theme[c[0]][c[1] ?? 'base'];
        cols.push(col);
      }
    }
  }
  return cols;
});


</script>

<style scoped>
.gauge {
  position: relative;
}
.value {
  display: flex;
  position: absolute;
  left:  15%;
  top:   15%;
  right: 15%;
  bottom: 20%;
  justify-content: center;
  align-items: center;
}
.gauge-title {
  display: flex;
  position: absolute;
  bottom: 0;
  left: 20%;
  right: 20%;
  text-align: center;
}
</style>
