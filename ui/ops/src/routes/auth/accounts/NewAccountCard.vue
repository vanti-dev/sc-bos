<template>
  <v-card min-width="440">
    <v-form @submit.prevent="onSave"
            v-model="formValid"
            ref="formRef"
            :disabled="formDisabled">
      <v-card-title>New Account</v-card-title>
      <v-card-text class="ga-2 d-flex flex-column">
        <v-btn-toggle v-model="accountType"
                      variant="outlined"
                      divided
                      mandatory
                      density="comfortable"
                      class="mb-2">
          <v-btn :value="Account.Type.USER_ACCOUNT" text="User" v-tooltip:bottom="'User account'"/>
          <v-btn :value="Account.Type.SERVICE_ACCOUNT" text="Service" v-tooltip:bottom="'Service account'"/>
        </v-btn-toggle>
        <template v-if="accountType === Account.Type.USER_ACCOUNT">
          <v-text-field label="Full Name"
                        v-model.trim="fullName"
                        :rules="fullNameRules"
                        required
                        :counter="100"
                        autocomplete="off"/>
          <v-text-field label="Username"
                        v-model.trim="username"
                        :rules="usernameRules"
                        required
                        :counter="100"
                        autocomplete="off"/>
          <v-text-field label="Password"
                        v-model.trim="password"
                        required
                        type="password"
                        autocomplete="off"/>
        </template>
        <template v-if="accountType === Account.Type.SERVICE_ACCOUNT">
          <v-text-field label="Account Name"
                        v-model.trim="accountName"
                        :rules="accountNameRules"
                        required
                        :counter="100"
                        autocomplete="off"/>
          <v-textarea label="Description"
                      v-model.trim="description"
                      :counter="500"
                      autocomplete="off"/>
        </template>
      </v-card-text>
      <v-card-actions>
        <v-btn text="Create"
               color="primary"
               type="submit"
               variant="flat"
               :disabled="!formValid"
               :loading="saveLoading"/>
        <v-btn text="Cancel" @click="onCancel"/>
      </v-card-actions>
    </v-form>
  </v-card>
</template>

<script setup>
import {createAccount, createServiceCredential} from '@/api/ui/account.js';
import {Account} from '@vanti-dev/sc-bos-ui-gen/proto/account_pb';
import {computed, ref, watch} from 'vue';

const props = defineProps({
  name: {
    type: String,
    default: undefined,
  }
});

const emit = defineEmits(['save', 'cancel', 'error']);

const formRef = ref(null);
const formValid = ref(false);
const accountType = ref(Account.Type.USER_ACCOUNT);
const fullName = ref('');
const fullNameRules = computed(() => {
  return [
    (v) => !!v || 'Full Name is required',
    (v) => v.length >= 3 || 'Full Name must be at least 3 characters',
    (v) => v.length <= 100 || 'Full Name must be less than 100 characters',
  ];
})
const username = ref('');
const usernameRules = computed(() => {
  return [
    (v) => !!v || 'Username is required',
    (v) => v.length >= 3 || 'Username must be at least 3 characters',
    (v) => v.length <= 100 || 'Username must be less than 100 characters',
    (v) => /^[a-zA-Z0-9.-_@]+$/.test(v) || 'Username must use only alphanumerics with .-_@',
  ];
})
const password = ref('');
// guess a suitable username based on the entered full name
const nameToUsername = (name) => {
  if (!name) return '';
  return name.toLowerCase().replace(/[^a-zA-Z0-9.-_@]+/g, '.');
};
watch(fullName, (value, oldValue) => {
  if (username.value && username.value !== nameToUsername(oldValue)) {
    return; // manual username has been entered
  }
  username.value = nameToUsername(value);
});

const accountName = ref('');
const accountNameRules = computed(() => {
  return [
    (v) => !!v || 'Account Name is required',
    (v) => v.length >= 3 || 'Account Name must be at least 3 characters',
    (v) => v.length <= 100 || 'Account Name must be less than 100 characters',
  ];
});
const description = ref('');

const formDisabled = ref(false);
const saveLoading = ref(false);
const reset = () => {
  const form = formRef.value;
  if (!form) return;
  form.reset();
  accountType.value = Account.Type.USER_ACCOUNT;
}
const onSave = async () => {
  saveLoading.value = true;
  formDisabled.value = true;
  try {
    switch (accountType.value) {
      case Account.Type.USER_ACCOUNT: {
        const createAccountReq = {
          name: props.name,
          account: {
            type: Account.Type.USER_ACCOUNT,
            displayName: fullName.value,
            username: username.value,
          },
          password: password.value,
        }
        const account = await createAccount(createAccountReq);
        emit('save', {account});
        break;
      }
      case Account.Type.SERVICE_ACCOUNT: {
        const createAccountReq = {
          name: props.name,
          account: {
            type: Account.Type.SERVICE_ACCOUNT,
            displayName: accountName.value,
            description: description.value,
          }
        };
        const account = await createAccount(createAccountReq);
        const createServiceCredentialReq = {
          name: props.name,
          serviceCredential: {
            accountId: account.id,
            displayName: 'Credential 1',
          }
        };
        const serviceCredential = await createServiceCredential(createServiceCredentialReq);
        emit('save', {account, serviceCredential});
        break;
      }
    }
  } catch (e) {
    console.warn(e);
  } finally {
    saveLoading.value = false;
    formDisabled.value = false;
  }
};
const onCancel = () => {
  emit('cancel');
}

defineExpose({
  reset,
})
</script>

<style scoped>

</style>