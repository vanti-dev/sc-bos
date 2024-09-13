import {listHubNodes, pullHubNodes} from '@/api/sc/traits/hub.js';
import useCollection from '@/composables/collection.js';
import {computed, toValue} from 'vue';

/**
 * @param {MaybeRefOrGetter<Partial<ListHubNodesRequest.AsObject>>} request
 * @param {MaybeRefOrGetter<Partial<UseCollectionOptions>>?} options
 * @return {UseCollectionResponse<HubNode.AsObject>}
 */
export function useHubNodesCollection(request, options) {
  const normOptions = computed(() => {
    const optArg = toValue(options);
    return /** @type {UseCollectionOptions<HubNode>} */ {
      cmp: (a, b) => a.address.localeCompare(b.address),
      idFn: (item) => item.address,
      ...optArg
    };
  });
  const client = {
    async listFn(req, tracker) {
      const res = await listHubNodes(req, tracker);
      return {
        items: res.nodesList,
        nextPageToken: res.nextPageToken,
        totalSize: res.totalSize
      };
    },
    pullFn(req, resource) {
      pullHubNodes(req, resource);
    }
  };
  return useCollection(request, client, normOptions);
}
