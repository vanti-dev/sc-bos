<template>
  <div class="pinch-zoom" :class="{debug: debug, center}" v-bind="attrs">
    <div class="event-target" ref="el">
      <pan-zoom
          :options="panZoomOptions"
          class="zoomable"
          @init="panZoomInit"
          @change="updateTransform"
          @start="currentTransition = false"
          v-bind="listeners">
        <slot v-bind="{...currentTransform, transition: internalTransition}"/>
      </pan-zoom>
    </div>
    <div class="controls" v-if="!hideControls">
      <slot
          name="controls"
          v-bind="{home, zoomIn, zoomOut, panUp, panDown, panLeft, panRight, isHome: isDefaultTransform}">
        <zoom-controls
            :actions="zoomActions"
            @home="home"
            @in="zoomIn"
            @out="zoomOut"
            @up="panUp"
            @down="panDown"
            @left="panLeft"
            @right="panRight">
          <template #home>
            <v-icon>
              {{ isDefaultTransform ? 'mdi-image-filter-center-focus' : 'mdi-image-filter-center-focus-weak' }}
            </v-icon>
          </template>
        </zoom-controls>
      </slot>
    </div>
    <pre v-if="debug" class="marquee"/>
    <pre v-if="debug" class="marquee"/>
    <pre v-if="debug" class="marquee"/>
    <pre v-if="debug" class="marquee"/>
    <pre v-if="debug" class="marquee"/>
    <pre v-if="debug" class="marquee"/>
    <pre v-if="debug" class="marquee"/>
    <pre v-if="debug" class="marquee"/>
    <pre v-if="debug" class="marquee"/>
    <div class="debug" v-if="debug">
      <pre>currentTransform {{ JSON.stringify(currentTransform) }}</pre>
      <slot name="debug"/>
    </div>
  </div>
</template>

<script>
import ZoomControls from '@/components/zoom/ZoomControls.vue';
import PanZoom from '@/thirdparty/panzoom/PanZoom.vue';

