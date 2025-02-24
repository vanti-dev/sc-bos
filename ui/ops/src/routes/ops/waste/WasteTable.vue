<template>
  <div class="ml-2">
    <v-row class="mt-0 ml-0 pl-0">
      <h3 class="text-h3 pt-2 pb-6">Waste Records</h3>
      <v-spacer/>
    </v-row>

    <content-card class="px-8 mt-8">
      <v-data-table-server
          :headers="allHeaders"
          v-bind="tableAttrs"
          disable-sort
          :items-length="queryTotalCount"
          class="pt-4">
        <template #item.createTime="{ item }">
          {{ timestampToDate(item.wasteCreateTime).toLocaleString() }}
        </template>
        <template #item.weight="{ item }">
          {{ item.weight.toFixed(2) }} {{ uiConfig.config.wasteRecordUnit ?? "kg" }}
        </template>
        <template #item.disposalMethod="{ item }">
          {{ getDisposalMethod(item.disposalMethod) }}
        </template>
      </v-data-table-server>
    </content-card>
  </div>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb.js';
import ContentCard from '@/components/ContentCard.vue';
import {useDataTableCollection} from '@/composables/table.js';
import {useWasteRecordsCollection} from '@/composables/waste.js';
import {useCohortStore} from '@/stores/cohort.js';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {computed, ref} from 'vue';

const uiConfig = useUiConfigStore();
const cohort = useCohortStore();
const name = computed(() => uiConfig.config.wasteRecordsSource ?? cohort.hubNode?.name ?? '');

const wasteRecordsRequest = computed(() => ({
  name: name.value
}));
const wantCount = ref(20);
const wasteRecordsOptions = computed(() => ({
  wantCount: wantCount.value
}));
const wasteRecordsCollection = useWasteRecordsCollection(wasteRecordsRequest, wasteRecordsOptions);
const tableAttrs = useDataTableCollection(wantCount, wasteRecordsCollection);

// Calculate the total number of items in the query
const queryTotalCount = computed(() => {
  return wasteRecordsCollection.totalItems.value;
});

const allHeaders = [
  {title: 'Waste Created', value: 'createTime', width: '15em'},
  {title: 'Area', value: 'area', width: '40%'},
  {title: 'Disposal Method', value: 'disposalMethod', width: '30%'},
  {title: 'Weight', value: 'weight', width: '10em', align: 'end'},
  {title: 'Stream', value: 'stream', width: '30%'},
  {title: 'System', value: 'system', width: '30%'}
];

const getDisposalMethod = (disposalMethod) => {
  switch (disposalMethod) {
    case 1:
      return 'General Waste';
    case 2:
      return 'Recycling';
    default:
      return 'Unknown';
  }
}

</script>

<style scoped>
:deep(table) {
  table-layout: fixed;
}

.hide-pagination {
  :deep(.v-data-table-footer__info),
  :deep(.v-pagination__last) {
    display: none;
  }

  :deep(.v-pagination__first) {
    margin-left: 16px;
  }
}
</style>