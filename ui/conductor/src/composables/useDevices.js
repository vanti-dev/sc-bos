import {useErrorStore} from '@/components/ui-error/error';
import useFloors from '@/composables/useFloors';
import {useDevicesStore} from '@/routes/devices/store';
import {computed, onMounted, onUnmounted, reactive, ref, toValue, watch} from 'vue';

const NO_FLOOR = '< no floor >';

/**
 * @typedef {Object} UseDevicesOptions
 * @property {string} subsystem
 * @property {(value: Device.AsObject, index: number, array: Device.AsObject[]) => boolean} filter
 */

/**
 *
 * @param {MaybeRefOrGetter<Partial<UseDevicesOptions>>} props
 * @return {{
 * floorList: import('vue').ComputedRef<Array>,
 * filterFloor: import('vue').Ref<string>,
 * search: import('vue').Ref<string>,
 * query: import('vue').ComputedRef<Object>,
 * devicesData: import('vue').ComputedRef<Array>
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
    return ['All', ...listOfFloors.value];
  });

  // Create reactive collection
  const collection = reactive(devicesStore.newCollection());
  collection.needsMorePages = true; // Todo: Connect with paging logic instead

  // Computed property for the query object
  const query = computed(() => {
    const q = {conditionsList: []};
    const opts = toValue(props);
    if (search.value) {
      const words = search.value.split(/\s+/);
      q.conditionsList.push(...words.map(word => ({stringContainsFold: word})));
    }
    if (opts.subsystem && opts.subsystem.toLowerCase() !== 'all') {
      q.conditionsList.push({field: 'metadata.membership.subsystem', stringEqualFold: opts.subsystem});
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
    const opts = toValue(props);
    const values = Object.values(collection.resources.value);
    if (!opts.filter) return values;
    return values.filter(opts.filter);
  });

  return {
    floorList,
    filterFloor,
    search,
    query,
    devicesData
  };
}
