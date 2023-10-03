<template>
  <div>
    <v-form @submit.prevent="login">
      <v-card-text>
        <p class="text-center">Sign in locally.</p>
      </v-card-text>
      <v-text-field
          autofocus
          label="Username"
          placeholder="Username"
          :rules="[rules.required]"
          outlined
          v-model="username"
          type="text"
          required/>
      <v-text-field
          label="Password"
          placeholder="Password"
          :rules="[rules.required]"
          outlined
          v-model="password"
          type="password"
          required/>
      <v-card-actions class="mx-2">
        <v-btn
            type="submit"
            color="primary"
            :disabled="disableSignIn"
            block
            large
            class="text-body-1 font-weight-bold mb-4">
          Sign In
        </v-btn>
      </v-card-actions>
      <v-card-text v-if="appConfig?.config?.keycloak" class="d-flex justify-center">
        <a @click="store.toggleLoginForm()" class="text-center">
          Use a different sign in method
        </a>
      </v-card-text>
    </v-form>

    <v-snackbar v-model="snackbar">
      Failed to sign in, please try again.

      <template #action="{ attrs }">
        <v-btn color="pink" text v-bind="attrs" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script setup>
import {useAccountStore} from '@/stores/account.js';
import {useAppConfigStore} from '@/stores/app-config';
import {storeToRefs} from 'pinia';
import {computed, ref} from 'vue';

const appConfig = useAppConfigStore();
const store = useAccountStore();
const password = ref('');
const username = ref('');
const {snackbar} = storeToRefs(store);
const rules = {
  required: (value) => !!value || 'Required.'
};
const login = () => {
  // check if username and password are entered
  if (username.value && password.value) {
    store
        .loginLocal(username.value, password.value)
        .catch((err) => console.error('unable to log in', err));
  } else {
    console.error('username and password are required');
  }
};
const disableSignIn = computed(() => {
  const hasNone = !username.value && !password.value;
  const hasOne = !username.value && password.value || username.value && !password.value;
  if (hasNone || hasOne) return true;
  else return false;
});
</script>

<style lang="scss" scoped>
.v-input {
  background-color: transparent;
}
</style>
