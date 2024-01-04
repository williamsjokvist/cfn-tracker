export const cn = (...args: Array<undefined | null | string | boolean>) => (
  args
    .flat()
    .filter(x => typeof x === "string")
    .join(" ")
)