<template>
  <v-menu
      location="bottom"
      :close-on-content-click="false"
      @update:model-value="reset"
      min-width="400">
    <template #activator="{ props: menuProps }">
      <v-tooltip location="bottom">
        <template #activator="{ props: tooltipProps }">
          <v-chip
              v-bind="{ ...menuProps, ...tooltipProps }"
              closable
              @click:close="clear">
            {{ text }}
          </v-chip>
        </template>
        Change {{ title.toLowerCase() }} filter
      </v-tooltip>
    </template>
    <page-chooser :title="title" :type="type" v-model:search="search">
      <slot :value="value" :items="items" :choose="choose"/>
    </page-chooser>
  </v-menu>
</template>

<script setup>
import PageChooser from '@/components/filter/PageChooser.vue';
import {toRefs} from 'vue';

const props = defineProps({
  ctx: {
    type: /** @type {import('./pageCtx.js')} */ Object,
    required: true
  }
});
const emits = defineEmits(['active']);

const {title, type, text, value, items, choose, clear, search} = toRefs(props.ctx);

const reset = (e) => {
  if (e) {
    search.value = '';
  }
  emits('active', e);
};
</script>
