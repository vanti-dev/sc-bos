import {closeResource, newResourceValue} from '@/api/resource';
import {pullDevicesMetadata} from '@/api/ui/devices';
import {computed, onMounted, onUnmounted, reactive} from 'vue';

const NO_FLOOR = '< no floor >';

/**
 *
 * @return {{
 * listOfFloors: ComputedRef<Array>
 * }}
 */
export default function() {
  // Create reactive resources and data
  const floorListResource = reactive(newResourceValue());

  // Fetch floor list on component mount
  onMounted(() => {
    const req = {includes: {fieldsList: ['metadata.location.floor']}, updatesOnly: false};
    pullDevicesMetadata(req, floorListResource);
  });

  // Close resource on component unmount
  onUnmounted(() => {
    closeResource(floorListResource);
  });

  // Computed property for the floor list
  const listOfFloors = computed(() => {
    const fieldCounts = floorListResource.value?.fieldCountsList || [];
    const floorFieldCounts = fieldCounts.find(v => v.field === 'metadata.location.floor');
    if (!floorFieldCounts || floorFieldCounts.countsMap.size <= 0) return [];

    const floors = floorFieldCounts.countsMap.map(([k]) => (k === '' ? NO_FLOOR : k));
    return floors;
  });

  return {
    listOfFloors
  };
}
