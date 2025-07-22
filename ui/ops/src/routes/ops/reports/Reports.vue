<template>
  <div>
    <h2>Reports</h2>
    <content-card class="px-8 mt-8">
      <v-data-table
          :headers="allHeaders"
          :items="reports"
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
      </v-data-table>
    </content-card>
  </div>
</template>

<script setup>
import {timestampToDate} from '@/api/convpb.js';
import {getDownloadReportUrl, listReports} from '@/api/ui/reports.js';
import ContentCard from '@/components/ContentCard.vue';
import {computed, onMounted, ref, watch} from 'vue';

const reports = ref([]);

const props = defineProps({
  source: {
    type: String,
    required: true
  }
});

const source = ref(props.source);
const name = computed(() => source.value);

const fetchReports = async () => {
  reports.value = [];
  const req = { name: name.value }
  try {
    const response = await listReports(req);
    reports.value = response.reportsList.map(report => ({
      id: report.id,
      title: report.title,
      description: report.description,
      createTime: report.createTime
    }));
  } catch (error) {
    console.error('Failed to fetch reports:', error);
  }
};

onMounted(() => {
  fetchReports();
});

// Watch for changes to source and fetch reports
watch(source, () => {
  fetchReports();
});

/**
 * Downloads a report by its ID.
 *
 * @param {string} id
 */
function downloadReport(id) {
  const request = { id, name: name.value };
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