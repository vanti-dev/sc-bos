<template>
  <tbody>
    <tr
        v-for="item in props.items"
        v-intersect="onIntersection(item)"
        :key="item[props.itemKey] + triggerRerender"
        :class="tableItemClass(item)"
        @click="emits('onClick:row', item)">
      <td v-if="props.showSelect">
        <v-checkbox
            class="ma-0 pa-0 mb-n5"
            color="white"
            :input-value="itemSelection"
            :value="item"
            @change="onSelect(item)"/>
      </td>

      <!-- Static / Text content -->
      <td
          v-for="(header, headerIndex) in [...headerCollection.staticDataHeaders]"
          :key="headerIndex"
          :class="header.value === 'active' ? [item.active ? 'success--text' : 'error--text', 'text--lighten-2'] : ''">
        {{
          collectStaticData(header, item)
        }}
      </td>

      <!-- Live / Dynamic content -->
      <!-- Deepest slot for hot points -->
      <td
          v-for="(slot, slotIndex) in slotsToGenerate"
          :key="slot.slotName + '_' + slotIndex"
          :class="slot.tdClass">
        <slot
            :name="slot.slotName"
            :slot-name="slot.slotName"
            :item="item"
            :values="slot.slotData"/>
      </td>
    </tr>
  </tbody>
</template>

<script setup>
import {computed, onUnmounted, ref, watch} from 'vue';
import {storeToRefs} from 'pinia';

// Store imports
import {usePageStore} from '@/stores/page';
import {useTableDataStore} from '@/stores/tableDataStore';
import {useTableHeaderStore} from './tableHeaderStore';

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
  },
  tableHeaders: {
    type: Array,
    default: () => []
  }
});

const emits = defineEmits(['onClick:row', 'onItemSelect']);

// Stores
const pageStore = usePageStore();
const tableDataStore = useTableDataStore();
const tableHeaderStore = useTableHeaderStore();


// Store values
const {pageType, requiredSlots} = pageStore;
const {
  findSensor, intersectionHandler
} = tableDataStore;
const {tableSelection, triggerRerender} = storeToRefs(tableDataStore);
const {headerCollection} = tableHeaderStore;


// Local data
const itemSelection = ref([]);

// Computeds
const slotsToGenerate = computed(() => {
  let slots = [];

  if (requiredSlots.length) {
    if (pageType.automations || pageType.system) {
      slots = [
        {
          slotName: 'actions'
        }
      ];
    } else if (!pageType.site || !pageType.automations) {
      slots = [
        {
          slotName: 'hotpoints',
          tdClass: 'd-flex justify-end align-center',
          slotData: {findSensor}
        }
      ];
    }
  }

  return slots;
});

// Methods
/**
 *
 * @param {*} header
 * @param {*} item
 * @return {*}
 */
function collectStaticData(header, item) {
  let staticData;

  if (header.text) {
    const headerValue = header.value;
    const nestedObjects = headerValue.split('.');

    let result = item;

    for (let i = 0; i < nestedObjects.length; i++) {
      if (result && result.hasOwnProperty(nestedObjects[i])) {
        staticData = {
          key: nestedObjects[i],
          value: result[nestedObjects[i]]
        };
        result = result[nestedObjects[i]];
      } else {
        staticData = null;
        break;
      }
    }

    if (staticData?.key === 'name' && item.metadata?.appearance?.title) {
      staticData.value = item.metadata.appearance.title;
    }

    if (staticData?.key === 'active') {
      if (item.active) {
        staticData.value = 'Running';
      } else staticData.value = 'Stopped';
    }
  }

  if (staticData) return staticData.value;
}

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
  const matchingName = pageType.automations ?
   pageStore.sidebarData?.id === item.id :
   pageStore.sidebarData?.name === item.name;

  if (
    pageStore.showSidebar && matchingName ||
    pageType.editorMode && itemSelection.value.includes(item)
  ) {
    return 'item-selected';
  }
  return '';
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
  // background-color: rgba(100, 100, 100, 0.7);
  background-color: var(--v-primary-darken4);
}
</style>
