import {closeResource, newResourceValue} from '@/api/resource.js';
import {listDevices, pullDevices, pullDevicesMetadata} from '@/api/ui/devices.js';
import useCollection from '@/composables/collection.js';
import {watchResource} from '@/util/traits.js';
import {computed, reactive, toRefs, toValue} from 'vue';

/** @import {MaybeRefOrGetter, ToRefs, Ref} from 'vue' */

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

/**
 * @param {import('vue').MaybeRefOrGetter<string|string[]|PullDevicesMetadataRequest.AsObject>} query
 * @param {import('vue').MaybeRefOrGetter<{paused?: boolean}>?} options
 * @return {import('vue').ToRefs<ResourceValue<DevicesMetadata.AsObject, PullDevicesMetadataResponse>>}
 */
export function usePullDevicesMetadata(query, options) {
  const normQuery = computed(() => {
    const queryArg = toValue(query);
    if (typeof queryArg === 'string') {
      return {includes: {fieldsList: [queryArg]}};
    }
    if (Array.isArray(queryArg)) {
      return {includes: {fieldsList: queryArg}};
    }
    // we could check for the correct type here, but lets assume people know what they're doing
    return queryArg;
  });

  const resource = reactive(
      /** @type {ResourceValue<DevicesMetadata.AsObject, PullDevicesMetadataResponse>} */
      newResourceValue());

  watchResource(normQuery, () => toValue(options)?.paused ?? false, (req) => {
    pullDevicesMetadata(req, resource);
    return () => closeResource(resource);
  });

  return toRefs(resource);
}

/**
 * @param {import('vue').MaybeRefOrGetter<DevicesMetadata.AsObject>} value
 * @param {import('vue').MaybeRefOrGetter<string>} field
 * @return {{
 *  counts: import('vue').Ref<Array<[string, number]>>,
 *  countsMap: import('vue').Ref<Record<string, number>>,
 *  keys: import('vue').Ref<string[]>
 * }}
 */
export function useDevicesMetadataField(value, field) {
  const counts = computed(() => {
    const _value = toValue(value);
    const _field = toValue(field);
    return _value?.fieldCountsList?.find(v => v.field === _field)?.countsMap;
  });
  const countMap = computed(() => {
    const mapArr = counts.value || [];
    if (mapArr.length === 0) return {};
    return mapArr.reduce((acc, [k, v]) => {
      acc[k] = v;
      return acc;
    }, {});
  });
  const keys = computed(() => {
    return (counts.value ?? []).map(([k]) => k);
  });

  return {
    counts,
    countMap,
    keys
  };
}
