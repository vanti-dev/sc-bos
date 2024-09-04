import {listDevices, pullDevices} from '@/api/ui/devices.js';
import useCollection from '@/composables/collection.js';
import {computed, toValue} from 'vue';

/**
 * @param {MaybeRefOrGetter<Partial<ListDevicesRequest.AsObject>>}request
 * @param {MaybeRefOrGetter<Partial<UseCollectionOptions>>?} options
 * @return {UseCollectionResponse<Device.AsObject>}
 */
export function useDevicesCollection(request, options) {
  const normOptions = computed(() => {
    const optArg = toValue(options);
    return {
      cmp: (a, b) => a.name.localeCompare(b.name),
      ...optArg
    };
  });
  const client = {
    async listFn(req, tracker) {
      const res = await listDevices(req, tracker);
      return {
        items: res.devicesList,
        nextPageToken: res.nextPageToken,
        totalSize: res.totalSize
      };
    },
    pullFn(req, resource) {
      pullDevices(req, resource);
    }
  };
  return useCollection(request, client, normOptions);
}
