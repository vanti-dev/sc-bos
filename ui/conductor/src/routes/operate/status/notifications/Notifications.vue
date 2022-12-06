<template>
  <v-container fluid>
    <v-card>
      <v-data-table
          :headers="headers"
          :items="tableData"
          :sort-by="['createTime']"
          :sort-desc="[true]"
          :loading="notifications.loading"
      >
        <template #item.createTime="{ item }">
          {{ item.createTime.toLocaleString() }}
        </template>
        <template #item.severity="{ item }">
          {{ notifications.severityData(item.severity) }}
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
import {useNotifications} from '@/routes/operate/status/notifications/notifications.js';
import {computed} from 'vue';

const notifications = useNotifications();

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
const tableData = computed(() => {
  return Object.values(notifications.alerts.value)
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

