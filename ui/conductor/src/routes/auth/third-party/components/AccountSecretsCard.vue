<template>
  <v-card flat tile class="pa-0">
    <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Tokens</v-subheader>
    <v-list two-line class="pt-0">
      <v-progress-linear color="primary" indeterminate :active="secretsTracker.loading"/>
      <v-hover v-slot="{hover}" v-for="secret of secretList" :key="secret.id">
        <v-list-item class="py-0">
          <v-list-item-content class="py-0">
            <v-list-item-title>{{ secret.note }}</v-list-item-title>
            <v-list-item-subtitle
                v-if="secret.expireTime"
                class="text-body-small neutral--text text--lighten-5">
              Expire{{ secret.expireTime > Date.now() ? 's' : 'd' }} {{ humanizeDate(secret.expireTime) }}
            </v-list-item-subtitle>
            <v-list-item-subtitle
                v-else
                class="text-body-small neutral--text text--lighten-5">
              This secret will not expire
            </v-list-item-subtitle>
          </v-list-item-content>
          <v-dialog v-model="deleteConfirmation[secret.id]" max-width="320">
            <v-card class="pa-2">
              <v-card-title class="text-h4 error--text text--lighten">Delete Account</v-card-title>
              <v-card-text>
                Are you sure you want to delete the secret "{{ secret.note }}"?<br><br>
                <span class="font-bold error--text">Note: This action cannot be undone</span>
              </v-card-text>
              <v-card-actions>
                <v-spacer/>
                <v-btn @click="deleteConfirmation[secret.id] = false" color="primary">Cancel</v-btn>
                <v-btn @click="delSecret(secret.id)" color="error">Delete</v-btn>
              </v-card-actions>
              <v-progress-linear color="primary" indeterminate :active="deleteSecretTracker.loading"/>
            </v-card>
            <template #activator="{on, attrs}">
              <v-list-item-action v-show="hover" class="my-0" v-bind="attrs">
                <v-btn icon small v-on="on"><v-icon color="neutral lighten-5">mdi-trash-can</v-icon></v-btn>
              </v-list-item-action>
            </template>
          </v-dialog>
        </v-list-item>
      </v-hover>
    </v-list>
    <v-card-actions class="px-4 pb-4">
      <new-secret-form
          :account-name="account.title"
          :account-id="account.id"
          @finished="refreshSecrets">
        <template #activator="{on}">
          <v-btn width="100%" color="primary" class="font-weight-bold" v-on="on">
            Create new token
            <v-icon right>mdi-key</v-icon>
          </v-btn>
        </template>
      </new-secret-form>
    </v-card-actions>
  </v-card>
</template>

<script setup>
import {deleteSecret, listSecrets, secretToObject} from '@/api/ui/tenant';
import {computed, reactive, ref, watch} from 'vue';
import {newActionTracker} from '@/api/resource';
import {compareDesc} from 'date-fns';
import dayjs from 'dayjs';
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
import NewSecretForm from './NewSecretForm.vue';

dayjs.extend(duration);
dayjs.extend(relativeTime);

const deleteConfirmation = ref({});

const secretsTracker = reactive(/** @type {ActionTracker<ListSecretsResponse.AsObject>} */ newActionTracker());
const deleteSecretTracker = reactive(/** @type {ActionTracker<DeleteSecretResponse.AsObject>} */ newActionTracker());

const props = defineProps({
  account: {
    type: Object,
    default: () => {}
  }
});

const secretList = computed(() => {
  // sorted by create time
  return secretsTracker.response?.secretsList
      .map(s => secretToObject(s))
      .sort((a, b) => compareDesc(a.createTime, b.createTime));
});

// fetch secret data for account
watch(() => props.account, () => refreshSecrets(), {immediate: true});

/**
 */
function refreshSecrets() {
  if (!props.account.id) {
    secretsTracker.response = null;
    return;
  }
  listSecrets({tenantId: props.account.id}, secretsTracker).catch(err => console.error(err));
}

/**
 * @param {Date} date
 * @return {string}
 */
function humanizeDate(date) {
  return dayjs(date).from(dayjs());
}

/**
 * @param {string} id
 */
async function delSecret(id) {
  await deleteSecret({id}, deleteSecretTracker);
  deleteConfirmation.value[id] = false;
  refreshSecrets();
}


</script>

<style scoped>
.v-list--two-line .v-list-item, .v-list-item--two-line {
  min-height: 48px;
}
</style>
