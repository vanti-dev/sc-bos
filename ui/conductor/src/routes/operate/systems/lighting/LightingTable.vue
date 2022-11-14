<template>
  <v-container fluid class="pa-7">
    <BulkAction />
    <main-card>
      <v-data-table
        v-model="selected"
        :headers="headers"
        :items="filteredLights"
        item-key="device_id"
        show-select
        class="table"
      >
        <template v-slot:top>
          <Filters />
        </template>
        <template v-slot:item.status="{ item }">
          <p
            :class="getColor(item.status)"
            class="font-weight-bold text-uppercase"
          >
            {{ item.status }}
          </p>
        </template>
      </v-data-table>
    </main-card>
  </v-container>
</template>
<script setup>
import MainCard from "@/components/SectionCard.vue";
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

const store = useLightingStore();

const { headers, selected, filteredLights } = storeToRefs(store);

const getColor = (status) => {
  if (status == "On") {
    return "green--text";
  } else if (status == "Off") {
    return "red--text";
  } else {
    return "orange--text";
  }
};
</script>

<style scoped>
.table {
  background-color: transparent;
}

::v-deep(.v-data-table-header__icon) {
  margin-left: 8px;
}

.table ::v-deep(tbody tr) {
  cursor: pointer;
}

/* This selector is the one used by vuetify to match hovered table rows. We are more specific because of the scoped styles */
.table
  > ::v-deep(.v-data-table__wrapper
    > table
    > tbody
    > tr:hover:not(.v-data-table__expanded__content):not(.v-data-table__empty-wrapper)) {
  background-color: #ffffff1a;
}
</style>
