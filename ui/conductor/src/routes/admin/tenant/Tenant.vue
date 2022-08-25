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
          <v-slide-y-reverse-transition>
            <theme-btn elevation="0" @click="addSecretBegin" v-if="!addingSecret">Generate new secret</theme-btn>
          </v-slide-y-reverse-transition>
        </v-card-title>
        <v-expand-transition>
          <new-secret-form v-if="addingSecret" @commit="addSecretCommit" @rollback="addSecretRollback"/>
        </v-expand-transition>
        <v-expand-transition>
          <v-alert type="info" tile v-if="createdSecret">
            Make sure to copy your secret token now. You won't be able to see it again!
          </v-alert>
        </v-expand-transition>
        <v-list color="transparent">
          <secret-token-list-item v-if="createdSecret"
                                  :secret="createdSecret"
                                  :key="createdSecret.id"
                                  @hideToken="hideToken"
                                  @delete="deleteSecretStart"/>
          <secret-list-item v-for="secret in secretList" :key="secret.id" :secret="secret" @delete="deleteSecretStart"/>
        </v-list>
      </section-card>
    </section-card>
    <delete-secret-dialog v-model="deleteSecretDialogOpen" @commit="deleteSecretCommit"/>
  </v-container>
</template>

<script setup>
import {createSecret, deleteSecret, getTenant, listSecrets} from '@/api/ui/tenant.js';
import SectionCard from '@/components/SectionCard.vue';
import ThemeBtn from '@/components/ThemeBtn.vue';
import DeleteSecretDialog from '@/routes/admin/tenant/DeleteSecretDialog.vue';
import NewSecretForm from '@/routes/admin/tenant/NewSecretForm.vue';
import SecretListItem from '@/routes/admin/tenant/SecretListItem.vue';
import SecretTokenListItem from '@/routes/admin/tenant/SecretTokenListItem.vue';
import {compareDesc} from 'date-fns';
import {computed, ref, watch} from 'vue';
import {useRoute} from 'vue-router/composables';

const route = useRoute();
const tenantId = computed(() => route?.params.tenantId);
const tenant = ref(null);
const secrets = ref([]);
const secretList = computed(() => {
  // sorted by create time, excluding the createdSecret
  let sorted = secrets.value.sort((a, b) => compareDesc(a.createTime, b.createTime));
  if (createdSecret.value) {
    sorted = sorted.filter(s => s !== createSecret.value)
  }
  return sorted;
})

const tenantTitle = computed(() => tenant.value?.title ?? '');
const tenantZones = computed(() => tenant.value?.zones ?? []);

const addingSecret = ref(false);
const createdSecret = ref(null);

function addSecretBegin() {
  addingSecret.value = true;
}

function addSecretRollback() {
  addingSecret.value = false;
}

async function addSecretCommit(secret) {
  addingSecret.value = false;
  secret.tenant = tenant.value;
  if (createdSecret.value) {
    await hideToken();
  }
  createdSecret.value = await createSecret({secret});
}

async function hideToken() {
  createdSecret.value = null;
  secrets.value = await listSecrets({tenantId: tenant.value.id});
}

// fetch data for the tenant
watch(tenantId, async (newVal, oldVal) => {
  if (!newVal) {
    tenant.value = null;
    secrets.value = [];
    return;
  }

  tenant.value = await getTenant({id: newVal});
  secrets.value = await listSecrets({tenantId: newVal});
}, {immediate: true});

const deleteSecretDialogOpen = ref(false);
const deleteSecretDialogSecret = ref(null);

function deleteSecretStart(secret) {
  deleteSecretDialogSecret.value = secret;
  deleteSecretDialogOpen.value = true;
}

async function deleteSecretCommit() {
  if (!deleteSecretDialogSecret.value) return;

  const id = deleteSecretDialogSecret.value.id;
  await deleteSecret({id});
  deleteSecretDialogSecret.value = null;
  deleteSecretDialogOpen.value = false;
  if (createdSecret.value?.id === id) {
    createdSecret.value = null;
  } else {
    secrets.value = await listSecrets({tenantId: tenant.value.id});
  }
}

</script>

<style scoped>
</style>
