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
              Expire{{ secret.expireTime > Date.now() ? 's' : 'd' }}
              <v-tooltip bottom>
                <template #activator="{on, attrs}">
                  <span v-bind="attrs" v-on="on">{{ humanizeDate(secret.expireTime) }}</span>
                </template>
                <span>{{ Intl.DateTimeFormat('en-GB').format(secret.expireTime) }}</span>
              </v-tooltip>
            </v-list-item-subtitle>
            <v-list-item-subtitle
                v-else
                class="text-body-small neutral--text text--lighten-5">
              This secret will not expire
            </v-list-item-subtitle>
          </v-list-item-content>
          <delete-confirmation-dialog
              title="Delete Token"
              :progress-bar="deleteSecretTracker.loading"
              @confirm="delSecret(secret.id)">
            Are you sure you want to delete the token "{{ secret.note }}"?
            <template #alert-content>
              This action will stop any integrations using this token from functioning.<br><br>
              This action cannot be undone.
            </template>
            <template #confirmBtn>Delete Token</template>
            <template #activator="{on, attrs}">
              <v-list-item-action v-show="hover" class="my-0" v-bind="attrs">
                <v-btn icon small v-on="on"><v-icon color="neutral lighten-5">mdi-trash-can</v-icon></v-btn>
              </v-list-item-action>
            </template>
          </delete-confirmation-dialog>
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
import {computed, reactive, watch} from 'vue';
import {newActionTracker} from '@/api/resource';
import {compareDesc} from 'date-fns';
import dayjs from 'dayjs';
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
import NewSecretForm from './NewSecretForm.vue';
import DeleteConfirmationDialog from '@/routes/auth/third-party/components/DeleteConfirmationDialog.vue';

dayjs.extend(duration);
dayjs.extend(relativeTime);

const secretsTracker = reactive(/** @type {ActionTracker<ListSecretsResponse.AsObject>} */ newActionTracker());
const deleteSecretTracker = reactive(/** @type {ActionTracker<DeleteSecretResponse.AsObject>} */ newActionTracker());

const props = defineProps({
  account: {
    type: Object,
    default: () => {
    }
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
  refreshSecrets();
}


</script>

<style scoped>
.v-list--two-line .v-list-item, .v-list-item--two-line {
  min-height: 48px;
}
</style>
