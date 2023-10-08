import {reactive, ref, computed, onMounted, onUnmounted, watch} from 'vue';

import {useErrorStore} from '@/components/ui-error/error';
import {useDevicesStore} from '@/routes/devices/store';

import useFloors from '@/composables/useFloors';

const NO_FLOOR = '< no floor >';

/**
 *
 * @param {Object} props
 * @return {{
 * floorList: ComputedRef<Array>,
 * filterFloor: Ref<string>,
 * search: Ref<string>,
 * query: ComputedRef<Object>,
 * devicesData: ComputedRef<Array>
 * }}
 */
export default function(props) {
  const devicesStore = useDevicesStore();
  const errorStore = useErrorStore();
  const {listOfFloors} = useFloors();

  const filterFloor = ref('All');
  const search = ref('');

  // Computed property for the floor list
  const floorList = computed(() => {
    const floors = listOfFloors.value;
    floors.unshift('All');
    return floors;
  });

  // Create reactive collection
  const collection = reactive(devicesStore.newCollection());
  collection.needsMorePages = true; // Todo: Connect with paging logic instead

  // Computed property for the query object
  const query = computed(() => {
    const q = {conditionsList: []};
    if (search.value) {
      const words = search.value.split(/\s+/);
      q.conditionsList.push(...words.map(word => ({stringContainsFold: word})));
    }
    if (props.subsystem.toLowerCase() !== 'all') {
      q.conditionsList.push({field: 'metadata.membership.subsystem', stringEqualFold: props.subsystem});
    }
    switch (filterFloor.value.toLowerCase()) {
      case 'all':
        // no filter
        break;
      case NO_FLOOR:
        q.conditionsList.push({field: 'metadata.location.floor', stringEqualFold: ''});
        break;
      default:
        q.conditionsList.push({field: 'metadata.location.floor', stringEqualFold: filterFloor.value});
        break;
    }
    return q;
  });

  // Watch for changes to the query object and fetch new device list
  watch(query, () => collection.query(query.value), {deep: true, immediate: true});

  // UI error handling
  let unwatchErrors;
  onMounted(() => {
    unwatchErrors = errorStore.registerCollection(collection);
  });

  onUnmounted(() => {
    if (unwatchErrors) unwatchErrors();
    collection.reset(); // stop listening when the component is unmounted
  });

  // Computed property for the filtered table data
  const devicesData = computed(() => {
    return Object.values(collection.resources.value).filter(props.filter);
  });

  return {
    floorList,
    filterFloor,
    search,
    query,
    devicesData
  };
}
