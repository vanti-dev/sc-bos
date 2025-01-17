<template>
  <v-menu
      v-model="menu"
      :close-on-content-click="false">
    <template #activator="{ props }">
      <v-icon
          icon="mdi-dots-vertical"
          v-bind="props"
          @click="hideSnackbar()">
      </v-icon>
    </template>
    <v-card class="history-card">
      <v-container>
        <v-list-subheader class="text-title-caps-small text-white">From</v-list-subheader>
        <v-text-field
            v-model="startDate"
            label="DATE"
            type="date"
            class="date-picker"/>
        <v-list-subheader class="text-title-caps-small text-white">To</v-list-subheader>
        <v-text-field
            v-model="endDate"
            label="DATE"
            type="date"
            class="date-picker"/>
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
        <v-snackbar v-model="snackbar">
          No history found.
        </v-snackbar>
      </v-container>
    </v-card>
  </v-menu>
</template>

<script setup>
import {listMeterReadingHistory, meterReadingRecordToObject} from '@/api/sc/traits/meter.js';
import {downloadCSVRows} from '@/util/downloadCSV.js';
import {ref} from 'vue';

const startDate = ref(Date);
const endDate = ref(Date);
const fetchingHistory = ref(false);
const menu = ref(false);
const snackbar = ref(false);

const dateTimeProp = (obj) => {
  return `${obj.toLocaleDateString()} ${obj.toLocaleTimeString()}`;
};

const historyHeaders = [
  {title: 'Record Time', val: (a) => dateTimeProp(a.recordTime)},
  {title: 'Usage', val: (a) => a.meterReading.usage}
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

function hideSnackbar() {
  snackbar.value = false;
}

async function downloadHistory(n) {
  fetchingHistory.value = true;
  const baseRequest = /** @type {ListMeterReadingHistoryRequest.AsObject} */ {
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
    try {
      const response = await listMeterReadingHistory(baseRequest, {});
      for (let record of response.meterReadingRecordsList) {
        record = meterReadingRecordToObject(record);
        csvRows.push(historyHeaders.map(h => h.val(record)));
      }
      if (!response.nextPageToken) {
        break;
      }
      baseRequest.pageToken = response.nextPageToken;
    } catch (error) {
      snackbar.value = true;
      console.error(error);
      break;
    }
  }
  if (csvRows.length > 1) {
    const filename = `${n}_meter_history.csv`;
    downloadCSVRows(filename, csvRows);
  } else {
    snackbar.value = true;
  }
  fetchingHistory.value = false;
}
</script>

<style scoped>

.date-picker {
  padding-left: 15px;
  padding-right: 15px;
}

#download-history {
  justify-self: center;
}
</style>
