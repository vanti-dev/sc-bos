<template>
  <content-card>
    <v-data-table
        :headers="headers"
        :items="alerts.pageItems"
        disable-sort
        :server-items-length="queryTotalCount"
        :item-class="rowClass"
        :options.sync="dataTableOptions"
        :footer-props="{itemsPerPageOptions}"
        :loading="alerts.loading"
        class="pt-4"
        @click:row="showNotification">
      <template #top>
        <filters
            :floor.sync="query.floor"
            :floor-items="floors"
            :zone.sync="query.zone"
            :zone-items="zones"
            :subsystem.sync="query.subsystem"
            :subsystem-items="subsystems"
            :acknowledged.sync="query.acknowledged"
            :resolved.sync="query.acknowledged"/>
      </template>
      <template #item.createTime="{ item }">
        {{ item.createTime.toLocaleString() }}
      </template>
      <template #item.subsystem="{ item }">
        <subsystem-icon size="20px" :subsystem="item.subsystem" no-default/>
      </template>
      <template #item.source="{ item }">
        <v-tooltip bottom>
          <template #activator="{ on }">
            <span v-on="on">{{ formatSource(item.source) }}</span>
          </template>
          <span>{{ item.source }}</span>
        </v-tooltip>
      </template>
      <template #item.severity="{ item }">
        <v-tooltip v-if="item.resolveTime" bottom>
          <template #activator="{ on }">
            <span v-on="on">RESOLVED</span>
          </template>
          Was:
          <span :class="notifications.severityData(item.severity).color">
            {{ notifications.severityData(item.severity).text }}
          </span>
        </v-tooltip>
        <span v-else :class="notifications.severityData(item.severity).color">
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
  </content-card>
</template>
<script setup>
import ContentCard from '@/components/ContentCard.vue';
import SubsystemIcon from '@/components/SubsystemIcon.vue';
import Acknowledgement from '@/routes/ops/notifications/Acknowledgement.vue';
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import Filters from '@/routes/ops/notifications/Filters.vue';
import {useNotifications} from '@/routes/ops/notifications/notifications.js';
import useAlertsApi from '@/routes/ops/notifications/useAlertsApi';
import {useHubStore} from '@/stores/hub';
import {usePageStore} from '@/stores/page';
import {computed, reactive, ref, watch} from 'vue';

const notifications = useNotifications();
const alertMetadata = useAlertMetadata();
const hubStore = useHubStore();
const pageStore = usePageStore();

const query = reactive({
  createdNotBefore: undefined,
  createdNotAfter: undefined,
  severityNotAbove: undefined,
  severityNotBelow: undefined,
  floor: undefined,
  zone: undefined,
  subsystem: undefined,
  source: undefined,
  acknowledged: undefined,
  resolved: false,
  resolvedNotBefore: undefined,
  resolvedNotAfter: undefined
});

const dataTableOptions = ref({
  itemsPerPage: 20,
  page: 1
});
const itemsPerPageOptions = [20, 50, 100];

const name = computed(() => hubStore.hubNode?.name ?? '');
const alerts = reactive(useAlertsApi(name, query));
watch(dataTableOptions, () => {
  alerts.pageSize = dataTableOptions.value.itemsPerPage;
  alerts.pageIndex = dataTableOptions.value.page - 1;
}, {deep: true, immediate: true});

const floors = computed(() => Object.keys(alertMetadata.floorCountsMap).sort());
const zones = computed(() => Object.keys(alertMetadata.zoneCountsMap).sort());
const subsystems = computed(() => Object.keys(alertMetadata.subsystemCountsMap).sort());

// How many query fields are not undefined.
const queryFieldCount = computed(() => Object.values(query).filter(value => value !== undefined).length);
// How many items are there using the current query.
// This isn't always accurate, but we do our best.
const queryTotalCount = computed(() => {
  const fieldCount = queryFieldCount.value;
  switch (fieldCount) {
    case 0:
      return alertMetadata.totalCount;
    case 1:
      if (query.floor !== undefined) return alertMetadata.floorCountsMap[query.floor];
      if (query.zone !== undefined) return alertMetadata.zoneCountsMap[query.zone];
      if (query.acknowledged !== undefined) return alertMetadata.acknowledgedCountMap[query.acknowledged];
      if (query.resolved !== undefined) return alertMetadata.resolvedCountMap[query.resolved];
      if (query.severityNotAbove !== undefined) {
        let total = 0;
        for (const [level, count] of Object.entries(alertMetadata.severityCountsMap)) {
          if (level <= query.severityNotAbove) total += count;
        }
        return total;
      }
      if (query.severityNotBelow !== undefined) {
        let total = 0;
        for (const [level, count] of Object.entries(alertMetadata.severityCountsMap)) {
          if (level >= query.severityNotBelow) total += count;
        }
        return total;
      }
      break;
    case 2:
      if (query.acknowledged !== undefined && query.resolved !== undefined) {
        const key = [
          query.acknowledged ? 'ack' : 'nack',
          query.resolved ? 'resolved' : 'unresolved'
        ].join('_');
        return alertMetadata.needsAttentionCountsMap[key];
      }
      if (query.severityNotBelow !== undefined && query.severityNotAbove !== undefined) {
        let total = 0;
        for (const [level, count] of Object.entries(alertMetadata.severityCountsMap)) {
          if (level <= query.severityNotAbove) total += count;
          if (level >= query.severityNotBelow) total += count;
        }
        return total;
      }
      break;
  }
  return undefined;
});

const allHeaders = [
  {text: 'Timestamp', value: 'createTime', width: '15em'},
  {value: 'subsystem', width: '20px', class: 'pl-2 pr-0', cellClass: 'pl-2 pr-0'},
  {text: 'Source', value: 'source', width: '15em'},
  {text: 'Floor', value: 'floor', width: '10em'},
  {text: 'Zone', value: 'zone', width: '10em'},
  {text: 'Severity', value: 'severity', width: '9em', align: 'center'},
  {text: 'Description', value: 'description', width: '100%'},
  {text: 'Acknowledged', value: 'acknowledged', align: 'center', width: '12em'}
];
// We don't include _some_ headers we're filtering out to avoid repetition,
// for example if we're filtering to show Floor1, then all rows would show Floor1 in that column which we don't need to
// see over and over.
const headers = computed(() => {
  return allHeaders.filter(header => {
    if (!['floor', 'zone', 'subsystem', 'source', 'acknowledged'].includes(header.value)) return true;
    return query[header.value] === undefined;
  });
});

const formatSource = (source) => {
  const parts = source.split('/');
  return parts[parts.length - 1];
};
const rowClass = (item) => {
  if (item.resolveTime) return 'resolved';
  return '';
};

/**
 * Shows the device in the sidebar
 *
 * @param {*} item
 */
async function showNotification(item) {
  pageStore.showSidebar = true;
  pageStore.sidebarTitle = item.source;
  pageStore.sidebarData = {name, item};
}
</script>

<style scoped>
:deep(table) {
  table-layout: fixed;
}

:deep(.resolved) {
  color: #FFF5 !important;
}
</style>

