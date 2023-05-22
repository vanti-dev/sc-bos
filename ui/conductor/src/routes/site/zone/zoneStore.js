import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useZoneStore = defineStore('zone', () => {
  //
  // State
  const zoneCollection = ref({});
  const activeZone = ref('');

  return {
    zoneCollection,
    activeZone
  };
});
