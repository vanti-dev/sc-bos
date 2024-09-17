<template>
  <json-editor-vue v-model="jsonData" v-bind="defaultEditorAttrs"/>
</template>

<script setup>
import JsonEditorVue from 'json-editor-vue';
import 'vanilla-jsoneditor/themes/jse-theme-dark.css';
import {computed} from 'vue';
import {useTheme} from 'vuetify';

const jsonData = defineModel({type: [Object, String], default: '{}'});

const {name: themeName} = useTheme();
const defaultEditorAttrs = computed(() => {
  return {
    mode: 'text',
    mainMenuBar: false,
    navigationBar: false,
    statusBar: false,
    indentation: 2,
    tabSize: 2,
    class: {'jse-theme-dark': themeName.value === 'dark'}
  };
});
</script>

<style scoped>

</style>
