import {defineStore} from 'pinia';
import {reactive, ref} from 'vue';

import {listDevices, getDevicesMetadata} from '@/api/ui/devices';
import {newResourceCollection} from '@/api/resource';
import {Collection} from '@/util/query';

export const useDevicesStore = defineStore('devices', () => {
  // holds all the devices we can show
  const deviceList = reactive(/** @type {ResourceCollection<Device.AsObject, Device>} */newResourceCollection());
  const subSystems = ref({});

  /**
   * @param {ActionTracker<GetDevicesMetadataResponse.AsObject>} tracker
   * @return {Collection}
   */
  async function fetchDeviceSubsystemCounts(tracker) {
    // Fetch devices data
    const devices = await getDevicesMetadata({includes: {fieldsList: ['metadata.membership.subsystem']}}, tracker);

    // Extract the countsMap array from the devices object and set it to a var
    const countsMap = devices?.fieldCountsList[0].countsMap;

    // Format countsMap array -> object with the keys/values
    const subs = countsMap.reduce((accumulator, [key, value]) => {
      if (key) accumulator[key] = value;
      else accumulator['noType'] = value;
      return accumulator;
    }, {});

    // Reconstruct object to include subs & totalCount
    subSystems.value = {
      subs,
      totalCount: devices.totalCount
    };
  }

  /**
   *
   * @return {Collection}
   */
  function newCollection() {
    const listFn = async (query, tracker, pageToken, recordFn) => {
      const page = await listDevices({query, pageToken, pageSize: 100}, tracker);
      for (const device of page.devicesList) {
        recordFn(device, device.name);
      }
      return page.nextPageToken;
    };
    return new Collection(listFn);
  }


  return {
    deviceList,
    subSystems,

    fetchDeviceSubsystemCounts,
    newCollection
  };
});
