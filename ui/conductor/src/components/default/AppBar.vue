<template>
  <v-app-bar app height="60" :clipped-left="hasNav" :clipped-right="hasSidebar" elevation="0" class="pr-7">
    <app-menu v-if="isLoggedIn"/>
    <brand-logo :theme="config.theme" outline="white" style="height: 35px; margin-left: 16px"/>
    <span class="heading">{{ appBarHeadingWithBrand }}</span>

    <v-divider vertical v-if="hasSections" class="mx-8 section-divider" inset/>

    <v-spacer/>

    <router-view name="actions"/>
    <smart-core-status-card/>
    <v-divider vertical class="mx-1" inset/>
    <account-btn btn-class="mr-0"/>
  </v-app-bar>
</template>
<script setup>
import AccountBtn from '@/components/AccountBtn.vue';
import AppMenu from '@/components/AppMenu.vue';
import BrandLogo from '@/components/BrandLogo.vue';
import smartCoreStatusCard from '@/components/smartCoreStatus/smartCoreStatusCard.vue';

import {computed} from 'vue';
import {storeToRefs} from 'pinia';

import {usePage} from '@/components/page';
import useAuthSetup from '@/composables/useAuthSetup';
import {useAccountStore} from '@/stores/account';
import {useAppConfigStore} from '@/stores/app-config';

const appConfigStore = useAppConfigStore();
const {config} = storeToRefs(appConfigStore);
const {isLoggedIn} = useAuthSetup();

const {pageTitle, hasSections, hasNav, hasSidebar} = usePage();

const store = useAccountStore();
store.loadLocalStorage();

const appBarHeadingWithBrand = computed(() => {
  const brandName = config.value.theme?.appBranding.brandName ?? 'Smart Core';

  return brandName + (isLoggedIn.value ? ' | ' + pageTitle.value : '');
});
</script>

<style lang="scss" scoped>
.v-app-bar.v-toolbar.v-sheet {
  background: var(--v-neutral-base);
}

.v-app-bar :deep(.v-toolbar__content) {
  padding-right: 0px;
}

.heading {
  font-size: 22px;
  font-weight: 300;
}

.section-divider {
  border-color: currentColor;
}
</style>