export default {
  name: 'PinchZoom',
  components: {PanZoom, ZoomControls},
  // attrs target the root node, listeners for the pan-zoom component
  inheritAttrs: false,
  props: {
    scaleMax: {
      type: Number,
      default: 10
    },
    scaleMin: {
      type: Number,
      default: 0.1
    },
    defaultTransform: {
      type: Object,
      default() {
        return {x: 0, y: 0, scale: 1};
      }
    },
    transition: {
      type: Boolean,
      default: false
    },
    transform: {
      type: Object,
      default: null
    },
    scale: {
      type: Number,
      default: null
    },
    hideControls: {
      type: Boolean,
      default: false
    },
    disablePanZoom: Boolean,
    debug: {
      type: Boolean,
      default: false
    },
    center: Boolean
  },
  data() {
    return {
      panZoomInstance: null,
      internalTransform: this.transform || {x: 0, y: 0, scale: 1},
      internalTransition: this.transition
    };
  },
  computed: {
    currentTransform: {
      get() {
        return this.internalTransform;
      },
      set(transform) {
        if (this.transformIsEqual(this.internalTransform, transform)) return; // no change
        this.internalTransform = transform;
        // note: this is an optimisation for the dev-tools. Transforms update a lot, if no code is listening the
        // dev-tools still does which causes a lot of event churn
        if (this.hasListener('onUpdate:transform')) {
          this.$emit('update:transform', transform);
        }
      }
    },
    currentTransition: {
      get() {
        return this.internalTransition;
      },
      set(v) {
        if (this.internalTransition === v) return; // no change
        this.internalTransition = v;
        this.$emit('update:transition', v);
      }
    },
    panZoomOptions() {
      return {
        overflow: 'visible',
        canvas: true,
        step: 0.1,
        minDistance: 10,
        maxScale: this.scaleMax,
        minScale: this.scaleMin,
        startX: this.cappedDefaultTransform.x,
        startY: this.cappedDefaultTransform.y,
        startScale: this.cappedDefaultTransform.scale,
        noBind: this.disablePanZoom,
        handleStartEvent(e) {
          // WARNING! in browsers that don't support PointerEvents this can cause the click event to not fire for
          // anything that is inside the PinchZoom control. Use touchstart/end instead of click to receive events.
          e.preventDefault();
          // omit e.stopPropagation() so idle listeners work
        }
      };
    },
    cappedDefaultTransform() {
      const res = {...this.defaultTransform};
      res.scale = this.capZoom(res.scale);
      return res;
    },
    isDefaultTransform() {
      const t = this.cappedDefaultTransform;
      const c = this.currentTransform;
      return this.transformIsEqual(t, c);
    },
    zoomActions() {
      if (this.disablePanZoom) {
        return ['home'];
      } else {
        return undefined; // all
      }
    },
    listeners() {
      return Object.entries(this.$attrs).reduce((acc, [k, v]) => {
        if (k.startsWith('on')) {
          acc[k] = v;
        }
        return acc;
      }, {});
    },
    attrs() {
      return Object.entries(this.$attrs).reduce((acc, [k, v]) => {
        if (!k.startsWith('on')) {
          acc[k] = v;
        }
        return acc;
      }, {});
    }
  },
  watch: {
    defaultTransform() {
      this.reset();
    },
    panZoomInstance() {
      this.reset();
    },
    transform: {
      immediate: true,
      handler: 'updateTransform'
    },
    ['currentTransform.scale'](v) {
      if (this.hasListener('update:scale')) {
        this.$emit('update:scale', v);
      }
    },
    scale(scale) {
      if (scale === this.currentTransform.scale) return;
      this.currentTransform = {...this.currentTransform, scale};
    }
  },
  methods: {
    panZoomInit(panZoomInstance) {
      this.panZoomInstance = panZoomInstance;
    },
    updateTransform(t) {
      if (!t) return;
      this.currentTransform = {...this.currentTransform, ...t};
    },
    bounds() {
      return this.$refs.el.getBoundingClientRect();
    },
    reset() {
      const {x, y, scale} = this.cappedDefaultTransform;
      this.panZoomInstance.pan(x, y, {animate: true});
      this.panZoomInstance.zoom(scale, {animate: true});
    },

    panZoom(panZoom) {
      if (this.panZoomInstance) {
        let {x, y, scale} = panZoom;
        scale = this.capZoom(scale);
        this.currentTransition = true;
        this.panZoomInstance.pan(x, y, {animate: true});
        this.panZoomInstance.zoom(scale, {animate: true});
      }
    },

    home() {
      this.currentTransition = true;
      this.reset();
    },

    zoomIn() {
      if (this.panZoomInstance) {
        this.currentTransition = true;
        this.panZoomInstance.zoomIn();
      }
    },

    zoomOut() {
      if (this.panZoomInstance) {
        this.currentTransition = true;
        this.panZoomInstance.zoomOut();
      }
    },

    panLeft() {
      if (this.panZoomInstance) {
        this.currentTransition = true;
        this.panZoomInstance.pan(30, 0, {relative: true, animate: true});
      }
    },

    panRight() {
      if (this.panZoomInstance) {
        this.currentTransition = true;
        this.panZoomInstance.pan(-30, 0, {relative: true, animate: true});
      }
    },

    panUp() {
      if (this.panZoomInstance) {
        this.currentTransition = true;
        this.panZoomInstance.pan(0, 30, {relative: true, animate: true});
      }
    },

    panDown() {
      if (this.panZoomInstance) {
        this.currentTransition = true;
        this.panZoomInstance.pan(0, -30, {relative: true, animate: true});
      }
    },

    capZoom(zoom) {
      return Math.max(Math.min(this.scaleMax, zoom), this.scaleMin);
    },

    transformIsEqual(t1, t2) {
      const fpError = 0.001;
      return t1 && t2 &&
          Math.abs(t1.x - t2.x) < fpError &&
          Math.abs(t1.y - t2.y) < fpError &&
          Math.abs(t1.scale - t2.scale) < fpError;
    },
    hasListener(l) {
      return Boolean(this.$props?.['on' + l[0].toUpperCase() + l.slice(1)]);
    }
  }
};
</script>

<style scoped>
.pinch-zoom {
  position: relative;
  display: flex;
  justify-content: flex-end;
  align-items: flex-end;
}

.event-target {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
}

.zoomable {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  outline: none;
}

.pinch-zoom.center .zoomable {
  display: flex;
  justify-content: center;
  align-items: center;
}

.pinch-zoom.center .zoomable >>> > * {
  flex-shrink: 0;
}

.marquee {
  position: absolute;
  width: 100px;
  height: 100px;
  left: calc(50% - 50px);
  top: calc(50% - 50px);
  border: 1px dashed #0005;
  pointer-events: none;
}

.marquee:nth-of-type(1) {
  transform: scale(0.25);
}

.marquee:nth-of-type(2) {
  transform: scale(1);
}

.marquee:nth-of-type(3) {
  transform: scale(2);
}

.marquee:nth-of-type(4) {
  transform: scale(3);
}

.marquee:nth-of-type(5) {
  transform: scale(4);
}

.marquee:nth-of-type(6) {
  transform: scale(5);
}

.marquee:nth-of-type(7) {
  transform: scale(6);
}

.marquee:nth-of-type(8) {
  transform: scale(7);
}

.marquee:nth-of-type(9) {
  transform: scale(8);
}

.pinch-zoom.debug {
  border: 1px dotted black;
}

.pinch-zoom > .debug {
  position: absolute;
  top: 0;
  left: 0;
  background: #ffe5becc;
  border: 1px solid rgba(255, 217, 152, 0.8);
  border-radius: 3px;
  padding: 3px 12px;
}
</style>
