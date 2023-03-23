<template>
  <content-card>
    <v-card-title>
      <h4 class="text-h4">Emergency Lighting</h4>
      <v-spacer/>
      <v-btn color="primary" @click="downloadCSV">Download CSV <v-icon right>mdi-download</v-icon></v-btn>
    </v-card-title>
    <v-data-table
        :headers="headers"
        :items="lightHealth">
      <template #item.faultsList="{value}">
        <span class="text-title-bold success--text text--lighten-3" v-if="value.length === 0">OK</span>
        <span class="text-title-bold error--text text--lighten-1" v-else>
          {{ value.map(v => faultToString(v)).join(", ") }}
        </span>
      </template>
      <template #item.updateTime="{value}">{{ parseDate(value.seconds) }}</template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import {faultToString, getReportCSV, listLightHealth} from '@/api/sc/traits/lighting-test';
import {computed, onMounted, reactive} from 'vue';
import {newActionTracker} from '@/api/resource';
import ContentCard from '@/components/ContentCard.vue';

const headers = [
  {text: 'Name', value: 'name'},
  {text: 'Status', value: 'faultsList'},
  {text: 'Updated', value: 'updateTime'}
];

const csvTracker = reactive(
    /** @type {ActionTracker<ReportCSV.AsObject>} */ newActionTracker()
);
const lightHealthTracker = reactive(
    /** @type {ActionTracker<ListLightHealthResponse.AsObject>}*/ newActionTracker()
);

const lightHealth = computed(() => {
  return lightHealthTracker.response?.emergencyLightsList ?? [];
});

onMounted(() => refreshLightHealth());

/**
 *
 */
function refreshLightHealth() {
  listLightHealth(lightHealthTracker);
}

/**
 *
 * @param seconds
 */
function parseDate(seconds) {
  const d = new Date(seconds*1000);
  return d.toLocaleDateString()+' '+d.toLocaleTimeString();
}

/**
 * download the CSV report, fetched from the server
 */
async function downloadCSV() {
  const csvObj = await getReportCSV(csvTracker);
  // base64 decode contents
  const csvContent = atob(csvObj.csv);
  // create fake file to serve
  const anchor = document.createElement('a');
  anchor.href = 'data:text/csv;charset=utf-8,' + encodeURIComponent(csvContent);
  anchor.target = '_blank';
  anchor.download = 'emergency-lighting-report.csv';
  anchor.click();
}

</script>

<style scoped>

</style>
