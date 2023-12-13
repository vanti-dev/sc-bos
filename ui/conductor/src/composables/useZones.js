import {computed, onUnmounted, ref, watch} from 'vue';
import {storeToRefs} from 'pinia';

import {ServiceNames} from '@/api/ui/services';
import {usePageStore} from '@/stores/page';
import {useServicesStore} from '@/stores/services';

/**
 * @return {{
 * zoneCollection: Collection,
 * zoneList: string[]
 * }}
 */
export default function() {
  const servicesStore = useServicesStore();
  const pageStore = usePageStore();

  const {sidebarNode} = storeToRefs(pageStore);
  const zoneCollection = ref({});

  const zoneListWithDetails = computed(() => {
    return Object.values(zoneCollection.value?.resources?.value ?? []);
  });

  const zoneList = computed(() => {
    return Object.values(zoneCollection.value?.resources?.value ?? [])
        .map((zone) => {
          return zone.id;
        })
        .sort();
  });

  // Watch for changes to the sidebar node
  watch(
      sidebarNode,
      async () => {
        zoneCollection.value = servicesStore.getService(
            ServiceNames.Zones,
            await sidebarNode.value.commsAddress,
            await sidebarNode.value.commsName
        ).servicesCollection;

        // todo: this causes us to load all pages, connect with paging logic instead
        // - although we might want it in this case
        zoneCollection.value.needsMorePages = true;
      },
      {immediate: true}
  );


  // Watch for changes to the zone collection
  watch(zoneCollection, () => {
    zoneCollection.value.query(ServiceNames.Zones);
  });

  // Clear the collection when the component is unmounted
  onUnmounted(() => zoneCollection.value.reset());

  return {
    zoneCollection,
    zoneList,
    zoneListWithDetails
  };
}
