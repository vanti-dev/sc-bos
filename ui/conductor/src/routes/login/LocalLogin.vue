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
      <v-card-text v-if="displayLoginSwitch" class="d-flex justify-center">
        <a
            @click="accountStore.loginFormVisible = !accountStore.loginFormVisible"
            class="text-center">
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
import {computed, ref} from 'vue';
import {useAccountStore} from '@/stores/account.js';
import {storeToRefs} from 'pinia';

const accountStore = useAccountStore();
const password = ref('');
const username = ref('');
const {snackbar} = storeToRefs(accountStore);
const rules = {
  required: (value) => !!value || 'Required.'
};
const login = () => {
  // check if username and password are entered
  if (username.value && password.value) {
    accountStore
        .loginWithLocalAuth({username: username.value, password: password.value})
        .catch((err) => console.error('unable to log in', err));
  } else {
    console.error('username and password are required');
  }
};

const disableSignIn = computed(() => {
  const hasNone = !username.value && !password.value;
  const hasOne = !username.value && password.value || username.value && !password.value;
  return !!(hasNone || hasOne);
});

// Show/Hide the login switch depending on whether KeyCloak is enabled or not
const displayLoginSwitch = computed(() => accountStore.availableAuthProviders.includes('keyCloakAuth'));
</script>

<style lang="scss" scoped>
.v-input {
  background-color: transparent;
}
</style>
