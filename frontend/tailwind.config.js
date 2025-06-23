/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    extend: {
      colors: {
        background: "#0B1014",
        "background-secondary": "#14181F",
        "background-light": "#1B2127",
        "accent-1": "#E7AF21",
        "accent-2": "#AF8417",
        text: "#FFFFFF",
        secondary: "#838E99",
      },
      padding: {
        form: "50px",
      },
    },
    borderRadius: {
      lg: "20px",
      md: "10px",
      sm: "7px",
      rounded: "100%",
    },
    fontSize: {
      sm: "16px",
      md: "20px",
      lg: "40px",
      "lg-1": "64px",
    },
    fontFamily: {
      sans: ["Onest", "sans-serif"],
    },
  },
  plugins: [require("@tailwindcss/forms")],
};
