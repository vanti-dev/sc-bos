/**
 * Converts SVG coordinates to percentage values relative to the viewBox dimensions.
 *
 * @param {{x: number, y: number, width: number, height: number}} viewBox
 * @param {DOMRect} svgRect
 * @return {{x: number, y: number, width: number, height: number}}
 */
export function convertSVGToPercentage(viewBox, svgRect) {
  return {
    x: svgRect.x / viewBox.width,
    y: svgRect.y / viewBox.height,
    width: svgRect.width / viewBox.width,
    height: svgRect.height / viewBox.height
  };
}

/**
 * Returns the native width and height of the SVG root element.
 * If the root element has a width and height attribute, those values are returned, otherwise this is closer to the
 * client bounding rect.
 *
 * @param {SVGGraphicsElement} el
 * @return {{rootWidth: number, rootHeight: number}}
 */
export function svgRootSize(el) {
  const svg = el.ownerSVGElement;
  return {rootWidth: svg.width.baseVal.value, rootHeight: svg.height.baseVal.value};
}

/**
 * Returns the bounding box of an SVG element in the coordinate space of the SVG root.
 *
 * @param {SVGGraphicsElement} el
 * @return {DOMRect}
 */
export function elementBounds(el) {
  const bBox = el.getBBox({stroke: true, markers: true});
  const ctm = el.getCTM();

  // note: browsers don't currently support {stroke: true} for getBBox so we fake it
  const style = window.getComputedStyle(el);
  const strokeWidth = style.getPropertyValue('stroke-width');
  if (strokeWidth) {
    const stroke = parseFloat(strokeWidth);
    const halfStroke = stroke / 2;
    bBox.x -= halfStroke;
    bBox.y -= halfStroke;
    bBox.width += stroke;
    bBox.height += stroke;
  }

  return matrixTransformRect(bBox, ctm);
}

/**
 * @param {DOMRect} rect
 * @param {DOMMatrix} matrix
 * @return {DOMRect}
 */
function matrixTransformRect(rect, matrix) {
  const tl = new DOMPoint(rect.x, rect.y).matrixTransform(matrix);
  const tr = new DOMPoint(rect.x + rect.width, rect.y).matrixTransform(matrix);
  const bl = new DOMPoint(rect.x, rect.y + rect.height).matrixTransform(matrix);
  const br = new DOMPoint(rect.x + rect.width, rect.y + rect.height).matrixTransform(matrix);

  const minx = Math.min(tl.x, tr.x, bl.x, br.x);
  const miny = Math.min(tl.y, tr.y, bl.y, br.y);
  const maxx = Math.max(tl.x, tr.x, bl.x, br.x);
  const maxy = Math.max(tl.y, tr.y, bl.y, br.y);

  return new DOMRect(minx, miny, maxx - minx, maxy - miny);
}
