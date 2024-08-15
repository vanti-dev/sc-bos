<template>
  <v-menu
      location="bottom left"
      :close-on-content-click="false"
      content-class="elevation-0"
      max-height="600px"
      max-width="550px"
      min-width="400px"
      tile>
    <template #activator="{props}">
      <v-btn
          class="py-1 px-3"
          style="text-align: center"
          variant="text"
          v-bind="props">
        <span class="text-title mr-1">Smart Core OS:</span>
        <span :class="`text-title-bold text-uppercase ${generalStatus}`">
          {{ statusText }}
        </span>
      </v-btn>
    </template>

    <v-card class="elevation-0 mt-4 pb-1" min-width="400px" style="border: 1px solid var(--v-neutral-lighten2)">
      <v-card-title class="text-subtitle-1 mb-0 pb-0 mt-n1 mb-2">
        Smart Core Status
        <span
            class="ml-2 font-weight-light"
            style="font-size: 10px; margin-bottom: -1px">
          {{ isLoading ? '- Checking...' : 'Updated ' + timeAgo }}
        </span>
        <v-spacer/>
        <v-tooltip location="left">
          <template #activator="{ props }">
            <v-btn

                v-bind="props"
                :class="['mb-0', {'rotate-icon': isRefreshing}]"
                rounded="circle"
                size="small"
                style="padding-left: 1px;"
                @click="triggerRefresh">
              <v-icon size="18">mdi-reload</v-icon>
            </v-btn>
          </template>
          <div class="d-flex flex-column">
            <span>Check Now</span>
          </div>
        </v-tooltip>
      </v-card-title>
      <v-card-text class="d-flex flex-row justify-center align-center mb-n1 mt-4">
        <!-- Display chips and status alerts -->
        <v-chip class="bg-neutral-lighten-1" size="small">UI</v-chip>
        <template v-for="(chip, index) in statusPopupSetup">
          <v-divider class="mx-2" style="width: 10px; max-width: 10px;" :key="index + '-divider'"/>
          <status-alert
              v-if="chip.id.includes('Status')"
              :key="chip.id + '-status'"
              :is-clickable="chip.isClickable"
              :color="chip.color"
              :icon="chip.icon"
              :resource="chip.resource"
              :single="chip.single"/>
          <v-chip
              v-else
              :key="chip.id + '-chip'"
              :class="chip.color"
              :disabled="chipDisabled"
              size="small"
              :to="navigateToNodes(chip.to)">
            {{ chip.label }}
          </v-chip>
        </template>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup>
import {DAY, HOUR, MINUTE, SECOND, useNow} from '@/components/now';
import StatusAlert from '@/components/StatusAlert.vue';
import useAuthSetup from '@/composables/useAuthSetup';
import useSmartCoreStatus from '@/composables/useSmartCoreStatus';
import {formatTimeAgo} from '@/util/date';
import {formatErrorMessage} from '@/util/error';
import {computed, ref, watch} from 'vue';

const isRefreshing = ref(false);
const {hasNoAccess, isLoggedIn} = useAuthSetup();
const {
  isLoading,
  listHubNodesValue,
  enrollmentStatus,
  getEnrollmentAndListHubNodes,
  nodeStatus,
  displayedChips
} = useSmartCoreStatus();

// ------ Computed Properties ------ //
/**
 * Returns the general status of the UI, nodes, and enrollment
 *
 * @type {import ('vue').ComputedRef<string>} generalStatus
 */
const generalStatus = computed(() => {
  const ui = uiStatus.value;
  const isUIError = ui.color === 'error';
  let hasError = false;

  for (const status of Object.values(nodeStatus.value)) {
    if (status.color === 'error') hasError = true;
  }

  if (!!networkIssue.value.error) {
    return 'error--text';
  } else if (isUIError || hasError) {
    return 'warning--text';
  } else {
    return 'success--text text--lighten-4';
  }
});

/**
 * Returns the text to be displayed in the status popup button
 *
 * @type {import ('vue').ComputedRef<string>} statusText
 */
