<template>
  <v-menu
      bottom
      :close-on-content-click="false"
      content-class="elevation-0"
      left
      max-height="600px"
      max-width="550px"
      min-width="400px"
      offset-y
      tile>
    <template #activator="{on, attrs}">
      <v-btn class="py-1" style="text-align: center" text v-bind="attrs" v-on="on">
        <span class="text-title mr-1">Smart Core OS:</span>
        <span :class="`text-title-bold text-uppercase ${generalStatus}--text`">{{ statusText }}</span>
      </v-btn>
    </template>

    <v-card class="elevation-0 mt-4 pb-1" min-width="400px" style="border: 1px solid var(--v-neutral-lighten2)">
      <div class="d-flex flex-row align-center">
        <v-card-title class="text-subtitle-1 mb-0 pb-0 mt-n1 mb-2">
          Smart Core Status
          <span
              class="ml-2 font-weight-light"
              style="font-size: 10px; margin-bottom: -1px">{{
                isLoading ? '- Checking...' : 'Last update: ' + timeAgo
              }}</span>
        </v-card-title>
        <v-spacer/>
        <v-tooltip left>
          <template #activator="{ on, attrs }">
            <v-btn
                v-bind="attrs"
                v-on="on"
                :class="['mr-2', {'rotate-icon': isRefreshing}]"
                icon
                small
                style="padding-left: 1px;"
                @click="triggerRefresh">
              <v-icon size="18">mdi-reload</v-icon>
            </v-btn>
          </template>
          <div class="d-flex flex-column">
            <span>Check Now</span>
          </div>
        </v-tooltip>
      </div>
      <v-row class="d-flex flex-row justify-center mx-4 my-3">
        <div class="d-flex flex-row align-center">
          <v-chip class="neutral lighten-1" small>UI</v-chip>
          <v-divider class="mx-2" style="min-width: 10px"/>

          <!-- Display chips and status alerts -->
          <div v-for="(chip, index) in statusPopupSetup" class="d-flex flex-row" :key="chip.id">
            <status-alert
                v-if="chip.id.includes('Status')"
                :is-clickable="chip.isClickable"
                :color="chip.color"
                :icon="chip.icon"
                :resource="chip.resource"
                :single="chip.single"/>
            <div v-if="!chip.id.includes('Status')" class="d-flex flex-row align-center">
              <v-divider class="mx-2" style="min-width: 10px"/>
              <v-chip :class="chip.color" :disabled="defaultConnection" small :to="navigateToNodes(chip.to)">
                {{ chip.label }}
              </v-chip>
              <v-divider
                  v-if="index !== statusPopupSetup.length - 1"
                  class="mx-2"
                  style="min-width: 10px"/>
            </div>
          </div>
        </div>
      </v-row>
    </v-card>
  </v-menu>
</template>

<script setup>
import {computed, ref, watch} from 'vue';
import useSmartCoreStatus from '@/composables/useSmartCoreStatus';
import useAuthSetup from '@/composables/useAuthSetup';
import StatusAlert from '@/components/StatusAlert.vue';
import {useNow, SECOND, MINUTE, HOUR, DAY} from '@/components/now';

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

const triggerRefresh = async () => {
  isRefreshing.value = true;
  await getEnrollmentAndListHubNodes();
  setTimeout(() => {
    isRefreshing.value = false;
  }, 1000);
};

const navigateToNodes = (to) => {
  if (to && isLoggedIn && !hasNoAccess(to)) {
    return to;
  }

  return null;
};

const statusText = computed(() => {
  if (generalStatus.value === 'error') {
    return 'offline';
  } else if (generalStatus.value === 'warning') {
    return 'online';
  } else if (generalStatus.value === 'success') {
    return 'online';
  } else {
    return 'offline';
  }
});

