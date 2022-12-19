<template>
  <v-card
    class="d-flex white--text rounded-0"
    max-width="100%"
    elevation-0
    flat
    :to="to"
    style="background: #ffffff1a"
  >
    <v-avatar tile class="my-4 mx-4" size="57">
      <img :src="icon" class="" />
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
  icon: String,
  to: [String, Object]
})
</script>

