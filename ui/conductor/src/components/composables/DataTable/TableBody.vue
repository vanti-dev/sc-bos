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

      <td
          v-for="(header, headerIndex) in props.tableHeaders"
          :key="header.slotName ? header.slotName + '_' + headerIndex : headerIndex"
          :class="[
            header.value === 'active' ?
              [item.active ? 'success--text' : 'error--text', 'text--lighten-2'] :
              ''
          ]">
        <!-- Live / Dynamic content -->
        <template v-if="!header.text && !pageType.site">
          <span
              v-for="dynamicSlot in slotsToGenerate()"
              :key="dynamicSlot.slotName"
              :class="dynamicSlot.tdClass">
            <slot
                :name="dynamicSlot.slotName"
                :slot-name="dynamicSlot.slotName"
                :item="item"
                :values="dynamicSlot.slotData"/>
          </span>
        </template>

        <!-- Static / Text content -->
        <template v-else>
          {{ collectStaticData(header, item) }}
        </template>
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

// Store values
const {pageType, requiredSlots} = pageStore;
const {findSensor, intersectionHandler} = tableDataStore;
const {tableSelection, triggerRerender} = storeToRefs(tableDataStore);


// Local data
const itemSelection = ref([]);

// Methods
/**
 * Within this function we are going to collect
 * the relevant data and styles (css classes) for the props being generated
 *
 * Keep in mind, this function MUST BE UPDATED if we require/add any new action/hotpoint
 *
 *  @return {Array.<{ slotName: string, tdClass: string, slotData: Object }>}
 */
function slotsToGenerate() {
  let slots;

  if (requiredSlots.length) {
    requiredSlots.forEach(slot => {
      //
      if (slot === 'actions' && (pageType.automations || pageType.system)) {
        slots = [
          {
            slotName: 'actions',
            tdClass: 'd-flex justify-end align-center mx-auto'
          }
        ];
        //
      } else if (slot === 'hotpoints' && pageType.devices) {
        slots = [
          {
            slotName: 'hotpoints',
            tdClass: 'd-flex justify-end align-center mx-auto',
            slotData: {findSensor}
          }
        ];
      }
    });
  }

  return slots;
};

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
    handler: (entries, observer) => intersectionHandler(entries, observer, item[props.itemKey]),
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
  const matchingName = pageType.automations || pageType.system ?
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
:deep(tr) {
  cursor: pointer;
}
.v-data-table:not(.selectable) :deep(.v-data-table__selected) {
  background: none;
}

:deep(tr:hover) {
  .automation-device__btn {
    &--red {
      background-color: red;
      .v-btn__content {
        color: white;
      }
    }
    &--green {
      background-color: green;
      .v-btn__content {
        color: white;
      }
    }
  }
}

:deep(.item-selected) {
  background-color: var(--v-primary-darken4);
  .automation-device__btn--red {
      background-color: red;
      .v-btn__content {
        color: white;
      }
    }
    .automation-device__btn--green {
      background-color: green;
      .v-btn__content {
        color: white;
      }
    }
  }
</style>
