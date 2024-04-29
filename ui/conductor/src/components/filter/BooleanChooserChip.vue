<template>
  <v-tooltip bottom>
    <template #activator="{ on }">
      <v-chip
          v-on="on"
          @click="toggle"
          close
          @click:close="clear"
          v-bind="$attrs">
        {{ text }}
      </v-chip>
    </template>
    Toggle {{ title.toLowerCase() }} state
  </v-tooltip>
</template>

<script setup>
const props = defineProps({
  ctx: {
    type: /** @type {import('./pageCtx.js')} */ Object,
    required: true,
    validator(ctx) {
      return ctx.type.value === 'boolean';
    }
  }
});

// eslint-disable-next-line vue/no-setup-props-destructure
const {title, choose, clear, value, text, defaultChoice} = props.ctx;
const nextVal = (val) => {
  if (val === true) return false;
  if (val === null || val === undefined) return true;
  return undefined;
};
const toggle = () => {
  // selecting the default value will clear the chip, make sure we don't do that.
  // The options are [null, true, false], but we don't want to select the default out of those.
  const def = defaultChoice.value?.value;
  const val = value.value;
  let next = nextVal(val);
  if (next === def) {
    next = nextVal(next);
  }
  choose(next);
};
</script>
