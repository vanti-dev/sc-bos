import {MINUTE} from '@/components/now.js';
import {usePullDevicesMetadata} from '@/composables/devices.js';
import {HealthCheck} from '@vanti-dev/sc-bos-ui-gen/proto/health_pb';
import {computed, onScopeDispose, ref, toValue, watch} from 'vue';

/**
 * Returns an older version of value that is no older than age.
 *
 * @param {import('vue').MaybeRefOrGetter<T>} value
 * @param {number} [age] - maximum age of history, default 5 minutes
 * @param {number} [resolution] - resolution to sample history, default 1 second. Updates to value more frequent than this will be ignored.
 * @return {import('vue').ComputedRef<T>} age
 * @template T
 */
export function useRollingHistory(value, age = 5 * MINUTE, resolution = MINUTE) {
  /**
   * @typedef {Object} Record
   * @property {number} t - timestamp in milliseconds
   * @property {T} v - value
   */
  /** @type {import('vue').Ref<Record[]>} */
  const oldValues = ref([]);
  const lastRecordedTime = ref(0);
  watch(() => toValue(value), (newValue) => {
    const now = Date.now();
    if (lastRecordedTime.value - now < resolution) {
      lastRecordedTime.value = now;
      oldValues.value.push({t: now, v: newValue});
    }
  }, {deep: true});
  let timer = 0;
  onScopeDispose(() => clearTimeout(timer));
  const processOldValues = () => {
    const now = Date.now();
    while (oldValues.value.length > 0 && (oldValues.value[0].t - now) < age) {
      oldValues.value.shift();
    }
  }
  watch(oldValues, (vs) => {
    clearTimeout(timer);
    if (vs.length <= 1) return;
    timer = setTimeout(() => {
      processOldValues();
    }, age - (vs[1].t - Date.now()))
  })

  return {
    oldValue: computed(() => {
      return oldValues.value?.[0]?.v ?? toValue(value);
    }),
    oldValues,
  };
}

/**
 * @typedef {Object} TableProps
 * @property {string} title
 * @property {number} totalCount
 * @property {string} color
 * @property {string} affectLabel
 * @property {boolean} hideAffected
 * @property {Array<{title: string, count: number, prevCount: number, affect: string}>} issues
 * @property {number} errorCount
 * @property {number} prevErrorCount
 */

/**
 * Returns properties that can be used for a ImpactTable component showing occupant impact health checks.
 *
 * @param {import('vue').MaybeRefOrGetter<import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject>} currentCounts - should include the field 'health_checks.occupant_impact'
 * @param {import('vue').MaybeRefOrGetter<import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject>} oldCounts - should include the field 'health_checks.occupant_impact'
 * @param {import('vue').MaybeRefOrGetter<Array<{field: string, stringIn?: {stringsList: string[]}}>>} [conditions] - additional conditions to filter devices
 * @return {{
 *   table: import('vue').ComputedRef<TableProps>,
 * }}
 */
export function useOccupantImpactTable(currentCounts, oldCounts, conditions) {
  return useImpactTable(currentCounts, oldCounts, conditions, {
    title: 'People',
    affectLabel: 'People affected',
    impactField: 'occupant_impact',
    fields: [
      {title: 'Life', key: 'LIFE'},
      {title: 'Health', key: 'HEALTH'},
      {title: 'Comfort', key: 'COMFORT'},
    ]
  });
}

/**
 * Returns properties that can be used for a ImpactTable component showing equipment impact health checks.
 *
 * @param {import('vue').MaybeRefOrGetter<import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject>} currentCounts - should include the field 'health_checks.equipment_impact'
 * @param {import('vue').MaybeRefOrGetter<import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject>} oldCounts - should include the field 'health_checks.equipment_impact'
 * @param {import('vue').MaybeRefOrGetter<Array<{field: string, stringIn?: {stringsList: string[]}}>>} [conditions] - additional conditions to filter devices
 * @return {{
 *   table: import('vue').ComputedRef<TableProps>,
 * }}
 */
export function useEquipmentImpactTable(currentCounts, oldCounts, conditions) {
  return useImpactTable(currentCounts, oldCounts, conditions, {
    title: 'Equipment',
    affectLabel: 'Units affected',
    impactField: 'equipment_impact',
    fields: [
      {title: 'Function', key: 'FUNCTION'},
      {title: 'Warranty', key: 'WARRANTY'},
      {title: 'Lifespan', key: 'LIFESPAN'},
    ]
  });
}

/**
 * Returns properties that can be used for a ImpactTable component showing compliance impact health checks.
 *
 * @param {import('vue').MaybeRefOrGetter<import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject>} currentCounts - should include the field 'health_checks.compliance_impacts.contribution'
 * @param {import('vue').MaybeRefOrGetter<import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject>} oldCounts - should include the field 'health_checks.compliance_impacts.contribution'
 * @param {import('vue').MaybeRefOrGetter<Array<{field: string, stringIn?: {stringsList: string[]}}>>} [conditions] - additional conditions to filter devices
 * @return {{
 *   table: import('vue').ComputedRef<TableProps>,
 * }}
 */
