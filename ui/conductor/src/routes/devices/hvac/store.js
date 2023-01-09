import {defineStore} from 'pinia';
import {computed, ref} from 'vue';

import hvacJson from './hvacData.json';

export const useHvacStore = defineStore('hvac', () => {
  const deviceList = ref(hvacJson.devices);
  const data = hvacJson.data;

  /**
   *
   * @param {string} deviceId
   * @return {number}
   */
  const getSetPoint = computed((state) => {
    return (deviceId) => {
      if (data.hasOwnProperty(deviceId)) {
        return data[deviceId].setPoint;
      }
      return 0;
    };
  });

  /**
   *
   * @param {string} deviceId
   * @return {number}
   */
  const getCurrentTemp = computed(() => {
    return (deviceId) => {
      if (data.value.hasOwnProperty(deviceId)) {
        return data[deviceId].currentTemp;
      }
      return 0;
    };
  });

  return {
    deviceList,
    getSetPoint,
    getCurrentTemp
  };
});
