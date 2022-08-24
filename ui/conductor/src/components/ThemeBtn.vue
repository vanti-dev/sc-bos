<template>
  <v-btn :color="themeColor" v-bind="$attrs" v-on="$listeners">
    <slot/>
  </v-btn>
</template>

<script setup>
import {usePage} from '@/components/page.js';

const {themeColor} = usePage();
</script>
