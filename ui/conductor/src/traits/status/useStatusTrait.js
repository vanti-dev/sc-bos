import {newResourceValue} from '@/api/resource';
import {pullCurrentStatus} from '@/api/sc/traits/status';
import {toQueryObject, watchResource} from '@/util/traits.js';
import {toValue} from '@/util/vue.js';
import {computed, reactive} from 'vue';
import {StatusLog} from '@sc-bos/ui-gen/proto/status_pb';
import {ref} from 'vue/src/v3/index.js';

/**
 * @param {MaybeRefOrGetter<string|PullCurrentStatusRequest.AsObject>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @return {{
 *  statusValue: ResourceValue<StatusLog.AsObject, StatusLog>,
 *  statusLogLevel: import('vue').ComputedRef<number>,
 *  statusLogDescription: import('vue').ComputedRef<string>,
 *  statusOk: import('vue').ComputedRef<boolean>,
 *  statusNotOk: import('vue').ComputedRef<boolean>,
 *  statusLogProblems: import('vue').ComputedRef<Array<StatusLog.Problem.AsObject>>,
 *  hasMoreProblem: import('vue').ComputedRef<boolean>,
 *  showMoreProblems: import('vue').Ref<boolean>,
 *  statusLevelString: (level: number, type?: string) => string,
 *  statusMap: import('vue').ComputedRef<{level: {[key: string]: string}, description: string}>,
 *  statusProblemsMap: import('vue').ComputedRef<Array<{level: string, name: string, description: string}>>,
 *  statusLogIconColor: import('vue').ComputedRef<string>,
 *  statusLogIconString: import('vue').ComputedRef<string>,
 *  error: import('vue').ComputedRef<ResourceError>,
 *  loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const statusValue = reactive(
      /** @type {ResourceValue<StatusLog.AsObject, StatusLog>} */
      newResourceValue()
  );

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullCurrentStatus(req, statusValue);
        return statusValue;
      }
  );

  // --------------------- Status Log --------------------- //
  /** @type {import('vue').ComputedRef<number>} */
  const statusLogLevel = computed(() => {
    return statusValue.value?.level || 0;
  });

  /** @type {import('vue').ComputedRef<string>} */
  const statusLogDescription = computed(() => {
    return statusValue.value?.description || '';
  });

  /** @type {import('vue').ComputedRef<boolean>} */
  const statusOk = computed(() => {
    return statusValue.value?.level === StatusLog.Level.NOMINAL;
  });

  /** @type {import('vue').ComputedRef<boolean>} */
  const statusNotOk = computed(() => {
    return statusValue.value?.level > StatusLog.Level.NOMINAL;
  });

  // --------------------- Status Log Problems --------------------- //
  /** @type {import('vue').ComputedRef<Array<StatusLog.Problem.AsObject>>} */
  const statusLogProblems = computed(() => {
    return statusValue.value?.problemsList || [];
  });

  /** @type {import('vue').ComputedRef<boolean>} */
  const hasMoreProblem = computed(() => {
    return statusLogProblems.value.length > 0;
  });

  /** @type {import('vue').Ref<boolean>} */
  const showMoreProblems = ref(false);

  /**
   * Returns a string representation of the status level
   *
   * @param {number} level
   * @param {string} [type]
   * @return {string}
   */
  const statusLevelString = (level, type = 'short') => {
    if (level === StatusLog.Level.LEVEL_UNDEFINED) return '';
    if (level <= StatusLog.Level.NOMINAL) return type === 'long' ? 'Status: Nominal' : 'Nominal';
    if (level <= StatusLog.Level.NOTICE) return type === 'long' ? 'Status: Notice' : 'Notice';
    if (level <= StatusLog.Level.REDUCED_FUNCTION) return 'Reduced Function';
    if (level <= StatusLog.Level.NON_FUNCTIONAL) return 'Non-Functional';
    if (level <= StatusLog.Level.OFFLINE) return type === 'long' ? 'Status: Offline' : 'Offline';

    return 'Custom Level ' + level;
  };

  /**
   * Returns an object with the level and description of the status
   *
   * @type {import('vue').ComputedRef<{level: {[key: string]: string}, description: string}>}
   */
  const statusMap = computed(() => {
    return {
      level: {
        long: statusLevelString(statusLogLevel.value, 'long'),
        short: statusLevelString(statusLogLevel.value)
      },
      description: statusLogDescription.value
    };
  });

  /**
   * Returns an array of objects with the level, name, and description of the problems
   *
   * @type {import('vue').ComputedRef<Array<{level: string, name: string, description: string}>>}
   */
  const statusProblemsMap = computed(() => {
    return statusLogProblems.value.map((problem) => {
      return {
        level: statusLevelString(problem.level, 'long'),
        name: problem.name,
        description: problem.description
      };
    });
  });

  // --------------------- Status Log Style --------------------- //
  /** @type {import('vue').ComputedRef<string>} */
  const statusLogIconColor = computed(() => {
    if (statusLogLevel.value <= StatusLog.Level.NOTICE) return 'info';
    if (statusLogLevel.value <= StatusLog.Level.REDUCED_FUNCTION) return 'warning';
    if (statusLogLevel.value <= StatusLog.Level.NON_FUNCTIONAL) return 'error';
    if (statusLogLevel.value <= StatusLog.Level.OFFLINE) return 'grey';
    if (statusLogLevel.value <= StatusLog.Level.NOMINAL) return 'success';
    return 'white';
  });

  /** @type {import('vue').ComputedRef<string>} */
  const statusLogIconString = computed(() => {
    if (statusLogLevel.value <= StatusLog.Level.NOMINAL) return 'mdi-check-circle-outline';
    if (statusLogLevel.value <= StatusLog.Level.NOTICE) return 'mdi-information-outline';
    if (statusLogLevel.value <= StatusLog.Level.REDUCED_FUNCTION) return 'mdi-progress-alert';
    if (statusLogLevel.value <= StatusLog.Level.NON_FUNCTIONAL) return 'mdi-alert-circle-outline';
    if (statusLogLevel.value <= StatusLog.Level.OFFLINE) return 'mdi-connection';
    return '';
  });

  // --------------------- Error --------------------- //

  /** @type {import('vue').ComputedRef<ResourceError>} */
  const error = computed(() => {
    return statusValue.streamError;
  });

  // --------------------- Loading --------------------- //

  /** @type {import('vue').ComputedRef<boolean>} */
  const loading = computed(() => {
    return statusValue.loading;
  });


  return {
    statusValue,

    statusLogLevel,
    statusLogDescription,
    statusOk,
    statusNotOk,

    statusLogProblems,
    hasMoreProblem,
    showMoreProblems,

    statusLevelString,
    statusMap,
    statusProblemsMap,

    statusLogIconColor,
    statusLogIconString,

    error,
    loading
  };
}
