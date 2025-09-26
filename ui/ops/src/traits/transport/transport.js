import {closeResource, newActionTracker, newResourceValue} from '@/api/resource.js';
import {describeTransport, pullTransport} from '@/api/sc/traits/transport.js';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {Transport} from '@vanti-dev/sc-bos-ui-gen/proto/transport_pb';
import {computed, onScopeDispose, reactive, toRefs, toValue} from 'vue';

/**
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/transport_pb').DescribeTransportRequest} DescribeTransportRequest
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/transport_pb').PullTransportRequest} PullTransportRequest
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/transport_pb').PullTransportResponse} PullTransportResponse
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/transport_pb').Transport} Transport
 * @typedef {import('@vanti-dev/sc-bos-ui-gen/proto/transport_pb').TransportSupport} TransportSupport
 * @typedef {import('vue').UnwrapNestedRefs} UnwrapNestedRefs
 * @typedef {import('vue').ToRefs} ToRefs
 * @typedef {import('vue').ComputedRef} ComputedRef
 * @typedef {import('vue').MaybeRefOrGetter} MaybeRefOrGetter
 * @typedef {import('@/api/resource').ResourceValue} ResourceValue
 */

/**
 * @param {MaybeRefOrGetter<string|PullTransportRequest.AsObject>} query
 * @param {MaybeRefOrGetter<boolean>=} paused
 * @return {ToRefs<ResourceValue<Transport.AsObject, PullTransportResponse>>}
 */
export function usePullTransport(query, paused = false) {
  const resource = reactive(
      /** @type {ResourceValue<Transport.AsObject, PullTransportResponse>} */
      newResourceValue()
  );
  onScopeDispose(() => closeResource(resource));

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullTransport(req, resource);
        return () => closeResource(resource);
      }
  );

  return toRefs(resource);
}

/**
 * @param {MaybeRefOrGetter<string|DescribeTransportRequest.AsObject>} query - The name of the device or a query
 *   object
 * @return {ToRefs<ActionTracker<TransportSupport.AsObject>>}
 */
export function useDescribeTransport(query) {
  const tracker = reactive(
      /** @type {ActionTracker<TransportSupport.AsObject>} */
      newActionTracker()
  );

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      false,
      (req) => {
        describeTransport(req, tracker);
        return () => closeResource(tracker);
      }
  );

  return toRefs(tracker);
}

/**
 * Convert a proto enum value to a display name. e.g. "DOOR_OPEN" -> "Door Open"
 * replace underscores with spaces, capitalize the first letter of each word
 *
 * @param {string} e
 * @return {string}
 */
function enumToDisplayName(e) {
  // replace underscores in with spaces, change the letters to lower case
  e = e.replace(/_/g, ' ').toLowerCase();
  // then capitalize the first letter eof each word
  e = e.replace(/\b\w/g, l => l.toUpperCase());
  return e;
}

/**
 * @param {MaybeRefOrGetter<Transport.AsObject|null>} value
 * @param {MaybeRefOrGetter<TransportSupport.AsObject|null>=} support
 * @return {{
 *   actualPosition: ComputedRef<string>,
 *   doorStatus: ComputedRef<Array<{label:string, value:string}>>,
 *   movingDirection: ComputedRef<string>,
 *   nextDestination: ComputedRef<string>,
 *   load: ComputedRef<string>,
 *   table: ComputedRef<Array<{label:string, value:string}>>
 * }}
 */
export function useTransport(value, support = null) {
  const _v = computed(() => toValue(value));
  const _s = computed(() => toValue(support));

  const movingDirectionById = Object.entries(Transport.Direction).reduce((all, [name, id]) => {
    all[id] = name;
    return all;
  }, {});
  
  const operatingModeById = Object.entries(Transport.OperatingMode).reduce((all, [name, id]) => {
    all[id] = name;
    return all;
  }, {});

  const doorStatusById = Object.entries(Transport.Door.DoorStatus).reduce((all, [name, id]) => {
    all[id] = name;
    return all;
  }, {});

  const actualPosition = computed(() => {
    const v = _v.value;
    if (!v) return '';
    if (v.actualPosition?.floor && v.actualPosition?.floor !== '') {
      return v.actualPosition?.floor;
    }
    if (v.actualPosition?.title && v.actualPosition?.title !== '') {
      return v.actualPosition?.title;
    }
    if (v.actualPosition?.id && v.actualPosition?.id !== '') {
      return v.actualPosition?.id;
    }
    return '';
  });

  const doorStatus = computed(() => {
    const v = _v.value;
    if (!v) return [];
    let res = [];
    for (const door of v.doorsList) {
      res.push({
        label: door.title !== '' ? door.title : 'Door',
        value: enumToDisplayName(doorStatusById[door.status] ?? '')
      });
    }
    return res;
  });

  const movingDirection = computed(() => {
    const v = _v.value;
    if (!v) return '';
    return enumToDisplayName(movingDirectionById[v.movingDirection] ?? '');
  });
  
  const operatingMode = computed(() => {
    const v = _v.value;
    if (!v) return '';
    return enumToDisplayName(operatingModeById[v.operatingMode] ?? '');
  });

  const nextDestination = computed(() => {
    const v = _v.value;
    if (!v) return '';
    if (v.nextDestinationsList.length > 0) {
      return v.nextDestinationsList[0]?.floor ?? 'N/A';
    } else {
      return 'N/A';
    }
  });

  const loadUnit = computed(() => {
    return _s.value?.loadUnit ?? 'kg';
  });

  const load = computed(() => {
    const v = _v.value;
    if (!v) return '';
    return v.load?.toFixed(2);
  });

  const loadStr = computed(() => {
    if (load.value)
    {
        return `${load.value} ${loadUnit.value}`;
    }
    return '-';
  });

  const table = computed(() => {
    let t = [{
      label: 'Actual Position',
      value: actualPosition.value
    },
    {
      label: 'Moving Direction',
      value: movingDirection.value
    },
    {
      label: 'Next Destination',
      value: nextDestination.value
    },
    {
      label: 'Operating Mode',
      value: operatingMode.value
    }
    ];

    for (const door of doorStatus.value) {
      t.push(door);
    }
    t.push({
      label: 'Load',
      value: loadStr.value
    });
    return t;
  });

  return {
    actualPosition,
    doorStatus,
    movingDirection,
    nextDestination,
    operatingMode,
    loadStr,
    table
  };
}