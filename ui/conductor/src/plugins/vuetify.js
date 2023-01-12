import '@mdi/font/css/materialdesignicons.css';
import Vue from 'vue';
import Vuetify from 'vuetify/lib';


Vue.use(Vuetify);

const opts = {
  theme: {
    dark: true,
    options: {
      customProperties: true
    },
    themes: {
      dark: {
        primary: {
          darken1: '#00AAC1',
          base: '#00BED6',
          lighten1: '#33CADE',
          lighten2: '#4DD1E2',
          lighten3: '#66D7E6',
          lighten4: '#80DEEB',
          lighten5: '#99E5EF'
        },
        secondary: {
          darken1: '#C5CC3C',
          base: '#DAE343',
          lighten1: '#E2E969',
          lighten2: '#E6EB7B',
          lighten3: '#E9EE8E',
          lighten4: '#EDF1A1',
          lighten5: '#F1F4B4'
        },
        secondaryTeal: {
          darken1: '#004651',
          base: '#195962',
          lighten1: '#336B74',
          lighten2: '#4D7E85',
          lighten3: '#669097',
          lighten4: '#80A3A8',
          lighten5: '#99B5B9'
        },
        accent: {
          darken1: '#FFA400',
          base: '#FFAF25',
          lighten1: '#FFB833',
          lighten2: '#FFC14D',
          lighten3: '#FFCA66',
          lighten4: '#FFD380',
          lighten5: '#FFDB99'
        },
        neutral: {
          darken1: '#101820',
          base: '#282F36',
          lighten1: '#40464D',
          lighten2: '#585D63',
          lighten3: '#707479',
          lighten4: '#888C90',
          lighten5: '#9FA3A6',
          lighten6: '#B7BABC',
          lighten7: '#CFD1D2',
          lighten8: '#ECEDED',
          lighten9: '#F7F7F7'
        },
        error: {
          darken1: '#BB0434',
          base: '#D0043C',
          lighten1: '#D93661',
          lighten2: '#DE4F75',
          lighten3: '#E36889',
          lighten4: '#E8829D',
          lighten5: '#EC9BB0'
        },
        success: {
          darken1: '#00613E',
          base: '#008052',
          lighten1: '#00955F',
          lighten2: '#00A76B',
          lighten3: '#00B674',
          lighten4: '#00C27C',
          lighten5: '#0CCE88'
        },
        info: '#00BED6',
        warning: '#FFA400'
      }
    }
  }
};

export default new Vuetify(opts);
