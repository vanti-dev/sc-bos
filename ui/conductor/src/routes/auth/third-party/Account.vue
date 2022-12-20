<template>
  <v-container>
    <section-card class="mt-7 pb-4" :loading="tenantTracker.loading">
      <v-alert v-if="tenantTracker.error" type="error" text class="rounded-b-0">
        Unable to fetch tenant information: {{ tenantTracker.error.message }}
      </v-alert>
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
      <section-card class="mx-4 mt-4" :loading="secretsTracker.loading">
        <v-alert v-if="secretsTracker.error" type="error" text class="rounded-b-0">
          Unable to fetch tenant secrets: {{ secretsTracker.error.message }}
        </v-alert>
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
    <delete-secret-dialog v-model="deleteSecretDialogOpen"
                          @commit="deleteSecretCommit"
                          @rollback="deleteSecretRollback"/>
  </v-container>
</template>

<script setup>
import {newActionTracker} from '@/api/resource.js';
import {createSecret, deleteSecret, getTenant, listSecrets, secretToObject} from '@/api/ui/tenant.js';
import SectionCard from '@/components/SectionCard.vue';
import ThemeBtn from '@/components/ThemeBtn.vue';
import DeleteSecretDialog from '@/routes/auth/third-party/DeleteSecretDialog.vue';
import NewSecretForm from '@/routes/auth/third-party/NewSecretForm.vue';
import SecretListItem from '@/routes/auth/third-party/SecretListItem.vue';
import SecretTokenListItem from '@/routes/auth/third-party/SecretTokenListItem.vue';
import {ListSecretsResponse, Secret, Tenant} from '@sc-bos/ui-gen/proto/tenants_pb';
import {compareDesc} from 'date-fns';
import {computed, reactive, ref, watch} from 'vue';
import {useRoute} from 'vue-router/composables';

const route = useRoute();
const tenantId = computed(() => route?.params.tenantId);

const tenantTracker = reactive(/** @type {ActionTracker<Tenant.AsObject>} */ newActionTracker());
const secretsTracker = reactive(/** @type {ActionTracker<ListSecretsResponse.AsObject>} */ newActionTracker());
const createSecretTracker = reactive(/** @type {ActionTracker<Secret.AsObject>} */ newActionTracker());

const secretList = computed(() => {
  // sorted by create time, excluding the createdSecret
  let sorted = secretsTracker.response?.secretsList
      .map(s => secretToObject(s))
      .sort((a, b) => compareDesc(a.createTime, b.createTime));
  if (createSecretTracker.response) {
    sorted = sorted.filter(s => s.id !== createSecretTracker.response.id)
  }
  return sorted;
})

const tenantTitle = computed(() => tenantTracker.response?.title ?? '');
const tenantZones = computed(() => tenantTracker.response?.zoneNamesList ?? []);

const addingSecret = ref(false);
const createdSecret = computed(() => secretToObject(createSecretTracker.response));

function addSecretBegin() {
  addingSecret.value = true;
}

function addSecretRollback() {
  addingSecret.value = false;
}

async function addSecretCommit(secret) {
  addingSecret.value = false;
  secret.tenant = tenantTracker.response;
  if (createSecretTracker.response) {
    await hideToken();
  }
  await createSecret({secret}, createSecretTracker);
}

function hideToken() {
  createSecretTracker.response = null;
  listSecrets({tenantId: tenantTracker.response.id}, secretsTracker);
}

// fetch data for the tenant
watch(tenantId, (newVal, oldVal) => {
  if (!newVal) {
    tenantTracker.response = null;
    secretsTracker.response = null;
    createSecretTracker.response = null;
    return;
  }

  getTenant({id: newVal}, tenantTracker).catch(err => console.error(err));
  listSecrets({tenantId: newVal}, secretsTracker).catch(err => console.error(err));
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
    await listSecrets({tenantId: tenantTracker.response.id}, secretsTracker);
  }
}

function deleteSecretRollback() {
  deleteSecretDialogSecret.value = null;
  deleteSecretDialogOpen.value = false;
}

</script>

<style scoped>
</style>
