<template>
  <v-data-table
      :class="tableClasses"
      fixed-header
      :headers="props.tableHeaders"
      hide-default-footer
      :hide-default-header="!props.tableItems.length"
      :items="props.tableItems"
      :items-per-page="itemsPerPage"
      :page.sync="activePage"
      :show-select="siteEditor.editMode"
      :value="tableSelection"
      @toggle-select-all="onToggleSelectAll($event)">
    <!-- Search and Filter bar -->
    <template #top>
      <TopBar
          :dropdown="props.dropdown"
          @onDropdownSelect="emits('update:dropdownValue', $event)"/>
    </template>

    <!-- Column labels -->
    <template #header="{headers}">
      <TableHeader :table-headers="headers"/>
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
          :show-select="siteEditor.editMode"
          @onItemSelect="emits('update:selectedItems', $event)">
        <!-- Middle - bridge slot -->
        <template #hotpoint="{item, intersectedItemNames}">
          <slot
              name="hotpoint"
              :find-sensor="findSensor"
              :item="item"
              :intersected-item-names="intersectedItemNames"/>
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
import TableHeader from '@/components/composables/DataTable/TableHeader.vue';
import TableBody from '@/components/composables/DataTable/TableBody.vue';
import TableFooter from '@/components/composables/DataTable/TableFooter.vue';

// Store imports
import {useTableDataStore} from '@/stores/tableDataStore';

// Stores
const tableDataStore = useTableDataStore();
const {siteEditor, findSensor} = tableDataStore;

const {activePage, itemsPerPage, tableSelection, search} = storeToRefs(tableDataStore);

const props = defineProps({
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
  selectedItems: {
    type: Array,
    default: () => []
  },
  rowSelect: {
    type: Boolean,
    default: true
  }
});

const emits = defineEmits(['update:dropdownValue', 'update:selectedItems']);

// Computeds
const tableClasses = computed(() => {
  const c = [];
  if (siteEditor.editMode) c.push('selectable');
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
