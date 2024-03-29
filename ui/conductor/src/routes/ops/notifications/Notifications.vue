<template>
  <div class="ml-3">
    <v-row v-if="!props.overviewPage" class="mt-0 ml-0 pl-0">
      <h3 class="text-h3 pt-2 pb-6">Notifications</h3>
      <v-spacer/>
      <v-btn class="mt-2 mr-3" color="neutral" @click="alerts.exportData('Notifications')">
        Export CSV...
      </v-btn>
    </v-row>

    <content-card :class="['px-8', {'mt-8 px-4': !props.overviewPage}]">
      <v-data-table
          :headers="headers"
          :items="alerts.pageItems"
          disable-sort
          :server-items-length="queryTotalCount"
          :item-class="rowClass"
          :options.sync="dataTableOptions"
          :footer-props="setFooterProps"
          :loading="alerts.loading"
          class="pt-4"
          :class="{ 'hide-pagination': modifyFooter }"
          @click:row="showNotification">
        <template #top>
          <v-row :class="['mt-1 mb-2 ml-0 pl-0', {'mt-n4 mb-2 mr-2': props.overviewPage}]">
            <h3 v-if="props.overviewPage" :class="['text-h3 pt-2 pb-6', {'text-h4': props.overviewPage}]">
              Notifications
            </h3>
            <v-spacer/>
            <filters
                v-if="!props.overviewPage"
                class="mb-4 mt-n2"
                :floor.sync="query.floor"
                :floor-items="floors"
                :zone.sync="query.zone"
                :zone-items="zones"
                :subsystem.sync="query.subsystem"
                :subsystem-items="subsystems"
                :acknowledged.sync="query.acknowledged"
                :resolved.sync="query.acknowledged"/>
            <v-btn v-if="props.overviewPage" class="mr-2" color="primary" @click="alerts.exportData('Notifications')">
              Export CSV...
            </v-btn>
            <v-tooltip bottom>
              <template #activator="{ on }">
                <v-btn
                    v-if="!props.overviewPage"
                    @click="toggleManualEntry"
                    color="primary"
                    width="30"
                    x-small
                    height="36"
                    :class="[{'mt-2 mr-5': !props.overviewPage, 'mt-0 ml-2 mr-0': props.overviewPage}]"
                    v-on="on">
                  <v-icon size="28">mdi-plus</v-icon>
                </v-btn>
              </template>
              <span>Add new entry</span>
            </v-tooltip>

            <v-expansion-panels v-if="!props.overviewPage" class="mt-n3 mb-3" flat v-model="manualEntryPanel">
              <v-expansion-panel>
                <v-expansion-panel-content>
                  <v-subheader class="pl-0 text-body-1">Manual Notification Entry Form</v-subheader>
                  <v-row class="align-center mr-2">
                    <v-col cols="self-align">
                      <v-text-field
                          v-model="manualEntryForm.description"
                          label="Description"
                          dense
                          outlined
                          hide-details
                          maxlength="160">
                        <template #append>
                          <span class="character-counter">
                            {{ manualEntryForm.description ? 160 - manualEntryForm.description.length : 160 }}
                          </span>
                        </template>
                      </v-text-field>
                    </v-col>
                    <v-col cols="2">
                      <v-select
                          v-model="manualEntryForm.severity"
                          dense
                          hide-details
                          :items="[
                            {label: 'INFO', value: 9},
                            {label:'WARN', value: 13},
                            {label:'ALERT', value:17},
                            {label:'DANGER', value: 21}
                          ]"
                          item-text="label"
                          item-value="value"
                          label="Severity"
                          outlined/>
                    </v-col>
                    <v-col cols="2">
                      <v-select
                          v-model="manualEntryForm.floor"
                          dense
                          hide-details
                          :items="floors"
                          label="Floor"
                          outlined/>
                    </v-col>
                    <v-col cols="2">
                      <v-select
                          v-model="manualEntryForm.zone"
                          dense
                          hide-details
                          :items="zones"
                          label="Zone"
                          outlined/>
                    </v-col>
                    <v-col cols="auto">
                      <v-btn @click="addManualEntry" :disabled="!manualEntryForm.description" color="primary">
                        Add
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-expansion-panel-content>
              </v-expansion-panel>
            </v-expansion-panels>
          </v-row>
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
  </div>
