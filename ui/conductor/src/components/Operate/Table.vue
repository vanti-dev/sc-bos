<template>
  <v-data-table
    v-model="selected"
    :headers="headers"
    :items="lights"
    item-key="device_id"
    show-select
    class="elevation-1"
  >
    <template v-slot:top>
      <Filters />
    </template>
    <template v-slot:item.status="{ item }">
      <p :class="getColor(item.status)" class="font-weight-bold text-uppercase">
        {{ item.status }}
      </p>
    </template>
  </v-data-table>
</template>

<script>
import { useOperateStore } from "@/stores/operate.js";
import { storeToRefs } from "pinia";
import Filters from "@/components/Operate/Filters.vue";

export default {
  setup() {
    const store = useOperateStore();

    const { headers, lights, selected } = storeToRefs(store);

    const getColor = (status) => {
      if (status == "On") {
        return "green--text";
      } else if (status == "Off") {
        return "red--text";
      } else {
        return "orange--text";
      }
    };

    return { store, headers, lights, selected, getColor };
  },
};
</script>

<style lang="scss" scoped>
.v-data-table {
  background-color: #283139;
  color: white;
  width: 100%;
  height: 93vh;
  position: relative;
}
.v-data-table ::v-deep(.v-data-footer) {
  background: rgb(62, 68, 77) !important;
  border-radius: 0px 0px 5px 5px;
  border: none;
  position: absolute;
  bottom: 0;
  width: 100%;
}

/* Using SCSS variables to store breakpoints */
$breakpoint-tablet: 768px;
@media (max-width: $breakpoint-tablet) {
  .v-data-table {
    background-color: #283139;
    color: white;
    width: 100%;
    height: 100%;
  }
  .v-data-table ::v-deep(.v-data-footer) {
    background: rgb(62, 68, 77) !important;
    border-radius: 0px 0px 5px 5px;
    border: none;
    width: 100%;
  }
}
</style>
