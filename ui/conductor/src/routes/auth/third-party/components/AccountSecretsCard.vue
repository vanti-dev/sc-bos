<template>
  <v-card flat tile class="pa-0">
    <v-subheader class="text-title-caps-large neutral--text text--lighten-3">Tokens</v-subheader>
    <v-list two-line class="pt-0">
      <v-list-item v-for="secret of secretList" :key="secret.id">
        <v-list-item-content class="pt-0">
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

dayjs.extend(duration);
dayjs.extend(relativeTime);

const secretsTracker = reactive(/** @type {ActionTracker<ListSecretsResponse.AsObject>} */ newActionTracker());
const createSecretTracker = reactive(/** @type {ActionTracker<Secret.AsObject>} */ newActionTracker());

const props = defineProps({
  accountId: {
    type: String,
    default: ''
  }
});

const secretList = computed(() => {
  // sorted by create time, excluding the createdSecret
  let sorted = secretsTracker.response?.secretsList
      .map(s => secretToObject(s))
      .sort((a, b) => compareDesc(a.createTime, b.createTime));
  if (createSecretTracker.response) {
    sorted = sorted.filter(s => s.id !== createSecretTracker.response.id);
  }
  return sorted;
});

const addingSecret = ref(false);
const createdSecret = computed(() => secretToObject(createSecretTracker.response));

// fetch secret data for account
watch(() => props.accountId, (id) => {
  if (!id) {
    secretsTracker.response = null;
    createSecretTracker.response = null;
    return;
  }
  listSecrets({tenantId: id}, secretsTracker).catch(err => console.error(err));
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

</style>