</template>
<script setup>
import {newActionTracker} from '@/api/resource.js';
import {createAlert} from '@/api/ui/alerts.js';
import ContentCard from '@/components/ContentCard.vue';
import SubsystemIcon from '@/components/SubsystemIcon.vue';
import Acknowledgement from '@/routes/ops/notifications/Acknowledgement.vue';
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import Filters from '@/routes/ops/notifications/Filters.vue';
import {useNotifications} from '@/routes/ops/notifications/notifications.js';
import useAlertsApi from '@/routes/ops/notifications/useAlertsApi';
import {useHubStore} from '@/stores/hub';
import {usePageStore} from '@/stores/page';
import {computed, onUnmounted, reactive, ref, watch} from 'vue';

const props = defineProps({
  overviewPage: {
    type: Boolean,
    default: false
  },
  zone: {
    type: String,
    default: undefined
  }
});

const notifications = useNotifications();
const alertMetadata = useAlertMetadata();
const hubStore = useHubStore();
const pageStore = usePageStore();
const activeZone = ref(props.zone);
watch(
    () => props.zone,
    (value) => {
      activeZone.value = value;
    },
    {immediate: true}
);

const manualEntryValue = reactive(newActionTracker());
const manualEntryPanel = ref(null);
const manualEntryForm = ref({
  source: 'manual',
  description: undefined,
  severity: undefined,
  floor: undefined,
  zone: undefined
});
const toggleManualEntry = () => {
  manualEntryPanel.value === null ? manualEntryPanel.value = 0 : manualEntryPanel.value = null;
};

const addManualEntry = async () => {
  await createAlert({alert: manualEntryForm.value}, manualEntryValue);
  manualEntryForm.value = {
    source: 'manual',
    description: undefined,
    severity: undefined,
    floor: undefined,
    zone: undefined
  };
};

