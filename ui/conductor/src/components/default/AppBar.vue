<template>
  <v-app-bar height="60" elevation="0" class="pr-7" color="neutral">
    <app-menu v-if="accountStore.isLoggedIn"/>
    <brand-logo :theme="config.theme" outline="white" style="height: 35px" class="ml-4 mr-2"/>
    <span class="heading">{{ appBarHeadingWithBrand }}</span>

    <v-divider vertical v-if="hasSections" class="mx-8 section-divider" inset/>

    <v-spacer/>

    <router-view name="actions"/>
    <smart-core-status-card/>
    <span
        v-if="!hideAccountBtn"
        class="d-flex flex-row">
      <v-divider vertical class="mx-1 my-1" inset/>
      <account-btn btn-class="mr-0"/>
    </span>
  </v-app-bar>
</template>
<script setup>
import AccountBtn from '@/components/AccountBtn.vue';
import AppMenu from '@/components/AppMenu.vue';
import BrandLogo from '@/components/BrandLogo.vue';

import {usePage} from '@/components/page';
import SmartCoreStatusCard from '@/components/smartCoreStatus/SmartCoreStatusCard.vue';
import {useAccountStore} from '@/stores/account';
import {useUiConfigStore} from '@/stores/uiConfig.js';
import {storeToRefs} from 'pinia';

import {computed} from 'vue';
import {useRoute} from 'vue-router';

const uiConfig = useUiConfigStore();
const {config} = storeToRefs(uiConfig);
const accountStore = useAccountStore();
const route = useRoute();

const {pageTitle, hasSections} = usePage();

const appBarHeadingWithBrand = computed(() => {
  const brandName = config.value.theme?.appBranding.brandName ?? 'Smart Core';

  const hasPageTitle = Boolean(pageTitle.value) && accountStore.isLoggedIn;
  return brandName + (hasPageTitle ? ' | ' + pageTitle.value : '');
});

const isLoginPage = computed(() => route.path === '/login');
const isAuthDisabled = computed(() => accountStore.isAuthenticationDisabled);

const hideAccountBtn = computed(() => isLoginPage.value || isAuthDisabled.value);
</script>

<style lang="scss" scoped>
.v-app-bar.v-toolbar.v-sheet {
  background: rgb(var(--v-theme-neutral));
}

.v-app-bar :deep(.v-toolbar__content) {
  padding-right: 0;
}

.heading {
  font-size: 22px;
  font-weight: 300;
}

.section-divider {
  border-color: currentColor;
}
</style>
