import {defineStore} from 'pinia';
import {computed, reactive, ref} from 'vue';

import {listDevices} from '@/api/ui/devices';
import {newActionTracker, newResourceCollection} from '@/api/resource';
import {Collection} from '@/util/query';

export const useDevicesStore = defineStore('devices', () => {
  // holds all the devices we can show
  const deviceList = reactive(/** @type {ResourceCollection<Device.AsObject, Device>} */newResourceCollection());
  // tracks the fetching of a single page
  const fetchingPage = reactive(/** @type {ActionTracker<ListDevicesResponse.AsObject>} */ newActionTracker());

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
    newCollection
  };
});
