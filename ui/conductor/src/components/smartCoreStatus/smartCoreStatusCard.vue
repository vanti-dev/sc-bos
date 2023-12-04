<template>
  <v-dialog content-class="elevation-0" max-width="500" overlay-opacity="0.4" v-model="showStatus">
    <template #activator="{ on, attrs }">
      <v-btn class="py-1" style="text-align: center" text v-bind="attrs" v-on="on" @click="showStatus = !showStatus">
        <span class="text-title mr-1">Smart Core OS:</span>
        <span :class="`text-title-bold text-uppercase ${generalStatus}--text`">{{ statusText }}</span>
      </v-btn>
    </template>

    <content-card class="popup__status">
      <v-card class="elevation-0" style="border: 2px solid var(--v-neutral-lighten2)" width="500px">
        <div class="d-flex flex-row align-center">
          <v-card-title class="text-subtitle-1 mb-0 pb-0 mt-n1 mb-2">Connection status</v-card-title>
          <v-spacer/>
          <v-tooltip left>
            <template #activator="{ on, attrs }">
              <v-btn
                  v-bind="attrs"
                  v-on="on"
                  :class="['mr-4', {'rotate-icon': isRefreshing}]"
                  icon
                  x-small
                  @click="triggerRefresh">
                <v-icon>mdi-reload</v-icon>
              </v-btn>
            </template>
            <span>Refresh</span>
          </v-tooltip>
        </div>
        <v-row class="d-flex flex-row justify-center mx-4 my-3">
          <!-- Test Hub to Nodes connection -->
          <div class="d-flex flex-row">
            <v-chip class="neutral lighten-1" small>UI</v-chip>
            <v-divider class="my-auto mx-2" style="min-width: 10px"/>

            <!-- Display chips and status alerts -->
            <div v-for="(chip, index) in statusPopupSetup" class="d-flex flex-row align-center" :key="chip.id">
              <status-alert
                  v-if="chip.id.includes('Status')"
                  :click-action="chip.clickAction"
                  :clickable="chip.isClickable"
                  :color="chip.color"
                  :icon="chip.icon"
                  :loading="chip.loading"
                  :resource="chip.resource"
                  :type="chip.type"/>
              <div v-if="!chip.id.includes('Status')" class="d-flex flex-row">
                <v-divider class="my-auto mx-2" style="min-width: 10px"/>
                <v-chip :class="chip.color" small>
                  {{ chip.label }}
                </v-chip>
                <v-divider
                    v-if="index !== statusPopupSetup.length - 1"
                    class="my-auto mx-2"
                    style="min-width: 10px"/>
              </div>
            </div>
          </div>
        </v-row>
      </v-card>
    </content-card>
  </v-dialog>
</template>

<script setup>
import {computed, ref} from 'vue';
import useSmartCoreStatus from '@/composables/useSmartCoreStatus';

import StatusAlert from '@/components/StatusAlert.vue';
import ContentCard from '@/components/ContentCard.vue';

const showStatus = ref(false);
const isRefreshing = ref(false);

const {
  listHubNodesValue,
  enrollmentStatus,
  triggerListHubNodesAction,
  getEnrollmentAndListHubNodes,
  nodeStatus,
  updatingNodeStatus,
  displayedChips
} = useSmartCoreStatus();

