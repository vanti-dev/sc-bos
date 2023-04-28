<template>
  <v-container fluid>
    <v-card>
      <v-data-table
          :headers="headers"
          :items="tableData"
          :sort-by="['createTime']"
          :sort-desc="[true]"
          :loading="notifications.loading"
          class="pt-4">
        <template #top>
          <filters
              :floor.sync="query.floor"
              :floor-items="floors"
              :zone.sync="query.zone"
              :zone-items="zones"/>
        </template>
        <template #item.createTime="{ item }">
          {{ item.createTime.toLocaleString() }}
        </template>
        <template #item.severity="{ item }">
          <span :class="notifications.severityData(item.severity).color">
            {{ notifications.severityData(item.severity).text }}
          </span>
        </template>
        <template #item.acknowledged="{ item }">
          <acknowledgement
              :ack="item.acknowledgement"
              @acknowledge="notifications.setAcknowledged(true, item, name)"
              @unacknowledge="notifications.setAcknowledged(false, item, name)"/>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>
<script setup>
import {timestampToDate} from '@/api/convpb.js';
import {useErrorStore} from '@/components/ui-error/error';
import Acknowledgement from '@/routes/ops/notifications/Acknowledgement.vue';
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import Filters from '@/routes/ops/notifications/Filters.vue';
import {useNotifications} from '@/routes/ops/notifications/notifications.js';
import {useAppConfigStore} from '@/stores/app-config';
import {useHubStore} from '@/stores/hub';
import {computed, onUnmounted, reactive, ref, watch} from 'vue';

const notifications = useNotifications();
const alertMetadata = useAlertMetadata();

const appConfig = useAppConfigStore();
const hubStore = useHubStore();
const errors = useErrorStore();

const query = reactive({
  createdNotBefore: undefined,
  createdNotAfter: undefined,
  severityNotAbove: undefined,
  severityNotBelow: undefined,
  floor: undefined,
  zone: undefined,
  source: undefined,
  acknowledged: undefined
});

const floors = computed(() => Object.keys(alertMetadata.floorCountsMap).sort());
const zones = computed(() => Object.keys(alertMetadata.zoneCountsMap).sort());

const allHeaders = [
  {text: 'Timestamp', value: 'createTime', width: '15em'},
  {text: 'Floor', value: 'floor', width: '10em'},
  {text: 'Zone', value: 'zone', width: '10em'},
  {text: 'Severity', value: 'severity', width: '9em'},
  {text: 'Description', value: 'description', width: '100%'},
  {text: 'Acknowledged', value: 'acknowledged', align: 'center', width: '12em'}
];
// We don't include _some_ headers we're filtering out to avoid repetition,
// for example if we're filtering to show Floor1, then all rows would show Floor1 in that column which we don't need to
// see over and over.
const headers = computed(() => {
  return allHeaders.filter(header => {
    if (!['floor', 'zone', 'source', 'acknowledged'].includes(header.value)) return true;
    return query[header.value] === undefined;
  });
});

const alertsCollection = ref({});
const name = computed(() => appConfig.config.proxy? hubStore.hubNode.name : '');

watch(() => appConfig.config, () => {
  init();
}, {immediate: true});
watch(() => hubStore.hubNode, () => {
  init();
}, {immediate: true});

let unwatchErrors;

/**
 *
 */
async function init() {
  // check config is loaded
  if (!appConfig.config) return;
  // check hubNode is loaded if proxy is enabled
  if (appConfig.config.proxy && !hubStore.hubNode) return;

  const name = appConfig.config.proxy? hubStore.hubNode.name : '';
  alertsCollection.value = notifications.newCollection(name);
  console.debug('Fetching alert metadata for', name, alertsCollection);
  alertsCollection.value.query(query);
}

watch(alertsCollection, () => {
  // todo: this causes all pages to be loaded, which is not ideal - connect with paging logic
  alertsCollection.value.needsMorePages = true;
  unwatchErrors = errors.registerCollection(alertsCollection);
});

watch(query, () => {
  alertsCollection.value.query(query);
}, {deep: true});


// UI error handling
onUnmounted(() => {
  if (unwatchErrors) unwatchErrors();
  alertsCollection.value.reset();
});

const tableData = computed(() => {
  return Object.values(alertsCollection.value.resources?.value ?? [])
      .map(alert => ({
        ...alert,
        createTime: timestampToDate(alert.createTime),
        acknowledged: notifications.isAcknowledged(alert)
      }));
});
</script>

<style scoped>
:deep(table) {
  table-layout: fixed;
}
</style>

