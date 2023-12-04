import {computed, onMounted, onUnmounted, reactive, ref, watch, watchEffect} from 'vue';

import {closeResource, newActionTracker} from '@/api/resource';
import {getEnrollment, testEnrollment} from '@/api/sc/traits/enrollment';
import {testHubNode} from '@/api/sc/traits/hub';
import {useHubStore} from '@/stores/hub';

/**
 * @return {{
 *   isLoading: import('vue').Ref<boolean>,
 *   enrollmentValue: ActionTracker<GetEnrollmentResponse.AsObject>,
 *   listHubNodesValue: ActionTracker<ListHubNodesResponse.AsObject>,
 *   testEnrollmentValue: ActionTracker<TestEnrollmentResponse.AsObject>,
 *   enrollmentStatus: import('vue').ComputedRef<{
 *      enrollment: GetEnrollmentResponse.AsObject | null,
 *      error: Error | null,
 *      loading: boolean,
 *      isTested: boolean | null,
 *      testLoading: boolean
 *   }>,
 *   triggerListHubNodesAction: function(): Promise<void>,
 *   getEnrollmentAndListHubNodes: function(): Promise<void>,
 *   nodeStatus: import('vue').Ref<Record<string, {
 *      color: string,
 *      icon: string,
 *      loading: boolean,
 *      resource: {status: string, message: string} | {error: Error},
 *      type: string
 *   }>>,
 *   updatingNodeStatus: import('vue').Ref<boolean>,
 *   serverChipType: import('vue').ComputedRef<string>,
 *   displayedChips: import('vue').ComputedRef<string[]>
 *  }}
 */
