import '@mdi/font/css/materialdesignicons.css';
import Vue from 'vue';
import Vuetify from 'vuetify/lib';


Vue.use(Vuetify);

const opts = {
  theme: {
    dark: true,
    options: {
      customProperties: true,
    },
    themes: {
      dark: {
        primary: "#00BED6",
        secondary: "#DAE343",
        accent: "#FFA400",

        bg: "#101820",
        card: "#283036",
        sectionText: "#3F454B",

        design: "#DAE343",
        commission: "#00BED6",
        operate: "#FFA400",
        admin: "#A2A2A2",

        error: "#FF6200",
        info: "#00BED6",
        success: "#DAE343",
        warning: "#FFA400",
      },
    },
  },
};

export default new Vuetify(opts);
