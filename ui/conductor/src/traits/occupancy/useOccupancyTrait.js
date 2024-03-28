import {timestampToDate} from '@/api/convpb.js';
import {newResourceValue} from '@/api/resource';
import {occupancyStateToString, pullOccupancy} from '@/api/sc/traits/occupancy';
import {toQueryObject, watchResource} from '@/util/traits';
import {toValue} from '@/util/vue';
import {Occupancy} from '@smart-core-os/sc-api-grpc-web/traits/occupancy_sensor_pb';
import {computed, reactive} from 'vue';

/**
 *
 * @param {MaybeRefOrGetter<string|PullOccupancyRequest.AsObject>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @return {{
 *  occupancyValue: ResourceValue<Occupancy.AsObject, Occupancy>,
 *  occupancyPeopleCount: import('vue').ComputedRef<number>,
 *  occupancyStateNumber: import('vue').ComputedRef<Occupancy.State>,
 *  occupancyStateString: import('vue').ComputedRef<string>,
 *  occupancyLastUpdateDate: import('vue').ComputedRef<Date>,
 *  occupancyColorString: import('vue').ComputedRef<string>,
 *  occupancyIconString: import('vue').ComputedRef<string>,
 *  occupancyIconColor: import('vue').ComputedRef<string|undefined>,
 *  error: import('vue').ComputedRef<ResourceError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const occupancyValue = reactive(
      /** @type {ResourceValue<Occupancy.AsObject, Occupancy>} */
      newResourceValue()
  );

  const queryObject = computed(() => toQueryObject(query));

  // Utility function to call the API with the query and the resource
  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullOccupancy(req, occupancyValue);
        return occupancyValue;
      }
  );


  // ------------------- Occupancy Values ------------------- //
  /** @type {import('vue').ComputedRef<number>} */
  const occupancyPeopleCount = computed(() => {
    return occupancyValue.value?.peopleCount || 0;
  });

  /** @type {import('vue').ComputedRef<number>} */
  const occupancyStateNumber = computed(() => {
    return occupancyValue.value?.state;
  });

  const occupancyStateString = computed(() => {
    return occupancyStateToString(occupancyValue.value?.state) || '';
  });

  /** @type {import('vue').ComputedRef<Date>} */
  const occupancyLastUpdateDate = computed(() => {
    return timestampToDate(occupancyValue.value?.stateChangeTime);
  });

  // ------------------- Occupancy Styles ------------------- //
  /** @type {import('vue').ComputedRef<string>} */
  const occupancyColorString = computed(() => {
    return occupancyStateString.value.toLowerCase();
  });

  /** @type {import('vue').ComputedRef<string>} */
  const occupancyIconString = computed(() => {
    if (occupancyStateNumber.value === Occupancy.State.OCCUPIED) {
      return 'mdi-crosshairs-gps';
    } else if (occupancyStateNumber.value === Occupancy.State.UNOCCUPIED) {
      return 'mdi-crosshairs';
    } else if (occupancyStateNumber.value === Occupancy.State.IDLE) {
      return 'mdi-crosshairs-gps';
    } else {
      return '';
    }
  });

  /** @type {import('vue').ComputedRef<string|undefined>} */
  const occupancyIconColor = computed(() => {
    if (occupancyStateNumber.value === Occupancy.State.OCCUPIED) {
      return 'success lighten1';
    } else if (occupancyStateNumber.value === Occupancy.State.UNOCCUPIED) {
      return 'warning';
    } else if (occupancyStateNumber.value === Occupancy.State.IDLE) {
      return 'info';
    } else {
      return undefined;
    }
  });

  // ------------------- Error ------------------- //
  /** @type {import('vue').ComputedRef<ResourceError>} */
  const error = computed(() => {
    return occupancyValue.streamError;
  });

  // ------------------- Loading ------------------- //
  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => {
    return occupancyValue.loading;
  });


  return {
    occupancyValue,

    occupancyPeopleCount,
    occupancyStateNumber,
    occupancyStateString,
    occupancyLastUpdateDate,

    occupancyColorString,
    occupancyIconString,
    occupancyIconColor,

    error,
    loading
  };
}