// Returns the actual UI status from the enrollment status
const uiStatus = computed(() => {
  const es = enrollmentStatus.value;

  if (!es.enrollment && !es.error && !es.loading) {
    // Unavailable state
    return createStatusObject(
        'error',
        'mdi-close',
        {
          error: {
            code: 14,
            message: 'Service is unavailable'
          }
        }
    );
  }

  if (es.loading || es.testLoading) {
    // Loading state
    return createStatusObject(
        'success',
        'mdi-check',
        {
          error: {
            code: 2,
            message: 'Loading'
          }
        }
    );
  }

  if (es.error) {
    // Error state
    return createStatusObject(
        'error',
        'mdi-close',
        {
          error: es.error.error
        }
    );
  }

  if (!es.isTested && !es.testLoading) {
    // Enrolled but not tested state
    return createStatusObject(
        'warning',
        'mdi-alert',
        {
          error: {
            code: 0,
            message: 'Enrollment successful, but not yet tested'
          }
        }
    );
  }

  if (!es.isTested && es.testLoading) {
    // Testing state
    return createStatusObject(
        'warning',
        'mdi-check',
        {
          error: {
            code: 2,
            message: 'Testing in progress'
          }
        }
    );
  }

  // Check if there's an error in the test
  if (es.isTested.error && es.isTested.error !== '') {
    // Enrolled and test failed state
    return createStatusObject(
        'error',
        'mdi-alert',
        {
          error: es.isTested
        }
    );
  }

  // Enrolled with pass state
  return createStatusObject(
      'success',
      'mdi-check',
      {
        error: {
          code: 0,
          message: 'Connection established'
        }
      }
  );
});
// Returns the general status of the UI - Online, Offline, or Error
const generalStatus = computed(() => {
  const ui = uiStatus.value;

  // Check if UI is in error
  const isUIError = ui.color === 'error';

  // Initialize flags
  let hasError = false;
  let allError = true;

  // Check the status of each node
  for (const status of Object.values(nodeStatus.value)) {
    if (status.color === 'error') {
      hasError = true;
    } else {
      allError = false;
    }
  }

  // Determine the general status
  if (isUIError && allError) {
    // Return 'error' only if both UI and all nodes are in error state
    return 'error';
  } else if (isUIError || hasError) {
    // Return 'warning' if UI is in error or some nodes are in error state
    return 'warning';
  } else {
    // Return 'success' if no errors are found
    return 'success';
  }
});

// Returns server/gateway to hub status
const serverToHubStatus = computed(() => {
  const isTested = enrollmentStatus.value.isTested;

  if (!isTested) {
    return createStatusObject(
        'warning',
        'mdi-alert',
        {
          error: {
            code: 2,
            message: 'Not tested'
          }
        }
    );
  }

  if (isTested.error !== '') {
    return createStatusObject(
        'error',
        'mdi-alert',
        {
          error: isTested.error
        }
    );
  }

  if (isTested.error === '') {
    return createStatusObject(
        'success',
        'mdi-check',
        {
          error: {
            code: 0,
            message: 'Connection successful'
          }
        }
    );
  }

  return createStatusObject(
      'warning',
      'mdi-alert',
      {
        error: {
          code: 2,
          message: 'Not tested'
        }
      }
  );
});

