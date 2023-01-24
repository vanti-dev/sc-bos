<template>
  <v-card flat tile class="pa-0">
    <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Tokens</v-subheader>
    <v-list two-line class="pt-0">
      <v-list-item v-for="secret of secretList" :key="secret.id" class="py-0">
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
        <!--        <v-list-item-avatar>{{ secret.lastUseTime }}</v-list-item-avatar>-->
      </v-list-item>
    </v-list>
    <v-card-actions class="px-4 pb-4">
      <new-secret-form
          :account-name="account.title"
          :account-id="account.id">
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
import {listSecrets, secretToObject} from '@/api/ui/tenant';
import {computed, reactive, ref, watch} from 'vue';
import {newActionTracker} from '@/api/resource';
import {compareDesc} from 'date-fns';
import dayjs from 'dayjs';
import duration from 'dayjs/plugin/duration';
import relativeTime from 'dayjs/plugin/relativeTime';
import NewSecretForm from '@/routes/auth/third-party/NewSecretForm.vue';

dayjs.extend(duration);
dayjs.extend(relativeTime);

const secretsTracker = reactive(/** @type {ActionTracker<ListSecretsResponse.AsObject>} */ newActionTracker());

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
watch(() => props.account, (account) => {
  if (!account.id) {
    secretsTracker.response = null;
    return;
  }
  listSecrets({tenantId: account.id}, secretsTracker).catch(err => console.error(err));
}, {immediate: true});

/**
 * @param {Date} date
 * @return {string}
 */
function humanizeDate(date) {
  return dayjs(date).from(dayjs());
}


</script>

<style scoped>
.v-list--two-line .v-list-item, .v-list-item--two-line {
  min-height: 48px;
}
</style>
