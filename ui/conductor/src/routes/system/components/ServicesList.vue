<template>
  <content-card>
    <v-row class="pa-4" v-if="configStore.config?.hub">
      <v-combobox
          v-model="node"
          :items="Object.values(hubStore.nodesList)"
          label="System Component"
          item-text="name"
          item-value="name"
          hide-details="auto"
          :loading="hubStore.nodesListCollection.loading ?? true"
          outlined/>
      <v-spacer/>
    </v-row>
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
import {newActionTracker} from '@/api/resource';
import {ServiceNames, startService, stopService} from '@/api/ui/services';
import ContentCard from '@/components/ContentCard.vue';
import {useErrorStore} from '@/components/ui-error/error';
import {useAppConfigStore} from '@/stores/app-config';
import {useHubStore} from '@/stores/hub';
import {usePageStore} from '@/stores/page';
import {useServicesStore} from '@/stores/services';
import {serviceName} from '@/util/proxy';
import {computed, onMounted, onUnmounted, reactive, ref, watch} from 'vue';

const serviceStore = useServicesStore();
const pageStore = usePageStore();
const errors = useErrorStore();
const configStore = useAppConfigStore();
const hubStore = useHubStore();

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

const node = computed({
  get() {
    return pageStore.sidebarNode;
  },
  set(val) {
    pageStore.sidebarNode = val;
  }
});
const search = ref('');

const headers = [
  {text: 'ID', value: 'id'},
  {text: 'Status', value: 'active'},
  {text: '', value: 'actions', align: 'end', width: '100'}
];

const serviceCollection = ref({});

// query watchers
watch(() => props.name, async () => {
  if (serviceCollection.value.reset) serviceCollection.value.reset();
  serviceCollection.value =
      serviceStore.getService(props.name, await node.value.commsAddress, await node.value.commsName).servicesCollection;
  // reinitialise in case this service collection has been previously reset;
  serviceCollection.value.init();
  serviceCollection.value.query(props.name);
}, {immediate: true});
watch(node, async () => {
  if (serviceCollection.value.reset) serviceCollection.value.reset();
  serviceCollection.value =
      serviceStore.getService(props.name, await node.value.commsAddress, await node.value.commsName).servicesCollection;
  serviceCollection.value.init();
  serviceCollection.value.query(props.name);
}, {immediate: true});


watch(serviceCollection, () => {
  // todo: this causes us to load all pages, connect with paging logic instead
  serviceCollection.value.needsMorePages = true;
});

// UI error handling
let unwatchErrors; let unwatchStartStopErrors;
onMounted(() => {
  unwatchErrors = errors.registerCollection(serviceCollection);
  unwatchStartStopErrors = errors.registerTracker(startStopTracker);
});
onUnmounted(() => {
  if (unwatchErrors) unwatchErrors();
  if (unwatchStartStopErrors) unwatchStartStopErrors();
  serviceCollection.value.reset();
});

const serviceList = computed(() => {
  return Object.values(serviceCollection.value?.resources?.value ?? []).filter(service => {
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
  console.debug('Starting:', serviceName(node.value.name, props.name), service.id);
  await startService({name: serviceName(await node.value.commsName, props.name), id: service.id}, startStopTracker);
}

/**
 *
 * @param {Service.AsObject} service
 */
async function _stopService(service) {
  console.debug('Stopping:', serviceName(node.value.name, props.name), service.id);
  await stopService({name: serviceName(await node.value.commsName, props.name), id: service.id}, startStopTracker);
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
