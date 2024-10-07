<template>
  <v-toolbar color="transparent">
    <v-toolbar-title>
      {{ readonly ? 'Viewing' : 'Editing' }} "{{ serviceID }}"
    </v-toolbar-title>
    <v-fade-transition>
      <v-btn v-if="refreshNeeded" v-bind="refreshAttrs"
             variant="text" class="mr-4"
             v-tooltip="refreshTooltip">
        Refresh
      </v-btn>
    </v-fade-transition>
    <span v-tooltip="saveTooltip">
      <v-btn v-bind="saveAttrs" variant="elevated"/>
    </span>
  </v-toolbar>
  <v-card :loading="loading">
    <v-expand-transition>
      <v-alert v-if="alertVisible" v-bind="alertAttrs" tile/>
    </v-expand-transition>
    <v-card-text>
      <v-data-table
          :headers="tableHeaders"
          :items="tableRows"
          item-key="name"
          :items-per-page="20">
        <template #item.traits="{item}">
          <v-chip
              v-for="trait in item.traits" :key="trait.name"
              v-tooltip="trait.name"
              class="mr-2">
            {{ traitTitle(trait) }}
          </v-chip>
        </template>
        <template v-for="slot in textItemSlots" #[slot]="{item, column, value}" :key="slot">
          <!-- The below looks like a v-text-field, but without the performance problems -->
          <!--<div v-if="!isActiveCell(item, column)"
               class="text-cell mx-n4 v-input v-input--horizontal v-input--density-compact"
               @click="setActiveCell(item, column)">
            <div class="v-input__control">
              <div class="v-field v-field--variant-outlined">
                <div class="v-field__overlay"/>
                <div class="v-field__field">
                  <input class="v-field__input" :value="value" readonly>
                </div>
                <div class="v-field__outline">
                  <div class="v-field__outline__start"/>
                  <div class="v-field__outline__end"/>
                </div>
              </div>
            </div>
          </div>-->
          <v-text-field
              class="text-cell mx-n4"
              :model-value="value" @update:model-value="doUpdateItemText(column.value, item, $event)"
              variant="outlined" density="compact"
              hide-details/>
        </template>
      </v-data-table>
    </v-card-text>
  </v-card>
</template>

<script setup>
import {useServiceConfig} from '@/dynamic/service/service.js';
import {usePopulatedFields} from '@/traits/metadata/metadata.js';
import {get as _get, set as _set, unset as _unset} from 'lodash';
import {computed} from 'vue';

/**
 * @typedef {Metadata.AsObject} Device
 */

const {
  saveAttrs, saveTooltip,
  refreshAttrs, refreshTooltip, refreshNeeded,
  alertVisible, alertAttrs,
  loading, readonly,
  configModel, serviceID
} = useServiceConfig();
const deviceList = computed(() => /** @type {Device[]} */ configModel.value?.devices ?? []);

const populatedMetadataFields = usePopulatedFields(deviceList);

const fieldToHeaderKey = (field) => field.replace(/\./g, '-');
const tableHeaders = computed(() => {
  const dst = /** @type {import('vuetify/lib/components/VDataTable').DataTableHeader[]} */ [
    {title: 'Name', key: 'name'},
    {title: 'Traits', key: 'traits', sortable: false}
  ];

  const partToTitle = (part) => {
    if (part.endsWith('Map')) part = part.slice(0, -3);
    if (part.endsWith('List')) part = part.slice(0, -4);
    return part[0].toUpperCase() + part.slice(1);
  };

  // for nested properties, create nested headers
  const childrenKey = Symbol('header');
  const headersByName = {
    [childrenKey]: dst
  };
  for (const field of populatedMetadataFields.value) {
    const parts = field.split('.');
    let parent = headersByName;
    for (let i = 0; i < parts.length - 1; i++) { // all but the last part
      const part = parts[i];
      if (part.endsWith('Map')) continue; // inline maps
      if (!parent[part]) {
        const children = parent[childrenKey];
        const header = {title: partToTitle(part), children: []};
        children.push(header);
        parent[part] = {[childrenKey]: header.children};
      }
      parent = parent[part];
    }

    // process the last part
    const lastPart = parts[parts.length - 1];
    if (!parent[lastPart]) {
      parent[lastPart] = true;
      parent[childrenKey].push({title: partToTitle(lastPart), key: fieldToHeaderKey(field), value: field});
    }
  }

  return dst;
});
const tableRows = computed(() => {
  return deviceList.value;
});

const textItemSlots = computed(() => {
  return populatedMetadataFields.value.map(field => 'item.' + fieldToHeaderKey(field));
});

const traitTitle = (trait) => {
  const lastDot = trait.name.lastIndexOf('.');
  return lastDot === -1 ? trait.name : trait.name.slice(lastDot + 1);
};

const doUpdateItemText = (key, item, newValue) => {
  const old = _get(item, key);
  if (old === newValue) return;
  if (newValue === '') _unset(item, key);
  else _set(item, key, newValue);
  // configEdited.value = {...configEdited.value};
};

// The below is only needed if we have performance issues with v-text-field.
// There is a bug where the vue dev tools will cause v-text-field to take 15ms to render,
// which causes issues when you have a lot on the page.
//
// const activeItem = ref(null);
// const isActiveCell = (item, column) => {
//   const {name, key} = activeItem.value ?? {};
//   return item.name === name && key === column.key;
// };
// const setActiveCell = (item, column) => {
//   activeItem.value = {name: item.name, key: column.key};
// };
</script>

<style scoped>
.text-cell:not(:hover) :deep(.v-field:not(.v-field--focused) .v-field__outline) {
  opacity: 0;
}

:deep(.v-input), :deep(.v-field), :deep(.v-field__field) {
  font-size: inherit;
}
</style>
