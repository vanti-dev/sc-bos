@@ -0,0 +1,147 @@
<template>
  <v-card
      color="#40464D"
      elevation="0"
      dark
      min-width="400px"
      height="100%"
      max-height="240px">
    <div
        :class="[
          `door-${accessPointCardData.notification.type}`, 'd-flex flex-row align-center']"
        style="height: 36px;">
      <v-card-title class="text-uppercase text-body-large">
        {{ accessPointCardData.notification.device }}
      </v-card-title>
      <v-spacer/>
      <v-card-title class="text-body-large font-weight-bold text-uppercase">
        {{ accessPointCardData.notification.type }}
      </v-card-title>
    </div>

    <v-card-text class="text-h6 white--text font-weight-regular d-flex flex-row pa-0 px-4 pt-4">
      <span>Last access:</span>
      <span :class="[`access-${accessPointCardData.access.type}` ,'ml-auto font-weight-bold text-uppercase']">
        {{ accessPointCardData.access.type }}
      </span>
    </v-card-text>
    <v-card-text class="text-h6 white--text d-flex flex-column pt-1">
      <span class="text-subtitle-1">
        {{ accessPointCardData.user.name }} ({{ accessPointCardData.user.userId }})
      </span>
      <span class="text-subtitle-2">{{ accessPointCardData.user.accessTime }}</span>
    </v-card-text>

    <v-card-actions v-if="!isAcknowledged" class="pt-0">
      <v-list-item class="px-2">
        <v-list-item-content>
          <v-list-item-title class="text-uppercase">
            {{ accessPointCardData.entry.type }}
          </v-list-item-title>
          <v-list-item-subtitle>
            {{ accessPointCardData.entry.time }}
          </v-list-item-subtitle>
        </v-list-item-content>

        <v-row class="justify-end pr-2">
          <Acknowledgement/>
        </v-row>
      </v-list-item>
    </v-card-actions>
  </v-card>
</template>

<script setup>
import {computed, reactive, ref} from 'vue';
import useAlertsApi from '../../notifications/useAlertsApi';

import {useHubStore} from '@/stores/hub';

import Acknowledgement from '@/routes/ops/notifications/Acknowledgement.vue';

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  isAcknowledged: {
    type: Boolean,
    default: false
  }
});

const hubStore = useHubStore();

const hubName = computed(() => hubStore.hubNode?.name ?? '');

const notification = ref({
  device: 'Test-Device-Name',
  type: 'alert',
  time: '2021-09-01 12:00:00'
});

const user = ref({
  name: 'Test User',
  userId: 'test-user-id',
  accessTime: '2021-09-01 12:00:00'
});

const access = ref({
  type: 'granted'
});

const entry = ref({
  type: 'Forced Entry',
  time: '2021-09-01 12:00:00'
});


const accessPointCardData = computed(() => {
  return {
    notification: notification.value,
    user: user.value,
    access: access.value,
    entry: entry.value
  };
});


const query = computed(() => {
  return {
    source: props.name
  };
});

const alerts = reactive(useAlertsApi(hubName, query));
alerts.pageSize = 10;
</script>

<style lang="scss" scoped>
.access {
  &-granted {
    color: green;
  }

  &-denied {
    color: #DE4F75;
  }
}

.door {
  &-open {
    background-color: #00AAC1;
  }

  &-closed {
    background-color: green;
  }

  &-alert {
    background-color: #D0043C;
  }

  &-offline {
    background-color: #C17C00;
  }
}
</style>
