<template>
  <v-container>
    <section-card class="mt-7 pb-4">
      <v-card-title>{{ tenantTitle }}</v-card-title>
      <v-card-text>
        <v-combobox v-model="tenantZones" :items="['L3', 'L4']" multiple label="Occupies zones" hide-details outlined>
          <template #selection="data">
            <v-chip :key="JSON.stringify(data.item)"
                    v-bind="data.attrs"
                    :input-value="data.selected"
                    :disabled="data.disabled"
                    close
                    outlined
                    @click:close="data.parent.selectItem(data.item)">
              {{ data.item }}
            </v-chip>
          </template>
        </v-combobox>
      </v-card-text>
      <section-card class="mx-4 mt-4">
        <v-card-title><span>Secrets</span>
          <v-spacer/>
          <theme-btn elevation="0">Generate new secret</theme-btn>
        </v-card-title>
        <v-list color="transparent">
          <secret-list-item v-for="secret in secrets" :key="secret.id" :secret="secret"/>
        </v-list>
      </section-card>
    </section-card>
  </v-container>
</template>

<script setup>
import {getTenant, listSecrets} from '@/api/ui/tenant.js';
import SectionCard from '@/components/SectionCard.vue';
import ThemeBtn from '@/components/ThemeBtn.vue';
import SecretListItem from '@/routes/admin/tenant/SecretListItem.vue';
import {computed, ref, watch} from 'vue';
import {useRoute} from 'vue-router/composables';

const route = useRoute();
const tenantId = computed(() => route?.params.tenantId);
const tenant = ref(null);
const secrets = ref([]);

const tenantTitle = computed(() => tenant.value?.title ?? '');
const tenantZones = computed(() => tenant.value?.zones ?? []);

watch(tenantId, async (newVal, oldVal) => {
  if (!newVal) {
    tenant.value = null;
    secrets.value = [];
    return;
  }

  tenant.value = await getTenant({tenantId: newVal});
  secrets.value = await listSecrets({tenantId: newVal});
}, {immediate: true})
</script>

<style scoped>

</style>
