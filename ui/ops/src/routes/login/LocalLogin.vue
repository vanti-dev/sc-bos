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
          variant="outlined"
          v-model="username"
          type="text"
          required/>
      <v-text-field
          label="Password"
          placeholder="Password"
          :rules="[rules.required]"
          variant="outlined"
          v-model="password"
          :type="showPassword ? 'text' : 'password'"
          :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
          @click:append-inner="showPassword = !showPassword"
          required/>
      <v-card-actions class="mx-2">
        <v-btn
            type="submit"
            color="primary"
            variant="elevated"
            :disabled="disableSignIn"
            block
            size="large"
            class="text-body-1 font-weight-bold mb-4">
          Sign In
        </v-btn>
      </v-card-actions>
      <v-expand-transition>
        <v-alert v-if="errorStr" type="error" :text="errorStr"/>
      </v-expand-transition>
    </v-form>
  </div>
</template>

<script setup>
import {useAccountStore} from '@/stores/account.js';
import {computed, ref} from 'vue';

const accountStore = useAccountStore();
const password = ref('');
const showPassword = ref(false);
const username = ref('');
const rules = {
  required: (value) => !!value || 'Required.'
};

const loginError = ref(null);
const errorStr = computed(() => {
  const err = loginError.value;
  if (!err) {
    return '';
  }
  if (err.error_description) return err.error_description;

  if (err.error) {
    switch (err.error) {
      case 'invalid_grant':
      case 'invalid_client':
        return 'Incorrect username or password';
      case 'account_locked':
        return 'Account is locked. Please contact support.';
      case 'account_disabled':
        return 'Account is disabled. Please contact support.';
      default:
        return 'Login failed: ' + err.error;
    }
  }
  return (err.message ?? err);
})

const login = async () => {
  // check if username and password are entered
  if (username.value && password.value) {
    const values = {username: username.value, password: password.value};

    // Forwards the login request to the account store with the values from the form
    try {
      await accountStore.loginWithLocalAuth(values);
      loginError.value = null;
    } catch (e) {
      loginError.value = e;
    }
  } else {
    console.error('username and password are required');
  }
};

const disableSignIn = computed(() => {
  const hasNone = !username.value && !password.value;
  const hasOne = !username.value && password.value || username.value && !password.value;
  return !!(hasNone || hasOne);
});
</script>

<style lang="scss" scoped>
.v-input {
  background-color: transparent;
}
</style>
