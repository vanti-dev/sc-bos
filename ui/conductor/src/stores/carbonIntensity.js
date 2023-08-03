import {MINUTE, useNow} from '@/components/now';
import {acceptHMRUpdate, defineStore} from 'pinia';
import {computed, ref, watchEffect} from 'vue';

export const useCarbonIntensity = defineStore('carbonIntensity', () => {
  const {now} = useNow(30 * MINUTE);
  const latestData = ref(null);
  watchEffect(async () => {
    const date = now.value;
    latestData.value = await carbonIntensityApi.intensity.pt24h(date);
  });

  const processedData = computed(() => {
    if (!latestData.value?.data) return [];
    return latestData.value.data.map(d => ({
      from: new Date(d.from),
      to: new Date(d.to),
      intensity: d.intensity
    }));
  });

  return {
    last24Hours: processedData
  };
});

const baseUrl = 'https://api.carbonintensity.org.uk';
const carbonIntensityApi = {
  intensity: {
    date(date) {
      let url = `${baseUrl}/intensity/date`;
      if (date) url += `/${date.toISOString().substring(0, 10)}`;
      return fetch(url)
          .then(response => response.json());
    },
    fw24h(date) {
      if (!date) {
        date = new Date();
        date.setHours(date.getHours() - 24);
      }
      const url = `${baseUrl}/intensity/${date.toISOString()}/fw24h`;
      return fetch(url)
          .then(response => response.json());
    },
    pt24h(date) {
      if (!date) {
        date = new Date();
      }
      const url = `${baseUrl}/intensity/${date.toISOString()}/pt24h`;
      return fetch(url)
          .then(response => response.json());
    }
  }
};

// make sure to pass the right store definition, `useAuth` in this case.
if (import.meta.hot) {
  import.meta.hot.accept(acceptHMRUpdate(useCarbonIntensity, import.meta.hot));
}
