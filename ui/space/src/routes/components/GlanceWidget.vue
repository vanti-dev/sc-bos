<template>
  <div class="glance-widget mt-2" @click="emits('adminClick')">
    <div class="timedate">
      <h1 class="time">{{ timeStr }}</h1>
      <span class="text-h5">{{ dateStr }}</span>
    </div>
    <div class="weather" v-if="showWeather">
      <span class="text-body-large">{{ outdoorTemperature }}&deg;</span>
      <v-icon end>mdi-weather-partly-cloudy</v-icon>
    </div>
  </div>
</template>

<script setup>
import {onMounted, onUnmounted, ref} from 'vue';

const timeStr = ref('');
const dateStr = ref('');
const outdoorTemperature = ref('18');

const emits = defineEmits(['adminClick']);

defineProps({
  showWeather: Boolean
});

/**
 */
function updateDateTime() {
  const now = new Date();
  timeStr.value = new Intl.DateTimeFormat('en-GB', {
    hour: 'numeric',
    minute: 'numeric'
  }).format(now);
  dateStr.value = new Intl.DateTimeFormat('en-GB', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric'
  }).format(now);
}

let timeInterval;
onMounted(() => {
  updateDateTime();
  timeInterval = setInterval(updateDateTime, 1000);
});
onUnmounted(() => {
  clearInterval(timeInterval);
});

</script>

<style scoped>
.glance-widget {
  position: relative;
}

.timedate {
  position: absolute;
  left: 0;
  top: 0;
}

.time {
  color: transparent;
  font-size: 4.5em;
  line-height: 1em;
  -webkit-text-stroke: 3px #fff;
  -webkit-font-smoothing: subpixel-antialiased;
  letter-spacing: 0.02em;
  font-weight: 900;
}

.weather {
  position: absolute;
  right: 0;
  top: -10px;
}

.weather .v-icon {
  font-size: 100px;
}
</style>
