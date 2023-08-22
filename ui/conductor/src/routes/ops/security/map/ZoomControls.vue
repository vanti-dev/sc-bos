<template>
  <div class="root">
    <v-btn
        v-for="b in activeBtns"
        :key="b.c"
        :class="[b.c, 'black--text']"
        @click="handleClick(b.c)"
        icon
        elevation="2">
      <slot :name="b.c">
        <v-icon>{{ b.i }}</v-icon>
      </slot>
    </v-btn>
  </div>
</template>

<script>
export default {
  name: 'ZoomControls',
  props: {
    actions: {
      type: Array,
      default: null
    }
  },
  data() {
    return {
      btns: [
        {c: 'up', i: 'mdi-chevron-up'},
        {c: 'down', i: 'mdi-chevron-down'},
        {c: 'left', i: 'mdi-chevron-left'},
        {c: 'right', i: 'mdi-chevron-right'},
        {c: 'home', i: 'mdi-image-filter-center-focus'},
        {c: 'in', i: 'mdi-magnify-plus-outline'},
        {c: 'out', i: 'mdi-magnify-minus-outline'}
      ]
    };
  },
  computed: {
    activeBtns() {
      if (!this.actions) return this.btns;
      return this.btns.filter((b) => this.actions.includes(b.c));
    }
  },
  methods: {
    handleClick(name) {
      this.$emit(name);
    }
  }
};
</script>

<style scoped>
.root {
  display: grid;
  grid-gap: 8px;
  grid-template-columns: repeat(3, auto);
  grid-template-rows: repeat(2, auto) 8px repeat(3, auto);
  margin-bottom: 1em;
}

.v-btn {
  background: white;
}

.in,
.out,
.up,
.down,
.home {
  grid-column: 2 / span 1;
}

.left {
  grid-column: 1 / span 1;
}

.right {
  grid-column: 3 / span 1;
}

.in {
  grid-row: 1 / span 1;
}

.out {
  grid-row: 2 / span 1;
}

.up {
  grid-row: -4 / span 1;
}

.left,
.home,
.right {
  grid-row: -3 / span 1;
}

.down {
  grid-row: -2 / span 1;
}
</style>
