<template>
  <v-menu
      offset-y
      bottom
      :close-on-content-click="false"
      @input="reset"
      min-width="400"
      nudge-bottom="4">
    <template #activator="{ on: onMenu, attrs: bindMenu }">
      <v-tooltip bottom>
        <template #activator="{ on: onTooltip, attrs: bindTooltip }">
          <v-chip
              v-on="{ ...onMenu, ...onTooltip }"
              v-bind="{ ...bindMenu, ...bindTooltip, ...$attrs }"
              close
              @click:close="clear">
            {{ text }}
          </v-chip>
        </template>
        Change {{ title.toLowerCase() }} filter
      </v-tooltip>
    </template>
    <page-chooser :title="title" :type="type" :search.sync="search">
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
