<template>
  <div>
    <v-btn
        text
        elevation="0"
        v-if="!loggedIn"
        :class="btnClass"
        to="/login">
      <v-icon left>mdi-account-circle-outline</v-icon>
      Sign in
    </v-btn>
    <v-menu v-else bottom left offset-y max-width="100%" tile>
      <template #activator="{ props }">
        <v-btn icon elevation="0" :class="btnClass" v-bind="props">
          <v-icon :left="!loggedIn">mdi-account-circle-outline</v-icon>
        </v-btn>
      </template>

      <v-card
          tile
          :light="$vuetify.theme.dark"
          class="text-center"
          min-width="300px">
        <v-avatar
            size="64"
            style="background: #eee; padding: 40px; margin-top: 24px">
          <v-icon size="64">mdi-account-circle-outline</v-icon>
        </v-avatar>
        <v-card-title class="justify-center">
          {{ accountStore.fullName }}
        </v-card-title>
        <v-card-subtitle>
          {{ accountStore.email }}
        </v-card-subtitle>
        <v-card-actions>
          <v-btn elevation="0" @click="logout" block>Sign out</v-btn>
        </v-card-actions>
      </v-card>
    </v-menu>
  </div>
</template>

<script setup>
import {useAccountStore} from '@/stores/account.js';
import {computed} from 'vue';

defineProps({
  btnClass: {
    type: [String, Object],
    default: ''
  }
});

const accountStore = useAccountStore();
const loggedIn = computed(() => accountStore.isLoggedIn);

const logout = async () => {
  await accountStore.logout();
};
</script>


