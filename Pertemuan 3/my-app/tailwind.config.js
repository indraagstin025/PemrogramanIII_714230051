import withMT from "@material-tailwind/react/utils/withMT.js";

/** @type {import('tailwindcss').Config} */
export default withMT({
  content: [
    "./index.html",
    "./src//*.{vue,js,ts,jsx,tsx}",
    "node_modules/@material-tailwind/react/components//*.{js,ts,jsx,tsx}",
    "node_modules/@material-tailwind/react/theme/components//*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
});