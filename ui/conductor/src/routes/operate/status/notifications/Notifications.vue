<template>
  <v-container fluid>
    <v-card>
      <v-data-table
          :headers="headers"
          :items="tableData"
          :sort-by="['createTime']"
          :sort-desc="[true]"
          :loading="notifications.loading"
          class="pt-4"
      >
        <template #top>
          <filters
              :floor.sync="query.floor"
              :floor-items="floors"
              :zone.sync="query.zone"
              :zone-items="zones"
          />
        </template>
        <template #item.createTime="{ item }">
          {{ item.createTime.toLocaleString() }}
        </template>
        <template #item.severity="{ item }">
          <span :class="notifications.severityData(item.severity).color">{{
              notifications.severityData(item.severity).text
            }}</span>
        </template>
        <template #item.acknowledged="{ item }">
          <v-simple-checkbox :value="item.acknowledged" @input="notifications.setAcknowledged($event, item)"/>
        </template>
      </v-data-table>
    </v-card>
  </v-container>
</template>
<script setup>
import {timestampToDate} from '@/api/convpb.js';
import Filters from '@/routes/operate/status/notifications/Filters.vue';
import {useNotifications} from '@/routes/operate/status/notifications/notifications.js';
import {computed, onUnmounted, reactive, watch} from 'vue';

const notifications = useNotifications();

const query = reactive({
  createdNotBefore: undefined,
  createdNotAfter: undefined,
  severityNotAbove: undefined,
  severityNotBelow: undefined,
  floor: undefined,
  zone: undefined,
  source: undefined,
  acknowledged: undefined,
});

// todo: get these from the server, or at least the data, or something more relevant
const floors = computed(() => ['L01', 'L02', 'L03', 'L04', '12']);
const zones = computed(() => ['Z01', 'Z02', 'R01', 'NE'])

const allHeaders = [
  {text: 'Timestamp', value: 'createTime', width: '14em'},
  {text: 'Floor', value: 'floor', width: '7em'},
  {text: 'Zone', value: 'zone', width: '7em'},
  {text: 'Severity', value: 'severity', width: '9em'},
  {text: 'Description', value: 'description', width: '100%'},
  {text: 'Acknowledged', value: 'acknowledged', align: 'center', width: '12em'},
];
// We don't include _some_ headers we're filtering out to avoid repetition,
// for example if we're filtering to show Floor1, then all rows would show Floor1 in that column which we don't need to
// see over and over.
const headers = computed(() => {
  return allHeaders.filter(header => {
    if (!['floor', 'zone', 'source', 'acknowledged'].includes(header.value)) return true;
    return query[header.value] === undefined;
  })
});

/** @type {Collection} */
const collection = notifications.newCollection();
collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead

watch(query, () => collection.query(query), {deep: true, immediate: true})

onUnmounted(() => {
  collection.reset(); // stop listening when the component is unmounted
})

const tableData = computed(() => {
  return Object.values(collection.resources.value)
      .map(alert => ({
        ...alert,
        createTime: timestampToDate(alert.createTime),
        acknowledged: notifications.isAcknowledged(alert)
      }))
})
</script>

<style scoped>
::v-deep(table) {
  table-layout: fixed;
}
</style>

