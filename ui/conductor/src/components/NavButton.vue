<template>
  <v-card class="d-flex flex-nowrap justify-space-between" max-width="300" :to="to">
    <v-avatar tile class="my-6 ml-4" size="57">
      <sc-logo :fill="color"/>
    </v-avatar>
    <div>
      <slot>
        <v-card-title class="font-weight-bold">
          <slot name="title">{{ title }}</slot>
        </v-card-title>
        <v-card-subtitle>
          <slot name="subtitle">{{ subtitle }}</slot>
        </v-card-subtitle>
      </slot>
    </div>
  </v-card>
</template>

<script setup>
import ScLogo from '@/components/ScLogo.vue';

const props = defineProps({
  title: String,
  subtitle: String,
  color: String,
  to: [String, Object]
})
</script>

