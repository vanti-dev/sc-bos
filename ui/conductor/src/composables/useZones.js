import {ServiceNames} from '@/api/ui/services';
import {useServicesCollection} from '@/composables/services.js';
import {useServicesStore} from '@/stores/services';
import {computed} from 'vue';

/**
 * @return {{
 * zoneCollection: Collection,
 * zoneList: string[]
 * }}
 */
export default function() {
  const servicesStore = useServicesStore();
  const serviceName = computed(() => `${servicesStore.node?.name}/${ServiceNames.Zones}`);
  const zoneCollection = useServicesCollection(serviceName, computed(() => ({
    paused: !servicesStore.node?.name,
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
