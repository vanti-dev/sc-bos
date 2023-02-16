<template>
  <content-card>
    <v-data-table
        :headers="headers"
        :items="automationsList"
        item-key="id"
        :search="search"
        :loading="automationsCollection.loading"
        @click:row="showAutomation">
      <template #item.active="{value}">
        <span :class="value?'success--text':'error--text'" class="text--lighten-2">
          {{ value?'Running':'Stopped' }}
        </span>
      </template>
      <template #item.actions="{item}">
        <v-btn outlined v-if="item.active" @click.stop="stopAutomation(item)">Stop</v-btn>
        <v-btn outlined v-else @click.stop="startAutomation(item)">Start</v-btn>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {computed, onMounted, onUnmounted, reactive, ref} from 'vue';
import {useServicesStore} from '@/stores/services';
import {ServiceNames, startService, stopService} from '@/api/ui/services';
import {usePageStore} from '@/stores/page';
import {newActionTracker} from '@/api/resource';

const serviceStore = useServicesStore();
const automationsCollection = ref(serviceStore.getService(ServiceNames.Automations).serviceCollection);
const pageStore = usePageStore();

const props = defineProps({
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

/** @type {Collection} */
const collection = serviceStore.newServicesCollection();
collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead

onMounted(() => collection.query(ServiceNames.Automations));
onUnmounted(() => collection.reset());

const automationsList = computed(() => {
  return Object.values(collection.resources.value).filter(automation => {
    return automation.type === props.type;
  });
});

/**
 *
 * @param {Service.AsObject} service
 * @param {*} row
 */
function showAutomation(service, row) {
  pageStore.showSidebar = true;
  pageStore.sidebarTitle = service.id;
  pageStore.sidebarData = service;
}

const startStopTracker = reactive(newActionTracker());
/**
 *
 * @param {Service.AsObject} service
 */
async function startAutomation(service) {
  console.debug('Starting:', service.id);
  await startService({name: ServiceNames.Automations, id: service.id}, startStopTracker);
}

/**
 *
 * @param {Service.AsObject} service
 */
async function stopAutomation(service) {
  console.debug('Stopping:', service.id);
  await stopService({name: ServiceNames.Automations, id: service.id}, startStopTracker);
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
