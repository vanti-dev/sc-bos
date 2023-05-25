<template>
  <v-data-table
      :class="tableClasses"
      fixed-header
      hide-default-footer
      :headers="props.tableHeaders"
      :loading="props.tableLoading"
      :items="props.tableItems"
      :items-per-page="itemsPerPage"
      :page.sync="activePage"
      :show-select="pageType.editorMode"
      :value="tableSelection"
      @toggle-select-all="onToggleSelectAll($event)">
    <!-- Search and Filter bar -->
    <template #top>
      <TopBar
          v-if="!pageType.automations"
          :dropdown="props.dropdown"
          @onDropdownSelect="emits('update:dropdownValue', $event)"/>
    </template>

    <!-- Table data -->
    <template #body="{items}">
      <tr v-if="!props.tableItems.length && !search.length">
        <td class="text-center py-4">No data to be displayed.</td>
      </tr>
      <TableBody
          v-else
          :items="items"
          :item-key="props.tableItemKey"
          :show-select="pageType.editorMode"
          @onClick:row="emits('onClick:row', $event)"
          @onItemSelect="emits('update:selectedItems', $event)">
        <!-- Middle - bridge slot -->
        <!-- Item row end with live data / possible user actions-->
        <template
            v-for="requiredSlot in props.requiredSlots"
            #[requiredSlot]="{slotName, item, values}">
          <slot
              :name="slotName"
              :item="item"
              :values="values"/>
        </template>
      </TableBody>
    </template>

    <!-- Pagination and items per page -->
    <template #footer>
      <TableFooter :table-item-length="props.tableItems.length"/>
    </template>
  </v-data-table>
</template>

<script setup>
import {computed, watch} from 'vue';
import {storeToRefs} from 'pinia';

// Component import
import TopBar from '@/components/composables/DataTable/TableTopBar.vue';
import TableBody from '@/components/composables/DataTable/TableBody.vue';
import TableFooter from '@/components/composables/DataTable/TableFooter.vue';

// Store imports
import {useTableDataStore} from '@/stores/tableDataStore';
import {usePageStore} from '@/stores/page';

// Stores
const {pageType} = usePageStore();
const tableDataStore = useTableDataStore();

const {activePage, itemsPerPage, tableSelection, search} = storeToRefs(tableDataStore);

const props = defineProps({
  colSpan: {
    type: String,
    default: ''
  },
  dropdown: {
    type: Object,
    default: () => {}
  },
  tableHeaders: {
    type: Array,
    default: () => []
  },
  tableItems: {
    type: Array,
    default: () => []
  },
  tableItemKey: {
    type: String,
    default: 'name'
  },
  tableItemsPerPage: {
    type: Number,
    default: 10
  },
  tableLoading: {
    type: Boolean,
    default: false
  },
  selectedItems: {
    type: Array,
    default: () => []
  },
  rowSelect: {
    type: Boolean,
    default: true
  },
  requiredSlots: {
    type: Array,
    default: () => []
  }
});

const emits = defineEmits(['onClick:row', 'update:dropdownValue', 'update:selectedItems']);

// Computeds
const tableClasses = computed(() => {
  const c = [];
  if (pageType.editorMode) c.push('selectable');
  if (props.rowSelect) c.push('rowSelectable');
  return c.join(' ');
});


// Methods

/**
 *
 * @param {*} event
 */
function onToggleSelectAll(event) {
  // If select all
  if (event.value) {
    tableSelection.value = [...props.selectedItems, ...event.items].reduce((acc, item) => {
      if (!acc.some((existingItem) => JSON.stringify(existingItem) === JSON.stringify(item))) {
        acc.push(item);
      }
      return acc;
    }, []);

    // If remove all
  } else if (!event.value) {
    tableSelection.value = tableSelection.value.filter((selectedItem) => !event.items.includes(selectedItem));
  }

  // update item selection
  emits('update:selectedItems', tableSelection.value);
}

// Watching selected item prop and updating pinia ref
watch(() => props.selectedItems, items => {
  tableSelection.value = items;
});
</script>

<style lang="scss" scoped>
:deep(.v-data-table-header__icon) {
  margin-left: 8px;
}
</style>
