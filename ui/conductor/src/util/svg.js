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
