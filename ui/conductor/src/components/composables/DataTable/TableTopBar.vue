<template>
  <!-- todo: bulk actions -->
  <!-- filters -->
  <v-container fluid style="width: 100%">
    <v-row dense>
      <v-col cols="12" md="5">
        <v-text-field
            :value="search"
            append-icon="mdi-magnify"
            label="Search devices"
            hide-details
            filled
            @input="search = $event"/>
      </v-col>
      <v-spacer/>
      <v-col cols="12" md="2">
        <v-select
            :disabled="props.dropdown.dropdownItems.length <= 1"
            :value="props.dropdown.dropdownValue || 'All'"
            :items="props.dropdown.dropdownItems"
            :label="props.dropdown.dropdownLabel"
            hide-details
            filled
            @change="emits('onDropdownSelect', $event)"/>
      </v-col>
      <!--            <v-col cols="12" md="2">
              <v-select
                  v-model="filterZone"
                  :items="zoneList"
                  label="Zone"
                  hide-details
                  filled/>
            </v-col>-->
    </v-row>
  </v-container>
</template>

<script setup>
import {storeToRefs} from 'pinia';

import {useTableDataStore} from '@/stores/tableDataStore';

// Incoming data
const props = defineProps({
  dropdown: {
    type: Object,
    default: () => {
      return {
        dropdownItems: [],
        dropdownLabel: '',
        dropdownValue: 'All'
      };
    }
  }
});
// Data mutation
const emits = defineEmits(['onDropdownSelect']);

const tableDataStore = useTableDataStore();

const {search} = storeToRefs(tableDataStore);

</script>
