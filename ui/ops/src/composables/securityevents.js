import {timestampToDate} from '@/api/convpb.js';
import {listSecurityEvents, pullSecurityEvents} from '@/api/ui/securityevents.js';
import useCollection from '@/composables/collection.js';
import {cmpDesc} from '@/util/date.js';
import {computed, toValue} from 'vue';

/**
 * @param {MaybeRefOrGetter<Partial<ListSecurityEventsRequest.AsObject>>} request
 * @param {MaybeRefOrGetter<Partial<UseCollectionOptions<SecurityEvent.AsObject>>>?} options
 * @return {UseCollectionResponse<SecurityEvent.AsObject>}
 */
export function useSecurityEventsCollection(request, options) {
  const normOpts = computed(() => {
    return {
      cmp: (a, b) => cmpDesc(timestampToDate(a.securityEventTime), timestampToDate(b.securityEventTime)),
      ...toValue(options)
    };
  });
  const client = {
    async listFn(req, tracker) {
      const res = await listSecurityEvents(req, tracker);
      return {
        items: res.securityEventsList,
        nextPageToken: res.nextPageToken,
        totalSize: res.totalSize
      };
    },
    pullFn(req, resource) {
      pullSecurityEvents(req, resource);
    }
  };
  return useCollection(request, client, normOpts);
}
