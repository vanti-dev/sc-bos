<script setup>
import VueDatePicker from '@vuepic/vue-datepicker';
import '@vuepic/vue-datepicker/dist/main.css';
</script>

<template>
  <v-card class="history-card">
    <v-list-subheader class="text-title-caps-large text-neutral-lighten-3">Air Quality History</v-list-subheader>
    <v-list-subheader class="text-title-caps-small text-white">Start Date</v-list-subheader>
    <VueDatePicker class="date-picker" v-model="startDate" :teleport="true" :format="format"/>
    <v-list-subheader class="text-title-caps-small text-white">End Date</v-list-subheader>
    <VueDatePicker class="date-picker" v-model="endDate" :teleport="true" :format="format"/>
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

<script>
import {airQualityRecordToObject, useListAirQualityHistory} from '@/traits/airQuality/airQuality.js';
import {downloadCSVRows} from '@/util/downloadCSV.js';
import {ref} from 'vue';

const format = (date) => {
  const day = date.getDate();
  const month = date.getMonth() + 1;
  const year = date.getFullYear();
  const hours = date.getHours();
  const minutes = date.getMinutes();

  return `${day}/${month}/${year} ${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}`;
};

const startDate = ref();
const endDate = ref();

const dateTimeProp = (obj) => {
  return `${obj.toLocaleDateString()} ${obj.toLocaleTimeString()}`;
};

const historyHeaders = [
  {title: 'Record Time', val: (a) => dateTimeProp(a.recordTime)},
  {title: 'C02', val: (a) => a.airQuality.carbonDioxideLevel},
  {title: 'VOC', val: (a) => a.airQuality.volatileOrganicCompounds}
];

export default {
  props: {
    name: {
      type: String,
      required: true
    }
  },
  data() {
    return {
      fetchingHistory: false
    };
  },
  methods: {
    async downloadHistory(n) {
      this.fetchingHistory = true;
      const baseRequest = /** @type {ListAirQualityHistoryRequest.AsObject} */ {
        name: n,
        period: {
          startTime: startDate.value,
          endTime: endDate.value
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
      this.fetchingHistory = false;
    }
  }
};
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
