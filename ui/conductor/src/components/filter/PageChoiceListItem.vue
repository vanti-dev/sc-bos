<template>
  <v-list-item @click="emits('click', props.filter.key)">
    <v-list-item-icon>
      <v-icon>{{ props.filter.icon }}</v-icon>
    </v-list-item-icon>
    <v-list-item-content class="flex-column align-start">
      <template v-if="props.defaultChoice">
        <v-list-item-title class="align-self-auto">{{ props.filter.title }}</v-list-item-title>
      </template>
      <template v-else>
        <v-list-item-subtitle class="text-overline">{{ props.filter.title }}</v-list-item-subtitle>
        <v-list-item-title class="align-self-auto">{{ props.choice.text ?? props.choice.value }}</v-list-item-title>
      </template>
    </v-list-item-content>
    <v-list-item-action>
      <v-icon v-if="props.defaultChoice">mdi-chevron-right</v-icon>
      <v-btn v-else icon @click.stop="emits('clear', props.filter.key)" class="mr-n1">
        <v-icon size="20">mdi-close</v-icon>
      </v-btn>
    </v-list-item-action>
  </v-list-item>
</template>

<script setup>

const props = defineProps({
  filter: {
    type: Object,
    required: true
  },
  choice: {
    type: Object,
    required: true
  },
  defaultChoice: {
    type: Boolean,
    default: false
  }
});
const emits = defineEmits(['click', 'clear']);
</script>

<style scoped>
</style>
