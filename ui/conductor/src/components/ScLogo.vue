<template>
  <!-- eslint-disable max-len -->
  <svg
      viewBox="0 0 226 226"
      style="fill-rule:evenodd;clip-rule:evenodd;stroke-linecap:round;stroke-linejoin:round;stroke-miterlimit:10;"
      :class="{noFill, colorFill}">
    <g class="fill">
      <circle class="tr" cx="158.748" cy="28.757" r="21.094"/>
      <circle class="bl" cx="32.185" cy="179.929" r="24.609"/>
      <circle class="c" cx="113.044" cy="113.132" r="35.156"/>
      <circle class="tl" cx="56.794" cy="49.851" r="14.063"/>
      <circle class="br" cx="193.904" cy="193.991" r="24.609"/>
      <path
          d="M179.033,34.484c-3.98,5.255 -10.203,8.349 -16.795,8.349c-11.556,-0 -21.065,-9.509 -21.065,-21.066c-0,-1.95 0.271,-3.892 0.805,-5.768c-2.831,3.685 -4.366,8.204 -4.366,12.851c-0,11.571 9.522,21.093 21.093,21.093c9.451,0 17.803,-6.351 20.328,-15.459Zm-164.426,134.898c0.003,-2.684 0.448,-5.349 1.315,-7.889c-5.33,4.67 -8.39,11.422 -8.39,18.508c-0,13.501 11.109,24.61 24.609,24.61c10.552,-0 19.983,-6.787 23.335,-16.793c-4.493,3.975 -10.289,6.171 -16.288,6.171c-13.485,-0 -24.581,-11.096 -24.581,-24.581c-0,-0.009 -0,-0.017 -0,-0.026Zm73.828,-70.312c-0.001,-4.584 0.897,-9.123 2.644,-13.36c-8.336,6.667 -13.197,16.78 -13.197,27.454c0,19.286 15.871,35.156 35.157,35.156c14.24,0 27.134,-8.651 32.532,-21.829c-6.237,5.009 -14.001,7.741 -22.001,7.741c-19.274,0 -35.135,-15.86 -35.135,-35.135c-0,-0.009 -0,-0.018 0,-0.027Zm-28.125,-42.188c-0.05,0.001 -0.1,0.001 -0.15,0.001c-7.633,-0 -13.914,-6.282 -13.914,-13.915c-0,-0.884 0.084,-1.766 0.251,-2.634c-2.407,2.595 -3.745,6.006 -3.745,9.544c-0,7.7 6.335,14.035 14.034,14.035c6.775,-0 12.631,-4.906 13.818,-11.577c-2.642,2.889 -6.379,4.54 -10.294,4.546Zm116.016,126.563c0.003,-2.684 0.447,-5.35 1.314,-7.889c-5.329,4.67 -8.39,11.422 -8.39,18.508c0,13.5 11.11,24.609 24.61,24.609c10.552,0 19.983,-6.787 23.335,-16.792c-4.493,3.974 -10.289,6.17 -16.288,6.17c-13.485,0 -24.581,-11.096 -24.581,-24.581c-0,-0.008 -0,-0.017 -0,-0.025Z"
          style="fill:#000064;fill-opacity:0.16;fill-rule:nonzero;"/>
      <path
          d="M128.626,94.035c-11.198,-10.772 -24.74,-3.164 -22.332,1.055c1.821,4.05 7.875,6.616 10.853,10.42c3.804,2.978 6.37,9.028 10.42,10.853c4.223,2.394 11.845,-11.131 1.059,-22.328Z"
          style="fill:#fff;fill-opacity:0.4;fill-rule:nonzero;"/>
    </g>
    <g class="stroke">
      <path
          d="M147.884,46.834l-18.358,35.238m-44.916,51.781l-32.903,31.096"
          style="fill:none;fill-rule:nonzero;;stroke-width:3.52px;"/>
      <circle cx="158.748" cy="28.757" r="21.094" style="fill:none;stroke-width:3.52px;"/>
      <circle cx="32.185" cy="179.929" r="24.609" style="fill:none;;stroke-width:3.52px;"/>
      <circle cx="113.044" cy="113.132" r="35.156" style="fill:none;;stroke-width:3.52px;"/>
      <path
          d="M113.044,92.038c11.572,0 21.094,9.522 21.094,21.094m4.142,24.472l37.223,40.085m-109.15,-117.544l24.149,26.008"
          style="fill:none;fill-rule:nonzero;stroke-width:3.52px;"/>
      <circle cx="56.794" cy="49.851" r="14.063" style="fill:none;;stroke-width:3.52px;"/>
      <circle cx="193.904" cy="193.991" r="24.609" style="fill:none;;stroke-width:3.52px;"/>
    </g>
  </svg>
  <!-- eslint-enable max-len -->
</template>

<script setup>
import {computed} from 'vue';

const props = defineProps({
  outline: {
    type: String,
    default: '#4c4c4c'
  },
  fill: {
    type: String,
    default: 'none'
  }
});
const noFill = computed(() => props.fill === 'none');
const colorFill = computed(() => Boolean(!noFill.value && props.fill));
</script>

<style scoped>
.noFill .fill {
  display: none;
}

.colorFill .fill circle {
  fill: v-bind(fill);
}

.fill .tl {
  fill: #ff9800;
}

.fill .c {
  fill: #00bcd4;
}

.fill .br {
  fill: #ff9800;
}

.fill .tr,
.fill .bl {
  fill: #cddc39;
}

.stroke > * {
  stroke: v-bind(outline);
}
</style>