const triggerRefresh = async () => {
  isRefreshing.value = true;
  await getEnrollmentAndListHubNodes();
  setTimeout(() => {
    isRefreshing.value = false;
  }, 1000);
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
    return {
      color: 'error',
      icon: 'mdi-close',
      loading: false,
      resource: {error: {code: 14, message: 'Service is unavailable'}},
      type: 'error'
    };
  }

  if (es.loading || es.testLoading) {
    // Loading state
    return {
      color: 'warning',
      icon: 'mdi-loading',
      loading: es.loading || es.testLoading,
      resource: {error: {code: 2, message: 'Loading'}},
      type: 'warning'
    };
  }

  if (es.error) {
    // Error state
    return {
      color: 'error',
      icon: 'mdi-close',
      loading: false,
      resource: {error: es.error.error},
      type: 'error'
    };
  }

  if (!es.isTested && !es.testLoading) {
    // Enrolled but not tested state
    return {
      color: 'warning',
      icon: 'mdi-alert',
      loading: false,
      resource: {error: {code: 0, message: 'Enrollment successful, but not yet tested'}},
      type: 'warning'
    };
  }

  if (!es.isTested && es.testLoading) {
    // Testing state
    return {
      color: 'warning',
      icon: 'mdi-',
      loading: true,
      resource: {error: {code: 2, message: 'Testing in progress'}},
      type: 'warning'
    };
  }

  if (es.isTested) {
    // Check if there's an error in the test
    if (es.isTested.error && es.isTested.error !== '') {
      // Enrolled and test failed state
      return {
        color: 'error',
        icon: 'mdi-alert',
        loading: false,
        resource: {error: es.isTested},
        type: 'error'
      };
    }
  }

  // Enrolled with pass state
  return {
    color: 'success',
    icon: 'mdi-check',
    loading: false,
    resource: {status: {code: 0, message: 'Enrolled'}},
    type: 'success'
  };
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
const serverToHubStatue = computed(() => {
  const isTested = enrollmentStatus.value.isTested;

  if (!isTested) {
    return {
      color: 'warning',
      icon: 'mdi-alert',
      loading: false,
      resource: {error: {code: 2, message: 'Not tested'}},
      type: 'warning'
    };
  }

  if (isTested.error) {
    return {
      color: 'error',
      icon: 'mdi-alert',
      loading: false,
      resource: {error: isTested.error},
      type: 'error'
    };
  }

  return {
    color: 'success',
    icon: 'mdi-check',
    loading: false,
    resource: {status: {code: 0, message: 'Tested'}},
    type: 'success'
  };
});

// Returns the overall status of the Node Status alert
const nodeOverallStatus = computed(() => {
  const nodes = nodeStatus.value;
  const isAvailable = listHubNodesValue.response?.nodesList?.length > 0 && nodeStatus.value !== {};
  const isUpdating = updatingNodeStatus.value;

  // Check if nodes are not available
  if (!isAvailable) {
    return {
      color: 'error',
      icon: 'mdi-information-outline',
      loading: false,
      resource: {error: {code: 14, message: 'No nodes available'}},
      type: 'error'
    };
  }

  // Check if the node status is being updated
  if (isUpdating) {
    return {
      color: 'warning',
      icon: 'mdi-loading',
      loading: true,
      resource: {error: {code: 2, message: 'Updating node status'}},
      type: 'warning'
    };
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
    return {
      color: 'error',
      icon: 'mdi-close',
      loading: false,
      resource: {error: {code: 14, message: 'All nodes are in error state'}},
      type: 'error'
    };
  } else if (hasError) {
    return {
      color: 'warning',
      icon: 'mdi-alert',
      loading: false,
      resource: {error: {code: 2, message: 'Some nodes have errors'}},
      type: 'warning'
    };
  } else if (allSuccess) {
    return {
      color: 'success',
      icon: 'mdi-check',
      loading: false,
      resource: {status: {code: 0, message: 'All nodes are operating normally'}},
      type: 'success'
    };
  }

  return {
    color: 'warning',
    icon: 'mdi-information-outline',
    loading: false,
    resource: {error: {code: 2, message: 'Node status is mixed or unknown'}},
    type: 'warning'
  };
});

const statusPopupSetup = computed(() => {
  const uiStatusObj = {
    id: 'uiStatus',
    ...uiStatus.value
  };

  const serverHubStatusObj = {
    isClickable: false,
    id: 'serverHubStatus',
    ...serverToHubStatue.value
  };

  const nodeStatusObj = {
    clickAction: async () => await triggerListHubNodesAction(),
    isClickable: true,
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

  const nodes = {
    color: 'neutral lighten-2',
    id: 'nodes',
    label: 'Nodes'
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

</script>

<style lang="scss" scoped>
.popup {
  height: 100%;

  &__status {
    position: absolute;
    width: 500px;
    height: auto;
    max-height: 600px;
    left: auto;
    right: 155px;
    top: 55px;

    &::after {
      content: '';
      position: absolute;
      top: -7px;
      right: 20px;
      width: 0;
      height: 0;
      border-style: solid;
      border-width: 0 15px 15px 0;
      border-color: transparent transparent var(--v-neutral-lighten2) transparent;
      transform: rotate(135deg);
    }
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
