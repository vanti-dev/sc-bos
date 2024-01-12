import {computed, onMounted, onUnmounted, reactive, ref, watch, watchEffect} from 'vue';

import {closeResource, newActionTracker} from '@/api/resource';
import {getEnrollment, testEnrollment} from '@/api/sc/traits/enrollment';
import {testHubNode} from '@/api/sc/traits/hub';
import {useHubStore} from '@/stores/hub';

/**
 * @typedef {Object} ComponentData
 * @property {import('vue').Ref<boolean>} isLoading - Indicates whether the component is currently loading data.
 * @property {ActionTracker<GetEnrollmentRequest.AsObject>} enrollmentValue - for enrollment requests.
 * @property {ActionTracker<ListHubNodesResponse.AsObject>} listHubNodesValue - for listing hub nodes.
 * @property {ActionTracker<TestEnrollmentResponse.AsObject>} testEnrollmentValue - for testing enrollment.
 * @property {import('vue').ComputedRef<{
 *     enrollment: GetEnrollmentResponse.AsObject | null,
 *     error: Error | null,
 *     loading: boolean,
 *     isTested: boolean | null,
 *     testLoading: boolean
 * }>} enrollmentStatus - Computed ref representing the current status of enrollment.
 * @property {function(): Promise<void>} triggerListHubNodesAction - triggers an action for listing hub nodes.
 * @property {function(): Promise<void>} getEnrollmentAndListHubNodes - triggers fetch enrollment and list hub nodes.
 * @property {import('vue').Ref<Record<string, {
 *     color: string,
 *     icon: string,
 *     loading: boolean,
 *     resource: {
 *         status: string,
 *         message: string
 *     } | {
 *         error: Error
 *     },
 *     type: string
 * }>>} nodeStatus - holding the status of each node.
 * @property {import('vue').Ref<boolean>} updatingNodeStatus - indicating whether the status is currently being updated.
 * @property {import('vue').ComputedRef<string>} serverChipType - representing the type of server chip.
 * @property {import('vue').ComputedRef<string[]>} displayedChips - representing the chips to be displayed in the UI.
 * @return {ComponentData}
 */
