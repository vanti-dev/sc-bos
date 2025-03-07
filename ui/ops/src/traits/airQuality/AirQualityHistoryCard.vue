<template>
  <v-menu :close-on-content-click="false">
    <template #activator="{ props }">
      <v-btn :icon="true" elevation="0" size="x-small" @click="resetMenu" v-bind="props" class="mt-n1 mr-n1">
        <v-icon size="20">mdi-dots-vertical</v-icon>
      </v-btn>
    </template>
    <v-card class="history-card" min-width="420">
      <v-card-text>
        <div class="d-flex align-start">
          <v-date-input
              v-model="dateRange" multiple="range" :readonly="fetchingHistory"
              label="Download History" placeholder="from - to" persistent-placeholder prepend-icon=""
              hint="Select a date range to download historical data." persistent-hint
              :error-messages="downloadError"/>
          <div v-tooltip="downloadBtnDisabled || 'Download CSV...'">
            <v-btn
                @click="downloadHistory(name)"
                icon="mdi-file-download" elevation="0" class="ml-2 mr-n2 mt-1"
                :loading="fetchingHistory" :disabled="!!downloadBtnDisabled"/>
          </div>
        </div>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup>
import {airQualityRecordToObject, listAirQualitySensorHistory} from '@/api/sc/traits/air-quality-sensor.js';
import {downloadCSVRows} from '@/util/downloadCSV.js';
import {addDays, startOfDay} from 'date-fns';
import {computed, ref} from 'vue';
import {VDateInput} from 'vuetify/labs/components';

defineProps({
  name: {
    type: String,
    required: true
  }
});

const dateRange = ref([]);
const startDate = computed(() => dateRange.value[0]);
const endDate = computed(() => dateRange.value[dateRange.value.length - 1]);
const downloadBtnDisabled = computed(() => {
  if (dateRange.value.length === 0) {
    return 'No date range selected';
  }
  return '';
});
const fetchingHistory = ref(false);
const downloadError = ref('');

/**
 * Resets the menu to its initial state.
 */
function resetMenu() {
  dateRange.value = [];
  downloadError.value = '';
}

const dateTimeProp = (obj) => {
  return `${obj.toLocaleDateString()} ${obj.toLocaleTimeString()}`;
};

const historyHeaders = [
  {title: 'Record Time', val: (a) => dateTimeProp(a.recordTime)},
  {title: 'C02', val: (a) => a.airQuality.carbonDioxideLevel},
  {title: 'VOC', val: (a) => a.airQuality.volatileOrganicCompounds}
];

/**
 * Downloads the history of the air quality sensor and saves it as a CSV file.
 *
 * @param {string} n
 * @return {Promise<void>}
 */
async function downloadHistory(n) {
  fetchingHistory.value = true;
  const baseRequest = /** @type {ListAirQualityHistoryRequest.AsObject} */ {
    name: n,
    period: {
      startTime: startOfDay(startDate.value),
      endTime: startOfDay(addDays(endDate.value, 1))
    },
    pageSize: 1000
  };

  const csvRows =
      /** @type {string[][]} */
      [historyHeaders.map(h => h.title)];
  while (true) {
    try {
      const response = await listAirQualitySensorHistory(baseRequest, {});
      for (let record of response.airQualityRecordsList) {
        record = airQualityRecordToObject(record);
        csvRows.push(historyHeaders.map(h => h.val(record)));
      }
      if (!response.nextPageToken) {
        break;
      }
      baseRequest.pageToken = response.nextPageToken;
    } catch (error) {
      downloadError.value = error.message;
      break;
    }
  }
  if (csvRows.length > 1) {
    const filename = `${n}_airquality_history.csv`;
    downloadCSVRows(filename, csvRows);
  } else {
    downloadError.value = 'No historical records found for these dates';
  }
  fetchingHistory.value = false;
}
</script>
