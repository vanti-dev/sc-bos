import {ServiceNames} from '@/api/ui/services';
import {useServicesCollection} from '@/composables/services.js';
import {useUserConfig} from '@/stores/userConfig.js';
import {computed} from 'vue';

/**
 * @return {{
 * zoneCollection: UseCollectionResponse,
 * zoneList: string[]
 * }}
 */
export default function() {
  const userConfig = useUserConfig();
  const serviceName = computed(() => `${userConfig.node?.name}/${ServiceNames.Zones}`);
  const zoneCollection = useServicesCollection(serviceName, computed(() => ({
    paused: !userConfig.node?.name,
    wantCount: -1 // there's no server search features, so we have to get them all and do it client side
  })));

  const zoneListWithDetails = computed(() => {
    return Object.values(zoneCollection.value?.resources?.value ?? []);
  });

  const zoneList = computed(() => {
    return Object.values(zoneCollection.items.value)
        .map((zone) => {
          return zone.id;
        })
        .sort();
  });

  return {
    zoneCollection,
    zoneList,
    zoneListWithDetails
  };
}
