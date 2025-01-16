<template>
  <v-card class="history-card">
    <v-list-subheader class="text-title-caps-large text-neutral-lighten-3">Download History</v-list-subheader>
    <v-container>
      <v-list-subheader class="text-title-caps-small text-white">From</v-list-subheader>
      <v-text-field
          v-model="startDate"
          label="DATE"
          type="date"
          class="date-picker"></v-text-field>
      <v-list-subheader class="text-title-caps-small text-white">To</v-list-subheader>
      <v-text-field
          v-model="endDate"
          label="DATE"
          type="date"
          class="date-picker"></v-text-field>
    </v-container>
    <v-container class="d-flex align-center justify-center">
      <v-container v-if="fetchingHistory">
        <v-row class="d-flex align-center justify-center">
          <v-progress-circular color="primary" indeterminate/>
        </v-row>
        <v-row class="d-flex align-center justify-center">
          <v-label>Downloading, please wait...</v-label>
        </v-row>
      </v-container>
      <v-btn v-else id="download-history" class="mt-2 mr-4 elevation-0" @click="downloadHistory(name)">
        Download History
      </v-btn>
    </v-container>
  </v-card>
</template>

<script setup>
import {airQualityRecordToObject, useListAirQualityHistory} from '@/traits/airQuality/airQuality.js';
import {downloadCSVRows} from '@/util/downloadCSV.js';
import {ref} from 'vue';

const startDate = ref();
const endDate = ref();
const fetchingHistory = ref(false);

const dateTimeProp = (obj) => {
  return `${obj.toLocaleDateString()} ${obj.toLocaleTimeString()}`;
};

const historyHeaders = [
  {title: 'Record Time', val: (a) => dateTimeProp(a.recordTime)},
  {title: 'C02', val: (a) => a.airQuality.carbonDioxideLevel},
  {title: 'VOC', val: (a) => a.airQuality.volatileOrganicCompounds}
];

defineProps({
  name: {
    type: String,
    required: true
  }
});

function addDay(date) {
  const result = new Date(date);
  result.setDate(result.getDate() + 1);
  return result;
}

async function downloadHistory(n) {
  fetchingHistory.value = true;
  const baseRequest = /** @type {ListAirQualityHistoryRequest.AsObject} */ {
    name: n,
    period: {
      startTime: startDate.value,
      endTime: addDay(endDate.value)
    },
    pageSize: 1000
  };

  const csvRows =
      /** @type {string[][]} */
      [historyHeaders.map(h => h.title)];
  while (true) {
    const response = await useListAirQualityHistory(baseRequest);
    for (let record of response.airQualityRecordsList) {
      record = airQualityRecordToObject(record);
      csvRows.push(historyHeaders.map(h => h.val(record)));
    }

    if (!response.nextPageToken) {
      break;
    }
    baseRequest.pageToken = response.nextPageToken;
  }
  const filename = `${n}_airquality_history.csv`;
  downloadCSVRows(filename, csvRows);
  fetchingHistory.value = false;
}
</script>

<style scoped>
.v-list-item {
  min-height: auto;
}

.v-progress-linear {
  width: auto;
}

.date-picker {
  padding-left: 15px;
  padding-right: 15px;
}

#download-history {
  justify-self: center;
  fill: #0d47a1;
}
</style>
