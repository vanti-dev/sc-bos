import {timestampToDate} from '@/api/convpb.js';
import {listWasteRecords, pullWasteRecords} from '@/api/ui/waste.js';
import useCollection from '@/composables/collection.js';
import {cmpDesc} from '@/util/date.js';
import {computed, toValue} from 'vue';

/**
 * @param {MaybeRefOrGetter<Partial<ListWasteRecordsRequest.AsObject>>} request
 * @param {MaybeRefOrGetter<Partial<UseCollectionOptions<WasteRecord.AsObject>>>?} options
 * @return {UseCollectionResponse<WasteRecord.AsObject>}
 */
export function useWasteRecordsCollection(request, options) {
  const normOpts = computed(() => {
    return {
      cmp: (a, b) => cmpDesc(timestampToDate(a.wasteCreateTime), timestampToDate(b.wasteCreateTime)),
      ...toValue(options)
    };
  });
  const client = {
    async listFn(req, tracker) {
      const res = await listWasteRecords(req, tracker);
      return {
        items: res.wasterecordsList,
        nextPageToken: res.nextPageToken,
        totalSize: res.totalSize
      };
    },
    pullFn(req, resource) {
      pullWasteRecords(req, resource);
    }
  };
  return useCollection(request, client, normOpts);
}