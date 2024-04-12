/**
 * @param {ChoiceRange=} val
 * @return {string}
 */
export default function choiceRangeStr(val) {
  const itemValue = (item) => item?.value ?? item;
  const from = val?.from;
  const to = val?.to;
  if (!from && !to) return 'All';
  if (from && !to) return `${from.title} and above`;
  if (!from && to) return `${to.title} and below`;
  if (itemValue(from) === itemValue(to)) return `${from.title} only`;
  return `${from.title} to ${to.title}`;
}
