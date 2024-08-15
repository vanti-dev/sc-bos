<template>
  <div class="ml-2">
    <v-row v-if="!props.overviewPage" class="mt-0 ml-0 pl-0">
      <h3 class="text-h3 pt-2 pb-6">Notifications</h3>
      <v-spacer/>
      <v-btn class="mt-2 mr-3" color="neutral" @click="alerts.exportData('Notifications')">
        Export CSV...
      </v-btn>
    </v-row>

    <content-card :class="['px-8', {'mt-8 px-4': !props.overviewPage}]">
      <v-data-table-server
          :headers="headers"
          :items="alerts.pageItems"
          disable-sort
          :items-length="queryTotalCount"
          :row-props="rowProps"
          :options.sync="dataTableOptions"
          :footer-props="setFooterProps"
          :loading="alerts.loading"
          class="pt-4"
          :class="{ 'hide-pagination': modifyFooter }"
          @click:row="showNotification">
        <template #top>
          <v-row
              :class="[
                'd-flex flex-row align-center mb-2 mt-1 ml-0 pl-0 mr-1',
                {'mb-4 mr-n2': props.overviewPage}
              ]">
            <h3 v-if="props.overviewPage" class="text-h4">
              Notifications
            </h3>
            <v-spacer/>
            <filter-choice-chips
                :ctx="filterCtx"
                :class="['mr-2 mt-n2', {'mt-n2': props.overviewPage}, {'mb-4': !props.overviewPage}]"/>
            <span :class="['mt-n2', {'mb-4 mt-2': !props.overviewPage}, {'mr-2': props.overviewPage}]">
              <filter-btn :ctx="filterCtx" tile/>
            </span>
            <v-tooltip location="top">
              <template #activator="{ props }">
                <v-btn
                    v-if="props.overviewPage"
                    class="mt-n2 rounded"
                    color="neutral"
                    elevation="0"
                    height="36"
                    size="small"
                    tile
                    v-bind="props"
                    width="34"
                    @click="alerts.exportData('Notifications')">
                  <v-icon>mdi-file-download</v-icon>
                </v-btn>
              </template>
              Export CSV...
            </v-tooltip>
            <v-tooltip location="top">
              <template #activator="{ props }">
                <v-btn
                    v-if="!props.overviewPage"
                    class="mt-n2 ml-2 mb-4 rounded"
                    color="neutral"
                    elevation="0"
                    height="36"
                    width="34"
                    size="small"
                    tile
                    v-bind="props"
                    @click="toggleManualEntry">
                  <v-icon size="30">mdi-plus</v-icon>
                </v-btn>
              </template>
              <span>Add New Entry</span>
            </v-tooltip>

            <v-expansion-panels v-if="!props.overviewPage" class="mt-n3 mb-3" flat v-model="manualEntryPanel">
              <v-expansion-panel>
                <v-expansion-panel-text>
                  <div class="text-subtitle-2 pl-0 text-body-1">Manual Notification Entry Form</div>
                  <v-row class="align-center mr-2">
                    <v-col cols="self-align">
                      <v-text-field
                          v-model="manualEntryForm.description"
                          label="Description"
                          density="compact"
                          variant="outlined"
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
                          density="compact"
                          hide-details
                          :items="[
                            {label: 'INFO', value: 9},
                            {label:'WARN', value: 13},
                            {label:'ALERT', value:17},
                            {label:'DANGER', value: 21}
                          ]"
                          item-title="label"
                          item-value="value"
                          label="Severity"
                          variant="outlined"/>
                    </v-col>
                    <v-col cols="2">
                      <v-select
                          v-model="manualEntryForm.floor"
                          density="compact"
                          hide-details
                          :items="floors"
                          label="Floor"
                          variant="outlined"/>
                    </v-col>
                    <v-col cols="2">
                      <v-select
                          v-model="manualEntryForm.zone"
                          density="compact"
                          hide-details
                          :items="zones"
                          label="Zone"
                          variant="outlined"/>
                    </v-col>
                    <v-col cols="auto">
                      <v-btn @click="addManualEntry" :disabled="!manualEntryForm.description" color="primary">
                        Add
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-expansion-panel-text>
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
          <v-tooltip location="bottom">
            <template #activator="{ props }">
              <span v-bind="props">{{ formatSource(item.source) }}</span>
            </template>
            <span>{{ item.source }}</span>
          </v-tooltip>
        </template>
        <template #item.severity="{ item }">
          <v-tooltip v-if="item.resolveTime" location="bottom">
            <template #activator="{ props }">
              <span v-bind="props">RESOLVED</span>
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
          <acknowledgement-btn
              :ack="item.acknowledgement"
              @acknowledge="notifications.setAcknowledged(true, item, name)"
              @unacknowledge="notifications.setAcknowledged(false, item, name)"/>
        </template>
      </v-data-table-server>
    </content-card>
  </div>
