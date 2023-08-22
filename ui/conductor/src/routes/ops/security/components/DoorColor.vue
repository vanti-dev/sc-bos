<template>
  <div/>
</template>
<script setup>

import {useStatus} from '@/routes/ops/security/components/access-point-card/useStatus';
import {watch} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: ''
  },
  accessAttempt: {
    type: Object,
    default: () => ({})
  },
  statusLog: {
    type: Object,
    default: () => ({})
  }
});

const emit = defineEmits(['update']);
const {color} = useStatus(() => props.accessAttempt, () => props.statusLog);

watch(color, (newColor) => {
  emit('update', {name: props.name, color: newColor});
});
</script>
