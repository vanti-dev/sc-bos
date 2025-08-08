import {timestampToDate} from '@/api/convpb.js';
import {listReports} from '@/api/ui/reports.js';
import useCollection from '@/composables/collection.js';
import {cmpDesc} from '@/util/date.js';
import {computed, toValue} from 'vue';


/**
 * @param {MaybeRefOrGetter<Partial<ListReportsRequest.AsObject>>} request
 * @param {MaybeRefOrGetter<Partial<UseCollectionOptions<Report.AsObject>>>?} options
 * @return {UseCollectionResponse<Report.AsObject>}
 */
export function useReportsCollection(request, options) {
  const normOpts = computed(() => {
    return {
      cmp: (a, b) => cmpDesc(timestampToDate(a.createTime), timestampToDate(b.createTime)),
      ...toValue(options)
    };
  });
  const client = {
    async listFn(req, tracker) {
      const res = await listReports(req, tracker);
      return {
        items: res.reportsList,
        nextPageToken: res.nextPageToken,
        totalSize: res.totalSize
      };
    },
    pullFn() {
      // No pull function for reports
    }
  };
  return useCollection(request, client, normOpts);
}