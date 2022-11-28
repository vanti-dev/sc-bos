import { defineStore } from "pinia";
import { computed, ref } from "vue";

export const useSmartCoreStore = defineStore("smartCoreStore", () => {

  //state for availability of Smart Core
  const smartCoreStatus = ref("Online");

  return {
    smartCoreStatus
  };
});
