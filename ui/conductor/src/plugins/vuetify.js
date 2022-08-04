import 'vuetify/dist/vuetify.min.css';
import Vue from 'vue';
import Vuetify from 'vuetify/lib/framework';
import '@mdi/font/css/materialdesignicons.css';

Vue.use(Vuetify);

const opts = {
  theme: {
    dark: true,
    themes: {
      dark: {
        primary: '#00BED6',
        secondary: '#DAE343',
        accent: '#FFA400',

        bg: '#101820',
        card: '#474747',
        sectionText: '#5C6165',

        design: '#DAE343',
        commission: '#00BED6',
        operate: '#FFA400',
        admin: '#00BED6',

        error: '#FF6200',
        info: '#00BED6',
        success: '#DAE343',
        warning: '#FFA400'
      }
    }
  }
};

export default new Vuetify(opts);
