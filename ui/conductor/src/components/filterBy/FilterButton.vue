<template>
  <div>
    <v-tooltip left>
      <template #activator="{ on }">
        <v-btn
            v-on="on"
            class="rounded"
            color="primary"
            elevation="0"
            fab
            height="36"
            tile
            small
            @click="showHideFilterList = !showHideFilterList">
          <v-icon>mdi-filter</v-icon>
        </v-btn>
      </template>
      <span>Set Filter</span>
    </v-tooltip>
  </div>
</template>

<script setup>
import {computed} from 'vue';

const emits = defineEmits(['update:showFilterList']);
const props = defineProps({
  showFilterList: {
    type: Boolean,
    default: false
  }
});

const showHideFilterList = computed({
  get: () => props.showFilterList,
  set: (value) => emits('update:showFilterList', value)
});
</script>

<style scoped>

</style>
