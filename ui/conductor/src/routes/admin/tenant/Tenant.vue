<template>
  <v-container>
    <section-card class="mt-7 pb-7">
      <v-card-title>{{tenantTitle}}</v-card-title>
      <section-card class="mx-7">
        <v-card-title>Zones</v-card-title>
        <v-list color="transparent">
          <v-list-item v-for="zone in tenantZones" :key="zone">{{ zone }}</v-list-item>
        </v-list>
      </section-card>
    </section-card>
  </v-container>
</template>

<script setup>
import {getTenant} from '@/api/ui/tenant.js';
import SectionCard from '@/components/SectionCard.vue';
import {computed, ref, watch} from 'vue';
import {useRoute} from 'vue-router/composables';

const route = useRoute();
const tenantId = computed(() => route?.params.tenantId);
const tenant = ref(null);

const tenantTitle = computed(() => tenant.value?.title ?? '');
const tenantZones = computed(() => tenant.value?.zones ?? []);

watch(tenantId, async (newVal, oldVal) => {
  if (!newVal) {
    tenant.value = null;
    return;
  }

  tenant.value = await getTenant({tenantId: newVal});
}, {immediate: true})
</script>

<style scoped>

</style>
