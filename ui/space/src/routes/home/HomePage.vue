<template>
  <div class="pa-8 d-flex flex-column fill-height">
    <glance-widget @admin-click="handle10Click"/>
    <v-spacer/>
    <light-card :name="zoneId"/>
    <air-temperature-card :name="zoneId" class="mt-6"/>
    <img src="/img/sc-fav.svg" height="115" class="pt-14 pr-3 align-self-end" alt="Smart Core logo">
    <NotificationToast :show-alert="alertMessage.show" :message="alertMessage.message"/>
  </div>
</template>

<script setup>
import NotificationToast from '@/components/NotificationToast.vue';
import AirTemperatureCard from '@/routes/components/AirTemperatureCard.vue';
import GlanceWidget from '@/routes/components/GlanceWidget.vue';
import LightCard from '@/routes/components/LightCard.vue';
import {useConfigStore} from '@/stores/config';
import {computed, ref} from 'vue';

const configStore = useConfigStore();
const zoneId = computed(() => configStore.zoneId);


// ----------------- 10 click safety feature ----------------- //
const clickCount = ref(0); // Define a ref to keep track of the click count
const alertMessage = computed(() => {
  if (clickCount.value >= 5 && clickCount.value < 10) {
    return {
      show: true,
      message: `${10 - clickCount.value} clicks left for admin menu.`
    };
  } else return {show: false, message: ''};
});
let clickTimeout; // Define a timeout to clear when the component is unmounted

// 10 click safety feature for logging the user out and returning to the setup page
const handle10Click = () => {
  clearTimeout(clickTimeout); // Clear any existing timeout

  clickCount.value += 1; // Increment the click count on each click

  if (clickCount.value === 10) {
    configStore.reconfigure();
    // Reset the click count
    clickCount.value = 0;
  }

  // Set a timeout to reset the click count after 1 second if no new click occurs
  clickTimeout = setTimeout(() => {
    clickCount.value = 0;
  }, 1000);
};
</script>
