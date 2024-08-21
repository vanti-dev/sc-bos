<template>
  <v-menu
      :close-on-content-click="false"
      v-model="active"
      location="bottom left"
      width="340">
    <template #activator="{ props: _props }">
      <v-btn
          v-bind="{..._props, ...$attrs}"
          icon="true"
          class="filter-btn">
        <v-badge dot v-if="badgeShown" :color="badgeColor">
          <v-icon size="24">mdi-filter</v-icon>
        </v-badge>
        <v-icon v-else size="24">mdi-filter</v-icon>
      </v-btn>
    </template>
    <filter-chooser :ctx="ctx"/>
  </v-menu>
</template>
<script setup>
import FilterChooser from '@/components/filter/FilterChooser.vue';
import useFilterCtx, {filterCtxSymbol} from '@/components/filter/filterCtx.js';
import {inject, provide} from 'vue';

defineOptions({inheritAttrs: false});
const props = defineProps({
  ctx: {
    type: Object, // import('./filterCtx.js')
    default: () => {}
  },
  filterOpts: {
    type: Object, // import('./filterCtx.js').Options
    default: () => {}
  }
});

const ctx = inject(filterCtxSymbol,
    () => props.ctx ?? useFilterCtx(() => props.filterOpts),
    true);
provide(filterCtxSymbol, ctx);
const {badgeColor, badgeShown, active} = ctx;
</script>
<style scoped>

</style>
