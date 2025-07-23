<template>
  <div>
    <h2>Reports</h2>
    <content-card class="px-8 mt-8">
      <v-data-table-server
          :headers="allHeaders"
          v-bind="tableAttrs"
          :items-length="queryTotalCount"
          item-key="id"
          class="pt-4">
        <template #item.created="{ item }">
          {{ timestampToDate(item.createTime).toLocaleString() }}
        </template>
        <template #item.download="{ item }">
          <v-btn @click="downloadReport(item.id)" title="Download" variant="flat">
            <v-icon icon="mdi-download" size="large"/>
          </v-btn>
        </template>
      </v-data-table-server>
    </content-card>
  </div>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb.js';
import {getDownloadReportUrl} from '@/api/ui/reports.js';
import ContentCard from '@/components/ContentCard.vue';
import {useReportsCollection} from '@/composables/reports.js';
import {useDataTableCollection} from '@/composables/table.js';
import {computed, ref} from 'vue';

const props = defineProps({
  source: {
    type: String,
    required: true
  }
});

const listReportsRequest = computed(() => ({
  name: props.source,
}))
const wantCount = ref(20);
const reportOptions = computed(() => ({
  wantCount: wantCount.value
}));
const reportsCollection = useReportsCollection(listReportsRequest, reportOptions);
const tableAttrs = useDataTableCollection(wantCount, reportsCollection);

const queryTotalCount = computed(() => {
  return reportsCollection.totalItems.value;
});

/**
 * Downloads a report by its ID.
 *
 * @param {string} id
 */
function downloadReport(id) {
  const request = { id, name: props.source };
  getDownloadReportUrl(request)
      .then(url => {
        const link = document.createElement('a');
        link.href = url.url;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      })
      .catch(e => {
        console.error('Failed to get download devices URL', e);
      });
}

const allHeaders = computed(() => [
  { title: 'Title', value: 'title', width: '30%' },
  { title: 'Description', value: 'description', width: '30%' },
  { title: 'Create Time', value: 'created', width: '15em', sortable: false },
  { title: 'Download', value: 'download', width: '15em', sortable: false }
]);

</script>

<style scoped>

button:hover .v-icon {
  color: #1976d2;
}
</style>