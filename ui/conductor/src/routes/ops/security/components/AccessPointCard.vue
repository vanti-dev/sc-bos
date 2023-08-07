<template>
  <v-card color="#40464D" elevation="0" dark min-width="400px" height="100%" max-height="240px">
    <div
        :class="[notifications.severityData(alert.severity).background, 'd-flex flex-row align-center']"
        style="height: 36px">
      <v-card-title class="text-uppercase text-body-large">
        {{ props.name }}
      </v-card-title>
      <v-spacer/>
      <v-card-title class="text-body-large font-weight-bold text-uppercase">
        {{ notifications.severityData(alert.severity).text }}
      </v-card-title>
    </div>

    <v-card-text class="text-h6 white--text font-weight-regular d-flex flex-row pa-0 px-4 pt-4">
      <span>Last access:</span>
      <span :class="[grantStates, 'ml-auto font-weight-bold']">
        {{ formatString(grantStates) }}
      </span>
    </v-card-text>
    <v-card-text class="text-h6 white--text d-flex flex-column pt-1">
      <span class="text-subtitle-1"> {{ user.name }} ({{ user.cardId }}) </span>
      <!-- <span class="text-subtitle-2">{{ accessPointCardData.user.accessTime }}</span> -->
    </v-card-text>

    <v-card-actions v-if="!isAcknowledged" class="pt-0">
      <v-col class="mx-0 px-0" cols="align-self" style="max-width: 350px">
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
      <v-col class="mx-0 px-0" cols="1">
        <Acknowledgement
            :ack="alert.acknowledgement"
            @acknowledge="notifications.setAcknowledged(true, alert, hubName)"
            @unacknowledge="notifications.setAcknowledged(false, alert, hubName)"/>
      </v-col>
    </v-card-actions>
  </v-card>
</template>

<script setup>
import {computed, reactive, ref} from 'vue';
import {AccessAttempt} from '@sc-bos/ui-gen/proto/access_pb';
import {useHubStore} from '@/stores/hub';
import {useNotifications} from '../../notifications/notifications';
import useAlertsApi from '../../notifications/useAlertsApi';

import Acknowledgement from '@/routes/ops/notifications/Acknowledgement.vue';

const props = defineProps({
  value: {
    type: Object,
    default: () => ({})
  },
  source: {
    type: String,
    default: ''
  },
  name: {
    type: String,
    default: ''
  },
  floor: {
    type: String,
    default: ''
  },
  isAcknowledged: {
    type: Boolean,
    default: false
  }
});

const notifications = useNotifications();

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
  source: props.source,
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
    alerts,
    cardId: props.value?.actor?.idsMap[0][1] ?? 'Unknown'
  };
});
</script>

<style lang="scss" scoped>
.granted {
  color: green;
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
