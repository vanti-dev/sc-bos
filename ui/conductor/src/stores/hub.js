import {closeResource, newActionTracker, newResourceCollection} from '@/api/resource';
import {listHubNodes, pullHubNodes} from '@/api/sc/traits/hub';
import {useAppConfigStore} from '@/stores/app-config';
import {defineStore} from 'pinia';
import {reactive, set, watch} from 'vue';

export const useHubStore = defineStore('hub', () => {
  const appConfig = useAppConfigStore();
  const nodesList = reactive(newResourceCollection());

  watch(() => appConfig.config, async config => {
    closeResource(nodesList);
    console.debug('config changed', config);
    if (config?.hub) {
      pullHubNodes(nodesList);
      try {
        const nodes = await listHubNodes(newActionTracker());
        for (const node of nodes.nodesList) {
          set(nodesList.value, node.name, node);
        }
      } catch (e) {
        console.warn('Error fetching first page', e);
      }
    }
  }, {immediate: true});

  return {
    nodesList
  };
});
