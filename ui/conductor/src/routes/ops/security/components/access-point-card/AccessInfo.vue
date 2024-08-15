<template>
  <div>
    <!-- Status bar -->
    <status-bar
        :title="props.device.title"
        :status-bar-color="statusColor"
        :door-status-text="doorStatusText"
        :show-close="props.showClose"
        @click:close="emit('click:close')"/>

    <!-- Grant area -->
    <div class="d-flex flex-column justify-space-between">
      <v-card-text class="text-h6 text-white font-weight-regular d-flex flex-row pa-0 px-4 pt-4">
        <span>Last access:</span>
        <span :class="[`${color}--text`, 'ml-auto font-weight-bold']">
          {{ formatString(grantStates) }}
        </span>
      </v-card-text>
      <v-card-text class="text-h6 text-white d-flex flex-column pt-1" style="max-width: 350px">
        <span class="text-subtitle-1"> {{ user.name }} {{ user.cardId }} </span>
      </v-card-text>

      <!-- Alert/Acknowledge area -->
      <v-card-actions v-if="alert.length" class="mt-4">
        <v-col class="mx-0 px-0" cols="align-self" style="max-width: 370px">
          <v-list-item class="px-2">
            <v-tooltip location="bottom">
              <template #activator="{ props }">
                <v-list-item-title class="text-uppercase" v-bind="props">
                  {{ alert.description }}
                </v-list-item-title>
              </template>
              {{ alert.description }}
            </v-tooltip>
            <v-list-item-subtitle>
              {{ alert.createTime }}
            </v-list-item-subtitle>
          </v-list-item>
        </v-col>
        <v-spacer/>
        <v-col class="mx-0 px-0" cols="1">
          <acknowledgement-btn
              :ack="alert.acknowledgement"
              @acknowledge="notifications.setAcknowledged(true, alert, hubName)"
              @unacknowledge="notifications.setAcknowledged(false, alert, hubName)"/>
        </v-col>
      </v-card-actions>
    </div>
  </div>
</template>

<script setup>
import {closeResource} from '@/api/resource';
import {grantNamesByID} from '@/api/sc/traits/access';
import AcknowledgementBtn from '@/routes/ops/notifications/AcknowledgementBtn.vue';
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import {useNotifications} from '@/routes/ops/notifications/notifications.js';
import useAlertsApi from '@/routes/ops/notifications/useAlertsApi';
import StatusBar from '@/routes/ops/security/components/access-point-card/StatusBar.vue';
import {useStatus} from '@/routes/ops/security/components/access-point-card/useStatus';
import {useHubStore} from '@/stores/hub';
import {computed, onBeforeUnmount, reactive} from 'vue';

const props = defineProps({
  accessAttempt: {
    type: Object,
    default: () => {
    }
  },
  openClose: {
    type: Object,
    default: () => {
    }
  },
  statusLog: {
    type: Object,
    default: () => {
    }
  },
  device: {
    type: Object,
    default: () => {
    }
  },
  paused: {
    type: Boolean,
    default: false
  },
  showClose: {
    type: Boolean,
    default: false
  }
});
const notifications = useNotifications();
const emit = defineEmits(['click:close']);

const {alertMetadata} = useAlertMetadata();
const {color} = useStatus(
    () => props.accessAttempt,
    () => props.statusLog
);

const statusColor = computed(() => {
  if (props.openClose) {
    const percentage = props.openClose.statesList[0].openPercent;

    if (percentage === 0) return 'success'; // closed
    if (percentage > 0 && percentage <= 100) return 'warning'; // moving and open
    else return 'grey'; // unknown
  } else return color.value;
});

const doorStatusText = computed(() => {
  if (props.openClose) {
    const percentage = props.openClose.statesList[0].openPercent;

    if (percentage === 0) return 'closed'; // closed
    if (percentage > 0 && percentage < 100) return 'moving'; // moving
    if (percentage === 100) return 'open'; // open
    else return 'unknown'; // unknown
  } else return undefined;
});

// ----------------- Access Attempt ----------------- //
const grantId = computed(() => props.accessAttempt?.grant);
const grantStates = computed(() => {
  return grantNamesByID[grantId.value || 0].toLowerCase();
});
const formatString = (str) => {
  return str.split('_').join(' ').charAt(0).toUpperCase() + str.split('_').join(' ').slice(1);
};

// ----------------- Alerts ----------------- //
const hubStore = useHubStore();
const hubName = computed(() => hubStore.hubNode?.name ?? '');
const query = reactive({
  createdNotBefore: undefined,
  createdNotAfter: undefined,
  severityNotAbove: undefined,
  severityNotBelow: undefined,
  floor: undefined,
  zone: undefined,
  subsystem: undefined,
  source: props.device.source === '' ? undefined : props.device.source,
  acknowledged: undefined,
  resolved: false,
  resolvedNotBefore: undefined,
  resolvedNotAfter: undefined
});

const alerts = reactive(useAlertsApi(hubName, query));
alerts.pageSize = 1;

const alert = computed(() => {
  if (!alerts.allItems) {
    return {};
  }

  const alertData = alerts.allItems[0];

  if (!alertData) {
    return {};
  }

  return {
    ...alertData,
    createTime: alertData.createTime
  };
});

const user = computed(() => {
  if (!props.accessAttempt?.actor?.displayName && !props.accessAttempt?.actor?.idsMap?.[0]?.[1]) {
    return {
      name: 'No card details provided'
    };
  }

  return {
    name: props.accessAttempt?.actor?.displayName ?? 'Unknown',
    cardId: `(${props.accessAttempt?.actor?.idsMap?.[0]?.[1] ?? 'Unknown'})`
  };
});

onBeforeUnmount(() => {
  closeResource(alertMetadata.value);
  closeResource(props.accessAttempt);
  closeResource(props.statusLog);
  closeResource(alerts.listPageTracker);
  closeResource(alerts.pullResource);
});
</script>
<style lang="scss" scoped>
.granted {
  color: var(--v-success-base);
  transition: color 0.5s ease-in-out;
}

.denied,
.forced,
.failed {
  color: var(--v-error-base);
}

.pending,
.aborted,
.tailgate {
  color: var(--v-warning-base);
}

.grant_unknown {
  color: grey;
}
</style>
