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
        {{ value?'Running':'Stopped' }}
      </template>
      <template #item.actions="{item}">
        <v-btn outlined v-if="item.active" @click="stopService(item)">Stop</v-btn>
        <v-btn outlined v-else @click="startService(item)">Start</v-btn>
      </template>
    </v-data-table>
  </content-card>
</template>

<script setup>
import ContentCard from '@/components/ContentCard.vue';
import {computed, onMounted, onUnmounted, ref} from 'vue';
import {useAutomationsStore} from '@/routes/automations/store';
import {storeToRefs} from 'pinia';
import {ServiceNames as ServiceTypes} from '@/api/ui/services';
import {usePageStore} from '@/stores/page';

const automationsStore = useAutomationsStore();
const {automationsCollection} = storeToRefs(automationsStore);
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
  {text: '', value: 'actions'}
];

/** @type {Collection} */
const collection = automationsStore.newAutomationsCollection();
collection.needsMorePages = true; // todo: this causes us to load all pages, connect with paging logic instead

onMounted(() => collection.query(ServiceTypes.Automations));
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

/**
 *
 * @param {Service.AsObject} service
 */
function startService(service) {
  console.debug('Starting:', service.id);
  // todo
}

/**
 *
 * @param {Service.AsObject} service
 */
function stopService(service) {
  console.debug('Stopping:', service.id);
  // todo
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
