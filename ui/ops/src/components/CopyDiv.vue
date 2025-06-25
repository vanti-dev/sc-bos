<template>
  <div class="d-flex align-center">
    <span ref="container" class="mr-auto">
      <slot>{{ props.text ?? '' }}</slot>
    </span>
    <v-btn :icon="copied ? 'mdi-check' : 'mdi-content-copy'" variant="text" @click="onCopy" v-bind="props.btnProps">
      <v-icon :size="iconSize"/>
      <v-menu activator="parent"
              :open-on-click="false"
              :model-value="copied"
              location="bottom"
              offset="6">
        <v-card :text="copiedText" color="success"/>
      </v-menu>
    </v-btn>
  </div>
</template>

<script setup>
import {onScopeDispose, ref, watch} from 'vue';

const props = defineProps({
  text: {
    type: String,
    default: undefined
  },
  iconSize: {
    type: [Number, String],
    default: 24
  },
  btnProps: {
    type: Object,
    default: () => ({})
  },
  copiedText: {
    type: String,
    default: 'Copied to clipboard'
  }
});

const container = ref(null);

const copied = ref(false);
const onCopy = () => {
  let _text = props.text;
  if (!_text) {
    if (!container.value) return;
    _text = container.value.innerText;
  }
  navigator.clipboard.writeText(_text);
  copied.value = true;
}

let copiedTimeout = 0;
watch(copied, (value) => {
  clearTimeout(copiedTimeout);
  if (!value) return;
  copiedTimeout = setTimeout(() => {
    copied.value = false;
  }, 5000);
});
onScopeDispose(() => clearTimeout(copiedTimeout));
</script>

<style scoped>

</style>