const query = reactive({
  createdNotBefore: undefined,
  createdNotAfter: undefined,
  severityNotAbove: undefined,
  severityNotBelow: undefined,
  floor: undefined,
  zone: computed({
    get: () => activeZone.value,
    set: (value) => {
      activeZone.value = value;
    }
  }),
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
const modifyFooter = computed(() => queryMetadataCount.value === undefined);

const name = computed(() => hubStore.hubNode?.name ?? '');
const alerts = reactive(useAlertsApi(name, query));
watch(
    dataTableOptions,
    () => {
      alerts.pageSize = dataTableOptions.value.itemsPerPage;
      alerts.pageIndex = dataTableOptions.value.page - 1;
    },
    {deep: true, immediate: true}
);

const floors = computed(() => Object.keys(alertMetadata.floorCountsMap).sort());
const zones = computed(() => Object.keys(alertMetadata.zoneCountsMap).sort());
const subsystems = computed(() => Object.keys(alertMetadata.subsystemCountsMap).sort());

const queryFieldCount = computed(() => Object.values(query).filter((value) => value !== undefined).length);

/**
 *  Calculate the total number of items in the query
 *
 * @return {number|undefined}
 */
function calculateQueryMetadataCount() {
  const fieldCount = queryFieldCount.value;

  /**
   * Get the total number of alerts for the given severity range
   *
   * @return {number|undefined}
   */
  function getSeverityTotal() {
    let total = 0;
    for (const [level, count] of Object.entries(alertMetadata.severityCountsMap)) {
      if (level <= query.severityNotAbove && level >= query.severityNotBelow) total += count;
    }
    return total;
  }

  /**
   * Get the total number of alerts for the given needs attention range
   *
   * @return {number|undefined}
   */
  function getNeedsAttentionTotal() {
    const key = [query.acknowledged ? 'ack' : 'nack', query.resolved ? 'resolved' : 'unresolved'].join('_');
    return alertMetadata.needsAttentionCountsMap[key];
  }

  // Switch on the number of query fields
  switch (fieldCount) {
    case 0: // If there are no query fields, then we can use the total count from the metadata
      return alertMetadata.totalCount;

    case 1: // If there is one query field, then we can use the count from the metadata
      for (const [key, value] of Object.entries(query)) {
        if (value !== undefined) {
          switch (key) {
            case 'subsystem':
              return alertMetadata.subsystemCountsMap[value];
            case 'floor':
              return alertMetadata.floorCountsMap[value];
            case 'zone':
              return alertMetadata.zoneCountsMap[value];
            case 'acknowledged':
              return alertMetadata.acknowledgedCountMap[value];
            case 'resolved':
              return alertMetadata.resolvedCountMap[value];
            case 'severityNotAbove':
              return getSeverityTotal();
            case 'severityNotBelow':
              return getSeverityTotal();
            default:
              return undefined;
          }
        }
      }
      break;
    case 2: // If there are two or more query fields, then we need to calculate the total ourselves
      if (query.acknowledged !== undefined && query.resolved !== undefined) {
        return getNeedsAttentionTotal();
      }
      if (query.severityNotBelow !== undefined && query.severityNotAbove !== undefined) {
        return getSeverityTotal();
      }
      break;
  }

  return undefined;
}

// Calculate the total number of items in the query
const queryMetadataCount = computed(() => calculateQueryMetadataCount());
const queryTotalCount = computed(() => {
  const totalCount = queryMetadataCount.value;

  // If the query metadata count is defined, then we can use it
  if (totalCount >= 0) return totalCount;
  // If there is a next page token, then we know there are more pages available.
  else if (alerts.nextPageToken) return alerts.allItems.length + 1;
  // If there is no next page token, then we know there are no more pages available.
  else return alerts.allItems.length;
});


// Set the footer props
const setFooterProps = computed(() => {
  // If there are more than 2 query fields, then we need to hide the pagination
  if (queryMetadataCount.value === undefined) {
    const nextPageToken = alerts.nextPageToken; // Get the next page token

    // If there is a next page token 'ready' to be used, then we know there are more pages available.
    if (nextPageToken) {
      // Keeping the item pp options and pagination object empty will show the next page button
      return {
        itemsPerPageOptions,
        pagination: {}
      };
    } else {
      // If there is no next page token, then we know there are no more pages available.
      // We can block the next button by setting the itemsLength to the total number of items.
      return {
        itemsPerPageOptions,
        pagination: {
          itemsLength: alerts.allItems.length
        }
      };
    }
  } else {
    // If there are less than 2 query fields, then we can use the default pagination options.
    return {itemsPerPageOptions};
  }
});

// Watch the query object for changes
watch(
    query,
    () => {
      // Reset the page to 1
      dataTableOptions.value = {
        ...dataTableOptions.value,
        page: 1
      };
    },
    {immediate: true, deep: true}
);

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
  return allHeaders.filter((header) => {
    if (!['floor', 'zone', 'subsystem', 'source'].includes(header.value)) return true;
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

onUnmounted(() => {
  pageStore.closeSidebar();
});
</script>

<style lang="scss" scoped>
:deep(table) {
  table-layout: fixed;
}

:deep(.resolved) {
  color: #fff5 !important;
}

.v-data-table :deep(tr:hover) {
  cursor: pointer;
}

.hide-pagination {
  :deep(.v-data-footer__pagination) {
    display: none;
  }
}

.character-counter {
  position: relative;
  text-align: center;
  width: 1.75em;
  height: auto;
  top: 6px;
  bottom: 0;
  right: -5px;
  font-size: 75%;
  color: var(--v-primary-base); /* Adjust color as needed */
}
</style>
