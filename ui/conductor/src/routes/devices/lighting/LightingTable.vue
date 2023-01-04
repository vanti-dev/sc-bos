<template>
  <main-card>
    <v-data-table
        v-model="selected"
        :headers="headers"
        :items="filteredLights"
        item-key="device_id"
        :search="search"
        @click:row="rowClick"
        :header-props="{ sortIcon: 'mdi-arrow-up-drop-circle-outline' }"
        show-select>
      <template #top>
        <Filters v-if="selected.length <= 1"/>
        <BulkAction v-else/>
      </template>
      <template #item.status="{ item }">
        <span
            :class="getColor(item.status)"
            class="font-weight-bold text-uppercase">
          {{ item.status }}
        </span>
      </template>
    </v-data-table>
  </main-card>
</template>
<script setup>
import MainCard from '@/components/ContentCard.vue';
import Filters from '@/routes/devices/lighting/components/Filters.vue';
import BulkAction from '@/routes/devices/lighting/components/BulkAction.vue';
import {useLightingStore} from '@/stores/devices/lighting.js';
import {storeToRefs} from 'pinia';
import {usePageStore} from '@/stores/page';

const store = useLightingStore();
const {headers, selected, filteredLights, search} = storeToRefs(store);

const pageStore = usePageStore();

const getColor = (status) => {
  if (status == 'On') {
    return 'green--text';
  } else if (status == 'Off') {
    return 'red--text';
  } else {
    return 'orange--text';
  }
};

const rowClick = (item, row) => {
  pageStore.showSidebar = true;
  store.setSelectedItem(item);
};
</script>

<style lang="scss" scoped>
::v-deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.table ::v-deep(tbody tr) {
  cursor: pointer;
}

.v-data-table ::v-deep(.v-data-footer) {
  background: var(--v-neutral-lighten1) !important;
  border-radius: 0px 0px $border-radius-root*2 $border-radius-root*2;
  border: none;
  margin: 0 -12px -12px;
}
</style>
