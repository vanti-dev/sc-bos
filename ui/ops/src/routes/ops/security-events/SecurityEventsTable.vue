<template>
  <div class="ml-2">
    <v-row class="mt-0 ml-0 pl-0">
      <h3 class="text-h3 pt-2 pb-6">Security Events</h3>
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
          {{ timestampToDate(item.securityEventTime).toLocaleString() }}
        </template>
      </v-data-table-server>
    </content-card>
  </div>
</template>
<script setup>
import {timestampToDate} from '@/api/convpb';
import ContentCard from '@/components/ContentCard.vue';
import {useSecurityEventsCollection} from '@/composables/securityevents.js';
import {useDataTableCollection} from '@/composables/table.js';
import {useCohortStore} from '@/stores/cohort.js';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {computed, ref} from 'vue';

const uiConfig = useUiConfigStore();
const cohort = useCohortStore();
const name = computed(() => uiConfig.config.securityEventsSource ?? cohort.hubNode?.name ?? '');

const securityEventsRequest = computed(() => ({
  name: name.value
}));
const wantCount = ref(20);
const securityEventsOptions = computed(() => ({
  wantCount: wantCount.value
}));
const securityEventsCollection = useSecurityEventsCollection(securityEventsRequest, securityEventsOptions);
const tableAttrs = useDataTableCollection(wantCount, securityEventsCollection);

// Calculate the total number of items in the query
const queryTotalCount = computed(() => {
  return securityEventsCollection.totalItems.value;
});

const allHeaders = [
  {title: 'Timestamp', value: 'createTime', width: '15em'},
  {title: 'Description', value: 'description', width: '60%'},
  {title: 'Priority', value: 'priority', width: '10em', align: 'end'},
  {title: 'Source', value: 'source.name', width: '30%'}
];

</script>

<style lang="scss" scoped>
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
