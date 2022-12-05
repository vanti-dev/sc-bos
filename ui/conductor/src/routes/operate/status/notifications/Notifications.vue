<template>
  <v-container fluid>
    <v-card>
      <v-data-table
          :headers="headers"
          :items="tableData"
          :sort-by="['createTime']"
          :sort-desc="[true]"
      >
        <template #item.createTime="{ item }">
          {{ item.createTime.toLocaleString() }}
        </template>
        <template #item.severity="{ item }">
          {{ severityData(item.severity) }}
        </template>
        <template #item.acknowledged="{ item }">
          <v-simple-checkbox :value="item.acknowledged" @input="setAcknowledged($event, item)"/>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>
<script setup>
import {timestampToDate} from '@/api/convpb.js';
import {newActionTracker, newResourceCollection} from '@/api/resource.js';
import {acknowledgeAlert, listAlerts, pullAlerts, unacknowledgeAlert} from '@/api/ui/alerts';
import {Alert, ListAlertsResponse} from '@bsp-ew/ui-gen/proto/alerts_pb.js';
import {computed, onBeforeMount, reactive, set} from 'vue';

const SeverityStrings = {};
for (const [name, val] of Object.entries(Alert.Severity)) {
  SeverityStrings[val] = name;
}

const name = 'test-ac';

const alerts = reactive(/** @type {ResourceCollection<Alert.AsObject, Alert>} */newResourceCollection()); // holds all the alerts we can show
const fetchingPage = reactive(/** @type {ActionTracker<ListAlertsResponse.AsObject>} */ newActionTracker()); // tracks the fetching of a single page
const tableData = computed(() => {
  return Object.values(alerts.value)
      .map(alert => ({
        ...alert,
        createTime: timestampToDate(alert.createTime),
        acknowledged: Boolean(alert.acknowledgement)
      }))
})

onBeforeMount(async () => {
  pullAlerts({name: 'test-ac'}, alerts);
  try {
    const firstPage = await listAlerts({name, pageSize: 100, pageToken: undefined}, fetchingPage);
    for (let alert of firstPage.alertsList) {
      set(alerts.value, alert.id, alert);
    }
    fetchingPage.response = null;
  } catch (e) {
    console.warn('Error fetching first page', e);
  }
});

function severityData(severity) {
  for (let i = severity; i > 0; i--) {
    if (SeverityStrings[i]) {
      let str = SeverityStrings[i];
      if (i < severity) {
        str += '+' + (severity - i);
      }
      return str;
    }
  }
  return 'unspecified';
}

function setAcknowledged(e, alert) {
  if (e) {
    acknowledgeAlert({name, id: alert.id})
        .catch(err => console.error(err));
  } else {
    unacknowledgeAlert({name, id: alert.id})
        .catch(err => console.error(err));
  }
}

const headers = computed(() => {
  return [
    {text: 'Timestamp', value: 'createTime', width: '14em'},
    {text: 'Floor', value: 'floor', width: '7em'},
    {text: 'Zone', value: 'zone', width: '7em'},
    {text: 'Severity', value: 'severity', width: '10em'},
    {text: 'Description', value: 'description', width: '100%'},
    {text: 'Acknowledged', value: 'acknowledged', align: 'center', width: '12em'},
  ]
});
</script>

<style scoped>
::v-deep(table) {
  table-layout: fixed;
}
</style>

