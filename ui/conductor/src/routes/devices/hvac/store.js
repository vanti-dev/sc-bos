import {defineStore} from 'pinia';
import {computed, ref} from 'vue';

import hvacJson from './hvacData.json';

export const useHvacStore = defineStore('hvac', () => {
  const deviceList = ref(hvacJson.devices);
  const data = hvacJson.data;

  const getDevice = computed((state) => {
    return (deviceId) => {
      for (const key in deviceList.value) {
        if (deviceList.value.hasOwnProperty(key)) {
          const d = deviceList.value[key];
          if (d.deviceId === deviceId) {
            return d;
          }
        }
      }
      return {};
    };
  });

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
      if (data.hasOwnProperty(deviceId)) {
        return data[deviceId].currentTemp;
      }
      return 0;
    };
  });

  return {
    deviceList,
    getSetPoint,
    getCurrentTemp,
    getDevice
  };
});
