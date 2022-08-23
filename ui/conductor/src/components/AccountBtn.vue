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

    <v-card tile :light="$vuetify.theme.dark" class="text-center" min-width="300px">
      <v-avatar size="64" style="background: #eee; padding: 40px; margin-top: 24px">
        <v-icon size="64">mdi-account-circle-outline</v-icon>
      </v-avatar>
      <v-card-title class="justify-center">
        {{ accountStore.fullName }}
      </v-card-title>
      <v-card-subtitle>
        {{ accountStore.email }}
      </v-card-subtitle>
      <v-card-actions>
        <v-btn elevation="0" @click="logout" block>Log out</v-btn>
      </v-card-actions>
    </v-card>
  </v-menu>
</template>

<script setup>

import {useAccountStore} from '@/stores/account.js';
import {computed} from 'vue';

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
