<template>
  <side-bar>
    <v-tabs grow v-model="tab">
      <v-tab value="securityevents">
        <v-icon>mdi-bell</v-icon>
      </v-tab>
      <v-tab value="control">
        <v-icon>mdi-devices</v-icon>
      </v-tab>
    </v-tabs>
  </side-bar>
</template>

<script setup>
import SideBar from '@/components/SideBar.vue';
import {ref} from 'vue';

const tab = ref(null);

</script>

<style scoped>
:deep(.v-slide-group__prev),
:deep(.v-slide-group__next) {
  display: none !important;
}

</style>
