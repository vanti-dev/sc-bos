<template>
  <v-menu v-bind="showMenu" bottom offset-y max-width="100%" tile transition="slide-y-transition">
    <template #activator="{on, attrs}">
      <v-btn icon tile v-bind="attrs" v-on="on" style="background: #FFFFFF1A" :class="btnClass">
        <menu-icon :color="themeColor"/>
      </v-btn>
    </template>
    <v-card width="100vw" :light="$vuetify.theme.dark" tile class="menu-card">
      <nav-button :to="{name:'design'}"
                  :color="$vuetify.theme.currentTheme.design"
                  title="Designer"
                  subtitle="Create our buildings, floors, floor plans, and spaces."/>
      <nav-button :to="{name:'commission'}"
                  :color="$vuetify.theme.currentTheme.commission"
                  title="Commissioner"
                  subtitle="Configure your devices and set their job roles within spaces."/>
      <nav-button :to="{name:'operate'}"
                  :color="$vuetify.theme.currentTheme.operate"
                  title="Operator"
                  subtitle="Get a holistic view of the health of your environment and access analytics."/>
      <nav-button :to="{name:'admin'}"
                  :color="$vuetify.theme.currentTheme.admin"
                  title="Administrator"
                  subtitle="Generate tenant virtual ID cards, create users and set their permission levels."/>
    </v-card>
  </v-menu>
</template>

<script setup>
import {ref} from 'vue';
import MenuIcon from './MenuIcon.vue';
import NavButton from './NavButton.vue';
import {usePage} from './page.js';

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