</template>
<script setup>
import {newActionTracker} from '@/api/resource.js';
import {createAlert} from '@/api/ui/alerts.js';
import ContentCard from '@/components/ContentCard.vue';
import FilterBtn from '@/components/filter/FilterBtn.vue';
import FilterChoiceChips from '@/components/filter/FilterChoiceChips.vue';
import useFilterCtx from '@/components/filter/filterCtx.js';
import SubsystemIcon from '@/components/SubsystemIcon.vue';
import AcknowledgementBtn from '@/routes/ops/notifications/AcknowledgementBtn.vue';
import {useAlertMetadata} from '@/routes/ops/notifications/alertMetadata';
import {severityData, useNotifications} from '@/routes/ops/notifications/notifications.js';
import useAlertsApi from '@/routes/ops/notifications/useAlertsApi';
import {useHubStore} from '@/stores/hub';
import {useSidebarStore} from '@/stores/sidebar';
import {Alert} from '@sc-bos/ui-gen/proto/alerts_pb';
import {computed, onUnmounted, reactive, ref, watch} from 'vue';
import deepEqual from 'fast-deep-equal';

const props = defineProps({
  overviewPage: {
    type: Boolean,
    default: false
  },
  forceQuery: {
    type: Object, /** @type {import('@sc-bos/ui-gen/proto/alerts_pb').Alert.Query.AsObject} */
    default: null
  }
});

const notifications = useNotifications();
const alertMetadata = useAlertMetadata();
const hubStore = useHubStore();
const sidebar = useSidebarStore();

const manualEntryValue = reactive(newActionTracker());
const manualEntryPanel = ref(null);
const manualEntryForm = ref({
  source: 'manual',
  description: undefined,
  severity: undefined,
  subsystem: undefined,
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
    subsystem: undefined,
    floor: undefined,
    zone: undefined
  };
};

const dataTableOptions = ref({
  itemsPerPage: 20,
  page: 1
});
const itemsPerPageOptions = [20, 50, 100];
const modifyFooter = computed(() => queryMetadataCount.value === undefined);

const floors = computed(() => Object.keys(alertMetadata.floorCountsMap).sort());
const zones = computed(() => Object.keys(alertMetadata.zoneCountsMap).sort());
const subsystems = computed(() => Object.keys(alertMetadata.subsystemCountsMap).sort());

const filterOpts = computed(() => {
  // we only add filters that can affect the output, i.e. no floor filter if nothing has a floor.
  const filters = [];
  const defaults = [];
  // Acknowledged
  if (!props.forceQuery?.hasOwnProperty('acknowledged')) {
    filters.push({
      key: 'acknowledged',
      icon: 'mdi-checkbox-marked-circle-outline', title: 'Acknowledgement', type: 'boolean',
      valueToString(value) {
        switch (value) {
          case true:
            return 'Acknowledged';
          case false:
            return 'Unacknowledged';
          default:
            return 'All';
        }
      }
    });
  }
  // Floor
  if (!props.forceQuery?.hasOwnProperty('floor')) {
    const items = floors.value
        // we can't query for empty strings anyway.
        .filter(s => Boolean(s));
    if (items.length > 0) {
      filters.push({key: 'floor', icon: 'mdi-layers-triple-outline', title: 'Floor', type: 'list', items});
    }
  }
  // Severity
  if (!props.forceQuery?.hasOwnProperty('severityNotAbove') &&
      !props.forceQuery?.hasOwnProperty('severityNotBelow')) {
    filters.push({
      key: 'severity', // maps to severityNotBelow and severityNotAbove
      icon: 'mdi-alert-box-outline', title: 'Severity', type: 'range',
      items: Object.entries(Alert.Severity)
          .map(([, v]) => ({value: v, title: severityData(v).text}))
          .filter((item) => item.value !== 0) // skip UNSPECIFIED
    });
  }
  // Subsystem
  if (!props.forceQuery?.hasOwnProperty('subsystem')) {
    const items = subsystems.value
        // we can't query for empty strings anyway.
        .filter(s => Boolean(s));
    if (items.length > 0) {
      filters.push({key: 'subsystem', icon: 'mdi-file-tree', title: 'Subsystem', type: 'list', items});
    }
  }
  // Zone
  if (!props.forceQuery?.hasOwnProperty('zone')) {
    const items = zones.value
        // we can't query for empty strings anyway.
        .filter(s => Boolean(s));
    if (items.length > 0) {
      filters.push({key: 'zone', icon: 'mdi-select-all', title: 'Zone', type: 'list', items});
    }
  }

  if (!props.forceQuery?.hasOwnProperty('resolved')) {
    filters.push({
      key: 'resolved',
      icon: 'mdi-checkbox-marked-circle-outline', title: 'Resolution', type: 'boolean',
      valueToString(value) {
        switch (value) {
          case true:
            return 'Resolved';
          case false:
            return 'Unresolved';
          default:
            return 'All';
        }
      }
    });
    defaults.push({filter: 'resolved', value: false});
  }
  return {filters, defaults};
});
const filterCtx = useFilterCtx(filterOpts);

