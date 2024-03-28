import {newResourceValue} from '@/api/resource';
import {pullAccessAttempts} from '@/api/sc/traits/access';
import {toQueryObject, watchResource} from '@/util/traits';
import {toValue} from '@/util/vue';
import {AccessAttempt} from '@sc-bos/ui-gen/proto/access_pb';
import {computed, reactive} from 'vue';

/**
 * @param {MaybeRefOrGetter<string|PullAccessAttemptsRequest.AsObject>} query - The name of the device or a query object
 * @param {MaybeRefOrGetter<boolean>} paused - Whether to pause the data stream
 * @return {{
 *   accessAttemptValue: ResourceValue<AccessAttempt.AsObject, PullAccessAttemptsResponse>,
 *   accessAttemptGrantId: import('vue').ComputedRef<AccessAttempt.Grant>,
 *   accessAttemptGrantNamesByID: {[key: number]: string},
 *   accessAttemptGrantState: import('vue').ComputedRef<string>,
 *   accessAttemptInformation: import('vue').ComputedRef<Array<AccessAttemptInfo.info, AccessAttemptInfo.subInfo>>,
 *   error: import('vue').ComputedRef<ResourceError>,
 *   loading: import('vue').ComputedRef<boolean>
 * }}
 */
export default function(query, paused) {
  const accessAttemptValue = reactive(
      /** @type {ResourceValue<AccessAttempt.AsObject, PullAccessAttemptsResponse>} */
      newResourceValue()
  );

  const queryObject = computed(() => toQueryObject(query));

  watchResource(
      () => toValue(queryObject),
      () => toValue(paused),
      (req) => {
        pullAccessAttempts(req, accessAttemptValue);
        return accessAttemptValue;
      }
  );

  /**
   * Returns the grant id of the access attempt
   *
   * @type {import('vue').ComputedRef<AccessAttempt.Grant>} accessAttemptGrantId
   */
  const accessAttemptGrantId = computed(() => accessAttemptValue.value?.grant);

  /**
   * Returns all the grant names sorted by id
   *
   * @type {{[key: number]: string}} accessAttemptGrantNamesByID
   */
  const accessAttemptGrantNamesByID = Object.entries(AccessAttempt.Grant).reduce((all, [name, id]) => {
    all[id] = name.toLowerCase();
    return all;
  }, {});

  /**
   * Returns the grant state of the access attempt as a string based on the grant id
   *
   * @type {import('vue').ComputedRef<string>} accessAttemptGrantState
   */
  const accessAttemptGrantState = computed(() => {
    return accessAttemptGrantNamesByID[accessAttemptGrantId.value || 0];
  });

  /**
   * A computed property that processes and returns access attempt information.
   *
   * @typedef {Object} AccessAttemptSubInfoValue
   * @property {Object.<string, string|number|Object[]|string[]|number[]>} [subKey]
   * @typedef {Object} AccessAttemptInfo
   * @property {Object.<string, string|number>} info
   * @property {Object.<string, AccessAttemptSubInfoValue>} subInfo
   * @return {[AccessAttemptInfo.info, AccessAttemptInfo.subInfo]} accessAttemptInformation
   */
  const accessAttemptInformation = computed(() => {
    // Initialize objects to hold processed information
    const info = {};
    const subInfo = {};

    // Check if accessAttemptValue has a value
    if (accessAttemptValue.value) {
      // Iterate over each entry in the accessAttemptValue object
      Object.entries(accessAttemptValue.value).forEach(([key, value]) => {
        // Process non-empty and non-object values directly
        if (value && typeof value !== 'object') {
          // Special handling for 'grant' key to format the grant name
          if (key === 'grant') {
            // Assign the formatted grant name to the info object
            info[key] = accessAttemptGrantNamesByID[value]
                .split('_')
                .join(' ')
                .replace(/^./, match => match.toUpperCase()); // Capitalize the first letter
          } else {
            // Directly assign other non-object values
            info[key] = value;
          }
        } else if (value) { // Ensure value is not null or undefined for object type
          // Process object type values
          Object.entries(value).forEach(([subKey, subValue]) => {
            // Only process non-empty subValues
            if (subValue) {
              // Initialize subInfo[key] as an object if it doesn't exist
              subInfo[key] = subInfo[key] || {};

              // Handle array values by mapping them to objects, otherwise, assign directly
              subInfo[key][subKey] = Array.isArray(subValue) ?
                  subValue.map(([innerKey, innerValue]) => ({[innerKey]: innerValue})) :
                  subValue;
            }
          });
        }
      });
    }

    // Return the structured information
    return [info, subInfo];
  });


  const error = computed(() => accessAttemptValue.streamError);

  const loading = computed(() => accessAttemptValue.loading);

  return {
    accessAttemptValue,

    accessAttemptGrantId,
    accessAttemptGrantNamesByID,
    accessAttemptGrantState,
    accessAttemptInformation,

    error,
    loading
  };
}
