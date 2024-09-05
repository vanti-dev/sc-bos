import {useDevicesCollection, useDevicesMetadataField, usePullDevicesMetadata} from '@/devices/devices.js';
import {computed, toValue} from 'vue';

const NO_FLOOR = '< no floor >';

/**
 * @typedef {Object} UseDevicesOptions
 * @property {number} wantCount
 * @property {string} subsystem
 *   - if present and not 'all', adds the condition {field: "metadata.membership.subsystem", stringEqualFold: subsystem}
 * @property {string} floor
 *   - if present and not 'all', adds the condition {field: "metadata.location.floor", stringEqualFold: floor}
 * @property {string} search
 *   - if present adds a condition for each word {stringContainsFold: word}
 * @property {Device.Query.Condition.AsObject[]} conditions
 * @property {(value: Device.AsObject, index: number, array: Device.AsObject[]) => boolean} filter
 */

/**
 *
 * @param {MaybeRefOrGetter<Partial<UseDevicesOptions>>} props
 * @return {UseCollectionResponse<Device.AsObject> & {
 *   floorList: import('vue').ComputedRef<Array>,
 *   query: import('vue').ComputedRef<Object>,
 * }}
 */
export default function(props) {
  const opts = computed(() => /** @type {Partial<UseDevicesOptions>} */ toValue(props));

  const {value: md} = usePullDevicesMetadata('metadata.location.floor');
  const {keys: listOfFloors} = useDevicesMetadataField(md, 'metadata.location.floor');

  const floorList = computed(() => {
    return ['All', ...listOfFloors.value.map(v => v === '' ? NO_FLOOR : v)];
  });

  const conditions = computed(() => {
    const _opts = opts.value;
    const conditionsList = [..._opts.conditions ?? []];
    if (_opts.search) {
      const words = _opts.search.split(/\s+/);
      conditionsList.push(...words.map(word => ({stringContainsFold: word})));
    }
    if (_opts.subsystem && _opts.subsystem.toLowerCase() !== 'all') {
      conditionsList.push({field: 'metadata.membership.subsystem', stringEqualFold: _opts.subsystem});
    }
    if (_opts.floor) {
      switch (_opts.floor.toLowerCase()) {
        case 'all':
          // no filter
          break;
        case NO_FLOOR:
          conditionsList.push({field: 'metadata.location.floor', stringEqualFold: ''});
          break;
        default:
          conditionsList.push({field: 'metadata.location.floor', stringEqualFold: _opts.floor});
          break;
      }
    }
    return conditionsList;
  });
  const query = computed(() => {
    return {conditionsList: conditions.value};
  });
  const request = computed(() => {
    return {query: query.value};
  });
  const deviceCollectionOptions = computed(() => {
    return {
      wantCount: opts.value.wantCount ?? 20
    };
  });

  const collection = useDevicesCollection(request, deviceCollectionOptions);

  // Computed property for the filtered table data
  const items = computed(() => {
    const values = collection.items.value;
    if (!opts.value.filter) return values;
    return values.filter(opts.value.filter);
  });

  return {
    ...collection,
    floorList,
    query,
    items
  };
}