export default function() {
  // ------------------- UI Status ------------------- //
  /**
   * Returns an object with enrollment status
   *
   * @type {import ('vue').ComputedRef<{
   *  enrollment: GetEnrollmentResponse.AsObject | null,
   *  error: Error | null,
   *  loading: boolean,
   *  isTested: boolean | null,
   *  testLoading: boolean
   * }>} enrollmentStatus
   */
  const enrollmentStatus = computed(() => {
    return {
      enrollment: enrollmentValue.response || null,
      error: enrollmentValue.error || testEnrollmentValue.error || null,
      loading: enrollmentValue.loading,
      isTested: testEnrollmentValue.response || testEnrollmentValue.error || null,
      testLoading: testEnrollmentValue.loading
    };
  });

  /**
   * Returns a type for the chip representing the server/hub/gateway depending on the successful checks
   *
   * @type {import ('vue').ComputedRef<string>} serverChipType
   */
  const serverChipType = computed(() => {
    const enrollment = enrollmentStatus.value.enrollment;
    if (enrollment) return 'gateway'; // If enrollment and hub nodes exist, return 'gateway'
    if (hubNodesList.value) return 'hub'; // If hub nodes exist, return 'hub'
    return 'server'; // Otherwise, return 'server'
  });

  /**
   * Returns an array of chips to display in the UI
   *
   * @type {import ('vue').ComputedRef<string[]>} displayedChips
   */
  const displayedChips = computed(() => [
    'ui',
    serverChipType.value,
    ...(enrollmentStatus.value.enrollment ? ['hub'] : []),
    ...(hubNodesList.value && hubNodesList.value.length > 0 ? ['nodes'] : [])
  ]);

  // ----------------------------------------------------------- //
  // ------------------- API Calls & Loading ------------------- //
  // const isLoading = ref(false); // Returns a boolean for whether the page is loading
  const isLoading = computed(() => {
    return enrollmentValue.loading || listHubNodesValue.loading;
  });
  const lastFetch = ref(null); // Returns a now timestamp for the last successful check

  //
  //
  /**
   * // ------------------- Main Function ------------------- //
   *
   * Asynchronously fetches the enrollment status and lists hub nodes.
   * Main function for determining the status of the server, hub, and nodes.
   * This is called on mount and every 15 seconds.
   *
   * @return {Promise<void>}
   */
  const getEnrollmentAndListHubNodes = async () => {
    await Promise.all([
      getEnrollmentValue(),
      triggerListHubNodesAction()
    ]);
    lastFetch.value = Date.now(); // Update the last execution time
  };

  /**
   * // ------------------- Enrollment & Test Enrollment ------------------- //
   *
   * This section manages the enrollment and testing of enrollment status in the application.
   * It includes action trackers for tracking enrollment and test enrollment status,
   * and provides functions to fetch and test enrollment data.
   *
   * - `enrollmentValue`: action tracker that tracks the enrollment status,
   *   including response, error, and loading states.
   * - `getEnrollmentValue()`: async function that fetches the enrollment status.
   *   It resets the error state on success and tests the enrollment status. In case of an error,
   *   it logs the error, resets the response, and also resets the node status.
   * - `testEnrollmentValue`: Similar to `enrollmentValue` - but for testing enrollment.
   * - `getTestEnrollmentValue()`: async function that tests the enrollment status
   *   and updates the `testEnrollmentValue`. It clears any existing error on success,
   *   and in case of an error, logs the error and resets the response state.
   *
   * These functions are initially called when the component is mounted and are subsequently
   * called every 15 seconds as part of the `getEnrollmentAndListHubNodes` routine.
   */

  // ------------------- Enrollment ------------------- //
  // Initialize the enrollmentValue as a reactive reference
  const enrollmentValue = reactive(
      /** @type {ActionTracker<GetEnrollmentRequest.AsObject>} */ newActionTracker()
  );

  // Function to get enrollment
  // This is called on mount and every 15 seconds as part of getEnrollmentAndListHubNodes
  const getEnrollmentValue = async () => {
    try {
      await getEnrollment(enrollmentValue); // Get the enrollment
      enrollmentValue.error = null; // Reset the error
      await getTestEnrollmentValue(); // Test the enrollment
    } catch (e) {
      console.error('Error fetching enrollment', e);
      enrollmentValue.response = null; // Reset the response
      nodeStatus.value = {}; // Reset the node status
    }
  };

  // // ------------------- Test Enrollment ------------------- //
  // Initialize the testEnrollmentValue as a reactive reference
  const testEnrollmentValue = reactive(
      /** @type {ActionTracker<TestEnrollmentResponse.AsObject>} */ newActionTracker()
  );

  /**
   * Asynchronously tests the enrollment status and updates the testEnrollmentValue state.
   * If the test is successful, any existing error is cleared. If an error occurs during the test,
   * the function logs the error and resets the response state.
   * This is called on mount and every 15 seconds as part of getEnrollmentAndListHubNodes -> getEnrollmentValue
   *
   * @return {Promise<void>} A promise that resolves when the test enrollment process is complete.
   * @async
   */
  const getTestEnrollmentValue = async () => {
    try {
      await testEnrollment(testEnrollmentValue); // Test the enrollment
      testEnrollmentValue.error = null; // Reset the error
    } catch (e) {
      console.error('Error fetching test enrollment', e);
      testEnrollmentValue.response = null; // Reset the response
    }
  };

  // // ------------------- List Hub Nodes ------------------- //
  /**
   * Asynchronously triggers an action to list hub nodes. It sets the updatingNodeStatus to true at the start
   * and sets it back to false upon completion. If the action is successful, it updates the listHubNodesValue state.
   * In case of an error, it logs the error and resets the listHubNodesValue response - this will trigger a re-render
   * of the UI to hide certain chips (e.g. nodes).
   * Can be used to manually trigger the action
   *
   * @return {Promise<void>} A promise that resolves when the action to list hub nodes is completed.
   * @async
   */
  const triggerListHubNodesAction = async () => {
    updatingNodeStatus.value = true;
    try {
      await listHubNodesAction(listHubNodesValue);
    } catch (e) {
      console.error('Error fetching hub nodes', e);
      listHubNodesValue.response = null;
    }
    updatingNodeStatus.value = false;
  };

  /**
   * // ------------------- Test Hub Nodes ------------------- //
   *
   * This section is focused on testing and managing the status of hub nodes in the network.
   * It involves initializing reactive references to track hub nodes, their test results, and any errors encountered.
   *
   * - `listHubNodesAction`: function from the `hub store` to fetch the list of hub nodes from the hub API.
   * - `listHubNodesValue`: tracker that tracks the response, error, and loading status for the `listHubNodesAction`.
   * - `testHubNodeValue`: Similar to `listHubNodesValue`, but for the testHubNode action.
   * - `trackedHubNodeErrors`: ref that stores errors encountered during testing of each node.
   * - `trackedHubNodeResults`: ref that stores the test results of each node.
   * - `hubNodesList`: computed property that derives the list of hub nodes from the `listHubNodesValue` response.
   * - `getTestHubNodeValue(nodesList)`: an async function that tests each node in the provided list
   *    and keeps track of the results and errors, maintaining only the last 5 results or errors for each node.
   * - `watch`: watcher that listens for changes in the `hubNodesList` and triggers the testing of nodes accordingly.
   *
   * This setup allows for efficient monitoring and updating of the status of each hub node in the network.
   */

  // Get the listHubNodesAction from the hub store which pulls from the hub API
  const {listHubNodesAction} = useHubStore();

  // Initialize the listHubNodesValue as a reactive reference
  // This is the action tracker for the listHubNodesAction
  // It contains the response, error, and loading status
  const listHubNodesValue = reactive(
      /** @type {ActionTracker<ListHubNodesResponse.AsObject>} */ newActionTracker()
  );
  // Initialize the trackedHubNodeResults as a reactive reference
  // Similar to the listHubNodesValue, this is the action tracker for the testHubNode action
  const testHubNodeValue = reactive(
      /** @type {ActionTracker<TestHubNodeResponse.AsObject>} */ newActionTracker()
  );

  // Initialize the trackedHubNodeErrors and trackedHubNodeResults as reactive references
  // Store the past and most recent results and errors for each node
  const trackedHubNodeErrors = ref({}); // Object to store errors for each node
  const trackedHubNodeResults = ref({}); // Object to store results for each node
  const hubNodesList = computed(() => listHubNodesValue.response?.nodesList);

  /**
   * Tests each hub node in the provided list and tracks the results and errors.
   * If a node test is successful, the result is stored; if it fails, the error is captured.
   * Only keeps the last 5 results or errors for each node to manage memory and reactivity.
   *
   * @param {Array} nodesList - An array of nodes to be tested. Each node should have an 'address' and 'name' property.
   * @return {Promise<Object>} A promise that resolves to an object containing two arrays: 'results' and 'errors'.
   *                            'results' is an array of the latest test results for each node,
   *                            and 'errors' is an array of the latest test errors for each node.
   * @async
   */
  const getTestHubNodeValue = async (nodesList) => {
    if (nodesList.length === 0) return {}; // Escaping early if the list is empty

    for (const node of nodesList) { // Iterate over each node in the list
      // Create a timestamp for this test, so we can keep track of the most recent results and errors
      const timestamp = {
        seconds: Math.floor(Date.now() / 1000), // Current time in seconds
        nanos: (Date.now() % 1000) * 1000000 // Current time in nanoseconds
      };

      try {
        const request = {address: node.address}; // Create a request object for the testHubNode action
        await testHubNode(request, testHubNodeValue); // Test the node

        // Initialize array for node results if it doesn't exist
        if (!trackedHubNodeResults.value[node.name]) {
          trackedHubNodeResults.value[node.name] = [];
        }

        // Add the new result to the node's result array
        trackedHubNodeResults.value[node.name].push({
          address: node.address, error: null, name: node.name, timestamp: timestamp
        });

        // Keep only the last 5 results for this node
        // Replace the entire object to trigger reactivity
        trackedHubNodeResults.value = {
          ...trackedHubNodeResults.value, [node.name]: trackedHubNodeResults.value[node.name].slice(-5)
        };
      } catch (e) {
        console.error('Error fetching test hub node', e);

        // Initialize array for node errors if it doesn't exist
        if (!trackedHubNodeErrors.value[node.name]) {
          trackedHubNodeErrors.value[node.name] = [];
        }

        // Add the new error to the node's error array
        trackedHubNodeErrors.value[node.name].push({
          address: node.address, error: e, name: node.name, timestamp: timestamp
        });

        // Keep only the last 5 errors for this node
        // Replace the entire object to trigger reactivity
        trackedHubNodeErrors.value = {
          ...trackedHubNodeErrors.value, [node.name]: trackedHubNodeErrors.value[node.name].slice(-5)
        };
      }
    }

    // Return flattened arrays of all results and errors
    return {
      results: Object.values(trackedHubNodeResults.value).flat(),
      errors: Object.values(trackedHubNodeErrors.value).flat()
    };
  };

  // Watch for changes in the list of hub nodes and test them
  watch(hubNodesList, async (nodesList) => {
    if (!nodesList || nodesList.length === 0) return; // Escaping early if the list is empty
    await getTestHubNodeValue(nodesList); // Test the nodes
  }, {deep: true});

  /**
   * // ------------------- Node Status ------------------- //
   *
   * This section manages the status of nodes in the network. It includes reactive references for node status and
   * a flag indicating whether the node status is being updated. The `updateNodeStatus` function asynchronously updates
   * the status of each node based on the latest test results and errors. It involves checking the most recent results
   * and errors, comparing timestamps to determine the current status, and updating the reactive `nodeStatus` object.
   * A `watchEffect` is set up to trigger `updateNodeStatus` whenever there are changes in
   * tracked node results or errors.
   * Additionally, a computed property `erroredNodes` is defined to return an array of nodes that are
   * currently in an error state.
   *
   * - `nodeStatus`: a ref object holding the current status of each node, including color, icon, and error information.
   * - `updatingNodeStatus`: a ref flag indicating whether the node status is currently being updated.
   * - `updateNodeStatus(nodeResults, nodeErrors)`: an async function that updates the status of each node.
   * - `watchEffect`: watcher that triggers `updateNodeStatus` on changes in node results or errors.
   * - `erroredNodes`: computed property that provides an array of nodes currently in an error state.
   */


  const nodeStatus = ref({}); // Initialize nodeStatus as a reactive reference
  const updatingNodeStatus = ref(false); // Initialize updatingNodeStatus as a reactive reference

  /**
   * Updates the status of each node based on the latest results and errors.
   * It checks the most recent result and error for each node, compares their timestamps,
   * and determines the current status of the node. The status includes a color indication,
   * an icon, and a message or error object.
   *
   * @param {Object} nodeResults - An object containing arrays of test results for each node.
   *                               The keys are node names, and the values are arrays of result objects.
   * @param {Object} nodeErrors - An object containing arrays of test errors for each node.
   *                              The keys are node names, and the values are arrays of error objects.
   * @return {Promise<void>} A promise that resolves when the node statuses have been updated.
   * @async
   */
  const updateNodeStatus = async (nodeResults, nodeErrors) => {
    const status = {}; // Initialize status as an empty object
    updatingNodeStatus.value = true; // Set updatingNodeStatus to true
    // Collect all unique node names from both results and errors
    const allNodes = await new Set([...Object.keys(nodeResults), ...Object.keys(nodeErrors)]);

    if (allNodes.size === 0) {
      updatingNodeStatus.value = false;
      return;
    }

    allNodes.forEach(node => {
      const results = nodeResults[node] || [];
      const errors = nodeErrors[node] || [];

      // Get the last result and the last error for this node
      const lastResult = results.length > 0 ? results[results.length - 1] : null;
      const lastError = errors.length > 0 ? errors[errors.length - 1] : null;

      // Compare timestamps
      const resultTimestamp = lastResult ? (lastResult.timestamp.seconds * 1e9 + lastResult.timestamp.nanos) : 0;
      const errorTimestamp = lastError ? (lastError.timestamp.seconds * 1e9 + lastError.timestamp.nanos) : 0;

      // Determine the status
      // If the last result is newer than the last error, and the last result is successful, return 'success'
      if (lastResult && resultTimestamp >= errorTimestamp) {
        status[node] = {
          color: 'success',
          icon: 'mdi-check',
          name: node,
          resource: {status: 'Success', message: `Node ${node} operating normally`}
        };
        // If the last result is newer than the last error, and the last result is not successful, return 'error'
      } else if (lastError) {
        status[node] = {
          color: 'error',
          icon: 'mdi-close',
          name: node,
          resource: {error: lastError.error}
        };
      } else {
        status[node] = {
          color: 'warning',
          icon: 'mdi-alert',
          name: node,
          resource: {
            error: {
              code: 2,
              message: `Node ${node} status unknown`
            }
          }
        };
      }
    });

    nodeStatus.value = status;
    updatingNodeStatus.value = false;
  };

  // Watcher to trigger updateNodeStatus when trackedHubNodeResults or trackedHubNodeErrors change
  watchEffect(async () => {
    const newNodeResults = trackedHubNodeResults.value;
    const newNodeErrors = trackedHubNodeErrors.value;
    await updateNodeStatus(newNodeResults, newNodeErrors);
  });

  /**
   * Returns an array of errored nodes
   *
   * @type {import ('vue').ComputedRef<{
   *  color: string,
   *  icon: string,
   *  loading: boolean,
   *  resource: {
   *    error: Error
   *  },
   *  type: string
   * }[]>} erroredNodes
   */
  const erroredNodes = computed(() => {
    const nodeStatuses = Object.values(nodeStatus.value);
    return nodeStatuses.filter(node => node.color === 'error');
  });

  // ------------------- Lifecycle ------------------- //
  // Set up an interval to check every second
  let checkInterval;

  const setUpdatingInterval = () => {
    return () => {
      const currentTime = Date.now();
      if (currentTime - lastFetch.value >= 15000) { // Check if more than 15 seconds have passed
        getEnrollmentAndListHubNodes();
      }
    };
  };

  onMounted(async () => {
    await getEnrollmentAndListHubNodes();

    // Set up an interval to check every second
    checkInterval = setInterval(setUpdatingInterval(), 1000); // Checking every second
  });

  // Closing enrollment on unmount
  onUnmounted(() => {
    closeResource(enrollmentValue);
    closeResource(listHubNodesValue);
    closeResource(testEnrollmentValue);
    closeResource(testHubNodeValue);
    clearInterval(checkInterval);
  });

  return {
    isLoading,
    enrollmentValue,
    listHubNodesValue,
    testEnrollmentValue,
    enrollmentStatus,
    triggerListHubNodesAction,
    getEnrollmentAndListHubNodes,
    nodeStatus,
    updatingNodeStatus,
    serverChipType,
    displayedChips,
    erroredNodes
  };
}
