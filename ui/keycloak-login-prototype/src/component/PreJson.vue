<script setup>
import {computed} from 'vue';

const props = defineProps({
  value: {
    type: Object,
    required: true
  }
})

const jsonHtml = computed(() => {
  let jsonStr = JSON.stringify(props.value, null, 2)
  jsonStr = jsonStr.replaceAll(/("https?:\/\/[^"]+")/g, '<a href=$1>$1</a>');
  return jsonStr
})
</script>

<template>
<pre v-html="jsonHtml"/>
</template>
