<template>
  <div>
    <v-btn
        variant="text"
        elevation="0"
        v-if="!loggedIn"
        :class="btnClass"
        to="/login">
      <v-icon start size="24">mdi-account-circle-outline</v-icon>
      Sign in
    </v-btn>
    <v-menu v-else location="bottom left">
      <template #activator="{ props }">
        <v-btn icon="true" variant="text" :class="btnClass" v-bind="props">
          <v-icon :start="!loggedIn" size="24">mdi-account-circle-outline</v-icon>
        </v-btn>
      </template>

      <v-card
          tile
          theme="light"
          class="text-center"
          min-width="300px">
        <v-avatar
            size="64"
            style="background: #eee; padding: 40px; margin-top: 24px">
          <v-icon size="64">mdi-account-circle-outline</v-icon>
        </v-avatar>
        <v-card-title class="justify-center pb-0">
          {{ accountStore.fullName }}
        </v-card-title>
        <v-card-subtitle class="mb-4">
          {{ accountStore.email }}
        </v-card-subtitle>
        <v-card-actions>
          <v-btn variant="tonal" @click="logout" block>Sign out</v-btn>
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


