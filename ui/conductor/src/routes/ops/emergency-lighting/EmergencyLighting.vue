<template>
  <content-card>
    <v-card-title>
      <h4 class="text-h4">Emergency Lighting</h4>
      <v-spacer/>
      <v-btn color="primary" @click="downloadCSV">Download CSV <v-icon right>mdi-download</v-icon></v-btn>
    </v-card-title>
    <v-data-table/>
  </content-card>
</template>

<script setup>
import {getReportCSV} from '@/api/sc/traits/lighting-test';
import {reactive} from 'vue';
import {newActionTracker} from '@/api/resource';
import ContentCard from '@/components/ContentCard.vue';


const csvTracker = reactive(
    /** @type {ActionTracker<ReportCSV.AsObject>} */ newActionTracker()
);

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
