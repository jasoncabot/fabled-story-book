/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./src/**/*.{html,js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        harlequin: {
          50: "#eeffe4",
          100: "#d8ffc4",
          200: "#b4ff90",
          300: "#83ff50",
          400: "#41ff00",
          500: "#35e600",
          600: "#26b800",
          700: "#1c8b00",
          800: "#1b6d07",
          900: "#195c0b",
          950: "#073400",
        },
      },
    },
  },
  future: {
    hoverOnlyWhenSupported: true,
  },
};
