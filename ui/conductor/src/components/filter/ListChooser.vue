<template>
  <v-list>
    <v-list-item v-for="item in _items" :key="item.value" @click="onClick(item)">
      <v-list-item-title>{{ item.title }}</v-list-item-title>
      <template #append>
        <v-fade-transition hide-on-leave>
          <v-list-item-action v-if="isSelectedItem(item)">
            <v-icon>mdi-check</v-icon>
          </v-list-item-action>
        </v-fade-transition>
      </template>
    </v-list-item>
  </v-list>
</template>
<script setup>
import {computed} from 'vue';

const props = defineProps({
  items: {
    type: null, // of string | {title: string, value: any}
    default: () => []
  },
  value: {
    type: null, // same type as items
    default: null
  }
});
const emit = defineEmits(['input']);

const _items = computed(() => props.items.map(item => {
  if (typeof item === 'object') return item;
  return {title: `${item}`, value: item, _src: item};
}));
const isSelectedItem = (item) => {
  if (item === null || item === undefined) return props.value === null || props.value === undefined;
  if (item.hasOwnProperty('_src')) return props.value === item._src;
  if (item.hasOwnProperty('value') && props.value?.hasOwnProperty('value')) return props.value.value === item.value;
  return props.value === item;
};
const onClick = (item) => {
  if (item === null || item === undefined) emit('input', null);
  else if (item.hasOwnProperty('_src')) emit('input', item._src);
  else emit('input', item);
};
</script>
<style scoped>

</style>
