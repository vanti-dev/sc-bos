<template>
  <div>
    <v-form @submit.prevent="login">
      <v-card class="pa-4">
        <div class="d-flex justify-end">
          <v-btn @click="store.toggleLoginDialog()" text dense>
            <v-icon> mdi-close </v-icon>
          </v-btn>
        </div>
        <v-card-title class="pt-8 justify-center text-h5 font-weight-semibold">
          Sign in to Smart Core
        </v-card-title>
        <v-card-text>
          <p class="text-center">Sign in locally.</p>
          <v-text-field
            label="Username"
            placeholder="Username"
            :rules="[rules.required]"
            outlined
            v-model="username"
            type="text"
            required
          ></v-text-field>
          <v-text-field
            label="Password"
            placeholder="Password"
            :rules="[rules.required]"
            outlined
            v-model="password"
            type="password"
            required
          ></v-text-field>
        </v-card-text>
        <v-card-actions class="mx-2">
          <v-btn
            type="submit"
            color="primary"
            block
            large
            class="font-weight-bold mb-4"
            >Sign In
          </v-btn>
        </v-card-actions>
        <v-card-text class="d-flex justify-center">
          <a @click="store.toggleLoginForm()" class="text-center">
            Use a different sign in method
          </a>
        </v-card-text>
      </v-card>
    </v-form>

    <v-snackbar v-model="snackbar">
      Failed to sign in, please try again.

      <template v-slot:action="{ attrs }">
        <v-btn color="pink" text v-bind="attrs" @click="snackbar = false">
          Close
        </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script setup>
import { useAccountStore } from "@/stores/account.js";
import { storeToRefs } from "pinia";
import { ref } from "vue";

const store = useAccountStore();
const password = ref("");
const username = ref("");
const { snackbar } = storeToRefs(store);

const rules = {
  required: (value) => !!value || "Required.",
};

const login = () => {
  //check if username and password are entered
  if (username.value && password.value) {
    store
      .loginLocal(username.value, password.value)
      .catch((err) => console.error("unable to log in", err));
  } else {
    console.error("username and password are required");
  }
};
</script>
