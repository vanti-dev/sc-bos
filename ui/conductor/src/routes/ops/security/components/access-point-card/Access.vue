<template>
  <div>
    <!-- Status bar -->
    <WithStatus v-slot="{ resource }" :name="props.device.name" :paused="props.paused">
      <Status v-bind="resource" :device="props.device" :status-bar-color="setBarColor(resource.live, grantStates)"/>
    </WithStatus>

    <!-- Grant area -->
    <div class="d-flex flex-column justify-space-between">
      <v-card-text class="text-h6 white--text font-weight-regular d-flex flex-row pa-0 px-4 pt-4">
        <span>Last access:</span>
        <span :class="[grantStates, 'ml-auto font-weight-bold']">
          {{ formatString(grantStates) }}
        </span>
      </v-card-text>
      <v-card-text class="text-h6 white--text d-flex flex-column pt-1" style="max-width: 350px">
        <span class="text-subtitle-1"> {{ user.name }} ({{ user.cardId }}) </span>
      </v-card-text>

      <!-- Alert/Acknowledge area -->
      <v-card-actions v-if="alert.source === props.device.source" class="mt-4">
        <v-col class="mx-0 px-0" cols="align-self" style="max-width: 370px">
          <v-list-item class="px-2">
            <v-list-item-content>
              <v-tooltip bottom>
                <template #activator="{ on }">
                  <v-list-item-title class="text-uppercase" v-on="on">
                    {{ alert.description }}
                  </v-list-item-title>
                </template>
                {{ alert.description }}
              </v-tooltip>
              <v-list-item-subtitle>
                {{ alert.createTime }}
              </v-list-item-subtitle>
            </v-list-item-content>
          </v-list-item>
        </v-col>
        <v-spacer/>
        <v-col class="mx-0 px-0" cols="1">
          <Acknowledgement
              :ack="alert.acknowledgement"
              @acknowledge="notifications.setAcknowledged(true, alert, hubName)"
              @unacknowledge="notifications.setAcknowledged(false, alert, hubName)"/>
        </v-col>
      </v-card-actions>
    </div>
  </div>
</template>

<script setup>
import {reactive, computed} from 'vue';
import {AccessAttempt} from '@sc-bos/ui-gen/proto/access_pb';
import {useHubStore} from '@/stores/hub';
import useAlertsApi from '@/routes/ops/notifications/useAlertsApi';
import {useStatusBarStore} from '@/routes/ops/security/components/access-point-card/statusBarStore';

import WithStatus from '@/routes/devices/components/renderless/WithStatus.vue';
import Status from '@/routes/ops/security/components/access-point-card/StatusBar.vue';
import Acknowledgement from '@/routes/ops/notifications/Acknowledgement.vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => {}
  },
  floor: {
    type: String,
    default: ''
  },
  device: {
    type: Object,
    default: () => {}
  },
  paused: {
    type: Boolean,
    default: false
  }
});

const statusBarStore = useStatusBarStore();
const {setBarColor} = statusBarStore;

// ----------------- Access Attempt ----------------- //
const grantId = computed(() => props.value?.grant);
const grantNamesByID = Object.entries(AccessAttempt.Grant).reduce((all, [name, id]) => {
  all[id] = name.toLowerCase();
  return all;
}, {});
const grantStates = computed(() => {
  return grantNamesByID[grantId.value || 0];
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
  floor: props.floor === 'All' ? undefined : props.floor,
  zone: undefined,
  subsystem: undefined,
  source: props.device.source === '' ? undefined : props.device.source,
  acknowledged: undefined,
  resolved: false,
  resolvedNotBefore: undefined,
  resolvedNotAfter: undefined
});

const alerts = reactive(useAlertsApi(hubName, query));
alerts.pageSize = 10;

const alert = computed(() => {
  const alertData = alerts.allItems[0];

  if (!alertData) {
    return {};
  }

  return {
    ...alertData,
    createTime: alertData.createTime.toLocaleString()
  };
});

const user = computed(() => {
  return {
    name: props.value?.actor?.displayName ?? 'Unknown',
    cardId: props.value?.actor?.idsMap?.[0]?.[1] ?? 'Unknown'
  };
});
</script>
<style lang="scss" scoped>
.granted {
  color: green;
  transition: color 0.5s ease-in-out;
}
.denied,
.forced,
.failed {
  color: red;
}
.pending,
.aborted,
.tailgate {
  color: orange;
}
.grant_unknown {
  color: grey;
}
</style>
