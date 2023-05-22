<template>
  <tbody>
    <tr
        v-for="item in props.items"
        v-intersect="onIntersection(item)"
        :key="item[props.itemKey] + triggerRerender"
        :class="tableItemClass(item)"
        @click="showDevice(item)">
      <td v-if="props.showSelect">
        <v-checkbox
            class="ma-0 pa-0 mb-n5"
            color="white"
            :input-value="itemSelection"
            :value="item"
            @change="onSelect(item)"/>
      </td>
      <td>
        {{ item.metadata.appearance ?
          item.metadata.appearance.title :
          item.name }}
      </td>
      <td>{{ item.metadata?.location?.floor ?? '' }}</td>
      <td>{{ item.metadata?.location?.title ?? '' }}</td>

      <!-- Deepest slot for hot points -->
      <td v-if="!siteEditor.zone" class="d-flex justify-end align-center">
        <slot
            name="hotpoint"
            :find-sensor="findSensor"
            :item="item"
            :intersected-item-names="intersectedItemNames"/>
      </td>
    </tr>
  </tbody>
</template>

<script setup>
import {onUnmounted, ref, watch} from 'vue';
import {storeToRefs} from 'pinia';

// Store imports
import {usePageStore} from '@/stores/page';
import {useTableDataStore} from '@/stores/tableDataStore';

const props = defineProps({
  items: {
    type: Array,
    default: () => []
  },
  itemKey: {
    type: String,
    default: 'name'
  },
  showSelect: {
    type: Boolean,
    default: false
  }
});

const emits = defineEmits(['onItemSelect']);

// Stores
const pageStore = usePageStore();
const tableDataStore = useTableDataStore();


// Store values
const {
  findSensor, intersectionHandler, intersectedItemNames, siteEditor
} = tableDataStore;
const {tableSelection, triggerRerender} = storeToRefs(tableDataStore);

const itemSelection = ref([]);

// Methods
/**
 *
 * @param {*} item
 * @return {{
 *   handler: (entries: IntersectionObserverEntry[], observer: IntersectionObserver) => void,
 *   options: {
 *     rootMargin: string,
 *     threshold: number
 *   }
 * }}
 */
function onIntersection(item) {
  return {
    handler: (entries, observer) => intersectionHandler(entries, observer, item.name),
    options: {
      rootMargin: '-50px 0px 0px 0px',
      threshold: 0.75,
      trackVisibility: true,
      delay: 100
    }
  };
}

/**
 * @param {*} item
 * @return {string}
 */
function tableItemClass(item) {
  if (
    pageStore.showSidebar && pageStore.sidebarData?.name === item.name ||
    siteEditor.editMode && itemSelection.value.includes(item)
  ) {
    return 'item-selected';
  }
  return '';
}


/**
 *
 * @param {*} item
 */
function showDevice(item) {
  pageStore.showSidebar = true;
  pageStore.sidebarTitle = item.metadata.appearance ? item.metadata.appearance.title : item.name;
  pageStore.sidebarData = item;
}

/**
 *
 * @param {*} item
 */
function onSelect(item) {
  if (!itemSelection.value.includes(item)) {
    itemSelection.value.push(item);
  } else if (itemSelection.value.includes(item)) {
    itemSelection.value = itemSelection.value.filter(selection => selection !== item);
  }

  emits('onItemSelect', itemSelection.value);
}

// Resetting local selection on leave
onUnmounted(() => {
  itemSelection.value = [];
});

// Setting changes to local selection
watch(tableSelection, items => {
  itemSelection.value = items;
});


</script>

<style lang="scss" scoped>
.v-data-table:not(.selectable) :deep(.v-data-table__selected) {
  background: none;
}

.v-data-table.rowSelectable :deep(.item-selected) {
  background-color: var(--v-primary-darken4);
}

.item-selected {
  background-color: rgba(100, 100, 100, 0.7);
}
</style>
