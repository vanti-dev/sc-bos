import ChildOnlyPage from '@/components/ChildOnlyPage.vue';
import thirdParty from '@/routes/auth/third-party/route.js';
import {route} from '@/util/router.js';

export default {
  name: 'auth',
  path: '/auth',
  redirect: '/auth/third-party',
  components: {
    default: ChildOnlyPage,
    nav: () => import('./AdminNav.vue')
  },
  children: [
    {path: 'users', component: () => import('./users/Users.vue')},
    ...route(thirdParty)
  ],
  meta: {
    title: 'Auth'
  }
}
