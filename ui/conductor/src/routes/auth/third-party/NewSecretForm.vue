<template>
  <v-form @submit.prevent="addSecretCommit" v-model="formValid">
    <v-card-text class="pt-0">
      <v-text-field
          label="Note"
          v-model="newSecret.note"
          style="max-width: 400px"
          required
          hint="Easily recognisable: 'Read only', 'Dev access'"
          :rules="noteRules"/>
      <div class="d-flex justify-start align-baseline">
        <v-select
            v-model="newSecret.expiresIn"
            :items="Object.values(suggestedExpiresIn)"
            hide-details
            class="expires-in"
            label="Expiration"/>
        <v-menu
            v-if="newSecret.expiresIn === suggestedExpiresIn.custom"
            v-model="customExpiryMenuVisible"
            :close-on-content-click="false"
            offset-y>
          <template #activator="{on, attrs}">
            <v-text-field
                v-model="newSecret.expiresAt"
                placeholder="yyyy-mm-dd"
                readonly
                v-bind="attrs"
                v-on="on"
                class="expires-at"
                hide-details="auto"
                :rules="expiresAtRules"/>
          </template>
          <v-date-picker v-model="newSecret.expiresAt" @input="customExpiryMenuVisible = false" no-title/>
        </v-menu>
        <span
            v-else-if="newSecret.expiresIn === suggestedExpiresIn.noExpiry"
            class="expires-in-message never">The token will never expire!</span>
        <span v-else class="expires-in-message">The token will expire <relative-date
            :date="computedExpiresAt"
            no-relative/></span>
      </div>
    </v-card-text>
    <v-card-actions class="justify-end">
      <theme-btn type="submit" depressed :disabled="!formValid">Create Secret</theme-btn>
      <v-btn type="cancel" text @click.prevent="addSecretRollback">Cancel</v-btn>
    </v-card-actions>
    <v-divider class="mt-4"/>
  </v-form>
</template>

<script setup>
import {DAY, useNow} from '@/components/now.js';
import ThemeBtn from '@/components/ThemeBtn.vue';
import {add} from 'date-fns';
import {computed, reactive, ref} from 'vue';

const emit = defineEmits(['commit', 'revert']);

const suggestedExpiresIn = {
  week: '7 days',
  month: '30 days',
  twoMonths: '60 days',
  quarter: '90 days',
  custom: 'Custom...',
  noExpiry: 'No expiration'
};

const newSecret = reactive({
  note: '',
  expiresIn: suggestedExpiresIn.month,
  expiresAt: null
});
const customExpiryMenuVisible = ref(false);

const {now} = useNow(DAY);
const computedExpiresAt = computed(() => {
  const duration = {};
  switch (newSecret.expiresIn) {
    case suggestedExpiresIn.week:
      duration.days = 7;
      break;
    case suggestedExpiresIn.month:
      duration.days = 30;
      break;
    case suggestedExpiresIn.twoMonths:
      duration.days = 60;
      break;
    case suggestedExpiresIn.quarter:
      duration.days = 90;
      break;
    case suggestedExpiresIn.custom:
      return newSecret.expiresAt;
    default:
      return null;
  }

  return add(now.value, duration);
});

/**
 *
 */
function addSecretReset() {
  newSecret.note = '';
  newSecret.expiresIn = suggestedExpiresIn.month;
  newSecret.expiresAt = null;
  customExpiryMenuVisible.value = false;
}

/**
 *
 */
function addSecretRollback() {
  emit('rollback');
  addSecretReset();
}

/**
 *
 */
function addSecretCommit() {
  emit('commit', mintSecret());
  addSecretReset();
}

/**
 *
 * @return {{note: UnwrapRef<string>, expireTime: unknown}}
 */
function mintSecret() {
  return {
    note: newSecret.note,
    expireTime: computedExpiresAt.value
  };
}

const formValid = ref(false);
const noteRules = [
  v => Boolean(v) || Boolean(v.trim()) || 'Notes can\'t be blank'
];
const expiresAtRules = [
  v => Boolean(v) || 'Expiry date can\'t be blank'
];
</script>

<style scoped>
.expires-in {
  width: 180px;
  flex-grow: 0;
}

.expires-at {
  min-width: 150px;
  flex-grow: 0;
  margin-left: 8px;
}

.expires-in-message {
  margin-left: 8px;
  opacity: .7;
}
</style>
