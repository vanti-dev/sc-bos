<template>
  <div>
    <v-divider/>
    <v-row class="ma-0 pa-0 pt-2 mb-n2">
      <v-spacer/>
      <v-col cols="auto">
        <v-pagination
            v-show="itemsPerPage !== -1"
            :length="pageCount"
            show-current-page
            total-visible="5"
            :value="activePage"
            @input="activePage = $event"/>
      </v-col>
      <v-col cols="auto">
        <v-text-field
            v-show="itemsPerPage !== -1"
            v-model="activePage"
            label="Go to page"
            type="number"
            min="1"
            :max="pageCount"
            outlined
            hide-details
            dense
            style="width: 100px"
            @input="activePage = parseInt($event, 10)"/>
      </v-col>
      <v-col cols="auto">
        <v-select
            dense
            outlined
            hide-details
            :value="itemsPerPage"
            label="Devices per page"
            :items="perPageChoices"
            style="width: 150px; cursor: pointer;"
            @change="itemsPerPage = parseInt($event, 10);"/>
      </v-col>
    </v-row>
  </div>
</template>

<script setup>
import {computed} from 'vue';
import {storeToRefs} from 'pinia';

// Store imports
import {useTableDataStore} from '@/stores/tableDataStore';

const props = defineProps({
  tableItemLength: {
    type: Number,
    default: 0
  }
});

// Stores
const tableDataStore = useTableDataStore();

// Store values
const {activePage, itemsPerPage, perPageChoices} = storeToRefs(tableDataStore);

// Computeds
const pageCount = computed(() => Math.ceil(props.tableItemLength / itemsPerPage.value));
</script>
<style lang="scss" scoped>
:deep(.v-pagination .v-pagination__item) {
  background-color: var(--v-neutral-lighten1);
}
</style>
