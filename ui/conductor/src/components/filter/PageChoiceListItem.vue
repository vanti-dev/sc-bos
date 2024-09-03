<template>
  <v-list-item @click="emits('click', props.filter.key)">
    <template #prepend>
      <v-icon>{{ props.filter.icon }}</v-icon>
    </template>
    <template v-if="props.defaultChoice">
      <v-list-item-title class="align-self-auto">{{ props.filter.title }}</v-list-item-title>
    </template>
    <template v-else>
      <v-list-item-subtitle class="text-overline mb-n2">{{ props.filter.title }}</v-list-item-subtitle>
      <v-list-item-title class="align-self-auto">{{ props.choice.text ?? props.choice.value }}</v-list-item-title>
    </template>
    <template #append>
      <v-list-item-action :class="{'mr-n2': !props.defaultChoice}">
        <v-icon v-if="props.defaultChoice" end size="24">mdi-chevron-right</v-icon>
        <v-btn
            v-else
            icon="mdi-close"
            size="small"
            variant="text"
            @click.stop="emits('clear', props.filter.key)"
            class="">
          <v-icon size="20">mdi-close</v-icon>
        </v-btn>
      </v-list-item-action>
    </template>
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
