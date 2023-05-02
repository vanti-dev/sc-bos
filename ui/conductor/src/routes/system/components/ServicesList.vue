<template>
  <content-card>
    <v-data-table
        :headers="headers"
        :items="serviceList"
        item-key="id"
        :search="search"
        :loading="serviceCollection.loading"
        @click:row="showService">
      <template #item.active="{value}">
        <span :class="value?'success--text':'error--text'" class="text--lighten-2">
          {{ value ? 'Running' : 'Stopped' }}
        </span>
      </template>
      <template #item.actions="{item}">
        <v-btn outlined v-if="item.active" class="mr-1" color="red" @click.stop="_stopService(item)">Stop</v-btn>
        <v-btn outlined v-else color="green" @click.stop="_startService(item)">Start</v-btn>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';
import {useServicesStore} from '@/stores/services';
import {ServiceNames, startService, stopService} from '@/api/ui/services';
import {usePageStore} from '@/stores/page';
import {newActionTracker} from '@/api/resource';
import {useErrorStore} from '@/components/ui-error/error';

const serviceStore = useServicesStore();
const pageStore = usePageStore();
const errors = useErrorStore();

const serviceCollection = reactive(serviceStore.getService(props.name).servicesCollection);
const startStopTracker = reactive(newActionTracker());

const props = defineProps({
  name: {
    type: String,
    default: ServiceNames.Systems
  },
  // optional type filter for services list
  type: {
    type: String,
    default: ''
  }
});

const search = ref('');

const headers = [
  {text: 'ID', value: 'id'},
  {text: 'Status', value: 'active'},
  {text: '', value: 'actions', align: 'end', width: '100'}
];

// todo: this causes us to load all pages, connect with paging logic instead
serviceCollection.needsMorePages = true;

// setup query
watch(() => props.name, (value) => {
  serviceCollection.query(props.name);
}, {immediate: true});
// UI error handling
let unwatchErrors; let unwatchStartStopErrors;
onMounted(() => {
  unwatchErrors = errors.registerCollection(serviceCollection);
  unwatchStartStopErrors = errors.registerTracker(startStopTracker);
});
onUnmounted(() => {
  if (unwatchErrors) unwatchErrors();
  if (unwatchStartStopErrors) unwatchStartStopErrors();
  serviceCollection.reset();
});

const serviceList = computed(() => {
  return Object.values(serviceCollection.resources.value).filter(service => {
    return props.type === '' || props.type === 'all' || service.type === props.type;
  });
});

/**
 *
 * @param {Service.AsObject} service
 * @param {*} row
 */
function showService(service, row) {
  pageStore.showSidebar = true;
  pageStore.sidebarTitle = service.id;
  pageStore.sidebarData = {...service, config: JSON.parse(service.configRaw)};
}


/**
 *
 * @param {Service.AsObject} service
 */
async function _startService(service) {
  console.debug('Starting:', service.id);
  await startService({name: props.name, id: service.id}, startStopTracker);
}

/**
 *
 * @param {Service.AsObject} service
 */
async function _stopService(service) {
  console.debug('Stopping:', service.id);
  await stopService({name: props.name, id: service.id}, startStopTracker);
}

</script>

<style lang="scss" scoped>
:deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.v-data-table :deep(.v-data-footer) {
  background: var(--v-neutral-lighten1) !important;
  border-radius: 0px 0px $border-radius-root*2 $border-radius-root*2;
  border: none;
  margin: 0 -12px -12px;
}

.v-data-table :deep(.item-selected) {
  background-color: var(--v-primary-darken4);
}
</style>
