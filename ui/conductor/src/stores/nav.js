import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useNavStore = defineStore('nav', () => {
  const drawer = ref(true);
  const miniVariant = ref(true);
  const drawerWidth = ref(60);
  const pinDrawer = ref(false);
  return {
    drawer,
    miniVariant,
    drawerWidth,
    pinDrawer
  };
});
