/** @type {import('prettier').Config} */
module.exports = {
  printWidth: 100,
  arrowParens: 'avoid',
  jsxSingleQuote: true,
  singleQuote: true,
  tabWidth: 2,
  semi: false,
  trailingComma: 'none',
  plugins: ['prettier-plugin-tailwindcss'],
  tailwindFunctions: ['cn'],
  tailwindConfig: 'tailwind.config.cjs'
}
