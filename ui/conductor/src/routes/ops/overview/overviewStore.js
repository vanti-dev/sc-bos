import {defineStore} from 'pinia';
import {computed, ref} from 'vue';

export const useOverviewStore = defineStore('activeOverview', () => {
  const activeOverview = ref(null);

  const getActiveOverview = computed(() => activeOverview.value);

  return {activeOverview, getActiveOverview};
});
