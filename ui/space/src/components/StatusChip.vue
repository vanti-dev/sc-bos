<template>
  <div class="status-chip">
    <span class="indicator d-flex align-center justify-center" :class="backgroundColorClasses" :style="backgroundColorStyles">
      <v-icon :icon="indicatorIcon" size="1em"/>
    </span>
    <span class="label"><slot>{{ props.label }}</slot></span>
  </div>
</template>

<script setup>
import {computed} from 'vue';
import {useBackgroundColor} from 'vuetify/lib/composables/color';

const props = defineProps({
  label: {
    type: String,
    default: ''
  },
  status: {
    type: String,
    default: 'success' // one of 'success', 'warning', 'error', 'info'
  }
});


const indicatorColor = computed(() => props.status)
const indicatorIcon = computed(() => {
  switch (props.status) {
    case 'success':
      return 'mdi-check';
    case 'warning':
      return 'mdi-exclamation';
    case 'error':
      return 'mdi-close';
    case 'info':
      return 'mdi-information-symbol';
    default:
      return '';
  }
})

const {backgroundColorClasses, backgroundColorStyles} = useBackgroundColor(indicatorColor)
</script>

<style scoped>
.status-chip {
  display: flex;
  align-items: center;
  column-gap: .5em;
}

.indicator {
  display: inline-block;
  width: 1.5em;
  height: 1.5em;
  border-radius: 100%;
  border: .15em solid currentColor;
}

.label {
  line-height: 1;
}
</style>