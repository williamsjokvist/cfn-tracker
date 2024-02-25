/** @type {import('prettier').Config} */
export default {
  printWidth: 100,
  arrowParens: 'avoid',
  jsxSingleQuote: true,
  singleQuote: true,
  tabWidth: 2,
  semi: false,
  trailingComma: 'none',
  plugins: ['prettier-plugin-tailwindcss'],
  tailwindFunctions: ['cn'],
  tailwindConfig: './.config/tailwind.config.js'
}
