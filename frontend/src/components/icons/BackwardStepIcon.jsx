function BackwardStepIcon({ ...props }) {
  return (
    <svg
      {...props}
      aria-hidden="true"
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      fill="none"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="2"
        d="M8 6v12m8-12v12l-8-6 8-6Z"
      />
    </svg>
  );
}

export default BackwardStepIcon;
