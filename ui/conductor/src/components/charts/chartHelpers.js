/**
 * This function generates the background color for the chart based on the segments
 *
 * @see https://www.chartjs.org/chartjs-plugin-annotation/latest/samples/box/basic.html
 * @param {Array<[number,number,string]>} segments
 * @param {string} conditionValue
 * @param {string} mainColor
 * @param {string} secondaryColor
 * @return {{
 *   annotations: Record<string, {
 *     backgroundColor: string,
 *     drawTime: string,
 *     type: string,
 *     xMin: number,
 *     xMax: number
 *   }>
 * }}
 */
export const generateAreaBackground = (segments, conditionValue, mainColor, secondaryColor) => {
  const annotations = {};

  segments.forEach((segment, index) => {
    const annotationKey = 'box' + index; // Generate a unique key for each annotation
    const annotation = {
      // Generate background color based on the segment type - given by the conditionValue
      backgroundColor: segment[2] === conditionValue ? mainColor : secondaryColor,
      // Set the draw time to after the dataset is drawn
      drawTime: 'beforeDatasetsDraw',
      //
      type: 'box', // https://www.chartjs.org/docs/latest/developers/annotations.html#types
      xMin: segment[0] - 0.5, // -0.5 to make it start at the beginning of the bar
      xMax: segment[1] + 0.5 // +0.5 to make it end at the end of the bar
    };

    annotations[annotationKey] = annotation;
  });

  return {
    annotations: {...annotations}
  };
};