const statusText = computed(() => {
  if (generalStatus.value.includes('error')) {
    if (networkIssue.value) {
      return 'offline';
    } else {
      return 'online';
    }
  } else if (generalStatus.value.includes('warning')) {
    return 'online';
  } else if (generalStatus.value.includes('success')) {
    return 'online';
  } else {
    return 'offline';
  }
});

/**
 * Check if there is a network issue (focusing on Unknown status and network related errors)
 *
 * @type {import ('vue').ComputedRef<{error: {code: number, message: string}} | boolean>} networkIssue
 */
const networkIssue = computed(() => {
  const {error} = enrollmentStatus.value;


  if (error) {
    const statusUnknown = error.error.code === 2;
    const networkRelated = error.error.message.toLowerCase().includes('http');

    if (statusUnknown && networkRelated) return error;
    else return false;
  }

  return false;
});

const chipDisabled = computed(() => {
  return !!networkIssue.value.error;
});

/**
 * Returns the actual UI status from the enrollment status
 *
 * @type {import ('vue').ComputedRef<{
 *   color: string,
 *   icon: string,
 *   resource: {
 *     error: {
 *       code: number,
 *       message: string
 *     },
 *     name?: string
 *   }
 * }>} uiStatus
 */
const uiStatus = computed(() => {
  if (networkIssue.value) {
    // Error state
    return createStatusObject(
        'error',
        'mdi-close',
        {
          error: networkIssue.value.error,
          name: null
        }
    );
  }

  // Active connection state
  return createStatusObject(
      'success',
      'mdi-check',
      {
        error: {
          code: 0,
          message: 'Connection active'
        }
      }
  );
});

/**
 * Returns the status of the server/gateway to hub connection
 *
 * @type {import ('vue').ComputedRef<{
 *  color: string,
 *  icon: string,
 *  resource: {
 *    error: {
 *      code: number,
 *      message: string
 *    }
 *  }
 * }>} serverToHubStatus
 */
const serverToHubStatus = computed(() => {
  const {isTested, testLoading} = enrollmentStatus.value;

  if (!isTested) {
    return testLoading ?
        createStatusObject('success', 'mdi-check', {error: {code: 2, message: 'Testing in progress'}}) :
        createStatusObject('error', 'mdi-close', {error: {code: 14, message: 'Unable to test connection'}});
  }

  return isTested.error !== '' ?
      createStatusObject('warning', 'mdi-alert', {
        error: {
          code: isTested.code,
          message: formatErrorMessage(isTested.error)
        }
      }) :
      createStatusObject('success', 'mdi-check', {error: {code: 0, message: 'Connection test successful'}});
});


/**
 * Returns the overall node status
 *
 * @type {import ('vue').ComputedRef<{
 *  color: string,
 *  icon: string,
 *  resource: {
 *    error: {
 *      code: number,
 *      message: string
 *    }
 *  }
 * }>} nodeOverallStatus
 */
const nodeOverallStatus = computed(() => {
  const nodes = nodeStatus.value;
  const nodesList = listHubNodesValue.response?.nodesList?.length > 0;
  const nodesStatuses = Object.values(nodes);

  if (!nodesList || !nodesStatuses.length) {
    return createStatusObject('success', null, {error: {code: 14, message: 'No nodes available'}});
  }

  let hasError = false;
  let hasSuccess = false;

  for (const status of nodesStatuses) {
    if (status.color === 'error') {
      hasError = true;
    }
    if (status.color === 'success') {
      hasSuccess = true;
    }
  }

  if (hasError && !hasSuccess) {
    return createStatusObject('error', null, {error: {code: 14, message: 'All nodes are in error state'}});
  } else if (hasError) {
    const erroredNodes = nodesStatuses.filter(node => node.color === 'error');
    return createStatusObject('warning', null, {errors: erroredNodes}, false);
  } else {
    return createStatusObject('success', null, {error: {code: 0, message: 'All nodes are operating normally'}});
  }
});

