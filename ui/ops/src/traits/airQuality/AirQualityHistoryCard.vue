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
                @click="onDownloadClick()"
                icon="mdi-file-download" elevation="0" class="ml-2 mr-n2 mt-1"
                :loading="fetchingHistory" :disabled="!!downloadBtnDisabled"/>
          </div>
        </div>
      </v-card-text>
    </v-card>
  </v-menu>
</template>

<script setup>
import {triggerDownload} from '@/components/download/download.js';
import {addDays, startOfDay} from 'date-fns';
import {computed, ref} from 'vue';
import {VDateInput} from 'vuetify/labs/components';

const p = defineProps({
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

const onDownloadClick = async () => {
  if (!p.name || p.name === '') {
    downloadError.value = 'No device name provided';
    return;
  }
  const names = [p.name];
  await triggerDownload(
      'air-quality',
      {conditionsList: [{stringIn: {stringsList: names}}]},
      {startTime: startOfDay(startDate.value), endTime: startOfDay(addDays(endDate.value, 1))},
      {

      }
  )
}
</script>
