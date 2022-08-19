<template>
  <v-btn v-if="!loggedIn" tile elevation="0" :class="btnClass" @click="login">
    <v-icon left>mdi-account-circle-outline</v-icon>
    Log in
  </v-btn>
  <v-menu v-else bottom left offset-y max-width="100%" tile>
    <template #activator="{on, attrs}">
      <v-btn tile elevation="0" :class="btnClass" v-bind="attrs" v-on="on">
        <v-icon :left="!loggedIn">mdi-account-circle-outline</v-icon>
        <template v-if="!loggedIn">Log in</template>
      </v-btn>
    </template>

    <v-card tile :light="$vuetify.theme.dark">
      <v-card-text>
        This is where the account menu would be
      </v-card-text>
      <v-card-actions>
        <v-btn elevation="0" @click="logout">Log out</v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script setup>

import {computed} from 'vue';
import {useAccountStore} from '../stores/account.js';

defineProps({
  btnClass: [String, Object]
});

const accountStore = useAccountStore();
const loggedIn = computed(() => accountStore.loggedIn);

function login() {
  accountStore.login(['profile'])
      .catch(err => console.error('error during login', err));
}

function logout() {
  accountStore.logout()
      .catch(err => console.error('error during logout', err));
}

</script>

<style scoped>

</style>