export default function() {
  const {listHubNodesAction} = useHubStore();
  const enrollmentValue = reactive(
      /** @type {ActionTracker<GetEnrollmentResponse.AsObject>} */ newActionTracker()
  );
  const listHubNodesValue = reactive(
      /** @type {ActionTracker<ListHubNodesResponse.AsObject>} */ newActionTracker()
  );
  const testEnrollmentValue = reactive(
      /** @type {ActionTracker<TestEnrollmentResponse.AsObject>} */ newActionTracker()
  );
  const testHubNodeValue = reactive(
      /** @type {ActionTracker<TestHubNodeResponse.AsObject>} */ newActionTracker()
  );

  const isLoading = ref(false);

  // Check if we are enrolled
  const getEnrollmentValue = async () => {
    try {
      console.debug('Getting enrollment...');
      await getEnrollment(enrollmentValue);
      enrollmentValue.error = null;
      await getTestEnrollmentValue();
    } catch (e) {
      console.error('Error fetching enrollment', e);
      enrollmentValue.response = null;
      nodeStatus.value = {};
    }

    return enrollmentValue;
  };

  const getTestEnrollmentValue = async () => {
    try {
      console.debug('Getting test enrollment...');
      await testEnrollment(testEnrollmentValue);
      testEnrollmentValue.error = null;
    } catch (e) {
      console.error('Error fetching test enrollment', e);
      testEnrollmentValue.response = null;
    }

    return testEnrollmentValue;
  };

  const trackedHubNodeErrors = ref({}); // Object to store errors for each node
  const trackedHubNodeResults = ref({}); // Object to store results for each node
  const getTestHubNodeValue = async (nodesList) => {
    if (nodesList.length === 0) return; // Escaping early if the list is empty

    for (const node of nodesList) { // Iterate over each node in the list
      const timestamp = {
        seconds: Math.floor(Date.now() / 1000), // Current time in seconds
        nanos: (Date.now() % 1000) * 1000000 // Current time in nanoseconds
      };

      try {
        console.debug('Testing hub node...');
        const request = {address: node.address};
        await testHubNode(request, testHubNodeValue);

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


  // Watch for changes to the list of hub nodes and test them
  watch(() => listHubNodesValue.response?.nodesList, async (nodesList) => {
    if (!nodesList || nodesList.length === 0) return;
    await getTestHubNodeValue(nodesList);
  }, {deep: true});


  // ------------------- UI Status ------------------- //
  // Generating statuses for the UI
  // Returns enrollment status
  const enrollmentStatus = computed(() => {
    return {
      enrollment: enrollmentValue.response || null,
      error: enrollmentValue.error || testEnrollmentValue.error || null,
      loading: enrollmentValue.loading,
      isTested: testEnrollmentValue.response || testEnrollmentValue.error || null,
      testLoading: testEnrollmentValue.loading
    };
  });


  const nodeStatus = ref({}); // Initialize nodeStatus as a reactive reference
  const updatingNodeStatus = ref(false); // Initialize updatingNodeStatus as a reactive reference
  // Function to update node status
  const updateNodeStatus = async (nodeResults, nodeErrors) => {
    const status = {}; // Initialize status as an empty object
    updatingNodeStatus.value = true; // Set updatingNodeStatus to true
    // Collect all unique node names from both results and errors
    const allNodes = await new Set([...Object.keys(nodeResults), ...Object.keys(nodeErrors)]);

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
          loading: false,
          resource: {status: 'Success', message: `Node ${node} operating normally`},
          type: 'success'
        };
        // If the last result is newer than the last error, and the last result is not successful, return 'error'
      } else if (lastError) {
        status[node] = {
          color: 'error', icon: 'mdi-cross', loading: false, resource: {error: lastError.error}, type: 'error'
        };
      } else {
        status[node] = {
          color: 'warning',
          icon: 'mdi-alert',
          loading: false,
          resource: {error: {message: `Node ${node} status unknown`}},
          type: 'warning'
        };
      }
    });

    nodeStatus.value = status;
    updatingNodeStatus.value = false;
  };

  // Watcher to react to changes in trackedHubNodeResults and trackedHubNodeErrors
  watchEffect(async () => {
    const newNodeResults = trackedHubNodeResults.value;
    const newNodeErrors = trackedHubNodeErrors.value;
    await updateNodeStatus(newNodeResults, newNodeErrors);
  });


  // Returns a type for the chip representing the server/hub/gateway depending on the successful checks
  const serverChipType = computed(() => {
    if (enrollmentStatus.value.enrollment && listHubNodesValue.response?.nodesList) return 'gateway';
    if (listHubNodesValue.response) return 'hub';
    return 'server';
  });

  // Returns an array of chips to be displayed
  const displayedChips = computed(() => {
    const es = enrollmentStatus.value;

    // default chips
    const chips = ['ui', serverChipType.value];
    // add hub chip if enrolled
    const isEnrolled = es.enrollment && es.isTested?.error === '';
    // add nodes chip if nodes are present
    const hubNodeResponse = listHubNodesValue.response !== null;

    if (isEnrolled) chips.push('hub');

    if (hubNodeResponse) chips.push('nodes');

    return chips;
  });

  // Returns a now timestamp for the last successful check
  const lastFetch = ref(Date.now());

  // Execute both GetEnrollment and ListHubNodes (or pull) against the server
  // Can be used to manually trigger the action
  const triggerListHubNodesAction = async () => {
    updatingNodeStatus.value = true;
    listHubNodesValue.response = null;
    await listHubNodesAction(listHubNodesValue);
    updatingNodeStatus.value = false;
  };
  const getEnrollmentAndListHubNodes = async () => {
    isLoading.value = true;
    await getEnrollmentValue();
    await triggerListHubNodesAction();
    lastFetch.value = Date.now(); // Update the last execution time
    isLoading.value = false;
  };


  // ------------------- Lifecycle ------------------- //
  // Set up an interval to check every second
  let checkInterval;

  onMounted(async () => {
    await getEnrollmentAndListHubNodes();

    // Set up an interval to check every second
    checkInterval = setInterval(() => {
      const currentTime = Date.now();
      if (currentTime - lastFetch.value >= 15000) { // Check if more than 15 seconds have passed
        getEnrollmentAndListHubNodes();
      }
    }, 1000); // Checking every second
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
    displayedChips
  };
}