export function useComplianceImpactTable(currentCounts, oldCounts, conditions) {
  return useImpactTable(currentCounts, oldCounts, conditions, {
    title: 'Compliance',
    affectLabel: 'Standards affected',
    impactField: 'compliance_impacts.contribution',
    fields: [
      {title: 'Fail', key: 'FAIL'},
      {title: 'Warning', key: 'WARNING'},
      {title: 'Rating', key: 'RATING'},
      {title: 'Note', key: 'NOTE'},
    ]
  });
}

/**
 * Generic function to create an impact table for a specific type of health check.
 *
 * @param {import('vue').MaybeRefOrGetter<import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject>} currentCounts
 * @param {import('vue').MaybeRefOrGetter<import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject>} oldCounts
 * @param {import('vue').MaybeRefOrGetter<Array<{field: string, stringIn?: {stringsList: string[]}}>>} [conditions]
 * @param {Object} opts
 * @param {string} opts.title - title of the table
 * @param {string} opts.color - color of the score (e.g. 'primary', 'secondary')
 * @param {string} opts.affectLabel - label for the affect column
 * @param {string} opts.impactField - HealthCheck field that represents the impact (e.g. 'equipment_impact')
 * @param {Array<{title: string, key: string}>} opts.fields - fields to include in the table
 * @return {{
 *   table: import('vue').ComputedRef<TableProps>,
 * }}
 */
function useImpactTable(currentCounts, oldCounts, conditions, opts) {
  const impactFieldPath = 'health_checks.' + opts.impactField;
  const mdQuery = computed(() => {
    const moreConditions = toValue(conditions) ?? [];
    return {
      query: {
        conditionsList: [
          {field: impactFieldPath, stringIn: {stringsList: opts.fields.map(f => f.key)}},
          ...moreConditions
        ],
      },
      includes: {
        fieldsList: [
          'health_checks.check.state',
          'health_checks.reliability.state'
        ],
      }
    }
  })
  const {value: md} = usePullDevicesMetadata(mdQuery);
  const {oldValue: oldMd} = useRollingHistory(() => md.value);
  const table = computed(() => {
    const totalCount = md.value?.totalCount ?? 0;
    const oldTotals = getMetadataField(oldCounts.value, impactFieldPath)
    const newTotals = getMetadataField(currentCounts.value, impactFieldPath)
    const issues = opts.fields.map(f => {
      const count = getCountField(newTotals, f.key);
      return {
        title: f.title,
        count,
        prevCount: getCountField(oldTotals, f.key) ?? count,
        affect: '-',
      }
    });
    const errorCount = totalUnreliableCount(getMetadataField(md.value, 'health_checks.reliability.state'));
    const prevErrorCount = totalUnreliableCount(getMetadataField(oldMd.value, 'health_checks.reliability.state'));
    return {
      title: opts.title,
      totalCount,
      affectLabel: opts.affectLabel,
      hideAffected: true,
      issues,
      errorCount,
      prevErrorCount,
    }
  });
  return {
    table
  }
}


/**
 * Returns the counts map for a specific field from the metadata object.
 *
 * @param {import('@vanti-dev/sc-bos-ui-gen/proto/devices_pb').DevicesMetadata.AsObject | null} m
 * @param {string} f
 * @return {Array<[string, number]> | null}
 */
export function getMetadataField(m, f) {
  return m?.fieldCountsList?.find(r => r.field === f)?.countsMap ?? []
}

/**
 * Returns the count for a specific field from the metadata object.
 *
 * @param {Array<[string, number]> | null} m
 * @param {string} f
 * @return {number}
 */
export function getCountField(m, f) {
  for (const [k, v] of m ?? []) {
    if (k === f) {
      return v;
    }
  }
  return 0;
}

/**
 * Returns the total count of unreliable states from a list of [state, count] pairs.
 * The argument is typically the result of a GetDevicesMetadata call.
 *
 * @param {Array<[string, number]>} counts
 * @return {number}
 */
export function totalUnreliableCount(counts) {
  return counts?.reduce((acc, [k, v]) => {
    if (k === 'RELIABLE' || k === 'STATE_UNSPECIFIED') return acc;
    return acc + v;
  }, 0);
}

/**
 * Counts the number of checks in a specific state.
 *
 * @param {Array<import('@vanti-dev/sc-bos-ui-gen/proto/health_pb').HealthCheck.AsObject>} checks
 * @param {import('@vanti-dev/sc-bos-ui-gen/proto/health_pb').HealthCheck.Check.State} state
 * @return {number}
 */
export function countChecksByState(checks, state) {
  return checks?.reduce((acc, check) => {
    if (check.check.state === state) acc++;
    return acc;
  }, 0);
}

/**
 * Counts the number of normal and abnormal checks.
 *
 * @param {Array<import('@vanti-dev/sc-bos-ui-gen/proto/health_pb').HealthCheck.AsObject>} checks
 * @return {{normalCount: number, abnormalCount: number, totalCount: number}}
 */
export function countChecks(checks) {
  const normalCount = countChecksByState(checks, HealthCheck.Check.State.NORMAL);
  const abnormalCount = checks?.reduce((acc, check) => {
    if (check.check.state > HealthCheck.Check.State.NORMAL) acc++;
    return acc;
  }, 0);
  return {
    normalCount,
    abnormalCount,
    totalCount: normalCount + abnormalCount,
  }
}