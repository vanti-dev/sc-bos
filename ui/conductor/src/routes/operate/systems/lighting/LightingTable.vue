<template>
  <v-container fluid class="pa-7">
    <Header />
    <v-sheet class="rounded-lg">
      <main-card>
        <v-data-table
          v-model="selected"
          :headers="headers"
          :items="filteredLights"
          item-key="device_id"
          :search="search"
          @click:row="rowClick"
          :header-props="{ sortIcon: 'mdi-arrow-up-drop-circle-outline' }"
          show-select

        >
          <template v-slot:top>
            <Filters v-if="selected.length <= 1" />
            <BulkAction class="mb-2" v-else />
          </template>
          <template v-slot:item.status="{ item }">
            <span
              :class="getColor(item.status)"
              class="font-weight-bold text-uppercase"
            >
              {{ item.status }}
            </span>
          </template>
        </v-data-table>
      </main-card>
      <RowMenu />
    </v-sheet>
  </v-container>
</template>
<script setup>
import MainCard from "@/components/SectionCard.vue";
import Header from "@/components/operate/Header.vue";
import Filters from "@/components/operate/Filters.vue";
import BulkAction from "@/components/operate/BulkAction.vue";
import RowMenu from "./RowMenu.vue";
import { useLightingStore } from "@/stores/operate/lighting.js";
import { storeToRefs } from "pinia";

const store = useLightingStore();

const { headers, selected, filteredLights, search } = storeToRefs(store);

const getColor = (status) => {
  if (status == "On") {
    return "green--text";
  } else if (status == "Off") {
    return "red--text";
  } else {
    return "orange--text";
  }
};

const rowClick = (item, row) => {
  store.toggleDrawer();
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

/* This selector is the one used by vuetify to match hovered table rows. We are more specific because of the scoped styles */
.table
  > ::v-deep(.v-data-table__wrapper
    > table
    > tbody
    > tr:hover:not(.v-data-table__expanded__content):not(.v-data-table__empty-wrapper)) {
  background-color: #292f36;
}

.v-data-table ::v-deep(.v-data-footer) {
  background: #3F454B !important;
  border-radius: 0px 0px 5px 5px;
  border: none;
  width: 100%;
}
</style>