/**
 * Returns the chips and alerts to be displayed in the status popup
 *
 * @type {import ('vue').ComputedRef<{
 *  color: string,
 *  icon: string,
 *  resource: {
 *    error: {
 *      code: number,
 *      message: string
 *    }
 *  }
 * }[]>} statusPopupSetup
 */
const statusPopupSetup = computed(() => {
  // Create an array of chips to be displayed in the status popup
  // Set a default chip for the UI status
  const chips = [
    {id: 'uiStatus', ...uiStatus.value}
  ];

  // Create a chip for the server/gateway to hub status
  const serverHubStatusObj = {id: 'serverHubStatus', ...serverToHubStatus.value};
  // Check if the enrollment status is enrolled and tested
  const isEnrolledTested = enrollmentStatus.value.enrollment && enrollmentStatus.value.isTested;

  // Function to add a chip to the chips array
  const addChip = (chipId, chipData, includeStatusObj = false) => {
    if (displayedChips.value.includes(chipId)) { // Check if the chip should be displayed
      chips.push(chipData); // Add the chip to the chips array
      if (isEnrolledTested && includeStatusObj) { // Check if the enrollment status is enrolled and tested
        chips.push(serverHubStatusObj); // Add the server/gateway to hub status chip to the chips array
      }
    }
  };

  // Populate the chips array with the appropriate chips
  addChip('server', {color: 'neutral lighten-4', id: 'server', label: 'Server'}, true);
  addChip('gateway', {color: 'accent', id: 'gateway', label: 'Gateway'}, true);
  addChip('hub', {color: 'primary', id: 'hub', label: 'Hub'});

  // Check if the nodes chip should be displayed
  if (displayedChips.value.includes('nodes')) {
    const nodeStatusColor = nodeOverallStatus.value.color;
    const nodeError = nodeStatusColor === 'error';
    const nodeWarning = nodeStatusColor === 'warning';

    // Add the node status chip to the chips array first, so it is displayed between the Hub and Nodes chips
    chips.push({id: 'nodeStatus', ...nodeOverallStatus.value});
    // Then add the nodes chip to the chips array
    chips.push({
      color: nodeError ? 'error' : nodeWarning ? 'red darken-2' : 'neutral lighten-2',
      id: 'nodes',
      label: 'Nodes',
      to: '/system/components'
    });
  }

  return chips;
});


// ------ Methods ------ //

// Trigger a refresh of the enrollment and node status
const triggerRefresh = async () => {
  try {
    isRefreshing.value = true;
    await getEnrollmentAndListHubNodes();
  } finally {
    setTimeout(() => {
      isRefreshing.value = false;
    }, 1000); // Wait 1 second before setting isRefreshing to false, this allows the icon to spin
  }
};


// Navigate to the nodes page if the user has access
const navigateToNodes = (to) => (to && isLoggedIn && !hasNoAccess(to)) ? to : null;


// ------ Helpers ------ //
const createStatusObject = (color, icon, resource, single = true) => {
  if (icon) {
    return {
      color,
      icon,
      resource,
      single
    };
  } else {
    let icn = 'mdi-check';

    if (color === 'error') {
      icn = 'mdi-close';
    } else if (color === 'warning') {
      icn = 'mdi-alert';
    }

    return {
      color,
      icon: icn,
      resource,
      single
    };
  }
};


// Create a lastChecked timestamp (for second to be used in the status popup
const {now} = useNow(SECOND);
const lastChecked = ref(null);

// Update lastChecked timestamp when isLoading changes
watch(isLoading, (isLoading) => {
  if (!isLoading) {
    lastChecked.value = new Date();
  }
}, {immediate: true});

// Create a timeAgo computed property to display time in words
const timeAgo = computed(() => {
  if (!lastChecked.value) return 'Never';
  return formatTimeAgo(lastChecked.value, now.value, MINUTE, HOUR, DAY);
});

</script>

<style lang="scss" scoped>
.popup {
  height: 100%;
  width: 100%;
  overflow: auto;

  &__status {
    top: 10px;
    min-height: 100%;
    max-height: 600px;
  }
}

.rotate-icon {
  animation: rotation 1s infinite linear;
}

@keyframes rotation {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
