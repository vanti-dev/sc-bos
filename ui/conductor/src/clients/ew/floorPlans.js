// Collect the floor plans here so we can import them all at once
// and can set them depending on the floor selected in the UI

// ?raw required at the end of the import to get the raw svg data
// with this we eliminate the svg conversion to .vue component
// but still can access the svg data as string
import level0 from './level0.svg?raw';

export const floorPlans = {
  level0
};
