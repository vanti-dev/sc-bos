<template>
  <v-menu bottom left offset-y max-width="100%" tile>
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
    </v-card>
  </v-menu>
</template>

<script setup>

import {useAccountStore} from '../stores/account.js';

defineProps({
  btnClass: [String, Object]
});

const accountStore = useAccountStore();
const loggedIn = accountStore.loggedIn;

</script>

<style scoped>

</style>