const nonFilterableQueryFields = computed(() => {
  return /** @type {import('@sc-bos/ui-gen/proto/alerts_pb').Alert.Query.AsObject} */ props.forceQuery ?? {};
});
const queryFields = computed(() => {
  const res = /** @type {import('@sc-bos/ui-gen/proto/alerts_pb').Alert.Query.AsObject} */ {};
  for (const choice of filterCtx.sortedChoices.value) {
    if (choice.value === undefined || choice.value === null) continue;
    switch (choice.filter) {
      case 'acknowledged':
        res.acknowledged = choice.value;
        break;
      case 'floor':
        res.floor = choice.value?.value ?? choice.value;
        break;
      case 'severity':
        const {from, to} = choice.value;
        if (from) {
          res.severityNotBelow = from.value;
        }
        if (to) {
          res.severityNotAbove = to.value;
        }
        break;
      case 'subsystem':
        res.subsystem = choice.value?.value ?? choice.value;
        break;
      case 'zone':
        res.zone = choice.value?.value ?? choice.value;
        break;
      case 'resolved':
        res.resolved = choice.value;
        break;
    }
  }
  return res;
});
const query = computed(() => {
  return {...nonFilterableQueryFields.value, ...queryFields.value};
});

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

const queryFieldCount = computed(() => Object.values(query).filter((value) => value !== undefined).length);

/**
 *  Calculate the total number of items in the query
 *
 * @param {import('@sc-bos/ui-gen/proto/alerts_pb').Alert.Query.AsObject} query
 * @return {number|undefined}
 */
function calculateQueryMetadataCount(query) {
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
const queryMetadataCount = computed(() => calculateQueryMetadataCount(query.value));
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
        showCurrentPage: true,
        itemsPerPageOptions,
        pagination: {}
      };
    } else {
      // If there is no next page token, then we know there are no more pages available.
      // We can block the next button by setting the itemsLength to the total number of items.
      return {
        showCurrentPage: true,
        itemsPerPageOptions,
        pagination: {
          itemsLength: alerts.allItems.length
        }
      };
    }
  } else {
    // If there are less than 2 query fields, then we can use the default pagination options.
    return {
      showCurrentPage: true,
      itemsPerPageOptions
    };
  }
});

// Watch the query object for changes
watch(
    query,
    (oldQuery, newQuery) => {
      if (deepEqual(oldQuery, newQuery)) return; // avoid reactivity churn
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
  const q = query.value;
  return allHeaders.filter((header) => {
    if (!['floor', 'zone', 'subsystem', 'source'].includes(header.value)) return true;
    return q[header.value] === undefined;
  });
});

const formatSource = (source) => {
  const parts = source.split('/');
  return parts[parts.length - 1];
};
const rowProps = ({item}) => {
  if (item.resolveTime) return {class: 'resolved'};
  return {};
};

/**
 * Shows the device in the sidebar
 *
 * @param {*} item
 */
async function showNotification(item) {
  sidebar.visible = true;
  sidebar.title = item.source;
  sidebar.data = {metadata: sidebar.data?.metadata, notification: {name, item}};
}

onUnmounted(() => {
  sidebar.closeSidebar();
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