// Returns the overall status of the Node Status alert
const nodeOverallStatus = computed(() => {
  const nodes = nodeStatus.value;
  const isAvailable = listHubNodesValue.response?.nodesList?.length > 0 && nodeStatus.value !== {};

  // Check if nodes are not available
  if (!isAvailable) {
    return createStatusObject(
        'success',
        null,
        {
          error: {
            code: 14,
            message: 'No nodes available'
          }
        }
    );
  }

  let hasError = false;
  let allError = true;
  let allSuccess = true;

  for (const status of Object.values(nodes)) {
    if (status.color === 'error') {
      hasError = true;
    } else {
      allError = false;
    }
    if (status.color !== 'success') {
      allSuccess = false;
    }
  }

  if (allError) {
    return createStatusObject(
        'error',
        null,
        {
          error: {
            code: 14,
            message: 'All nodes are in error state'
          }
        }
    );
  } else if (hasError) {
    const erroredNodes = Object.values(nodes).filter((node) => {
      if (node.color === 'error') {
        return {
          name: node.name,
          error: node.resource.error
        };
      }
    });

    return createStatusObject(
        'warning',
        null,
        {
          errors: erroredNodes
        },
        false
    );
  } else if (allSuccess) {
    return createStatusObject(
        'success',
        null,
        {
          error: {
            code: 0,
            message: 'All nodes are operating normally'
          }
        }
    );
  }

  return createStatusObject(
      'warning',
      null,
      {
        error: {
          code: 2,
          message: 'Node status is mixed or unknown'
        }
      }
  );
});

const defaultConnection = computed(() => {
  return displayedChips.value.length === 2; // Default to UI / Server
});

// Returns the chips and alerts to be displayed in the status popup
const statusPopupSetup = computed(() => {
  const uiStatusObj = {
    id: 'uiStatus',
    ...uiStatus.value
  };

  const serverHubStatusObj = {
    id: 'serverHubStatus',
    ...serverToHubStatus.value
  };

  const nodeStatusObj = {
    id: 'nodeStatus',
    ...nodeOverallStatus.value
  };

  const server = {
    color: 'neutral lighten-4',
    id: 'server',
    label: 'Server'
  };

  const gateway = {
    color: 'accent',
    id: 'gateway',
    label: 'Gateway'
  };
  const hub = {
    color: 'primary',
    id: 'hub',
    label: 'Hub'
  };

  const nodeStatusColor = nodeOverallStatus.value.color;
  const nodes = {
    color: nodeStatusColor === 'error' ? 'error' : nodeStatusColor === 'warning' ? 'warning' : 'neutral lighten-2',
    id: 'nodes',
    label: 'Nodes',
    to: '/system/components'
  };


  const chips = [uiStatusObj];
  const isEnrolledTested = enrollmentStatus.value.enrollment && enrollmentStatus.value.isTested;

  if (displayedChips.value.includes('server')) {
    chips.push(server);
    if (isEnrolledTested) chips.push(serverHubStatusObj);
  }

  if (displayedChips.value.includes('gateway')) {
    chips.push(gateway);
    if (isEnrolledTested) chips.push(serverHubStatusObj);
  }

  if (displayedChips.value.includes('hub')) {
    chips.push(hub);
  }

  if (displayedChips.value.includes('nodes')) {
    chips.push(nodeStatusObj);
    chips.push(nodes);
  }

  return chips;
});


// ------ Helpers ------ //
const createStatusObject = (color, icon, resource, single = true) => {
  let setIcon;

  if (!icon) {
    if (color === 'error') {
      setIcon = 'mdi-close';
    } else if (color === 'warning') {
      setIcon = 'mdi-alert';
    } else {
      setIcon = 'mdi-check';
    }
  } else {
    setIcon = icon;
  }

  return {
    color,
    icon: setIcon,
    resource,
    single
  };
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
  return formatTimeAgo(lastChecked.value, now.value);
});

// Format time ago using Intl.RelativeTimeFormat API to display time in words
const formatTimeAgo = (date, now) => {
  const diffInSeconds = (now - date) / 1000;
  const rtf = new Intl.RelativeTimeFormat('en', {numeric: 'auto'});

  if (Math.abs(diffInSeconds) < MINUTE) {
    return rtf.format(-Math.floor(diffInSeconds), 'second');
  } else if (Math.abs(diffInSeconds) < HOUR) {
    return rtf.format(-Math.floor(diffInSeconds / MINUTE), 'minute');
  } else if (Math.abs(diffInSeconds) < DAY) {
    return rtf.format(-Math.floor(diffInSeconds / HOUR), 'hour');
  }
};
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
