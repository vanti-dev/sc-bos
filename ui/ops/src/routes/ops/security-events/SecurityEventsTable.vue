<template>
  <div class="ml-2">
    <v-row v-if="!props.overviewPage" class="mt-0 ml-0 pl-0">
      <h3 class="text-h3 pt-2 pb-6">Security Events</h3>
      <v-spacer/>
    </v-row>

    <content-card :class="['px-8', {'mt-8 px-4': !props.overviewPage}]">
      <v-data-table-server
          :headers="allHeaders"
          v-bind="tableAttrs"
          disable-sort
          :items-length="queryTotalCount"
          :row-props="rowProps"
          class="pt-4">
        <template #top>
          <v-row
              :class="[
                'd-flex flex-row align-center mb-2 mt-1 ml-0 pl-0 mr-1'
              ]">
            <h3 v-if="props.overviewPage" class="text-h4">
              Security Events
            </h3>
            <v-spacer/>
          </v-row>
        </template>
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
import {useSidebarStore} from '@/stores/sidebar';
import {computed, onUnmounted, ref} from 'vue';

const props = defineProps({
  overviewPage: {
    type: Boolean,
    default: false
  }
});

const sidebar = useSidebarStore();

const cohort = useCohortStore();
const name = computed(() => cohort.hubNode?.name ?? '');

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
  {title: 'Source', value: 'source.name', width: '30%'},
  {title: 'Priority', value: 'priority', width: '10em'}
];

const rowProps = ({item}) => {
  if (item.resolveTime) return {class: 'resolved'};
  return {};
};
onUnmounted(() => {
  sidebar.closeSidebar();
});
</script>

<style lang="scss" scoped>
:deep(table) {
  table-layout: fixed;
}

:deep(.resolved) {
  color: #fff5 !important;
}

.v-data-table :deep(tr:hover) {
  cursor: pointer;
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

.character-counter {
  font-size: 75%;
  color: rgb(var(--v-theme-primary)); /* Adjust color as needed */
}
</style>
