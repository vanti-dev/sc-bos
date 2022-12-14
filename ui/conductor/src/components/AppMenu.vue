<template>
  <v-menu v-bind="showMenu" bottom offset-y max-width="100%" tile transition="slide-y-transition">
    <template #activator="{on, attrs}">
      <v-btn icon tile v-bind="attrs" v-on="on" style="background: #FFFFFF1A" :class="btnClass">
        <menu-icon :color="themeColor"/>
      </v-btn>
    </template>
    <v-card :light="$vuetify.theme.dark" tile width="100vw">
      <nav-button :to="{name:'auth'}"
                  subtitle=""
                  title="Auth"/>
      <nav-button :to="{name:'devices'}"
                  subtitle=""
                  title="Devices"/>
      <nav-button :to="{name:'ops'}"
                  subtitle=""
                  title="Operations"/>
      <nav-button :to="{name:'site-settings'}"
                  subtitle=""
                  title="Site Settings"/>
      <nav-button :to="{name:'workflows'}"
                  subtitle=""
                  title="Workflows & Automations"/>
    </v-card>
  </v-menu>
</template>

<script setup>
import MenuIcon from '@/components/MenuIcon.vue';
import NavButton from '@/components/NavButton.vue';
import {usePage} from '@/components/page.js';
import {ref} from 'vue';

defineProps({
  btnClass: [String, Object]
});

const {themeColor} = usePage();

const showMenu = ref(false);

</script>

<style scoped>
.icon circle {
  stroke: none;
  fill: v-bind(logoColor);
}

.menu-card {
  display: grid;
  grid-auto-flow: column;
  grid-gap: 30px;
  justify-content: start;
  padding: 30px 42px;
  background: #ECEDED;
}
</style>
